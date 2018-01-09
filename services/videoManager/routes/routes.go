package routes

import (
	"thaThrowdown/services/videoManager/controllers"
	"thaThrowdown/services/videoManager/wrappers"
)

const (
	GET_VIDEOS_ROUTE           = ""
	GET_SINGLE_VIDEO_ROUTE     = "/:videoId"
	GET_VIDEOS_BY_GENRE_ROUTE  = "/genre"
	GET_VIDEOS_BY_ARTIST_ROUTE = "/artist"
	GET_VIDEOS_BY_TEXT_ROUTE   = "/text"
	GET_UPLOAD_TOKEN_ROUTE    = "/:videoId/token"
	UPLOAD_ARTWORK_ROUTE      = "/:videoId/artwork"
	UPLOAD_NEW_VIDEO_ROUTE     = "/new"
	UPDATE_VIDEO_ROUTE         = "/:videoId"
	DELETE_VIDEO_ROUTE         = "/:videoId"

	GET_GENRES_ROUTE       = "/genres"
	UPLOAD_NEW_GENRE_ROUTE = "/genres"
	DELETE_GENRE_ROUTE     = "/genres/:genreId"
)

// DefaultRouter is the default implementation of Router interface
type DefaultRouter struct {
	group wrappers.GroupWrapper
}

// NewRouter acts as a factory for creating the default router
func NewRouter(grp wrappers.GroupWrapper) Router {
	return DefaultRouter{group: grp}
}

// SetVideoRoutes sets the endpoints for managing video details
func (r DefaultRouter) SetVideoRoutes(c controllers.VideoController) {
	r.group.POST(GET_VIDEOS_ROUTE, c.GetVideos)
	r.group.POST(GET_SINGLE_VIDEO_ROUTE, c.GetSingleVideo)
	r.group.POST(GET_VIDEOS_BY_GENRE_ROUTE, c.GetVideosByGenre)
	r.group.POST(GET_VIDEOS_BY_ARTIST_ROUTE, c.GetVideosByArtist)
	r.group.POST(GET_VIDEOS_BY_TEXT_ROUTE, c.GetVideosForText)
	r.group.POST(GET_UPLOAD_TOKEN_ROUTE, c.GetUploadToken)
	r.group.POST(UPLOAD_ARTWORK_ROUTE, c.UploadArtwork)
	r.group.POST(UPLOAD_NEW_VIDEO_ROUTE, c.UploadNewVideo)
	r.group.PUT(UPDATE_VIDEO_ROUTE, c.UpdateVideo)
	r.group.DELETE(DELETE_VIDEO_ROUTE, c.DeleteVideo)
}

// SetGenreRoutes sets the endpoints  for managing the video genres
func (r DefaultRouter) SetGenreRoutes(c controllers.GenreController) {
	r.group.GET(GET_GENRES_ROUTE, c.GetGenres)
	r.group.POST(UPLOAD_NEW_GENRE_ROUTE, c.UploadGenre)
	r.group.DELETE(DELETE_GENRE_ROUTE, c.DeleteGenre)
}
