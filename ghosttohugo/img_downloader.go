package ghosttohugo

import (
	"crypto/tls"
	"errors"
	"io"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	jww "github.com/spf13/jwalterweatherman"
)

type imageDownloader struct {
	siteUrl  string
	basePath string
	insecure bool
}

var ImgDownloader *imageDownloader

func NewImgDownloader(path string, url string, inscure bool) {
	ImgDownloader = &imageDownloader{
		siteUrl:  url,
		basePath: path,
		insecure: inscure,
	}
}

func (imgD imageDownloader) Download(img string) (string, error) {

	URL := strings.Replace(img, "__GHOST_URL__", imgD.siteUrl, -1)
	imageName := filepath.Base(img)
	jww.DEBUG.Printf("Downloading image: %s", imageName)

	if filepath.Ext(imageName) == "" || len(filepath.Ext(imageName)) > 5 {
		imageName = imageName + ".jpg"
	}

	fileName := path.Join("images", filepath.Base(imageName))
	savePath := path.Join(imgD.basePath, "static", fileName)

	transCfg := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: imgD.insecure}, // ignore expired SSL certificates
	}
	client := &http.Client{Transport: transCfg}

	response, err := client.Get(URL)
	if err != nil {
		return "", err
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", errors.New("received non 200 response code")
	}
	//Create a empty file
	file, err := os.Create(savePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	//Write the bytes to the fiel
	_, err = io.Copy(file, response.Body)
	if err != nil {
		return "", err
	}

	return "/" + fileName, nil
}
