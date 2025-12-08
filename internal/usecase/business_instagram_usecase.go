package usecase

import (
	"context"
	"time"

	"github.com/zuxt268/homing/internal/config"
	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/adapter"
	"github.com/zuxt268/homing/internal/interface/dto/req"
	"github.com/zuxt268/homing/internal/interface/dto/res"
	"github.com/zuxt268/homing/internal/interface/repository"
	"github.com/zuxt268/homing/internal/interface/util"
)

type BusinessInstagramUsecase interface {
	FetchGoogleBusinesses(ctx context.Context) error
	GetGoogleBusinesses(ctx context.Context, limit, offset int) ([]*domain.GoogleBusinesses, int64, error)

	GetBusinessInstagram(ctx context.Context, id int) (*res.BusinessInstagram, error)
	GetBusinessInstagramList(ctx context.Context, params req.GetBusinessInstagram) (*res.BusinessInstagramList, error)
	CreateBusinessInstagram(ctx context.Context, body req.BusinessInstagram) (*res.BusinessInstagram, error)
	UpdateBusinessInstagram(ctx context.Context, id int, body req.BusinessInstagram) (*res.BusinessInstagram, error)
	DeleteBusinessInstagram(ctx context.Context, id int) error
}

type businessInstagramUsecase struct {
	googleBusinessRepo    repository.GoogleBusinessRepository
	tokenRepo             repository.TokenRepository
	businessInstagramRepo repository.BusinessInstagramRepository
	instagramAdapter      adapter.InstagramAdapter
	gbpAdapter            adapter.GbpAdapter
}

func NewBusinessInstagramUsecase(
	googleBusinessRepo repository.GoogleBusinessRepository,
	tokenRepo repository.TokenRepository,
	businessInstagramRepo repository.BusinessInstagramRepository,
	instagramAdapter adapter.InstagramAdapter,
	gbpAdapter adapter.GbpAdapter,
) BusinessInstagramUsecase {
	return &businessInstagramUsecase{
		googleBusinessRepo:    googleBusinessRepo,
		tokenRepo:             tokenRepo,
		businessInstagramRepo: businessInstagramRepo,
		instagramAdapter:      instagramAdapter,
		gbpAdapter:            gbpAdapter,
	}
}

func (u *businessInstagramUsecase) FetchGoogleBusinesses(ctx context.Context) error {
	businesses, err := u.gbpAdapter.GetAllBusinesses(ctx, config.Env.GoogleBusinessAccountName)
	if err != nil {
		return err
	}
	for _, business := range businesses {
		exists, err := u.googleBusinessRepo.Exists(ctx, repository.GoogleBusinessFilter{
			Name: &business.Name,
		})
		if err != nil {
			return err
		}
		if exists {
			continue
		}
		if err := u.googleBusinessRepo.Create(ctx, &domain.GoogleBusinesses{
			Name:  business.Name,
			Title: business.Title,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (u *businessInstagramUsecase) GetGoogleBusinesses(ctx context.Context, limit, offset int) ([]*domain.GoogleBusinesses, int64, error) {
	// 総数を取得
	total, err := u.googleBusinessRepo.Count(ctx, repository.GoogleBusinessFilter{
		All: util.Pointer(true),
	})
	if err != nil {
		return nil, 0, err
	}

	// ページングして取得
	businesses, err := u.googleBusinessRepo.FindAll(ctx, repository.GoogleBusinessFilter{
		Limit:         &limit,
		Offset:        &offset,
		OrderByIDDesc: util.Pointer(true),
	})
	if err != nil {
		return nil, 0, err
	}

	return businesses, total, nil
}

func (u *businessInstagramUsecase) GetBusinessInstagramList(ctx context.Context, params req.GetBusinessInstagram) (*res.BusinessInstagramList, error) {
	biList, err := u.businessInstagramRepo.FindAll(ctx, repository.BusinessInstagramFilter{
		InstagramID: params.InstagramID,
		Limit:       params.Limit,
		Offset:      params.Offset,
		PartialName: params.Name,
		Status:      params.Status,
	})
	if err != nil {
		return nil, err
	}
	total, err := u.googleBusinessRepo.Count(ctx, repository.GoogleBusinessFilter{
		All: util.Pointer(true),
	})
	if err != nil {
		return nil, err
	}

	resBusinessInstagram := make([]res.BusinessInstagram, len(biList))
	for i, business := range biList {
		resBusinessInstagram[i] = res.BusinessInstagram{
			ID:           business.ID,
			Name:         business.Name,
			BusinessName: business.BusinessName,
			InstagramID:  business.InstagramID,
			Memo:         business.Memo,
			StartDate:    business.StartDate,
			Status:       int(business.Status),
			CreatedAt:    business.CreatedAt,
			UpdatedAt:    business.UpdatedAt,
		}
	}
	return &res.BusinessInstagramList{
		BusinessInstagramList: resBusinessInstagram,
		Paginate: res.Paginate{
			Total: total,
			Count: len(biList),
		},
	}, nil
}

func (u *businessInstagramUsecase) GetBusinessInstagram(ctx context.Context, id int) (*res.BusinessInstagram, error) {
	bi, err := u.businessInstagramRepo.Get(ctx, repository.BusinessInstagramFilter{
		ID: &id,
	})
	if err != nil {
		return nil, err
	}
	return &res.BusinessInstagram{
		ID:           bi.ID,
		Name:         bi.Name,
		BusinessName: bi.BusinessName,
		InstagramID:  bi.InstagramID,
		Memo:         bi.Memo,
		StartDate:    bi.StartDate,
		Status:       int(bi.Status),
		CreatedAt:    bi.CreatedAt,
		UpdatedAt:    bi.UpdatedAt,
	}, nil
}

func (u *businessInstagramUsecase) CreateBusinessInstagram(ctx context.Context, body req.BusinessInstagram) (*res.BusinessInstagram, error) {
	token, err := u.tokenRepo.First(ctx)
	if err != nil {
		return nil, err
	}
	instagram, err := u.instagramAdapter.GetAccount(ctx, token, body.InstagramID)
	if err != nil {
		return nil, err
	}
	if instagram.InstagramAccountUserName == "" {
		return nil, domain.ErrInstagramConnection
	}

	business, err := u.gbpAdapter.GetBusiness(ctx, body.BusinessName)
	if err != nil {
		return nil, err
	}
	if business.Title == "" {
		return nil, domain.ErrBusinessConnection
	}

	bi := &domain.BusinessInstagram{
		Name:          body.Name,
		Memo:          body.Memo,
		InstagramID:   instagram.InstagramAccountID,
		InstagramName: instagram.InstagramAccountName,
		BusinessName:  business.Name,
		BusinessTitle: business.Title,
		StartDate:     body.StartDate,
		Status:        domain.Status(body.Status),
	}

	if err := u.businessInstagramRepo.Create(ctx, bi); err != nil {
		return nil, err
	}

	return &res.BusinessInstagram{
		ID:           bi.ID,
		Name:         bi.Name,
		BusinessName: bi.BusinessName,
		InstagramID:  bi.InstagramID,
		Memo:         bi.Memo,
		StartDate:    bi.StartDate,
		Status:       int(bi.Status),
		CreatedAt:    bi.CreatedAt,
		UpdatedAt:    bi.UpdatedAt,
	}, nil
}

func (u *businessInstagramUsecase) UpdateBusinessInstagram(ctx context.Context, id int, body req.BusinessInstagram) (*res.BusinessInstagram, error) {
	token, err := u.tokenRepo.First(ctx)
	if err != nil {
		return nil, err
	}
	instagram, err := u.instagramAdapter.GetAccount(ctx, token, body.InstagramID)
	if err != nil {
		return nil, err
	}

	business, err := u.gbpAdapter.GetBusiness(ctx, body.BusinessName)
	if err != nil {
		return nil, err
	}

	bi, err := u.businessInstagramRepo.Get(ctx, repository.BusinessInstagramFilter{
		ID: &id,
	})
	if err != nil {
		return nil, err
	}
	bi.Name = body.Name
	bi.Memo = body.Memo
	bi.InstagramID = instagram.InstagramAccountID
	bi.InstagramName = instagram.InstagramAccountName
	bi.BusinessName = business.Name
	bi.BusinessTitle = business.Title
	bi.StartDate = body.StartDate
	bi.Status = domain.Status(body.Status)
	bi.UpdatedAt = time.Now()

	if err := u.businessInstagramRepo.Update(ctx, bi, repository.BusinessInstagramFilter{
		ID: &id,
	}); err != nil {
		return nil, err
	}

	return &res.BusinessInstagram{
		ID:           bi.ID,
		Name:         bi.Name,
		BusinessName: bi.BusinessName,
		InstagramID:  bi.InstagramID,
		Memo:         bi.Memo,
		StartDate:    bi.StartDate,
		Status:       int(bi.Status),
		CreatedAt:    bi.CreatedAt,
		UpdatedAt:    bi.UpdatedAt,
	}, nil
}

func (u *businessInstagramUsecase) DeleteBusinessInstagram(ctx context.Context, id int) error {
	return u.businessInstagramRepo.Delete(ctx, repository.BusinessInstagramFilter{
		ID: &id,
	})
}
