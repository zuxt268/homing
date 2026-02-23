package usecase

import (
	"context"
	"time"

	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/adapter"
	"github.com/zuxt268/homing/internal/interface/dto/req"
	"github.com/zuxt268/homing/internal/interface/dto/res"
	"github.com/zuxt268/homing/internal/interface/repository"
	"github.com/zuxt268/homing/internal/interface/util"
)

type WordpressGbpUsecase interface {
	GetWordpressGbpList(ctx context.Context, params req.GetWordpressGbp) (*res.WordpressGbpList, error)
	GetWordpressGbp(ctx context.Context, id int, params req.GetWordpressGbpDetail) (*res.WordpressGbpDetail, error)
	CreateWordpressGbp(ctx context.Context, body req.WordpressGbp) (*res.WordpressGbp, error)
	UpdateWordpressGbp(ctx context.Context, id int, body req.WordpressGbp) (*res.WordpressGbp, error)
	DeleteWordpressGbp(ctx context.Context, id int) error
}

type wordpressGbpUsecase struct {
	wordpressGbpRepo repository.WordpressGbpRepository
	googlePostRepo   repository.GooglePostRepository
	wordpressAdapter adapter.WordpressAdapter
	gbpAdapter       adapter.GbpAdapter
}

func NewWordpressGbpUsecase(
	wordpressGbpRepo repository.WordpressGbpRepository,
	googlePostRepo repository.GooglePostRepository,
	wordpressAdapter adapter.WordpressAdapter,
	gbpAdapter adapter.GbpAdapter,
) WordpressGbpUsecase {
	return &wordpressGbpUsecase{
		wordpressGbpRepo: wordpressGbpRepo,
		googlePostRepo:   googlePostRepo,
		wordpressAdapter: wordpressAdapter,
		gbpAdapter:       gbpAdapter,
	}
}

func (u *wordpressGbpUsecase) GetWordpressGbpList(ctx context.Context, params req.GetWordpressGbp) (*res.WordpressGbpList, error) {
	wgList, err := u.wordpressGbpRepo.FindAll(ctx, repository.WordpressGbpFilter{
		Limit:       params.Limit,
		Offset:      params.Offset,
		PartialName: params.Name,
		Status:      params.Status,
	})
	if err != nil {
		return nil, err
	}
	total, err := u.wordpressGbpRepo.Count(ctx, repository.WordpressGbpFilter{
		All: util.Pointer(true),
	})
	if err != nil {
		return nil, err
	}

	resWordpressGbp := make([]res.WordpressGbp, len(wgList))
	for i, wg := range wgList {
		resWordpressGbp[i] = res.WordpressGbp{
			ID:              wg.ID,
			Name:            wg.Name,
			WordpressDomain: wg.WordpressDomain,
			BusinessName:    wg.BusinessName,
			BusinessTitle:   wg.BusinessTitle,
			Memo:            wg.Memo,
			MapsURL:         wg.MapsURL,
			StartDate:       wg.StartDate,
			Status:          int(wg.Status),
			CreatedAt:       wg.CreatedAt,
			UpdatedAt:       wg.UpdatedAt,
		}
	}
	return &res.WordpressGbpList{
		WordpressGbpList: resWordpressGbp,
		Paginate: res.Paginate{
			Total: total,
			Count: len(wgList),
		},
	}, nil
}

func (u *wordpressGbpUsecase) GetWordpressGbp(ctx context.Context, id int, params req.GetWordpressGbpDetail) (*res.WordpressGbpDetail, error) {
	wg, err := u.wordpressGbpRepo.Get(ctx, repository.WordpressGbpFilter{
		ID: &id,
	})
	if err != nil {
		return nil, err
	}

	customerID := 300000 + wg.ID

	googlePostsCount, err := u.googlePostRepo.Count(ctx, repository.GooglePostFilter{
		CustomerID:    util.Pointer(customerID),
		OrderByIDDesc: util.Pointer(true),
		PostType:      util.Pointer(domain.PostTypePost),
	})
	if err != nil {
		return nil, err
	}

	googlePhotosCount, err := u.googlePostRepo.Count(ctx, repository.GooglePostFilter{
		CustomerID:    util.Pointer(customerID),
		OrderByIDDesc: util.Pointer(true),
		PostType:      util.Pointer(domain.PostTypePhoto),
	})
	if err != nil {
		return nil, err
	}

	return &res.WordpressGbpDetail{
		ID:                wg.ID,
		Name:              wg.Name,
		WordpressDomain:   wg.WordpressDomain,
		BusinessName:      wg.BusinessName,
		BusinessTitle:     wg.BusinessTitle,
		Memo:              wg.Memo,
		MapsURL:           wg.MapsURL,
		StartDate:         wg.StartDate,
		Status:            int(wg.Status),
		GooglePhotosCount: googlePhotosCount,
		GooglePostsCount:  googlePostsCount,
		CreatedAt:         wg.CreatedAt,
		UpdatedAt:         wg.UpdatedAt,
	}, nil
}

func (u *wordpressGbpUsecase) CreateWordpressGbp(ctx context.Context, body req.WordpressGbp) (*res.WordpressGbp, error) {
	// WordPress接続確認
	_, err := u.wordpressAdapter.GetTitle(ctx, body.WordpressDomain)
	if err != nil {
		return nil, domain.ErrWordpressConnection
	}

	// GBPビジネス存在確認
	business, err := u.gbpAdapter.GetBusiness(ctx, body.BusinessName)
	if err != nil {
		return nil, domain.ErrBusinessConnection
	}
	if business.Title == "" {
		return nil, domain.ErrBusinessConnection
	}

	wg := &domain.WordpressGbp{
		Name:            body.Name,
		Memo:            body.Memo,
		WordpressDomain: body.WordpressDomain,
		BusinessName:    business.Name,
		BusinessTitle:   business.Title,
		MapsURL:         business.MapsURL,
		StartDate:       body.StartDate,
		Status:          domain.Status(body.Status),
	}

	if err := u.wordpressGbpRepo.Create(ctx, wg); err != nil {
		return nil, err
	}

	return &res.WordpressGbp{
		ID:              wg.ID,
		Name:            wg.Name,
		WordpressDomain: wg.WordpressDomain,
		BusinessName:    wg.BusinessName,
		BusinessTitle:   wg.BusinessTitle,
		Memo:            wg.Memo,
		MapsURL:         wg.MapsURL,
		Status:          int(wg.Status),
		CreatedAt:       wg.CreatedAt,
		UpdatedAt:       wg.UpdatedAt,
	}, nil
}

func (u *wordpressGbpUsecase) UpdateWordpressGbp(ctx context.Context, id int, body req.WordpressGbp) (*res.WordpressGbp, error) {
	// WordPress接続確認
	_, err := u.wordpressAdapter.GetTitle(ctx, body.WordpressDomain)
	if err != nil {
		return nil, domain.ErrWordpressConnection
	}

	// GBPビジネス存在確認
	business, err := u.gbpAdapter.GetBusiness(ctx, body.BusinessName)
	if err != nil {
		return nil, domain.ErrBusinessConnection
	}

	wg, err := u.wordpressGbpRepo.Get(ctx, repository.WordpressGbpFilter{
		ID: &id,
	})
	if err != nil {
		return nil, err
	}
	wg.Name = body.Name
	wg.Memo = body.Memo
	wg.WordpressDomain = body.WordpressDomain
	wg.BusinessName = business.Name
	wg.BusinessTitle = business.Title
	wg.MapsURL = business.MapsURL
	wg.StartDate = body.StartDate
	wg.Status = domain.Status(body.Status)
	wg.UpdatedAt = time.Now()

	if err := u.wordpressGbpRepo.Update(ctx, wg, repository.WordpressGbpFilter{
		ID: &id,
	}); err != nil {
		return nil, err
	}

	return &res.WordpressGbp{
		ID:              wg.ID,
		Name:            wg.Name,
		WordpressDomain: wg.WordpressDomain,
		BusinessName:    wg.BusinessName,
		BusinessTitle:   wg.BusinessTitle,
		Memo:            wg.Memo,
		MapsURL:         wg.MapsURL,
		Status:          int(wg.Status),
		CreatedAt:       wg.CreatedAt,
		UpdatedAt:       wg.UpdatedAt,
	}, nil
}

func (u *wordpressGbpUsecase) DeleteWordpressGbp(ctx context.Context, id int) error {
	return u.wordpressGbpRepo.Delete(ctx, repository.WordpressGbpFilter{
		ID: &id,
	})
}
