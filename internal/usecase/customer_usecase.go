package usecase

import (
	"context"
	"sync"
	"time"

	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/adapter"
	"github.com/zuxt268/homing/internal/interface/dto/external"
	"github.com/zuxt268/homing/internal/interface/dto/model"
	"github.com/zuxt268/homing/internal/interface/dto/req"
	"github.com/zuxt268/homing/internal/interface/dto/res"
	"github.com/zuxt268/homing/internal/interface/repository"
	"github.com/zuxt268/homing/internal/interface/util"
)

type CustomerUsecase interface {
	SyncAll(ctx context.Context) error

	GetWordpressInstagramList(ctx context.Context, params req.GetWordpressInstagram) (*res.WordpressInstagramList, error)
	GetWordpressInstagram(ctx context.Context, id int) (*res.WordpressInstagram, error)
	CreateWordpressInstagram(ctx context.Context, body req.CreateWordpressInstagram) (*res.WordpressInstagram, error)
	UpdateWordpressInstagram(ctx context.Context, body req.UpdateWordpressInstagram) (*res.WordpressInstagram, error)
	DeleteWordpressInstagram(ctx context.Context, id int) error

	GetToken(ctx context.Context) (*res.Token, error)
	UpdateToken(ctx context.Context, body req.UpdateToken) error
}

type customerUsecase struct {
	fileDownloader         adapter.FileDownloader
	instagramAdapter       adapter.InstagramAdapter
	slack                  adapter.Slack
	wordpressAdapter       adapter.WordpressAdapter
	postRepo               repository.PostRepository
	wordpressInstagramRepo repository.WordpressInstagramRepository
	tokenRepo              repository.TokenRepository
}

func NewCustomerUsecase(
	fileDownloader adapter.FileDownloader,
	instagramAdapter adapter.InstagramAdapter,
	slack adapter.Slack,
	wordpressAdapter adapter.WordpressAdapter,
	postRepo repository.PostRepository,
	wordpressInstagramRepo repository.WordpressInstagramRepository,
	tokenRepo repository.TokenRepository,
) CustomerUsecase {
	return &customerUsecase{
		fileDownloader:         fileDownloader,
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
	wiList, err := u.wordpressInstagramRepo.FindAll(ctx, repository.WordpressInstagramFilter{})
	if err != nil {
		return err
	}

	// 20件の並列処理
	semaphore := make(chan struct{}, 20)
	var wg sync.WaitGroup

	for _, wi := range wiList {
		wg.Add(1)
		semaphore <- struct{}{} // セマフォを取得

		go func(wi *domain.WordpressInstagram) {
			defer wg.Done()
			defer func() { <-semaphore }() // セマフォを解放

			err := u.syncOne(ctx, wi)
			if err != nil {
				_ = u.slack.Alert(ctx, err.Error(), *wi)
			}
		}(wi)
	}

	wg.Wait()
	return nil
}

func (u *customerUsecase) syncOne(ctx context.Context, wi *domain.WordpressInstagram) error {

	/*
		トークンを取得する
	*/
	token, err := u.tokenRepo.First(ctx)
	if err != nil {
		return err
	}
	/*
		インスタグラムから投稿を一覧で取得する
	*/
	posts, err := u.instagramAdapter.GetPosts(ctx, token, wi.InstagramID)
	if err != nil {
		return err
	}

	/*
		まだ連携していない投稿をWordpressに連携する
	*/
	for _, post := range posts {
		err := u.transfer(ctx, wi, post)
		if err != nil {
			return err
		}
	}

	/*
		tempディレクトリを削除
	*/
	return u.fileDownloader.DeleteTempDirectory()
}

func (u *customerUsecase) transfer(ctx context.Context, wi *domain.WordpressInstagram, post domain.InstagramPost) error {

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
	localPath, err := u.fileDownloader.Download(ctx, post.MediaURL)
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

	/*
		アップロードしたファイルをFeaturedに指定して、記事を投稿
	*/
	postResp, err := u.wordpressAdapter.Post(ctx, external.WordpressPostInput{
		WordpressInstagram: *wi,
		FeaturedMediaID:    uploadResp.Id,
		Post:               post,
	})
	if err != nil {
		return err
	}

	/*
		投稿したことをDBに保存
	*/
	err = u.postRepo.SavePost(ctx, &model.Post{
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

func (u *customerUsecase) GetWordpressInstagramList(ctx context.Context, params req.GetWordpressInstagram) (*res.WordpressInstagramList, error) {
	filter := repository.WordpressInstagramFilter{
		Name:         params.Name,
		Wordpress:    params.Wordpress,
		InstagramID:  params.InstagramID,
		Status:       params.Status,
		DeleteHash:   params.DeleteHash,
		CustomerType: params.CustomerType,
	}

	wiList, err := u.wordpressInstagramRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	result := make([]res.WordpressInstagram, 0, len(wiList))
	for _, wi := range wiList {
		result = append(result, res.WordpressInstagram{
			ID:           wi.ID,
			Name:         wi.Name,
			Wordpress:    wi.Wordpress,
			InstagramID:  wi.InstagramID,
			Memo:         wi.Memo,
			StartDate:    wi.StartDate,
			Status:       int(wi.Status),
			DeleteHash:   wi.DeleteHash,
			CustomerType: int(wi.CustomerType),
		})
	}

	return &res.WordpressInstagramList{
		WordpressInstagramList: result,
	}, nil
}

func (u *customerUsecase) GetWordpressInstagram(ctx context.Context, id int) (*res.WordpressInstagram, error) {
	wi, err := u.wordpressInstagramRepo.Get(ctx, repository.WordpressInstagramFilter{
		ID: &id,
	})
	if err != nil {
		return nil, err
	}

	return &res.WordpressInstagram{
		ID:           wi.ID,
		Name:         wi.Name,
		Wordpress:    wi.Wordpress,
		InstagramID:  wi.InstagramID,
		Memo:         wi.Memo,
		StartDate:    wi.StartDate,
		Status:       int(wi.Status),
		DeleteHash:   wi.DeleteHash,
		CustomerType: int(wi.CustomerType),
	}, nil
}

func (u *customerUsecase) CreateWordpressInstagram(ctx context.Context, req req.CreateWordpressInstagram) (*res.WordpressInstagram, error) {
	wi := &domain.WordpressInstagram{
		Name:         req.Name,
		Wordpress:    req.Wordpress,
		InstagramID:  req.InstagramID,
		Memo:         req.Memo,
		StartDate:    req.StartDate,
		Status:       domain.Status(req.Status),
		DeleteHash:   req.DeleteHash,
		CustomerType: domain.CustomerType(req.CustomerType),
	}

	err := u.wordpressInstagramRepo.Save(ctx, wi)
	if err != nil {
		return nil, err
	}

	return &res.WordpressInstagram{
		ID:           wi.ID,
		Name:         wi.Name,
		Wordpress:    wi.Wordpress,
		InstagramID:  wi.InstagramID,
		Memo:         wi.Memo,
		StartDate:    wi.StartDate,
		Status:       int(wi.Status),
		DeleteHash:   wi.DeleteHash,
		CustomerType: int(wi.CustomerType),
	}, nil
}

func (u *customerUsecase) UpdateWordpressInstagram(ctx context.Context, req req.UpdateWordpressInstagram) (*res.WordpressInstagram, error) {
	wi, err := u.wordpressInstagramRepo.Get(ctx, repository.WordpressInstagramFilter{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}

	if req.Name != nil {
		wi.Name = *req.Name
	}
	if req.Wordpress != nil {
		wi.Wordpress = *req.Wordpress
	}
	if req.InstagramID != nil {
		wi.InstagramID = *req.InstagramID
	}
	if req.Memo != nil {
		wi.Memo = *req.Memo
	}
	if req.StartDate != nil {
		wi.StartDate = *req.StartDate
	}
	if req.Status != nil {
		wi.Status = domain.Status(*req.Status)
	}
	if req.DeleteHash != nil {
		wi.DeleteHash = *req.DeleteHash
	}
	if req.CustomerType != nil {
		wi.CustomerType = domain.CustomerType(*req.CustomerType)
	}

	err = u.wordpressInstagramRepo.Save(ctx, wi)
	if err != nil {
		return nil, err
	}

	return &res.WordpressInstagram{
		ID:           wi.ID,
		Name:         wi.Name,
		Wordpress:    wi.Wordpress,
		InstagramID:  wi.InstagramID,
		Memo:         wi.Memo,
		StartDate:    wi.StartDate,
		Status:       int(wi.Status),
		DeleteHash:   wi.DeleteHash,
		CustomerType: int(wi.CustomerType),
	}, nil
}

func (u *customerUsecase) DeleteWordpressInstagram(ctx context.Context, id int) error {
	return u.wordpressInstagramRepo.Delete(ctx, repository.WordpressInstagramFilter{
		ID: &id,
	})
}

func (u *customerUsecase) GetToken(ctx context.Context) (*res.Token, error) {
	token, err := u.tokenRepo.First(ctx)
	if err != nil {
		return nil, err
	}
	debug, err := u.instagramAdapter.DebugToken(ctx, token)
	if err != nil {
		return nil, err
	}
	expiredAt := time.Unix(debug.Data.ExpiresAt, 0)
	tenDaysLater := time.Now().AddDate(0, 0, 10)

	if tenDaysLater.After(expiredAt) {
		_ = u.slack.SendTokenExpired(ctx)
	}

	return &res.Token{
		Token:    token,
		ExpireAt: expiredAt,
	}, nil
}

func (u *customerUsecase) UpdateToken(ctx context.Context, req req.UpdateToken) error {
	return u.tokenRepo.DeleteInsert(ctx, req.Token)
}
