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
	ID        string
	Permalink string
	Caption   string
	Timestamp string
	MediaType string
	MediaURL  string
	Children  []InstagramPostChildren

	DeleteHash bool
}

type InstagramPostChildren struct {
	MediaType string
	MediaURL  string
	ID        string
}

func (i *InstagramPost) GetTitle(deleteHash bool) string {
	caption := i.Caption
	if deleteHash {
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

func (i *InstagramPost) GetContent(deleteHash bool) string {
	switch i.MediaType {
	case "IMAGE":
		return i.getHTMLForImage(deleteHash)
	case "VIDEO":
		return i.getHTMLForVideo(deleteHash)
	case "CAROUSEL_ALBUM":
		return i.getHTMLForCarousel(deleteHash)
	default:
		return i.getContentsHTML(deleteHash)
	}
}

func (i *InstagramPost) GetPostDate() string {
	instagramPost, _ := time.Parse("2006-01-02T15:04:05-0700", i.Timestamp)
	return instagramPost.Format("2006-01-02 15:04:05")
}

func (i *InstagramPost) getContentsHTML(deleteHash bool) string {
	caption := i.Caption
	// captionから#のタグ（#からスペースまたは改行まで）を削除
	if deleteHash {
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

func (i *InstagramPost) getHTMLForImage(deleteHash bool) string {
	imageHTML := fmt.Sprintf("<div style='text-align: center;'><img src='%s' style='margin: 0 auto;' width='500px' height='500px'/></div>", i.MediaURL)
	imageHTML += i.getContentsHTML(deleteHash)
	return imageHTML
}

func (i *InstagramPost) getHTMLForVideo(deleteHash bool) string {
	videoHTML := fmt.Sprintf("<div style='text-align: center;'><video src='%s' style='margin: 0 auto;' width='500px' height='500px' controls>Sorry, your browser does not support embedded videos.</video></div>", i.MediaURL)
	videoHTML += i.getContentsHTML(deleteHash)
	return videoHTML
}

func (i *InstagramPost) getHTMLForCarousel(deleteHash bool) string {
	html := "<div class='a-root-wordpress-instagram-slider'>"
	for _, child := range i.Children {
		if child.MediaType == "IMAGE" {
			html += fmt.Sprintf("<div style='text-align: center;'><img src='%s' style='margin: 0 auto;' width='500px' height='500px'/></div>", child.MediaURL)
		} else if child.MediaType == "VIDEO" {
			html += fmt.Sprintf("<div style='text-align: center;'><video src='%s' style='margin: 0 auto;' width='500px' height='500px' controls>Sorry, your browser does not support embedded videos.</video></div>", child.MediaURL)
		}
	}
	html += "</div>"
	html += i.getContentsHTML(deleteHash)
	return html
}
