package models

import (
	"bytes"
	"fmt"
	"strings"
	"thaThrowdown/common/database/dgraph"
	"time"
)

const (
	deleteItem = `<0x%x> * * . `

	addNewGenre = `
		_:genre <name> "%s" .
		_:genre <type> "video.genre" .
	`

	addNewVideo = `
		_:video <name> "%s" .
		_:video <description> "%s" .
		_:video <mediaArtist> "%s" .		
		_:video <mediaLength> "%f" .
		_:video <artworkUrl> "%s" .
		_:video <playUrl> "%s" .
		_:video <downloadUrl> "%s" .
		_:video <price> "%f" .
		_:video <purchaseEnabled> "%t" .
		_:video <downloadEnabled> "%t" .		
		_:video <uploadedBy>  <0x%x> .
		_:video <isActive> "%t" .
		_:video <isDeleted> "false" .
		_:video <isFree> "%t" .
		_:video <uploadedAt> "%s" .
		_:video <type> "video" .

		%s
	`
)

// ToDeleteMutation converts the id into delete NQuad
func ToDeleteMutation(id dgraph.UID) []byte {
	return []byte(fmt.Sprintf(deleteItem, id))
}

// ToMutation converts GenreRequest to Mutation to be used by DGraph
func (r *UploadRequest) ToMutation() []byte {

	m := func(ids []dgraph.UID) []string {
		r := make([]string, len(ids))

		for i, id := range ids {
			r[i] = fmt.Sprintf(`_:video <mediaGenre> <0x%x> . `, id)
		}

		return r
	}

	genres := strings.Join(m(r.Genres), "\n")

	return []byte(fmt.Sprintf(addNewVideo,
		r.Name,
		r.Description,
		r.MediaArtist,
		r.MediaLength,
		r.ArtworkURL,
		r.PlayURL,
		r.DownloadURL,
		r.Price,
		r.PurchaseEnabled,
		r.DownloadEnabled,
		r.UploadedBy,
		r.IsActive,
		r.IsFree,
		time.Now().Format(dgraph.DATE_FORMAT),
		genres))
}

func (r UpdateRequest) ToMutation(id dgraph.UID) ([]byte, []byte, error) {

	var updateBuffer bytes.Buffer
	var deleteBuffer bytes.Buffer

	for k, v := range r {

		switch k {
		case "name", "description", "mediaArtist", "mediaLength", "artworkUrl", "playUrl", "downloadUrl", "price",
			"purchaseEnabled", "downloadEnabled", "isActive", "isFree":
			updateBuffer.WriteString(fmt.Sprintf(`<0x%x> <%s> "%v" . `+"\n", id, k, v))

		case "mediaGenres":
			// check if mediaGenres is an array
			if a, ok := v.([]interface{}); ok {
				deleteBuffer.WriteString(fmt.Sprintf("<0x%x> <mediaGenre> * . ", id))
				for _, i := range a {
					if g, ok := i.(string); ok {
						uid, err := dgraph.ToUID(g)
						if err != nil {
							return nil, nil, err
						}
						updateBuffer.WriteString(fmt.Sprintf("<0x%x> <mediaGenre> <0x%x> . \n", id, uid))
					}
				}
			} else if g, ok := v.(string); ok {
				// if mediaGenre has a single value
				uid, err := dgraph.ToUID(g)
				if err != nil {
					return nil, nil, err
				}
				updateBuffer.WriteString(fmt.Sprintf("<0x%x> <mediaGenre> <0x%x> . \n", id, uid))
				deleteBuffer.WriteString(fmt.Sprintf("<0x%x> <mediaGenre> * . ", id))
			}
		}
	}

	return updateBuffer.Bytes(), deleteBuffer.Bytes(), nil

}
