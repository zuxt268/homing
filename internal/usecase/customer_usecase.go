package usecase

import (
	"context"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/adapter"
	"github.com/zuxt268/homing/internal/interface/dto/external"
	"github.com/zuxt268/homing/internal/interface/dto/model"
	"github.com/zuxt268/homing/internal/interface/repository"
	"github.com/zuxt268/homing/internal/interface/util"
)

type CustomerUsecase interface {
	SyncAll(ctx context.Context) error
	SyncOne(ctx context.Context, id int) error
}

type customerUsecase struct {
	instagramAdapter       adapter.InstagramAdapter
	slack                  adapter.Slack
	wordpressAdapter       adapter.WordpressAdapter
	postRepo               repository.PostRepository
	wordpressInstagramRepo repository.WordpressInstagramRepository
	tokenRepo              repository.TokenRepository
	customerLocks          sync.Map
}

func NewCustomerUsecase(
	instagramAdapter adapter.InstagramAdapter,
	slack adapter.Slack,
	wordpressAdapter adapter.WordpressAdapter,
	postRepo repository.PostRepository,
	wordpressInstagramRepo repository.WordpressInstagramRepository,
	tokenRepo repository.TokenRepository,
) CustomerUsecase {
	return &customerUsecase{
		instagramAdapter:       instagramAdapter,
		slack:                  slack,
		wordpressAdapter:       wordpressAdapter,
		postRepo:               postRepo,
		wordpressInstagramRepo: wordpressInstagramRepo,
		tokenRepo:              tokenRepo,
	}
}

const template = `<@U04P797HYPM>
[%s]
顧客 id=%d, name=%s`

func (u *customerUsecase) SyncAll(ctx context.Context) error {
	wiList, err := u.wordpressInstagramRepo.FindAll(ctx, repository.WordpressInstagramFilter{
		Status: util.Pointer(1),
	})
	if err != nil {
		return err
	}

	// 20件の並列処理
	semaphore := make(chan struct{}, 20)
	var wg sync.WaitGroup

	for _, wi := range wiList {
		wg.Add(1)
		semaphore <- struct{}{}

		go func(wi *domain.WordpressInstagram) {
			defer wg.Done()
			defer func() { <-semaphore }()

			// 各goroutine専用のFileDownloaderを作成
			fd := adapter.NewFileDownloader()
			defer func() {
				_ = fd.DeleteTempDirectory()
			}()

			err := u.syncOne(ctx, wi, fd)
			if err != nil {
				_ = u.slack.Alert(ctx, err.Error(), *wi)
			}
		}(wi)
	}

	wg.Wait()
	return nil
}

func (u *customerUsecase) syncOne(ctx context.Context, wi *domain.WordpressInstagram, fd adapter.FileDownloader) error {
	// 顧客IDごとのロックを取得
	lockInterface, _ := u.customerLocks.LoadOrStore(wi.ID, &sync.Mutex{})
	mu := lockInterface.(*sync.Mutex)

	mu.Lock()
	defer mu.Unlock()

	backGroundCtx := context.Background()

	go func() {
		/*
			トークンを取得する
		*/
		token, err := u.tokenRepo.First(backGroundCtx)
		if err != nil {
			_ = u.slack.Alert(backGroundCtx, err.Error(), *wi)
		}
		/*
			インスタグラムから投稿を一覧で取得する
		*/
		posts, err := u.instagramAdapter.GetPostsAll(backGroundCtx, token, wi.InstagramID)
		if err != nil {
			_ = u.slack.Alert(backGroundCtx, err.Error(), *wi)
		}

		/*
			まだ連携していない投稿をWordpressに連携する
		*/
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Timestamp < posts[j].Timestamp
		})
		for _, post := range posts {
			err := u.transfer(backGroundCtx, wi, post, fd)
			if err != nil {
				_ = u.slack.Alert(backGroundCtx, err.Error(), *wi)
			}
		}
	}()

	return nil
}

func (u *customerUsecase) transfer(ctx context.Context, wi *domain.WordpressInstagram, post domain.InstagramPost, fd adapter.FileDownloader) error {

	/*
		メディアのリンクがない場合はスキップ
	*/
	if post.MediaURL == "" {
		return nil
	}

	/*
		すでに投稿しているものかどうかをチェック
	*/
	exist, err := u.postRepo.ExistPost(ctx, repository.PostFilter{
		CustomerID: util.Pointer(100000 + wi.ID),
		MediaID:    &post.ID,
	})
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	/*
		連携開始日前のデータは連携しない
	*/
	instagramPost, _ := time.Parse("2006-01-02T15:04:05-0700", post.Timestamp)
	if instagramPost.Before(wi.StartDate) {
		return nil
	}

	/*
		インスタグラムの投稿の画像、動画を一時ディレクトリにダウンロード
	*/
	localPath, err := fd.Download(ctx, post.MediaURL)
	if err != nil {
		return err
	}

	/*
		ダウンロードしたファイルをWordpressにアップロード
	*/
	uploadResp, err := u.wordpressAdapter.FileUpload(ctx, external.WordpressFileUploadInput{
		Path:               localPath,
		WordpressInstagram: *wi,
	})
	if err != nil {
		return err
	}
	post.SetFeaturedMediaID(uploadResp.Id)
	post.AppendSourceURL(uploadResp.SourceUrl)
	post.SetDeleteHashFlag(wi.DeleteHash)

	for _, child := range post.Children {
		/*
			インスタグラムの投稿の画像、動画を一時ディレクトリにダウンロード
		*/
		childLocalPath, err := fd.Download(ctx, child.MediaURL)
		if err != nil {
			return err
		}

		/*
			ダウンロードしたファイルをWordpressにアップロード
		*/
		childUploadResp, err := u.wordpressAdapter.FileUpload(ctx, external.WordpressFileUploadInput{
			Path:               childLocalPath,
			WordpressInstagram: *wi,
		})
		if err != nil {
			return err
		}
		post.AppendSourceURL(childUploadResp.SourceUrl)
	}

	/*
		アップロードしたファイルをFeaturedに指定して、記事を投稿
	*/
	postResp, err := u.wordpressAdapter.Post(ctx, external.WordpressPostInput{
		WordpressInstagram: *wi,
		Post:               post,
	})
	if err != nil {
		return err
	}

	/*
		投稿したことをDBに保存
	*/
	err = u.postRepo.CreatePost(ctx, &model.Post{
		MediaID:       post.ID,
		CustomerID:    100000 + wi.ID,
		Timestamp:     post.Timestamp,
		MediaURL:      post.MediaURL,
		Permalink:     post.Permalink,
		WordpressLink: postResp.WordpressURL,
		CreatedAt:     time.Now(),
	})
	if err != nil {
		return err
	}

	/*
		Slackに通知
	*/
	_ = u.slack.Success(ctx, wi, postResp.WordpressURL, post.Permalink)

	/*
		ダウンロードファイルを都度削除
	*/
	_ = os.Remove(localPath)

	return nil
}

func (u *customerUsecase) SyncOne(ctx context.Context, id int) error {
	wi, err := u.wordpressInstagramRepo.Get(ctx, repository.WordpressInstagramFilter{
		ID: util.Pointer(id),
	})
	if err != nil {
		return err
	}

	fd := adapter.NewFileDownloader()
	defer func() {
		_ = fd.DeleteTempDirectory()
	}()

	err = u.syncOne(ctx, wi, fd)
	if err != nil {
		_ = u.slack.Alert(ctx, err.Error(), *wi)
		return err
	}
	return nil
}
