package domain

import (
	"fmt"
	"strings"
)

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
	return strings.Split(i.Caption, "\n")[0]
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

func (i *InstagramPost) getContentsHTML() string {
	caption := i.Caption
	contents := "<p>"
	lines := strings.Split(caption, "\n")
	for _, line := range lines {
		contents += line + "<br>"
	}
	contents += "</p>"
	return contents
}

func (i *InstagramPost) getHTMLForImage() string {
	imageHTML := fmt.Sprintf("<div style='text-align: center;'><img src='%s' style='margin: 0 auto;' width='500px' height='500px'/></div>", i.MediaURL)
	imageHTML += i.getContentsHTML()
	return imageHTML
}

func (i *InstagramPost) getHTMLForVideo() string {
	videoHTML := fmt.Sprintf("<div style='text-align: center;'><video src='%s' style='margin: 0 auto;' width='500px' height='500px' controls>Sorry, your browser does not support embedded videos.</video></div>", i.MediaURL)
	videoHTML += i.getContentsHTML()
	return videoHTML
}

func (i *InstagramPost) getHTMLForCarousel() string {
	html := "<div class='a-root-wordpress-instagram-slider'>"
	for _, child := range i.Children {
		if child.MediaType == "IMAGE" {
			html += fmt.Sprintf("<div style='text-align: center;'><img src='%s' style='margin: 0 auto;' width='500px' height='500px'/></div>", child.MediaURL)
		} else if child.MediaType == "VIDEO" {
			html += fmt.Sprintf("<div style='text-align: center;'><video src='%s' style='margin: 0 auto;' width='500px' height='500px' controls>Sorry, your browser does not support embedded videos.</video></div>", child.MediaURL)
		}
	}
	html += "</div>"
	html += i.getContentsHTML()
	return html
}
