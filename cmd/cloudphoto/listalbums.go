package cloudphoto

import (
	"github.com/spf13/cobra"
	"github.com/ynsssss/cloudphoto/pkg/yc"
)

func ListAlbumsCmdFunc() *cobra.Command {
	var listAlbumsCmd = &cobra.Command{
		Use: "list-albums",
		Short: "List all albums in bucket",
		Run: func(cmd *cobra.Command, args []string) {
			yc.ListAlbums()
		},
	}	
	return listAlbumsCmd
}
