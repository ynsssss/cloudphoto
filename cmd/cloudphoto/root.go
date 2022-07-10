package cloudphoto 
import (
	"fmt"
	"os"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use: "cloudphoto",
	Short: "cloudphoto is a cli tool to manage photos in your Yandex Cloud",

	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error executing CLI: %s", err)
		os.Exit(1)
	}
}


func init() {
	rootCmd.AddCommand(UploadCmdFunc())
	rootCmd.AddCommand(DownloadCmdFunc())
	rootCmd.AddCommand(ListAlbumsCmdFunc())
	rootCmd.AddCommand(ListPhotosInAlbumCmdFunc())
	rootCmd.AddCommand(GenerateSiteCmdFunc())
	cobra.OnInitialize(initConfig)
}


func initConfig() {

}
