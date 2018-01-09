package routes

import (
	"thaThrowdown/services/videoManager/controllers"
)

// Router is the interface that defines the methods for the router
type Router interface {
	SetVideoRoutes(c controllers.VideoController)
	SetGenreRoutes(c controllers.GenreController)
}
