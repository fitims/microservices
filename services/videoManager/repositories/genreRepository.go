package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"thaThrowdown/common/api"
	"thaThrowdown/common/database/dgraph"
	"thaThrowdown/services/videoManager/models"

	dgraphApi "github.com/dgraph-io/dgraph/protos/api"
)

const (
	getGenresByName = `{
		genres(func: eq(name, "%s")) @filter(eq(type, "video.genre")) {
		    uid
		    name
		}
	}`

	getAllGenres = `{
		genres(func: eq(type, "video.genre")) {
		    uid
		    name
		}
	}`
)

var (
	GenreAlreadyExistsErr = errors.New("Genre already exists")
)

type defaultGenreRepository struct{}

// NewGenreRepository acts as a factory for creating default GenreRepository implementation
func NewGenreRepository() GenreRepository {
	return defaultGenreRepository{}
}

// GetGenres return all the genres for the videos
func (r defaultGenreRepository) GetGenres() ([]api.J, error) {

	resp, err := dgraph.Client.NewTxn().Query(context.Background(), getAllGenres)

	if err != nil {
		log.Println("Error in getting response from server when reading genres : ", err)
		return nil, err
	}

	jsonResp := resp.GetJson()
	fmt.Println("Json Response : ", string(jsonResp))

	var qr GenresQueryResponse
	err = json.Unmarshal(jsonResp, &qr)
	if err != nil {
		log.Println("Error unmarshalling genres : ", err)
		return nil, err
	}

	return qr.Genres, nil
}

// AddGenre inserts a new genre into graph database
func (r defaultGenreRepository) AddGenre(g models.GenreRequest) (dgraph.UID, error) {

	ctx := context.Background()

	resp, err := dgraph.Client.NewTxn().Query(ctx, fmt.Sprintf(getGenresByName, g.Name))
	if err != nil {
		log.Println("Error in getting response from server when reading genres : ", err)
		return 0, err
	}

	var qry GenresQueryResponse
	err = json.Unmarshal(resp.GetJson(), &qry)
	if err != nil {
		log.Println("Error unmarshalling genres : ", err)
		return 0, err
	}

	if len(qry.Genres) > 0 {
		return 0, GenreAlreadyExistsErr
	}

	mutation := &dgraphApi.Mutation{
		CommitNow: true,
		SetNquads: g.ToMutation(),
	}

	assigned, err := dgraph.Client.NewTxn().Mutate(ctx, mutation)
	if err != nil {
		log.Println("Error in getting response from server when saving genre : ", err)
		return 0, err
	}

	return dgraph.ToUID(assigned.Uids["genre"])
}

// DeleteGenre deletes genre by its ID
func (r defaultGenreRepository) DeleteGenre(id dgraph.UID) error {

	mutation := &dgraphApi.Mutation{
		CommitNow: true,
		DelNquads: models.ToDeleteMutation(id),
	}

	_, err := dgraph.Client.NewTxn().Mutate(context.Background(), mutation)
	if err != nil {
		log.Println("Error in getting response from server when reading genres : ", err)
		return err
	}

	return nil
}
