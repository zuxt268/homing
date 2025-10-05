package adapter

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/zuxt268/homing/internal/domain"
	"github.com/zuxt268/homing/internal/infrastructure/driver"
	"github.com/zuxt268/homing/internal/interface/dto/external"
)

func TestWordpressAdapter_Post(t *testing.T) {
	post := domain.InstagramPost{
		ID:        "18287747659283707",
		Permalink: "https://www.instagram.com/p/DOdb8pEjyLU/",
		Caption:   "DDDDD",
		Timestamp: "2025-09-29T10:19:24+0000",
		MediaType: "CAROUSEL_ALBUM",
		MediaURL:  "https://scontent-nrt1-1.cdninstagram.com/v/t51.82787-15/545510899_17858706267498031_8817069271596986900_n.webp?stp=dst-jpg_e35_tt6&_nc_cat=106&ccb=1-7&_nc_sid=18de74&_nc_ohc=LyGPz35EmhwQ7kNvwEi04Dl&_nc_oc=Adkj6qPf6CGsJl7HSFGjiz3Nmodo2mtYH7lYEjrythPr0r-3L3RkYK-FD1isFLIqLug&_nc_zt=23&_nc_ht=scontent-nrt1-1.cdninstagram.com&edm=AL-3X8kEAAAA&_nc_gid=7xCyPmu21M4s5JlAQ9a9nw&oh=00_AfbXhH9SlsBe6LewpLQL2X4-fVJ-7cDTahiZxRBXpI8Bug&oe=68D056EE",
		Children: []domain.InstagramPostChildren{
			{
				MediaType: "IMAGE",
				MediaURL:  "https://scontent-nrt1-1.cdninstagram.com/v/t51.82787-15/545510899_17858706267498031_8817069271596986900_n.webp?stp=dst-jpg_e35_tt6&_nc_cat=106&ccb=1-7&_nc_sid=18de74&_nc_ohc=LyGPz35EmhwQ7kNvwEi04Dl&_nc_oc=Adkj6qPf6CGsJl7HSFGjiz3Nmodo2mtYH7lYEjrythPr0r-3L3RkYK-FD1isFLIqLug&_nc_zt=23&_nc_ht=scontent-nrt1-1.cdninstagram.com&edm=AL-3X8kEAAAA&_nc_gid=7xCyPmu21M4s5JlAQ9a9nw&oh=00_AfbXhH9SlsBe6LewpLQL2X4-fVJ-7cDTahiZxRBXpI8Bug&oe=68D056EE",
				ID:        "18057626774382203",
			},
		},
	}
	customer := domain.Customer{WordpressUrl: "hp-standard.moe"}

	httpClient := &http.Client{}
	client := driver.NewClient(httpClient)
	adapter := NewWordpressAdapter(client)
	posted, err := adapter.Post(context.Background(), external.WordpressPostInput{
		Post:            post,
		FeaturedMediaID: 0,
		Customer:        customer,
	})
	if err != nil {
		t.Errorf("error posting %v", err)
	}
	fmt.Println(posted)
}

func TestUploadFile(t *testing.T) {
	customer := domain.Customer{WordpressUrl: "hp-standard.moe"}
	httpClient := &http.Client{}
	client := driver.NewClient(httpClient)
	adapter := NewWordpressAdapter(client)

	resp, err := adapter.FileUpload(context.Background(), external.WordpressFileUploadInput{
		Path:     "/var/folders/3t/gfwjqksn6tqfj5kvg70dzlwr0000gn/T/homing_download_1656907455/548865242_17916787776176467_8381450328613983170_n.jpg",
		Customer: customer,
	})
	if err != nil {
		t.Errorf("error uploading %v", err)
	}
	fmt.Println(resp)
}
