package domain

import (
	"fmt"
	"regexp"
	"strings"
	"time"
)

type InstagramAccount struct {
	InstagramAccountName     string
	InstagramAccountID       string
	InstagramAccountUserName string
}

type InstagramPost struct {
	ID              string
	Permalink       string
	Caption         string
	Timestamp       string
	MediaType       string
	MediaURL        string
	Children        []InstagramPostChildren
	SourceURLs      []string
	FeaturedMediaID int

	DeleteHash bool
}

func (i *InstagramPost) SetFeaturedMediaID(mediaID int) {
	i.FeaturedMediaID = mediaID
}

func (i *InstagramPost) AppendSourceURL(imageUrl string) {
	i.SourceURLs = append(i.SourceURLs, imageUrl)
}

func (i *InstagramPost) SetDeleteHashFlag(deleteHash bool) {
	i.DeleteHash = deleteHash
}

type InstagramPostChildren struct {
	MediaType string
	MediaURL  string
	ID        string
}

func (i *InstagramPost) GetTitle() string {
	caption := i.Caption
	if i.DeleteHash {
		caption = removeHashtags(caption)
	}
	for _, w := range strings.Split(caption, "\n") {
		if strings.TrimSpace(w) == "" {
			continue
		}
		return w
	}
	return " "
}

func (i *InstagramPost) GetContent() string {
	switch i.MediaType {
	case "IMAGE":
		return i.getHTMLForImage()
	case "VIDEO":
		return i.getHTMLForVideo()
	case "CAROUSEL_ALBUM":
		return i.getHTMLForCarousel()
	default:
		return i.getContentsHTML()
	}
}

func (i *InstagramPost) GetPostDate() string {
	instagramPost, _ := time.Parse("2006-01-02T15:04:05-0700", i.Timestamp)
	return instagramPost.Format("2006-01-02 15:04:05")
}

func (i *InstagramPost) getContentsHTML() string {
	caption := i.Caption
	// captionから#のタグ（#からスペースまたは改行まで）を削除
	if i.DeleteHash {
		caption = removeHashtags(caption)
	}

	contents := "<p>"
	lines := strings.Split(caption, "\n")
	for _, line := range lines {
		// 空行はスキップしない（改行を保持）
		contents += line + "<br>"
	}
	contents += "</p>"
	return contents
}

// removeHashtags はキャプションからハッシュタグを削除する
func removeHashtags(text string) string {
	re := regexp.MustCompile(`#\S+`)
	result := re.ReplaceAllString(text, "")
	return strings.TrimSpace(result)
}

func (i *InstagramPost) getHTMLForImage() string {
	imageHTML := fmt.Sprintf("<div style='text-align: center;'><img src='%s' style='margin: 0 auto;' width='500px' height='500px'/></div>", i.SourceURLs[0])
	imageHTML += i.getContentsHTML()
	return imageHTML
}

func (i *InstagramPost) getHTMLForVideo() string {
	videoHTML := fmt.Sprintf("<div style='text-align: center;'><video src='%s' style='margin: 0 auto;' width='500px' height='500px' controls>Sorry, your browser does not support embedded videos.</video></div>", i.SourceURLs[0])
	videoHTML += i.getContentsHTML()
	return videoHTML
}

func (i *InstagramPost) getHTMLForCarousel() string {
	html := "<div class='a-root-wordpress-instagram-slider'>"
	for idx, child := range i.Children {
		if child.MediaType == "IMAGE" {
			html += fmt.Sprintf("<div style='text-align: center;'><img src='%s' style='margin: 0 auto;' width='500px' height='500px'/></div>", i.SourceURLs[idx+1])
		} else if child.MediaType == "VIDEO" {
			html += fmt.Sprintf("<div style='text-align: center;'><video src='%s' style='margin: 0 auto;' width='500px' height='500px' controls>Sorry, your browser does not support embedded videos.</video></div>", i.SourceURLs[idx+1])
		}
	}
	html += "</div>"
	html += i.getContentsHTML()
	return html
}
