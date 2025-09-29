package res

import "github.com/zuxt268/homing/internal/domain"

type InstagramAccounts struct {
	InstagramAccounts []InstagramAccount `json:"instagram_accounts"`
}

type InstagramAccount struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func GetInstagramAccounts(accounts []domain.InstagramAccount) *InstagramAccounts {
	resp := make([]InstagramAccount, len(accounts))
	for i, account := range accounts {
		resp[i] = InstagramAccount{
			ID:   account.InstagramAccountID,
			Name: account.InstagramAccountName,
		}
	}
	return &InstagramAccounts{resp}
}
