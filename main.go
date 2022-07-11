package main 

import (
	"github.com/ynsssss/cloudphoto/cmd/cloudphoto"
	"github.com/ynsssss/cloudphoto/pkg/yc"
	cpconfig "github.com/ynsssss/cloudphoto/config"
	"github.com/ynsssss/cloudphoto/internal/pagecons"
	"os"
	"embed"
	"fmt"
	"context"
)

//go:embed template
var templateFs embed.FS

func main() {
	pagecons.InitFs(templateFs)
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("could not obtain user home folder")
		fmt.Println("make sure that you can edit the contents of your home folder")
		fmt.Println("try this: sudo chown -R your_login /home")
		return
	}

	var cpconfFile = fmt.Sprintf("%s/.config/cloudphoto/cloudphotorc", homeDir)
	
	cpConf, err:= cpconfig.LoadConfigIni(cpconfFile)

	yc.ConnectAWSS3(context.TODO(), cpConf, cpconfFile)


	cloudphoto.Execute()
}
