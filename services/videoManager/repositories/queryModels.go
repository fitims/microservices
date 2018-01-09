package repositories

import (
	"thaThrowdown/common/api"
)

type GenresQueryResponse struct {
	Genres []api.J `json:"genres"`
}

type CountItemsResponse struct {
	Count int `json:"count"`
}

type VideosExistsResponse struct {
	Total []CountItemsResponse `json:"total"`
}

type VideosQueryResponse struct {
	Total []CountItemsResponse `json:"total"`
	Videos []api.J              `json:"videos"`
}

type VideoDetailQueryResponse struct {
	Video []api.J `json:"video"`
}
