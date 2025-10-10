package usecase

import (
	"context"

	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/interface/adapter"
	"github.com/zuxt268/homing/internal/interface/dto/req"
	"github.com/zuxt268/homing/internal/interface/dto/res"
	"github.com/zuxt268/homing/internal/interface/repository"
	"github.com/zuxt268/homing/internal/interface/util"
)

type WordpressInstagramUsecase interface {
	GetWordpressInstagramList(ctx context.Context, params req.GetWordpressInstagram) (*res.WordpressInstagramList, error)
	GetWordpressInstagram(ctx context.Context, id int) (*res.WordpressInstagramDetail, error)
	CreateWordpressInstagram(ctx context.Context, body req.CreateWordpressInstagram) (*res.WordpressInstagram, error)
	UpdateWordpressInstagram(ctx context.Context, body req.UpdateWordpressInstagram) (*res.WordpressInstagram, error)
	DeleteWordpressInstagram(ctx context.Context, id int) error
}

type wordpressInstagramUsecase struct {
	wordpressInstagramRepo repository.WordpressInstagramRepository
	tokenRepo              repository.TokenRepository
	postRepo               repository.PostRepository
	instagramAdapter       adapter.InstagramAdapter
	wordpressAdapter       adapter.WordpressAdapter
}

func NewWordpressInstagramUsecase(
	wordpressInstagramRepo repository.WordpressInstagramRepository,
	tokenRepo repository.TokenRepository,
	postRepo repository.PostRepository,
	instagramAdapter adapter.InstagramAdapter,
	wordpressAdapter adapter.WordpressAdapter,
) WordpressInstagramUsecase {
	return &wordpressInstagramUsecase{
		wordpressInstagramRepo: wordpressInstagramRepo,
		tokenRepo:              tokenRepo,
		postRepo:               postRepo,
		instagramAdapter:       instagramAdapter,
		wordpressAdapter:       wordpressAdapter,
	}
}

func (u *wordpressInstagramUsecase) GetWordpressInstagramList(ctx context.Context, params req.GetWordpressInstagram) (*res.WordpressInstagramList, error) {
	filter := repository.WordpressInstagramFilter{
		PartialName:            params.Name,
		PartialWordpressDomain: params.WordpressDomain,
		PartialInstagramName:   params.InstagramName,
		WordpressSiteTitle:     params.WordpressSiteTitle,
		InstagramID:            params.InstagramID,
		Status:                 params.Status,
		DeleteHash:             params.DeleteHash,
		Limit:                  params.Limit,
		Offset:                 params.Offset,
	}

	wiList, err := u.wordpressInstagramRepo.FindAll(ctx, filter)
	if err != nil {
		return nil, err
	}

	total, err := u.wordpressInstagramRepo.Count(ctx, filter)
	if err != nil {
		return nil, err
	}

	result := make([]res.WordpressInstagram, 0, len(wiList))
	for _, wi := range wiList {
		result = append(result, res.WordpressInstagram{
			ID:                 wi.ID,
			Name:               wi.Name,
			WordpressDomain:    wi.WordpressDomain,
			WordpressSiteTitle: wi.WordpressSiteTitle,
			InstagramID:        wi.InstagramID,
			InstagramName:      wi.InstagramName,
			Memo:               wi.Memo,
			StartDate:          wi.StartDate,
			Status:             int(wi.Status),
			DeleteHash:         wi.DeleteHash,
		})
	}

	return &res.WordpressInstagramList{
		WordpressInstagramList: result,
		Paginate: res.Paginate{
			Total: total,
			Count: len(wiList),
		},
	}, nil
}

func (u *wordpressInstagramUsecase) GetWordpressInstagram(ctx context.Context, id int) (*res.WordpressInstagramDetail, error) {
	wi, err := u.wordpressInstagramRepo.Get(ctx, repository.WordpressInstagramFilter{
		ID: &id,
	})
	if err != nil {
		return nil, err
	}

	filter := repository.PostFilter{
		CustomerID:           util.Pointer(100000 + wi.ID),
		OrderByCreatedAtDesc: util.Pointer(true),
	}

	posts, err := u.postRepo.GetPosts(ctx, filter)
	if err != nil {
		return nil, err
	}

	total, err := u.postRepo.CountPosts(ctx, filter)
	if err != nil {
		return nil, err
	}

	respPosts := make([]res.Post, len(posts))
	for i, post := range posts {
		respPosts[i] = res.Post{
			WordpressUrl: post.WordpressURL,
			InstagramUrl: post.InstagramURL,
			CreatedAt:    post.CreatedAt,
		}
	}

	return &res.WordpressInstagramDetail{
		ID:                 wi.ID,
		Name:               wi.Name,
		WordpressDomain:    wi.WordpressDomain,
		WordpressSiteTitle: wi.WordpressSiteTitle,
		InstagramID:        wi.InstagramID,
		InstagramName:      wi.InstagramName,
		Memo:               wi.Memo,
		StartDate:          wi.StartDate,
		Status:             int(wi.Status),
		DeleteHash:         wi.DeleteHash,
		Posts: res.Posts{
			Posts: respPosts,
			Paginate: res.Paginate{
				Total: total,
				Count: len(posts),
			},
		},
	}, nil
}

func (u *wordpressInstagramUsecase) CreateWordpressInstagram(ctx context.Context, req req.CreateWordpressInstagram) (*res.WordpressInstagram, error) {

	token, err := u.tokenRepo.First(ctx)
	if err != nil {
		return nil, domain.ErrTokenNotFound
	}

	// システムユーザーで取得できるか確認
	account, err := u.instagramAdapter.GetAccount(ctx, token, req.InstagramID)
	if err != nil {
		return nil, domain.ErrInstagramConnection
	}
	if account.InstagramAccountUserName == "" {
		return nil, domain.ErrInstagramConnection
	}

	// ワードプレスと疎通できるか
	title, err := u.wordpressAdapter.GetTitle(ctx, req.WordpressDomain)
	if err != nil {
		return nil, domain.ErrWordpressConnection
	}

	wi := &domain.WordpressInstagram{
		Name:               req.Name,
		WordpressDomain:    req.WordpressDomain,
		WordpressSiteTitle: title,
		InstagramID:        req.InstagramID,
		InstagramName:      account.InstagramAccountUserName,
		Memo:               req.Memo,
		StartDate:          req.StartDate,
		Status:             domain.Status(req.Status),
		DeleteHash:         req.DeleteHash,
	}

	if err := u.wordpressInstagramRepo.Create(ctx, wi); err != nil {
		return nil, err
	}

	return &res.WordpressInstagram{
		ID:                 wi.ID,
		Name:               wi.Name,
		WordpressDomain:    wi.WordpressDomain,
		WordpressSiteTitle: wi.WordpressSiteTitle,
		InstagramID:        wi.InstagramID,
		InstagramName:      wi.InstagramName,
		Memo:               wi.Memo,
		StartDate:          wi.StartDate,
		Status:             int(wi.Status),
		DeleteHash:         wi.DeleteHash,
	}, nil
}

func (u *wordpressInstagramUsecase) UpdateWordpressInstagram(ctx context.Context, req req.UpdateWordpressInstagram) (*res.WordpressInstagram, error) {
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
		wi.WordpressDomain = *req.Wordpress
		title, err := u.wordpressAdapter.GetTitle(ctx, wi.WordpressDomain)
		if err != nil {
			return nil, domain.ErrWordpressConnection
		}
		wi.WordpressSiteTitle = title
	}
	if req.InstagramID != nil {
		wi.InstagramID = *req.InstagramID
		token, err := u.tokenRepo.First(ctx)
		if err != nil {
			return nil, err
		}
		account, err := u.instagramAdapter.GetAccount(ctx, token, wi.InstagramID)
		if err != nil {
			return nil, domain.ErrInstagramConnection
		}
		if account.InstagramAccountUserName == "" {
			return nil, domain.ErrInstagramConnection
		}
		wi.InstagramName = account.InstagramAccountUserName
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

	err = u.wordpressInstagramRepo.Update(ctx, wi, repository.WordpressInstagramFilter{
		ID: req.ID,
	})
	if err != nil {
		return nil, err
	}

	return &res.WordpressInstagram{
		ID:                 wi.ID,
		Name:               wi.Name,
		WordpressDomain:    wi.WordpressDomain,
		WordpressSiteTitle: wi.WordpressSiteTitle,
		InstagramID:        wi.InstagramID,
		InstagramName:      wi.InstagramName,
		Memo:               wi.Memo,
		StartDate:          wi.StartDate,
		Status:             int(wi.Status),
		DeleteHash:         wi.DeleteHash,
	}, nil
}

func (u *wordpressInstagramUsecase) DeleteWordpressInstagram(ctx context.Context, id int) error {
	return u.wordpressInstagramRepo.Delete(ctx, repository.WordpressInstagramFilter{
		ID: &id,
	})
}
