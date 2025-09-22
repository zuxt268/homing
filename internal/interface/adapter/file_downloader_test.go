package adapter

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileDownloader_Download(t *testing.T) {
	u := "https://scontent-nrt1-1.cdninstagram.com/v/t51.82787-15/548865242_17916787776176467_8381450328613983170_n.jpg?stp=dst-jpg_e35_tt6&_nc_cat=109&ccb=1-7&_nc_sid=18de74&_nc_ohc=7WNWypEfsI4Q7kNvwF7t7E5&_nc_oc=AdmcjKDmwgqq_amRmGjCJLc2DD6J6zNptG7Z-M0qsbw-ehCvHF8mSWKEQUREdiot4Gk&_nc_zt=23&_nc_ht=scontent-nrt1-1.cdninstagram.com&edm=AL-3X8kEAAAA&_nc_gid=EgLlLvyMzVppYx0cYVA0kg&oh=00_AfauDgaShf6StUNFAp_wi8R2b_RijYdxCAV2INPJJNrmeA&oe=68D694BA"
	client := &http.Client{}
	downloader := NewFileDownloader(client)
	path, err := downloader.Download(context.Background(), u)
	assert.NoError(t, err)
	fmt.Println(path)

	//err = downloader.DeleteTempDirectory()
	assert.NoError(t, err)
}
