package controllers

import (
	"errors"
	"log"
	"os"
	"path"
	"thaThrowdown/common/api"
	"thaThrowdown/common/database/dgraph"
	"thaThrowdown/common/infrastructure"
	"thaThrowdown/lib/imageManipulator/images"
	"thaThrowdown/services/videoManager/models"
)

const (

	// AMAZON_S3_BASE_URL is the base URL for amazon S3
	AMAZON_S3_BASE_URL = "https://s3-eu-west-1.amazonaws.com/thathrowdownmedia"
)

func getS3Token(videoID dgraph.UID, fileType infrastructure.MediaType, fileName string) (string, string, error) {
	if len(fileName) > 0 {
		s3key := infrastructure.BuildKey(videoID, fileType, fileName)
		url := AMAZON_S3_BASE_URL + "/" + s3key

		token, err := infrastructure.GetS3SignedRequest(s3key)
		if err != nil {
			return "", "", err
		}
		return token, url, nil
	}
	return "", "", errors.New("Filename is empty")
}

func resizeImageAndDeleteSource(originalImage, resizedImage string, width, height uint) error {
	err := images.ResizeImage(originalImage, resizedImage, width, height)
	if err != nil {
		return err
	}
	log.Println("Artwork - file resized successfully. Removing source file : ", originalImage)

	return os.Remove(originalImage)
}

func buildFilePath(filename string) string {
	return path.Join("/home/artwork", filename)
}

func buildGetVideoResponse(page int, total int, videos []api.J) models.VideosResponse {
	r := models.VideosResponse{
		Paging: models.Paging{
			Page:  page,
			Total: total / models.PAGESIZE,
		},
		Videos: videos,
	}
	return r
}
