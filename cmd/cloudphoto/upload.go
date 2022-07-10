package cloudphoto

import (
	"github.com/spf13/cobra"
	"github.com/ynsssss/cloudphoto/pkg/yc"
	"github.com/ynsssss/cloudphoto/pkg/utils"
)

func UploadCmdFunc() *cobra.Command {
	var albumName string
	var pathToPhotos string

	var uploadCmd = &cobra.Command{
		Use: "upload",
		Aliases: []string{"up"},
		Short: "Upload photos to a specified album",
		Run: func(cmd *cobra.Command, args []string) {
			utils.CheckIfIsDir(pathToPhotos)
			_ = yc.UploadPhotosToAlbum(albumName, pathToPhotos)
		},
	}	
	uploadCmd.Flags().StringVar(&albumName, "album", "", "Specify what album do you want to upload photos to")
	uploadCmd.MarkFlagRequired("album")
	uploadCmd.Flags().StringVar(&pathToPhotos, "path", "", "Specify where your album is located")
	uploadCmd.MarkFlagRequired("path")
	return uploadCmd
}


