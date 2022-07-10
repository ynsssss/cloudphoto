package cloudphoto

import (
	"github.com/spf13/cobra"
	"github.com/ynsssss/cloudphoto/pkg/yc"
)

func ListPhotosInAlbumCmdFunc() *cobra.Command {
	var albumNameList string
	var listPhotosInAlbumCmd = &cobra.Command{
		Use: "list-photos",
		Short: "List all photos in album",
		Run: func(cmd *cobra.Command, args []string) {
			yc.ListPhotosInAlbum(albumNameList)
		},
	}

	listPhotosInAlbumCmd.Flags().StringVar(&albumNameList, "album", "", "Specify photos from what album do you want")
	listPhotosInAlbumCmd.MarkFlagRequired("album")	
	return listPhotosInAlbumCmd
}
