package controllers

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"thaThrowdown/common/api"
	"thaThrowdown/common/database/dgraph"
	"thaThrowdown/common/infrastructure"
	"thaThrowdown/services/videoManager/models"
	"thaThrowdown/services/videoManager/repositories"

	"github.com/labstack/echo"
)

const (
	VIDEO_ID_PARAM     = "videoId"
	GENRE_NAME_PARAM  = "genreName"
	ARTIST_NAME_PARAM = "artistName"
	FILENAME          = "video_%d_%s"
)

// DefaultVideoController is the implementation fo the VideoController interface
type DefaultVideoController struct {
	repo repositories.VideoRepository
}

// NewVideoController is used as a factory to create a video an DefaultVideoController that iplements VideoController interface
func NewVideoController(r repositories.VideoRepository) VideoController {
	return DefaultVideoController{repo: r}
}

// GetVideos returns videos paginated
func (ctl DefaultVideoController) GetVideos(c echo.Context) error {

	var r models.VideosRequest
	err := c.Bind(&r)
	if err != nil {
		log.Println("Error binding the request : ", err)
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid data"))
	}

	total, videos, err := ctl.repo.GetVideos(r)
	if err != nil {
		log.Println("Error getting videos from database : ", err)
		return c.JSON(http.StatusInternalServerError, api.Fail("Internal Error"))
	}

	resp := buildGetVideoResponse(r.Page, total, videos)
	return c.JSON(http.StatusOK, resp)
}

// GetSingleVideo returns video details for the video id specified
func (ctl DefaultVideoController) GetSingleVideo(c echo.Context) error {

	idParam := c.Param(VIDEO_ID_PARAM)
	videoID, err := dgraph.ToUID(idParam)

	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid request"))
	}

	var r models.DetailsRequest
	err = c.Bind(&r)
	if err != nil {
		log.Println("Error binding the request : ", err)
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid data"))
	}

	log.Println("Id : ", videoID, ", request : ", r)

	video, err := ctl.repo.GetVideoDetails(videoID, r.Fields)
	if err != nil {
		if err == repositories.VideoNotFoundError {
			log.Println("Error getting video details from database : ", err)
			return c.JSON(http.StatusBadRequest, api.Fail("video does not exist"))
		}

		log.Println("Error getting video details from database : ", err)
		return c.JSON(http.StatusInternalServerError, api.Fail("Internal Error"))
	}

	return c.JSON(http.StatusOK, video)
}

// GetVideosByGenre  returns videos for the specified genre paginated
func (ctl DefaultVideoController) GetVideosByGenre(c echo.Context) error {

	var r models.VideosForGenreRequest
	err := c.Bind(&r)
	if err != nil {
		log.Println("Error binding the request : ", err)
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid data"))
	}

	total, videos, err := ctl.repo.GetVideosForGenre(r)
	if err != nil {
		log.Println("Error getting videos from database : ", err)
		return c.JSON(http.StatusInternalServerError, api.Fail("Internal Error"))
	}

	resp := buildGetVideoResponse(r.Page, total, videos)
	return c.JSON(http.StatusOK, resp)
}

// GetVideosByArtist  returns videos for the specified artist paginated
func (ctl DefaultVideoController) GetVideosByArtist(c echo.Context) error {

	var r models.VideosForArtistRequest
	err := c.Bind(&r)
	if err != nil {
		log.Println("Error binding the request : ", err)
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid data"))
	}

	total, videos, err := ctl.repo.GetVideosForArtist(r)
	if err != nil {
		log.Println("Error getting videos from database : ", err)
		return c.JSON(http.StatusInternalServerError, api.Fail("Internal Error"))
	}

	resp := buildGetVideoResponse(r.Page, total, videos)
	return c.JSON(http.StatusOK, resp)
}

// GetVideosForText  returns videos for the specified text paginated
func (ctl DefaultVideoController) GetVideosForText(c echo.Context) error {

	var r models.VideosForTextRequest
	err := c.Bind(&r)
	if err != nil {
		log.Println("Error binding the request : ", err)
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid data"))
	}

	total, videos, err := ctl.repo.GetVideosForText(r)
	if err != nil {
		log.Println("Error getting videos from database : ", err)
		return c.JSON(http.StatusInternalServerError, api.Fail("Internal Error"))
	}

	resp := buildGetVideoResponse(r.Page, total, videos)
	return c.JSON(http.StatusOK, resp)
}

// GetUploadToken will return the upload tokens to be used by Amazon S3
func (ctl DefaultVideoController) GetUploadToken(c echo.Context) error {
	idParam := c.Param(VIDEO_ID_PARAM)
	videoID, err := dgraph.ToUID(idParam)

	fmt.Println("id Param : ", idParam)
	fmt.Println("videoID : ", videoID)

	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid request"))
	}

	isValid, err := ctl.repo.IsVideoValid(videoID)
	if err != nil {
		log.Println("Error getting videos from database : ", err)
		return c.JSON(http.StatusInternalServerError, api.Fail("Internal Error"))
	}

	if isValid == false {
		return c.JSON(http.StatusBadRequest, api.Fail("video does not exist"))
	}

	// get token request
	var t models.TokenRequest

	err = c.Bind(&t)
	if err != nil {
		log.Println("Error binding the token request : ", err)
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid data"))
	}

	// get the token from Amazon S3
	token, err := ctl.getS3Tokens(videoID, t)
	if err != nil {
		log.Println("There was an error getting token : ", err)
		return c.JSON(http.StatusInternalServerError, api.Fail("Internal Error"))
	}
	return c.JSON(http.StatusOK, api.SuccessWithPayload("Token retrieved successfully", api.J{"token": token}))
}

// UploadArtwork uploads artwork for the video
func (ctl DefaultVideoController) UploadArtwork(c echo.Context) error {
	idParam := c.Param(VIDEO_ID_PARAM)
	videoID, err := dgraph.ToUID(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid request"))
	}

	file, err := c.FormFile("file")
	if err != nil {
		log.Println("Cannot get the file from the request. Error : ", err)
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid request. Cannot get file"))
	}

	src, err := file.Open()
	if err != nil {
		log.Println("Cannot open file. Error : ", err)
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid request. Cannot read file"))
	}
	defer src.Close()

	sourceFile := buildFilePath(file.Filename)
	dst, err := os.Create(sourceFile)
	if err != nil {
		log.Println("Cannot create a file. Error : ", err)
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid request. Cannot write file"))
	}
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		log.Println("Cannot copy file content. Error : ", err)
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid request. Cannot copy file content"))
	}

	resizedFile := path.Join("/home/artwork", fmt.Sprintf(FILENAME, videoID, file.Filename))
	go resizeImageAndDeleteSource(sourceFile, resizedFile, 360, 360)

	artworkURL := path.Join("/public/artwork/", fmt.Sprintf(FILENAME, videoID, file.Filename))
	err = ctl.repo.SetArtworkURL(videoID, artworkURL)
	if err != nil {
		log.Println("Cannot set artwork URL to database file. Error : ", err)
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid request. Cannot update video"))
	}
	return c.JSON(http.StatusOK, api.SuccessWithPayload("Artwork uploaded successfully", api.J{"artworkUrl": artworkURL}))
}

// UploadNewVideo uploads new video details
func (ctl DefaultVideoController) UploadNewVideo(c echo.Context) error {

	var s models.UploadRequest
	err := c.Bind(&s)
	if err != nil {
		log.Println("Error binding the video : ", err)
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid data"))
	}

	id, err := ctl.repo.UploadVideo(s)
	if err != nil {
		log.Println("Error saving the video :", err)
		return c.JSON(http.StatusBadRequest, api.Fail("Could not save the video"))
	}

	return c.JSON(http.StatusOK, api.SuccessWithID("Video uploaded successfully", id))
}

// UpdateVideo updates existing video details
func (ctl DefaultVideoController) UpdateVideo(c echo.Context) error {
	idParam := c.Param("videoId")
	videoID, err := dgraph.ToUID(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid request"))
	}

	// get the video
	var req models.UpdateRequest
	err = c.Bind(&req)
	if err != nil {
		log.Println("Error binding the video : ", err)
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid data"))
	}

	err = ctl.repo.UpdateVideo(videoID, req)
	if err != nil {
		log.Println("Error saving the video :", err)
		return c.JSON(http.StatusBadRequest, api.Fail("Could not update the video"))
	}

	return c.JSON(http.StatusOK, api.Success("Video updated successfully"))
}

// DeleteVideo deletes existing video (the video in reality the video details are never deleted, it si just marked as deleted)
func (ctl DefaultVideoController) DeleteVideo(c echo.Context) error {
	idParam := c.Param("videoId")
	videoID, err := dgraph.ToUID(idParam)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid request"))
	}

	err = ctl.repo.DeleteVideo(videoID)
	if err != nil {
		log.Println("Error storing purchase data : ", err)
		return c.JSON(http.StatusInternalServerError, api.Fail("Could not delete video"))
	}

	return c.JSON(http.StatusOK, api.Success("Video deleted successfully"))
}

func (ctl DefaultVideoController) getS3Tokens(videoID dgraph.UID, tr models.TokenRequest) (models.TokenResponse, error) {
	response := models.TokenResponse{}
	prevToken, url, err := getS3Token(videoID, infrastructure.MediaType("videos/preview"), tr.PreviewFile)
	if err != nil {
		return response, err
	}
	response.Preview = models.TokenDetails{
		Token: prevToken,
		URL:   url,
	}

	purchaseToken, url, err := getS3Token(videoID, infrastructure.MediaType("videos/purchase"), tr.PurchaseFile)
	if err != nil {
		return response, err
	}
	response.Purchase = models.TokenDetails{
		Token: purchaseToken,
		URL:   url,
	}

	return response, nil
}
