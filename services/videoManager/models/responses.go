package models

import "thaThrowdown/common/api"

type Paging struct {
	Page  int `json:"page"`
	Total int `json:"total"`
}

// VideosResponse is the response for
type VideosResponse struct {
	Paging Paging  `json:"paging"`
	Videos  []api.J `json:"videos"`
}

// TokenDetails encapsulates the response for each token
type TokenDetails struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

// TokenResponse is the response for S3 token
type TokenResponse struct {
	Preview  TokenDetails `json:"preview"`
	Purchase TokenDetails `json:"purchase"`
}


