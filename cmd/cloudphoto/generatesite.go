package cloudphoto

import (
	"github.com/spf13/cobra"
	"github.com/ynsssss/cloudphoto/pkg/yc"
	"github.com/ynsssss/cloudphoto/internal/pagecons"
	"fmt"
	"os"
)

func GenerateSiteCmdFunc() *cobra.Command {
	var generateSiteCmd = &cobra.Command{
		Use: "generate-site",
		Short: "generate site and return link",
		Run: func(cmd *cobra.Command, args []string) {
			albums := yc.ReturnAlbums()
			homeDir, _ := os.UserHomeDir()
			_ = os.MkdirAll(fmt.Sprintf("%s/tmp/cloudphoto", homeDir), os.ModePerm)
			for _, album := range albums {
				photos := yc.ReturnPhotosInAlbum(album)
				pagecons.ConstructPageWithArgs(photos, album, homeDir)
			}

			pagecons.ConstructIndexPage(albums, homeDir)
			yc.UploadSites(albums, homeDir)
			yc.UploadIndex(homeDir)
			fmt.Printf("https://storage.yandexcloud.net/%s/index.html\n", yc.Conf.BucketName)
		},
	}	

	return generateSiteCmd
}

