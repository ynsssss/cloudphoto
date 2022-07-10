package pagecons


import (
	"bufio"
	"os"
	"bytes"
	"html/template"
	"fmt"
	"embed"
)

var tplFolder embed.FS

func ConstructPageWithArgs(photos []string, albumName string, homeDir string) {
	file := "template/page.tmpl"

	template := template.Must(template.New("").ParseFS(tplFolder, file))

	var processed bytes.Buffer
	template.ExecuteTemplate(&processed, "page", photos)

	outputPath := fmt.Sprintf("%s/tmp/cloudphoto/%s.html", homeDir, albumName)
	f, _ := os.Create(outputPath)
	w := bufio.NewWriter(f)
	w.WriteString(string(processed.Bytes()))
	w.Flush()
}

func ConstructIndexPage(albums []string, homeDir string) {
	file := "template/index.tmpl"
	template := template.Must(template.New("").ParseFS(tplFolder, file))

	var processed bytes.Buffer
	template.ExecuteTemplate(&processed, "index", albums)

	outputPath := fmt.Sprintf("%s/tmp/cloudphoto/index.html", homeDir)
	f, _ := os.Create(outputPath)
	w := bufio.NewWriter(f)
	w.WriteString(string(processed.Bytes()))
	w.Flush()
}

func InitFs(fold embed.FS) {
	tplFolder = fold
}
