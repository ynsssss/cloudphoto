package cloudphoto

import (
	"github.com/spf13/cobra"
	"github.com/ynsssss/cloudphoto/pkg/yc"
	"github.com/ynsssss/cloudphoto/pkg/utils"
)

func DownloadCmdFunc() *cobra.Command {
	var albumNameDownload string
	var pathToPhotosDownload string

	var downloadCmd = &cobra.Command{
		Use: "download",
		Aliases: []string{"dl"},
		Short: "Download photos from a specified album",
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckIfIsDir(pathToPhotosDownload)
			yc.DownloadAlbum(albumNameDownload, pathToPhotosDownload)
		},
	}	

	downloadCmd.Flags().StringVar(&albumNameDownload, "album", "", "Source to your photo folder")
	downloadCmd.MarkFlagRequired("album")
	downloadCmd.Flags().StringVar(&pathToPhotosDownload, "path", "", "Source to your photo folder")
	downloadCmd.MarkFlagRequired("path")
	return downloadCmd
}
