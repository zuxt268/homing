package external

type GoogleBusinessMediaUploadResponse struct {
	Name         string `json:"name"`
	MediaFormat  string `json:"mediaFormat"`
	GoogleURL    string `json:"googleUrl"`
	SourceURL    string `json:"sourceUrl"`
	ThumbnailURL string `json:"thumbnailUrl"`
	CreateTime   string `json:"createTime"`
}