package repositories

import (
	"thaThrowdown/common/api"
	"thaThrowdown/common/database/dgraph"
	"thaThrowdown/services/videoManager/models"
)

// GenreRepository wraps the functions used for communicating with the database
type GenreRepository interface {
	GetGenres() ([]api.J, error)
	AddGenre(g models.GenreRequest) (dgraph.UID, error)
	DeleteGenre(id dgraph.UID) error
}

// VideoRepository wraps the functions used for communicating with the database
type VideoRepository interface {
	GetVideos(request models.VideosRequest) (int, []api.J, error)
	GetVideoDetails(videoID dgraph.UID, fields []string) (api.J, error)
	IsVideoValid(videoID dgraph.UID) (bool, error)
	GetVideosForGenre(request models.VideosForGenreRequest) (int, []api.J, error)
	GetVideosForArtist(request models.VideosForArtistRequest) (int, []api.J, error)
	GetVideosForText(request models.VideosForTextRequest) (int, []api.J, error)
	UploadVideo(video models.UploadRequest) (dgraph.UID, error)
	UpdateVideo(id dgraph.UID, video models.UpdateRequest) error
	DeleteVideo(id dgraph.UID) error
	SetPlayURL(videoID dgraph.UID, playURL string) error
	SetDownloadURL(videoID dgraph.UID, fileURL string) error
	SetArtworkURL(videoID dgraph.UID, fileURL string) error
}
