package yc

import (
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	cpconfig "github.com/ynsssss/cloudphoto/config"

	"context"
	"strings"
	"os"
	"bytes"
	"net/http"
	"path/filepath"
	"io/fs"
	"fmt"
)

var downloader *manager.Downloader

var uploader *manager.Uploader

var client *s3.Client

var Conf cpconfig.Config

func ConnectAWSS3(ctx context.Context, extCfg cpconfig.Config, pathToConfig string) {
	Conf = extCfg
	Conf.PathToConfig = pathToConfig
	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		if service == s3.ServiceID && region == extCfg.Region {
			return aws.Endpoint{
				PartitionID:   "yc",
				URL:           extCfg.EndpointUrl,
				SigningRegion: extCfg.Region,
			}, nil
		}
		return aws.Endpoint{}, fmt.Errorf("unknown endpoint requested")
	})
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithEndpointResolverWithOptions(customResolver), config.WithSharedCredentialsFiles([]string{pathToConfig}))
	if err != nil {
		panic(err)
	}

	client = s3.NewFromConfig(cfg)
	downloader = manager.NewDownloader(client)
	uploader = manager.NewUploader(client)
}


func UploadPhotosToAlbum(albumName string, pathToDir string) error {
	var jpgPhotos []string
	filepath.WalkDir(pathToDir, func(s string, d fs.DirEntry, e error) error {
		if e != nil { return e }
		if (filepath.Ext(d.Name()) == ".jpg" || filepath.Ext(d.Name()) == ".jpeg") {
			jpgPhotos = append(jpgPhotos, s)
		}
		return nil
	})
	CreateAlbum(albumName)
	for _, elem := range jpgPhotos {
		UploadPhotoToAlbum(albumName, elem)
	}

	return nil
}


func UploadPhotoToAlbum(albumName string, path string) {
	file, err := os.Open(path)
	if err != nil {
		panic(err)
	}

	fileInfo, _ := file.Stat() 
	var size int64 = fileInfo.Size()
	buffer := make([]byte, size)
  file.Read(buffer)
	asize := aws.Int64(size)
	fileBytes := bytes.NewReader(buffer)

	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:        aws.String(Conf.BucketName),
		Key: 					 aws.String(fmt.Sprintf("albums/%s/%s", albumName, filepath.Base(path))),
		Body: 				 fileBytes,
		ContentType:   aws.String(http.DetectContentType(buffer)),
		ContentDisposition: aws.String(fmt.Sprintf("inline; %s=%s.%s", filepath.Base(path), filepath.Base(path), filepath.Ext(path))),
		ACL: types.ObjectCannedACLPublicRead,
		ContentLength: *asize,
	})

	file.Close()

	if err != nil {
		panic(err.Error())
	}
}

func CreateAlbum(albumName string) {
	_, err := client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(Conf.BucketName),
		Key: aws.String(fmt.Sprintf("albums/%s/", albumName)),
	})
	if err != nil {
		panic(err.Error())
	}
}


func GetObjectsWithPrefix(prefix string) *s3.ListObjectsV2Output {
	input := &s3.ListObjectsV2Input{
		Bucket: &Conf.BucketName,
		Prefix: &prefix,
	}

	resp, err := client.ListObjectsV2(context.TODO(), input)
	if err != nil {
		panic(err)
	}
	return resp
}

func DownloadAlbum(albumName string, path string) {
	photoObjects := GetObjectsWithPrefix(fmt.Sprintf("albums/%s/", albumName))
	if len(photoObjects.Contents) < 1 {
		fmt.Println("no such album")
		os.Exit(1)
	}
	photos := photoObjects.Contents[1:]
	for _, photo := range photos {
		fileDest := fmt.Sprintf("%s/%s", path, filepath.Base(*photo.Key))
		newFile, err := os.Create(fileDest)
		if err != nil {
			panic(err)
		}
		defer newFile.Close()
		numBytes, err := downloader.Download(context.TODO(), newFile, &s3.GetObjectInput{
			Bucket: aws.String(Conf.BucketName),
			Key: 		aws.String(*photo.Key),
		})
		if err != nil {
			panic(err)
		}
		fmt.Printf("downloaded %v bytes", numBytes)
	}
}

func ListAlbums() {
	albumsObjects := GetObjectsWithPrefix("albums/")
	albums := albumsObjects.Contents
	for _, album := range albums {
		if strings.HasSuffix(*album.Key, "/") {
			fmt.Println(filepath.Base(*album.Key))
		}
	}
}

func ListPhotosInAlbum(albumName string) {
	photosObjects := GetObjectsWithPrefix(fmt.Sprintf("albums/%s/", albumName))
	if len(photosObjects.Contents) == 0 {
		fmt.Println("Album does not exist")
		return
	}
	photos := photosObjects.Contents[1:]
	for _, photo := range photos {
		fmt.Println(*photo.Key)
	}
}

func ReturnAlbums() []string {
	albumsObjects := GetObjectsWithPrefix("albums/")
	albums := albumsObjects.Contents
	var returnAlbums []string
	for _, album := range albums {
		if strings.HasSuffix(*album.Key, "/") {
			returnAlbums = append(returnAlbums, filepath.Base(*album.Key))
		}
	}
	return returnAlbums
}

func UploadSites(albums []string, homeDir string) {
	for _, album := range albums {
		file, err := os.Open(fmt.Sprintf("%s/tmp/cloudphoto/%s.html", homeDir, album))
		if err != nil {
			panic(err)
		}

		_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: 			 aws.String(Conf.BucketName),
		Key:    			 aws.String(fmt.Sprintf("sites/%s.html", album)),
		Body: 				 file,
		ContentType:   aws.String("text/html; charset=utf-8"),
		ContentDisposition: aws.String(fmt.Sprintf("inline; %s=%s.html", album, album)),
		ACL: types.ObjectCannedACLPublicRead,
		})

		file.Close()

		if err != nil {
			panic(err.Error())
		}
	}
}

func UploadIndex(homeDir string) {
	file, err := os.Open(fmt.Sprintf("%s/tmp/cloudphoto/index.html", homeDir))
	if err != nil {
		panic(err)
	}
	_, err = uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: 			 aws.String(Conf.BucketName),
		Key:    			 aws.String("index.html"),
		Body: 				 file,
		ContentType:   aws.String("text/html; charset=utf-8"),
		ContentDisposition: aws.String("inline; index=index.html"),
		ACL: types.ObjectCannedACLPublicRead,
		})
  file.Close()

	if err != nil {
		panic(err)
	}
}


func ReturnPhotosInAlbum(albumName string) []string {
	photosObjects := GetObjectsWithPrefix(fmt.Sprintf("albums/%s/", albumName))
	if len(photosObjects.Contents) == 0 {
		fmt.Println("Album does not exist")
		os.Exit(1)
	}
	photos := photosObjects.Contents[1:]
	var returnPhotos []string
	for _, photo := range photos {
		returnPhotos = append(returnPhotos, *photo.Key)
	}
	return returnPhotos
}






