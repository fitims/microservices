package models

import "thaThrowdown/common/database/dgraph"

// VideosRequest contains the fields for the videos to be returned
type VideosRequest struct {
	Page   int      `json:"page"`
	Fields []string `json:"fields"`
}

// VideosForGenreRequest contains the genres to filter and fields for the videos to be returned
type VideosForGenreRequest struct {
	Page    int        `json:"page"`
	GenreID dgraph.UID `json:"genreId"`
	Fields  []string   `json:"fields"`
}

// VideosForArtistRequest contains the artist to filter and fields for the videos to be returned
type VideosForArtistRequest struct {
	Page   int      `json:"page"`
	Artist string   `json:"mediaArtist"`
	Fields []string `json:"fields"`
}

// VideosForTextRequest contains the artist/title to filter and fields for the videos to be returned
type VideosForTextRequest struct {
	Page    int        `json:"page"`
	Text    string     `json:"text"`
	GenreID dgraph.UID `json:"genreId, omitempty"`
	Fields  []string   `json:"fields"`
}

// DetailsRequest contains the fields for the video to be returned
type DetailsRequest struct {
	Fields []string `json:"fields"`
}

// UploadRequest contains the details about the video to be uploaded
type UploadRequest struct {
	Name            string       `json:"name"`
	Description     string       `json:"description"`
	MediaArtist     string       `json:"mediaArtist"`
	MediaLength     float32      `json:"mediaLength"`
	ArtworkURL      string       `json:"artworkUrl"`
	PlayURL         string       `json:"playUrl"`
	DownloadURL     string       `json:"downloadUrl"`
	Price           float32      `json:"price"`
	PurchaseEnabled bool         `json:"purchaseEnabled"`
	DownloadEnabled bool         `json:"downloadEnabled"`
	Genres          []dgraph.UID `json:"mediaGenres"`
	UploadedBy      dgraph.UID   `json:"uploadedBy"`
	IsFree          bool         `json:"isFree"`
	IsActive        bool         `json:"isActive"`
}

// UpdateRequest contains the details about the video to be updated
type UpdateRequest map[string]interface{}

// Token contains aws S3 token
type Token string

// TokenRequest is the request for getting S3 token
type TokenRequest struct {
	PreviewFile  string `json:"previewFile"`
	PurchaseFile string `json:"purchaseFile"`
}

// GenreRequest is the request for uploading a new genre
type GenreRequest struct {
	Name string `json:"name"`
}

// PurchaseRequest is the request for the user purchasing the video
type PurchaseRequest struct {
	Amount float32 `json:"amount"`
}
