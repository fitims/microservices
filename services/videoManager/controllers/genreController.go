package controllers

import (
	"log"
	"net/http"
	"thaThrowdown/common/api"
	"thaThrowdown/common/database/dgraph"

	"thaThrowdown/services/videoManager/models"
	"thaThrowdown/services/videoManager/repositories"

	"fmt"

	"github.com/labstack/echo"
)

// DefaultGenreController is the default implementation of the GenreController interface
type DefaultGenreController struct {
	repo repositories.GenreRepository
}

// NewGenreController is a factory that creates the default instance of the genre controllers
func NewGenreController(r repositories.GenreRepository) GenreController {
	return DefaultGenreController{repo: r}
}

// GetGenres returns genres
func (ctl DefaultGenreController) GetGenres(c echo.Context) error {
	genres, err := ctl.repo.GetGenres()
	if err != nil {
		log.Println("Error getting videos from database : ", err)
		return c.JSON(http.StatusInternalServerError, api.Fail("Internal Error"))
	}

	// return Music genres view model

	fmt.Println("genres returned from database : ", genres)

	return c.JSON(http.StatusOK, genres)
}

// UploadGenre is a method of APIUser that is used to upload a genre into db
func (ctl DefaultGenreController) UploadGenre(c echo.Context) error {
	// get the genre
	var r models.GenreRequest
	err := c.Bind(&r)
	if err != nil {
		log.Println("Error binding the genre : ", err)
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid data"))
	}

	id, err := ctl.repo.AddGenre(r)
	if err != nil {
		if err == repositories.GenreAlreadyExistsErr {
			log.Println("Error saving the genre :", err)
			return c.JSON(http.StatusBadRequest, api.Fail("Genre already exists"))
		}

		log.Println("Error saving the genre :", err)
		return c.JSON(http.StatusInternalServerError, api.Fail("Could not save the genre"))
	}

	return c.JSON(http.StatusOK, api.SuccessWithID("Genre uploaded successfully", id))
}

// DeleteGenre is a method of APIUser that is used to delete video
func (ctl DefaultGenreController) DeleteGenre(c echo.Context) error {
	genreID := c.Param("genreId")
	id, err := dgraph.ToUID(genreID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, api.Fail("Invalid request"))
	}

	err = ctl.repo.DeleteGenre(id)
	if err != nil {
		log.Println("Error storing purchase data : ", err)
		return c.JSON(http.StatusInternalServerError, api.Fail("Could not delete genre"))
	}
	return c.JSON(http.StatusOK, api.Success("Genre deleted successfully"))
}
