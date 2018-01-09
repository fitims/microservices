package repositories

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"thaThrowdown/common/api"
	"thaThrowdown/common/database/dgraph"
	"thaThrowdown/services/videoManager/models"

	dgraphApi "github.com/dgraph-io/dgraph/protos/api"
)

const (
	getVideoDetails = `{ video(func: uid(0x%x)) @filter(eq(isDeleted, "false")) { %s } }`
	doesVideoExist  = `{ total(func: uid(0x%x)) @filter(eq(type, "video") AND eq(isActive, "true")) { count(uid) } }`
)

var (
	VideoNotFoundError = errors.New("video does not exist")
	UnmarshalError    = errors.New("cannot unmarshal data")
	ParseError        = errors.New("parsing error")
	DatabaseError     = errors.New("database error")
)

type defaultVideoRepository struct{}

// NewVideoRepository acts as a factory for the default implementation of VideoRepository interface
func NewVideoRepository() VideoRepository {
	return defaultVideoRepository{}
}

// GetVideoDetails returns the details of the specified video from the dgraph
func (r defaultVideoRepository) GetVideoDetails(videoID dgraph.UID, fields []string) (api.J, error) {

	query := fmt.Sprintf(getVideoDetails, videoID, models.BuildFieldsQuery(fields))

	resp, err := dgraph.Client.NewTxn().Query(context.Background(), query)
	if err != nil {
		log.Println("Error in getting response from server when reading genres : ", err)
		return api.J{}, DatabaseError
	}

	jsonResp := resp.GetJson()

	var qr VideoDetailQueryResponse
	err = json.Unmarshal(jsonResp, &qr)
	if err != nil {
		log.Println("Error unmarshalling genres : ", err)
		return api.J{}, UnmarshalError
	}

	if len(qr.Video) == 0 {
		return api.J{}, VideoNotFoundError
	}

	return qr.Video[0], nil
}

// IsVideoValid returns true if the video exists and is not deleted
func (r defaultVideoRepository) IsVideoValid(videoID dgraph.UID) (bool, error) {
	query := fmt.Sprintf(doesVideoExist, videoID)

	resp, err := dgraph.Client.NewTxn().Query(context.Background(), query)
	if err != nil {
		log.Println("Error in getting response from server when reading genres : ", err)
		return false, DatabaseError
	}

	jsonResp := resp.GetJson()

	var qr VideosExistsResponse
	err = json.Unmarshal(jsonResp, &qr)
	if err != nil {
		log.Println("Error unmarshalling genres : ", err)
		return false, UnmarshalError
	}

	return qr.Total[0].Count > 0, nil
}

// GetVideos returns all the videos from the dgraph
func (r defaultVideoRepository) GetVideos(request models.VideosRequest) (int, []api.J, error) {

	query := request.ToQuery()
	return r.getVideosFromRepo(query)
}

// GetVideosForGenre returns all the videos for the genres provided from the mongo
func (r defaultVideoRepository) GetVideosForGenre(request models.VideosForGenreRequest) (int, []api.J, error) {

	query := request.ToQuery()
	return r.getVideosFromRepo(query)
}

// GetVideosForArtist returns all the videos for the artist provided from the mongo
func (r defaultVideoRepository) GetVideosForArtist(request models.VideosForArtistRequest) (int, []api.J, error) {
	query := request.ToQuery()
	return r.getVideosFromRepo(query)
}

// GetVideosForText returns all the videos for the artist provided from the mongo
func (r defaultVideoRepository) GetVideosForText(request models.VideosForTextRequest) (int, []api.J, error) {
	query := request.ToQuery()
	return r.getVideosFromRepo(query)
}

// UploadVideo adds a new video to the mongo
func (r defaultVideoRepository) UploadVideo(video models.UploadRequest) (dgraph.UID, error) {

	mutation := &dgraphApi.Mutation{
		CommitNow: true,
		SetNquads: video.ToMutation(),
	}

	assigned, err := dgraph.Client.NewTxn().Mutate(context.Background(), mutation)
	if err != nil {
		log.Println("Error in getting response from server when saving video : ", err)
		return 0, DatabaseError
	}

	id, err := strconv.ParseUint(strings.TrimPrefix(assigned.Uids["video"], "0x"), 16, 64)
	if err != nil {
		log.Println("Error getting newly created UID for the saved video : ", err)
		return 0, ParseError
	}

	return dgraph.UID(id), nil
}

// UpdateVideo adds a new video to the mongo
func (r defaultVideoRepository) UpdateVideo(id dgraph.UID, video models.UpdateRequest) error {

	u, d, err := video.ToMutation(id)
	if err != nil {
		return err
	}

	fmt.Println("\n\nTo delete mutation : \n\n", string(d), "\n\n")
	fmt.Println("\n\nTo update mutation : \n\n", string(u), "\n\n")

	if len(d) > 0 {
		delMutation := &dgraphApi.Mutation{
			CommitNow: true,
			DelNquads: d,
		}

		_, err := dgraph.Client.NewTxn().Mutate(context.Background(), delMutation)
		if err != nil {
			log.Println("Error in getting response from server when deleting mediaGenres : ", err)
			return err
		}
	}

	mutation := &dgraphApi.Mutation{
		CommitNow: true,
		SetNquads: u,
	}

	_, err = dgraph.Client.NewTxn().Mutate(context.Background(), mutation)
	if err != nil {
		log.Println("Error in getting response from server when saving video : ", err)
		return err
	}

	return nil
}

// DeleteVideo deletes video by its ID
func (r defaultVideoRepository) DeleteVideo(id dgraph.UID) error {
	//ds := mongo.NewDataStore()
	//defer ds.Close()
	//
	//return ds.C("videos").RemoveId(id)
	return nil
}

// SetPlayURL sets play Url (preview file) for the video
func (r defaultVideoRepository) SetPlayURL(videoID dgraph.UID, playURL string) error {
	//ds := mongo.NewDataStore()
	//defer ds.Close()
	//
	//model := bson.M{"playUrl": playURL}
	//return ds.C("videos").Update(bson.M{"_id": videoID}, bson.M{"$set": model})
	return nil
}

// SetFileURL sets play Url (purchse file) for the video
func (r defaultVideoRepository) SetDownloadURL(videoID dgraph.UID, fileURL string) error {
	//ds := mongo.NewDataStore()
	//defer ds.Close()
	//
	//model := bson.M{"fileUrl": fileURL}
	//return ds.C("videos").Update(bson.M{"_id": videoID}, bson.M{"$set": model})
	return nil
}

// SetArtworkURL sets artwork Url for the video
func (r defaultVideoRepository) SetArtworkURL(videoID dgraph.UID, fileURL string) error {

	//handler := func(cl *client.Dgraph) error {
	//	// set artworkURL
	//	req := client.Req{}
	//	m := fmt.Sprintf(setArtworkURL, videoID, fileURL)
	//	req.SetQuery(m)
	//
	//	_, err := cl.Run(context.Background(), &req)
	//	if err != nil {
	//		log.Println("Error in getting response from server when saving video : ", err)
	//		return err
	//	}
	//	return nil
	//}
	//return dgraph.Execute(handler)
	return nil
}

func (r defaultVideoRepository) getVideosFromRepo(query string) (int, []api.J, error) {

	resp, err := dgraph.Client.NewTxn().Query(context.Background(), query)
	if err != nil {
		log.Println("Error in getting response from server when reading genres : ", err)
		return 0, nil, err
	}

	jsonResp := resp.GetJson()

	var qr VideosQueryResponse
	err = json.Unmarshal(jsonResp, &qr)
	if err != nil {
		log.Println("Error unmarshalling genres : ", err)
		return 0, nil, err
	}

	return qr.Total[0].Count, qr.Videos, nil
}
