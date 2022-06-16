package S3AWS

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/jasonbronson/kwik-cms-engine/config"
)

func S3FileUpload(filename, folder string) (string, error) {
	//create session s3
	session, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.Cfg.MediaBucketRegion),
		Credentials: credentials.NewStaticCredentials(config.Cfg.AWSAccessKey, config.Cfg.AWSSecret, ""),
	})
	if err != nil {
		log.Println(err)
		return "", err
	}

	//Upload to S3
	hashedFilename, err := uploadFile(session, filename, folder)
	if err != nil {
		log.Println("Cannot upload file to s3", err)
		return "", err
	}

	//Get file link url
	var folderPath = string(folder + "/")
	if folder == "" {
		folderPath = ""
	}
	s3URL := getS3FileLink(session, hashedFilename, folderPath)
	if err != nil {
		log.Println("S3AWS:", err)
		return "", err
	}

	return s3URL, nil
}

func uploadFile(session *session.Session, uploadFileDir, folder string) (string, error) {

	log.Println(uploadFileDir)
	upFile, err := os.Open(uploadFileDir)
	if err != nil {
		return "", err
	}
	defer upFile.Close()

	var folderPath = ""
	if folder != "" {
		folderPath = fmt.Sprintf("/%s/", folder)
	}

	upFileInfo, _ := upFile.Stat()
	var fileSize int64 = upFileInfo.Size()
	fileBuffer := make([]byte, fileSize)
	upFile.Read(fileBuffer)
	hashedFilename := HashFilename(uploadFileDir)
	_, err = s3.New(session).PutObject(&s3.PutObjectInput{
		ACL:           aws.String("bucket-owner-full-control"),
		Bucket:        aws.String(config.Cfg.MediaBucketName + folderPath),
		Key:           aws.String(hashedFilename),
		Body:          bytes.NewReader(fileBuffer),
		ContentLength: aws.Int64(fileSize),
		ContentType:   aws.String(http.DetectContentType(fileBuffer)),
	})
	return hashedFilename, err
}

func getS3FileLink(session *session.Session, key, folder string) string {
	return config.Cfg.MediaBucketURL + "/" + folder + key
}

func HashFilename(filename string) string {
	reg, err := regexp.Compile("[^a-zA-Z0-9]+")
	if err != nil {
		log.Fatal(err)
	}
	//remove extension from filename
	filenameSplit := strings.Split(filename, ".")
	newFilename := reg.ReplaceAllString(filenameSplit[0], "")
	newFilename = newFilename + "." + filenameSplit[1]
	return newFilename
}
