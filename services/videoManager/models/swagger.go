package models

import (
	"mime/multipart"
)

// VideoRequest contains the details about the video to be uploaded
// swagger:parameters UploadNewVideo
//type VideoRequestWrapper struct {
//	// in:body
//	// required: true
//	Request VideoRequest
//}

// UpdateRequest contains the details about the video to be updated
// swagger:parameters UpdateVideo
type UpdateRequestWrapper struct {
	// in:body
	// required: true
	Request UpdateRequest
}

// swagger:parameters GetVideosByArtist
type ArtistNameRequestWrapper struct {
	// in: path
	// required: true
	ArtistName string `json:"artistName"`
}

// swagger:parameters GetVideosByGenre
type GenreNameRequestWrapper struct {
	// in: path
	// required: true
	GenreId string `json:"genreName"`
}

// swagger:parameters GetSingleVideo UpdateVideo DeleteVideo
type VideoDetailRequestWrapper struct {
	// in: path
	// required: true
	VideoId string `json:"videoId"`
}

// swagger:parameters PurchaseVideo
//type SellVideoRequestWrapper struct {
//
//	// in: path
//	// required: true
//	VideoId string `json:"videoId"`
//
//	// in:body
//	// required: true
//	Request SellVideoRequest `json:"request"`
//}

// swagger:parameters GetUploadToken
type UploadTokenRequestWrapper struct {

	// in: path
	// required: true
	VideoId string `json:"videoId"`

	// in:body
	// required: true
	Request TokenRequest `json:"request"`
}

type UploadTokenWithPayloadResponse struct {
	Data struct {
		Token string `json:"token"`
	} `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// swagger:response UploadTokenSuccessWithPayloadResponseWrapper
type UploadTokenSuccessWithPayloadResponseWrapper struct {
	// in:body
	// required: true
	Response UploadTokenWithPayloadResponse
}

// VideoResponse is the response sent to the user and holds paging info
// swagger:response VideoResponseWrapper
//type VideoResponseWrapper struct {
//	// in:body
//	// required: true
//	Response VideoResponse `json:"response"`
//}

// VideoResponse is the response sent to the user and holds paging info
// swagger:response VideoSummaryWrapper
//type VideoSummaryWrapper struct {
//	// in:body
//	// required: true
//	Response VideoSummary `json:"response"`
//}

type VideoSuccessWithPayloadResponse struct {
	Data struct {
		GenreId string `json:"videoId"`
	} `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// File represents an uploaded file.
type File struct {
	Data   multipart.File
	Header *multipart.FileHeader
}

func (f File) Read(p []byte) (n int, err error) {
	return f.Data.Read(p)
}

// swagger:parameters UploadArtwork
type ArtworkRequestWrapper struct {

	// in: formData
	Description *string

	// in: formData
	File File

	// in: path
	// required: true
	VideoId string `json:"videoId"`
}

// swagger:response VideoSuccessWithPayloadResponseWrapper
type VideoSuccessWithPayloadResponseWrapper struct {
	// in:body
	// required: true
	Response VideoSuccessWithPayloadResponse
}

type ArtworkSuccessWithPayloadResponse struct {
	Data struct {
		ArtworkUrl string `json:"artworkUrl"`
	} `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// swagger:response ArtworkSuccessWithPayloadResponseWrapper
type ArtworkSuccessWithPayloadResponseWrapper struct {
	// in:body
	// required: true
	Response ArtworkSuccessWithPayloadResponse
}

// swagger:response VideoSuccessWithIdResponse
type VideoSuccessWithIdResponse struct {
}

// GenreRequest is the request for uploading a new genre
// swagger:parameters UploadGenre
type GenreRequestWrapper struct {
	// in: body
	// required: true
	Request GenreRequest `json:"request"`
}

// swagger:parameters DeleteGenre
type DeleteGenreRequestWrapper struct {
	// in: path
	// required: true
	GenreId string `json:"genreId"`
}

// MusicGenres encapsulates a collection of music genres
// swagger:response MusicGenresWrapper
//type MusicGenresWrapper struct {
//	// in:body
//	// required: true
//	Response MusicGenres `json:"response"`
//}

type MessageResponse struct {
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// swagger:response FailResponseWrapper
type FailResponseWrapper struct {
	// in:body
	// required: true
	Response MessageResponse `json:"response"`
}

// swagger:response SuccessResponseWrapper
type SuccessResponseWrapper struct {
	// in:body
	// required: true
	Response MessageResponse
}

type GenreSuccessWithPayloadResponse struct {
	Data struct {
		GenreId string `json:"genreId"`
	} `json:"data"`
	Message string `json:"message"`
	Success bool   `json:"success"`
}

// swagger:response GenreSuccessWithPayloadResponseWrapper
type GenreSuccessWithPayloadResponseWrapper struct {
	// in:body
	// required: true
	Response GenreSuccessWithPayloadResponse
}

// swagger:response GenreSuccessWithIdResponse
type GenreSuccessWithIdResponse struct {
}
