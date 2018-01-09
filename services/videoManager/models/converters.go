package models

import (
	"fmt"
	"strings"
)

// PAGESIZE shows the number of items per page
const PAGESIZE int = 15

const (
	getVideos = `{
		total(func: eq(type, "video")) @filter(eq(isActive, "true")) {
			count(uid)
		}

		videos(func: eq(type, "video"), first:%d) @filter(eq(isActive, "true")) {
			%s
		}
	}`

	getVideosOffset = `{
		total(func: eq(type, "video")) {
			count(uid)
		}

		videos(func: eq(type, "video"), first:%d, offset:%d) @filter(eq(isActive, "true")) {
			%s
		}
	}`

	getVideosForGenre = `{
		total(func: eq(type, "video")) @filter(uid_in(mediaGenre, 0x%x) AND eq(isActive, "true")) {
			count(uid)
		}

		videos(func: eq(type, "video"), first:%d) @filter(uid_in(mediaGenre, 0x%x) AND eq(isActive, "true")) {
			%s
		}
	}`

	getVideosForGenreOffset = `{
		total(func: eq(type, "video")) @filter(uid_in(mediaGenre, 0x%x) AND eq(isActive, "true"))  {
			count(uid)
		}

		videos(func: eq(type, "video"), first:%d, offset:%d) @filter(uid_in(mediaGenre, 0x%x) AND eq(isActive, "true"))  {
			%s
		}
	}`

	getVideosForArtist = `{
		total(func: eq(type, "video")) @filter(anyofterms(mediaArtist,"%s") AND eq(isActive, "true")) {
			count(uid)
		}

		videos(func: eq(type, "video"), first:%d) @filter(anyofterms(mediaArtist,"%s") AND eq(isActive, "true")) {
			%s
		}
	}`

	getVideosForArtistOffset = `{
		total(func: eq(type, "video")) @filter(anyofterms(mediaArtist,"%s") AND eq(isActive, "true")) {
			count(uid)
		}

		videos(func: eq(type, "video"), first:%d, offset:%d) @filter(anyofterms(mediaArtist,"%s") AND eq(isActive, "true")) {
			%s
		}
	}`

	getVideosTextSearch = `{
		total(func: eq(type, "video")) @filter((anyofterms(name, "%s") OR anyofterms(mediaArtist, "%s")) AND eq(isActive, "true") %s) {
			count(uid)
		}

		videos(func: eq(type, "video"), first:%d) @filter((anyofterms(name, "%s") OR anyofterms(mediaArtist, "%s")) AND eq(isActive, "true") %s) {
			%s
		}
	}`

	getVideosTextSearchOffset = `{
		total(func: eq(type, "video")) @filter((anyofterms(name, "%s") OR anyofterms(mediaArtist, "%s")) AND eq(isActive, "true") %s) {
			count(uid)
		}

		videos(func: eq(type, "video"), first:%d, offset:%d) @filter((anyofterms(name, "%s") OR anyofterms(mediaArtist, "%s")) AND eq(isActive, "true") %s) {
			%s
		}
	}`
)

// ToMutation converts GenreRequest to Mutation to be used by DGraph
func (r *GenreRequest) ToMutation() []byte {
	return []byte(fmt.Sprintf(addNewGenre, r.Name))
}

// ToQuery converts GetVideosRequest to Query to be used by DGraph
func (r *VideosRequest) ToQuery() string {

	fields := BuildFieldsQuery(r.Fields)

	if r.Page < 2 {
		return fmt.Sprintf(getVideos, PAGESIZE, fields)
	}
	return fmt.Sprintf(getVideosOffset, PAGESIZE, PAGESIZE*(r.Page-1), fields)
}

// ToQuery converts VideosForGenreRequest to Query to be used by DGraph
func (r *VideosForGenreRequest) ToQuery() string {

	fields := BuildFieldsQuery(r.Fields)

	if r.Page < 2 {

		return fmt.Sprintf(getVideosForGenre, r.GenreID, PAGESIZE, r.GenreID, fields)
	}
	return fmt.Sprintf(getVideosForGenreOffset, r.GenreID, PAGESIZE, PAGESIZE*(r.Page-1), r.GenreID, fields)
}

// ToQuery converts VideosForArtistRequest to Query to be used by DGraph
func (r *VideosForArtistRequest) ToQuery() string {
	fields := BuildFieldsQuery(r.Fields)

	if r.Page < 2 {
		return fmt.Sprintf(getVideosForArtist, r.Artist, PAGESIZE, r.Artist, fields)
	}
	return fmt.Sprintf(getVideosForArtistOffset, r.Artist, PAGESIZE, PAGESIZE*(r.Page-1), r.Artist, fields)
}

// ToQuery converts VideosForArtistRequest to Query to be used by DGraph
func (r *VideosForTextRequest) ToQuery() string {
	fields := BuildFieldsQuery(r.Fields)

	var genreFilter string
	if r.GenreID != 0 {
		genreFilter = fmt.Sprintf("AND (uid_in(mediaGenre, 0x%x))", r.GenreID)
	}

	if r.Page < 2 {
		return fmt.Sprintf(getVideosTextSearch, r.Text, r.Text, genreFilter, PAGESIZE, r.Text, r.Text, genreFilter, fields)
	}
	return fmt.Sprintf(getVideosTextSearchOffset, r.Text, r.Text, genreFilter, PAGESIZE, PAGESIZE*(r.Page-1), r.Text, r.Text, genreFilter, fields)
}


func BuildFieldsQuery(fields []string) string {
	var query string

	var noOfFields = len(fields)
	if noOfFields < 1 {
		query = "uid name"
	} else if noOfFields == 1 {
		query = fields[0]
	} else {
		f := make([]string, 0)
		for _, v := range fields {
			switch v {
			case "mediaGenre" :
				f = append(f, "mediaGenre { uid name }\n")
			default:
				f = append(f, v)
			}
		}

		query = strings.Join(f, " ")
	}

	return query
}