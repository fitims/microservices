package controllers

import "github.com/labstack/echo"

// VideoController captures actions for video
type VideoController interface {
	GetVideos(c echo.Context) error
	GetSingleVideo(c echo.Context) error
	GetVideosByGenre(c echo.Context) error
	GetVideosByArtist(c echo.Context) error
	GetVideosForText(c echo.Context) error
	GetUploadToken(c echo.Context) error
	UploadArtwork(c echo.Context) error
	UploadNewVideo(c echo.Context) error
	UpdateVideo(c echo.Context) error
	DeleteVideo(c echo.Context) error
}

// GenreController captures actions for video genre
type GenreController interface {
	GetGenres(c echo.Context) error
	UploadGenre(c echo.Context) error
	DeleteGenre(c echo.Context) error
}
