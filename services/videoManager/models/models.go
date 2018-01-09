package models

import "thaThrowdown/common/database/dgraph"

// VideoGenre contains the genre of video
type VideoGenre struct {
	ID   dgraph.UID `json:"uid, omitempty"`
	Name string     `json:"name, omitempty"`
}

// VideoGenres encapsulates a collection of video genres
type VideoGenres []VideoGenre

// Video encapsulates the details about the video
type Video struct {
	ID              dgraph.UID   `json:"uid, omitempty"`
	Name            string       `json:"name, omitempty"`
	Artist          string       `json:"mediaArtist, omitempty"`
	Description     string       `json:"description, omitempty"`
	Genres          []VideoGenre `json:"mediaGenre, omitempty"`
	ArtworkURL      string       `json:"artworkUrl, omitempty"`
	PlayURL         string       `json:"playUrl, omitempty"`
	DownloadURL     string       `json:"downloadUrl, omitempty"`
	Price           float32      `json:"price, omitempty"`
	DownloadEnabled bool         `json:"downloadEnabled, omitempty"`
	IsFree          bool         `json:"isFree, omitempty"`
	Length          float32      `json:"mediaLength, omitempty"`
	IsActive        bool         `json:"isActive, omitempty"`
}
