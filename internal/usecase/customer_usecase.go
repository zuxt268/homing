package usecase

import (
	"context"
	"log/slog"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/zuxt268/homing/internal/config"
	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/adapter"
	"github.com/zuxt268/homing/internal/interface/dto/external"
	"github.com/zuxt268/homing/internal/interface/dto/model"
	"github.com/zuxt268/homing/internal/interface/repository"
	"github.com/zuxt268/homing/internal/interface/util"
)

type CustomerUsecase interface {
	SyncAllWordpressInstagram(ctx context.Context) error
	SyncOneWordpressInstagram(ctx context.Context, id int) error

	SyncAllGoogleBusinessInstagram(ctx context.Context) error
	SyncOneGoogleBusinessInstagram(ctx context.Context, id int) error
}

type customerUsecase struct {
	instagramAdapter       adapter.InstagramAdapter
	slack                  adapter.Slack
	wordpressAdapter       adapter.WordpressAdapter
	gbpAdapter             adapter.GbpAdapter
	postRepo               repository.PostRepository
	wordpressInstagramRepo repository.WordpressInstagramRepository
	tokenRepo              repository.TokenRepository
	businessInstagramRepo  repository.BusinessInstagramRepository
	googlePostRepo         repository.GooglePostRepository
	customerLocks          sync.Map
}

func NewCustomerUsecase(
	instagramAdapter adapter.InstagramAdapter,
	slack adapter.Slack,
	wordpressAdapter adapter.WordpressAdapter,
	gbpAdapter adapter.GbpAdapter,
	postRepo repository.PostRepository,
	wordpressInstagramRepo repository.WordpressInstagramRepository,
	tokenRepo repository.TokenRepository,
	businessInstagramRepo repository.BusinessInstagramRepository,
	googlePostRepo repository.GooglePostRepository,
) CustomerUsecase {
	return &customerUsecase{
		instagramAdapter:       instagramAdapter,
		slack:                  slack,
		wordpressAdapter:       wordpressAdapter,
		gbpAdapter:             gbpAdapter,
		postRepo:               postRepo,
		wordpressInstagramRepo: wordpressInstagramRepo,
		tokenRepo:              tokenRepo,
		businessInstagramRepo:  businessInstagramRepo,
		googlePostRepo:         googlePostRepo,
	}
}

const template = `<@U04P797HYPM>
[%s]
顧客 id=%d, name=%s`

func (u *customerUsecase) SyncAllWordpressInstagram(ctx context.Context) error {
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

			u.syncOne(ctx, wi, fd)
		}(wi)
	}

	wg.Wait()
	return nil
}

func (u *customerUsecase) syncOne(ctx context.Context, wi *domain.WordpressInstagram, fd adapter.FileDownloader) {
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
			_ = u.slack.Error(ctx, "instagram => wordpress", err, wi.ID, wi.Name)
			return
		}
		/*
			インスタグラムから投稿を一覧で取得する
		*/
		posts, err := u.instagramAdapter.GetPostsAll(backGroundCtx, token, wi.InstagramID)
		if err != nil {
			_ = u.slack.Error(ctx, "instagram => wordpress", err, wi.ID, wi.Name)
			return
		}

		/*
			まだ連携していない投稿をWordpressに連携する
		*/
		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Timestamp < posts[j].Timestamp
		})
		for _, post := range posts {
			err := u.instagram2wordpress(backGroundCtx, wi, post, fd)
			if err != nil {
				_ = u.slack.Error(ctx, "instagram => wordpress", err, wi.ID, wi.Name)
				return
			}
		}
	}()
}

func (u *customerUsecase) instagram2wordpress(ctx context.Context, wi *domain.WordpressInstagram, post domain.InstagramPost, fd adapter.FileDownloader) error {

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

	if len(post.Children) == 0 {
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
		post.SetDeleteHashFlag(wi.DeleteHash)
		post.AppendSourceURL(uploadResp.SourceUrl)

		/*
			ダウンロードファイルを都度削除
		*/
		err = os.Remove(localPath)
		if err != nil {
			slog.Warn(err.Error())
		}

	} else {
		for i, child := range post.Children {
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
			if i == 0 {
				post.SetFeaturedMediaID(childUploadResp.Id)
				post.SetDeleteHashFlag(wi.DeleteHash)
			}
			post.AppendSourceURL(childUploadResp.SourceUrl)

			/*
				ダウンロードファイルを都度削除
			*/
			err = os.Remove(childLocalPath)
			if err != nil {
				slog.Warn(err.Error())
			}
		}
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

	return nil
}

func (u *customerUsecase) SyncOneWordpressInstagram(ctx context.Context, id int) error {
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

	u.syncOne(ctx, wi, fd)
	return nil
}

func (u *customerUsecase) SyncAllGoogleBusinessInstagram(ctx context.Context) error {
	biList, err := u.businessInstagramRepo.FindAll(ctx, repository.BusinessInstagramFilter{
		Status: util.Pointer(1),
	})
	if err != nil {
		return err
	}

	fd := adapter.NewFileDownloader()
	defer func() {
		_ = fd.DeleteTempDirectory()
	}()

	backGroundCtx := context.Background()
	for _, bi := range biList {
		token, err := u.tokenRepo.First(backGroundCtx)
		if err != nil {
			_ = u.slack.Error(ctx, "instagram => google business profile", err, bi.ID, bi.BusinessTitle)
			continue
		}
		/*
			インスタグラムから投稿を一覧で取得する
		*/
		posts, err := u.instagramAdapter.GetPosts25(backGroundCtx, token, bi.InstagramID)
		if err != nil {
			_ = u.slack.Error(ctx, "instagram => google business profile", err, bi.ID, bi.BusinessTitle)
			continue
		}

		sort.Slice(posts, func(i, j int) bool {
			return posts[i].Timestamp < posts[j].Timestamp
		})

		for _, post := range posts {
			if err := u.instagramToGbp(backGroundCtx, bi, post, fd); err != nil {
				_ = u.slack.Error(ctx, "instagram => google business profile", err, bi.ID, bi.BusinessTitle)
				continue
			}
		}
	}
	return nil
}

func (u *customerUsecase) SyncOneGoogleBusinessInstagram(ctx context.Context, id int) error {
	bi, err := u.businessInstagramRepo.Get(ctx, repository.BusinessInstagramFilter{
		ID:     util.Pointer(id),
		Status: util.Pointer(1),
	})
	if err != nil {
		return err
	}

	fd := adapter.NewFileDownloader()
	defer func() {
		_ = fd.DeleteTempDirectory()
	}()

	backGroundCtx := context.Background()
	token, err := u.tokenRepo.First(backGroundCtx)
	if err != nil {
		return err
	}
	/*
		インスタグラムから投稿を一覧で取得する
	*/
	posts, err := u.instagramAdapter.GetPostsAll(backGroundCtx, token, bi.InstagramID)
	if err != nil {
		return err
	}

	sort.Slice(posts, func(i, j int) bool {
		return posts[i].Timestamp < posts[j].Timestamp
	})

	for _, post := range posts {
		if err := u.instagramToGbp(backGroundCtx, bi, post, fd); err != nil {
			return err
		}
	}

	return nil
}

func (u *customerUsecase) instagramToGbp(ctx context.Context, bi *domain.BusinessInstagram, post domain.InstagramPost, fd adapter.FileDownloader) error {

	/*
		メディアのリンクがない場合はスキップ
	*/
	if post.MediaURL == "" {
		return nil
	}

	/*
		連携開始日前のデータは連携しない
	*/
	instagramPost, _ := time.Parse("2006-01-02T15:04:05-0700", post.Timestamp)
	if instagramPost.Before(bi.StartDate) {
		return nil
	}

	/*
		すでに投稿しているものかどうかをチェック
	*/
	exist, err := u.googlePostRepo.Exists(ctx, repository.GooglePostFilter{
		MediaID:    &post.ID,
		CustomerID: &bi.ID,
	})
	if err != nil {
		return err
	}
	if exist {
		return nil
	}

	if len(post.Children) == 0 {

		/*
			インスタグラムの投稿の画像、動画を一時ディレクトリにダウンロード
		*/
		localPath, err := fd.Download(ctx, post.MediaURL)
		if err != nil {
			return err
		}
		/*
			ダウンロードしたファイルをGoogleBusinessにアップロード
		*/
		uploadResp, err := u.gbpAdapter.UploadMedia(ctx, config.Env.GoogleBusinessAccountName, bi.BusinessName, localPath)
		if err != nil {
			return err
		}

		/*
			ダウンロードファイルを都度削除
		*/
		err = os.Remove(localPath)
		if err != nil {
			slog.Warn(err.Error())
		}

		/*
			投稿したことをDBに保存
		*/
		err = u.googlePostRepo.Create(ctx, &domain.GooglePost{
			GoogleBusinessURL: bi.BusinessName,
			InstagramURL:      post.Permalink,
			MediaID:           post.ID,
			CustomerID:        bi.ID,
			Name:              uploadResp.Name,
			MediaFormat:       uploadResp.MediaFormat,
			GoogleURL:         uploadResp.GoogleURL,
			CreateTime:        uploadResp.CreateTime,
		})
		if err != nil {
			return err
		}

	} else {
		for _, child := range post.Children {
			/*
				すでに投稿しているものかどうかをチェック（子要素単位）
			*/
			childExist, err := u.googlePostRepo.Exists(ctx, repository.GooglePostFilter{
				MediaID:    &child.ID,
				CustomerID: &bi.ID,
			})
			if err != nil {
				return err
			}
			if childExist {
				continue
			}

			/*
				インスタグラムの投稿の画像、動画を一時ディレクトリにダウンロード
			*/
			childLocalPath, err := fd.Download(ctx, child.MediaURL)
			if err != nil {
				return err
			}

			/*
				ダウンロードしたファイルをGoogleBusinessにアップロード
			*/
			uploadResp, err := u.gbpAdapter.UploadMedia(ctx, config.Env.GoogleBusinessAccountName, bi.BusinessName, childLocalPath)
			if err != nil {
				return err
			}

			/*
				ダウンロードファイルを都度削除
			*/
			err = os.Remove(childLocalPath)
			if err != nil {
				slog.Warn(err.Error())
			}

			/*
				投稿したことをDBに保存
			*/
			err = u.googlePostRepo.Create(ctx, &domain.GooglePost{
				GoogleBusinessURL: bi.BusinessName,
				InstagramURL:      post.Permalink,
				MediaID:           child.ID,
				CustomerID:        bi.ID,
				Name:              uploadResp.Name,
				MediaFormat:       uploadResp.MediaFormat,
				GoogleURL:         uploadResp.GoogleURL,
				CreateTime:        uploadResp.CreateTime,
			})
			if err != nil {
				return err
			}
		}
	}

	return nil
}
