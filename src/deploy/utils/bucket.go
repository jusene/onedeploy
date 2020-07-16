package utils

import (
	"github.com/spf13/viper"
	"deploy/hwcloud/obs"
	"io/ioutil"
	"os"
	"strings"
)

type Download struct {
	bucketName string
	objectName  string
	location   string
	pkg        string
	obsClient  *obs.ObsClient
}

func NewDownload(config *viper.Viper, pkg string) *Download {
	accessKey := config.Get("bucket").(map[string]interface{})["accesskey"].(string)
	secretKey := config.Get("bucket").(map[string]interface{})["secretkey"].(string)
	endPoint := config.Get("bucket").(map[string]interface{})["endpoint"].(string)
	directory := config.Get("bucket").(map[string]interface{})["directory"].(string)
	bucketName := config.Get("bucket").(map[string]interface{})["bucketname"].(string)

	obsClient, err := obs.New(accessKey, secretKey, endPoint)
	if err != nil {
		panic(err)
	}

	return &Download{
		bucketName: bucketName,
		objectName:  strings.Join([]string{directory, pkg}, "/"),
		location:   endPoint,
		pkg: pkg,
		obsClient:  obsClient,
	}
}

func (ob *Download) Down() {
	input := &obs.GetObjectInput{}
	input.Bucket = ob.bucketName
	input.Key = ob.objectName

	output, err := ob.obsClient.GetObject(input)
	if err != nil {
		panic(err)
	}
	defer func() {
		errMsg := output.Body.Close()
		if errMsg != nil {
			panic(errMsg)
		}
	}()

	body, err := ioutil.ReadAll(output.Body)
	if err != nil {
		panic(err)
	}

	file, err := os.OpenFile("file/"+ob.pkg, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		panic(err)
	}
	file.Write(body)
}
