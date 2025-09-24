package usecase

import (
	"context"
	"time"

	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/adapter"
	"github.com/zuxt268/homing/internal/interface/dto/external"
	"github.com/zuxt268/homing/internal/interface/dto/model"
	"github.com/zuxt268/homing/internal/interface/repository"
)

type CustomerUsecase interface {
	SyncAll(ctx context.Context) error
	SyncOne(ctx context.Context, customerID int) error
}

type customerUsecase struct {
	fileDownloader   adapter.FileDownloader
	instagramAdapter adapter.InstagramAdapter
	slack            adapter.Slack
	wordpressAdapter adapter.WordpressAdapter
	customerRepo     repository.CustomerRepository
	postRepo         repository.PostRepository
}

func NewCustomerUsecase(
	fileDownloader adapter.FileDownloader,
	instagramAdapter adapter.InstagramAdapter,
	slack adapter.Slack,
	wordpressAdapter adapter.WordpressAdapter,
	customerRepo repository.CustomerRepository,
	postRepo repository.PostRepository,
) CustomerUsecase {
	return &customerUsecase{
		fileDownloader:   fileDownloader,
		instagramAdapter: instagramAdapter,
		slack:            slack,
		wordpressAdapter: wordpressAdapter,
		customerRepo:     customerRepo,
		postRepo:         postRepo,
	}
}

const template = `<@U04P797HYPM>
[%s]
顧客 id=%d, name=%s`

func (u *customerUsecase) SyncAll(ctx context.Context) error {
	customers, err := u.customerRepo.FindAllCustomers(ctx, repository.CustomerFilter{})
	if err != nil {
		return err
	}
	for _, customer := range customers {
		err := u.syncOne(ctx, customer)
		if err != nil {
			_ = u.slack.Alert(ctx, err.Error(), *customer)
		}
	}
	return nil
}

func (u *customerUsecase) SyncOne(ctx context.Context, customerID int) error {
	customer, err := u.customerRepo.GetCustomer(ctx, customerID)
	if err != nil {
		return err
	}
	err = u.syncOne(ctx, customer)
	if err != nil {
		_ = u.slack.Alert(ctx, err.Error(), *customer)
		return err
	}
	return nil
}

func (u *customerUsecase) syncOne(ctx context.Context, customer *domain.Customer) error {
	/*
		インスタグラムから投稿を一覧で取得する
	*/
	posts, err := u.instagramAdapter.GetPosts(ctx, customer.AccessToken, customer.InstagramAccountID)
	if err != nil {
		return err
	}

	/*
		まだ連携していない投稿をWordpressに連携する
	*/
	for _, post := range posts {
		err := u.transfer(ctx, customer, post)
		if err != nil {
			return err
		}
	}

	/*
		tempディレクトリを削除
	*/
	return u.fileDownloader.DeleteTempDirectory()
}

func (u *customerUsecase) transfer(ctx context.Context, customer *domain.Customer, post domain.InstagramPost) error {

	/*
		すでに投稿しているものかどうかをチェック
	*/
	exist, err := u.postRepo.ExistPost(ctx, repository.PostFilter{
		MediaID: &post.ID,
	})
	if err != nil {
		return err
	}
	if exist {
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
		Path:     localPath,
		Customer: *customer,
	})
	if err != nil {
		return err
	}

	/*
		アップロードしたファイルをFeaturedに指定して、記事を投稿
	*/
	postResp, err := u.wordpressAdapter.Post(ctx, external.WordpressPostInput{
		Customer:        *customer,
		FeaturedMediaID: uploadResp.Id,
		Post:            post,
	})
	if err != nil {
		return err
	}

	/*
		投稿したことをDBに保存
	*/
	err = u.postRepo.SavePost(ctx, &model.Post{
		MediaID:       post.ID,
		CustomerID:    customer.ID,
		Timestamp:     post.Timestamp,
		MediaURL:      post.MediaURL,
		CreatedAt:     time.Now(),
		Permalink:     post.Permalink,
		WordpressLink: postResp.WordpressURL,
	})
	if err != nil {
		return err
	}
	return nil
}
