package entity

type InstagramAccount struct {
	InstagramAccountName     string
	InstagramAccountID       string
	InstagramAccountUsername string
}

type InstagramPost struct {
	ID        string
	Permalink string
	Caption   string
	Timestamp string
	MediaType string
	MediaURL  string
	Children  []InstagramPostChildren
}

type InstagramPostChildren struct {
	MediaType string
	MediaURL  string
	ID        string
}

func (i *InstagramPost) GetTitle() string {
	return i.Caption
}

func (i *InstagramPost) GetContent() string {
	return "<h1>Instagram</h1>"
}
