//     Schemes: http
//     Host: localhost:7050
//     BasePath: /
//     Version: 1.0.0
//     License: MIT http://opensource.org/licenses/MIT
//     Contact: API Thathrowdown<api@thathrowdown.com> http://thathrowdown.com
//     Consumes:
//     - application/json
//     Produces:
//     - application/json
// swagger:meta
package main

import (
	"log"
	"thaThrowdown/common/database/dgraph"
	"thaThrowdown/common/infrastructure"
	"thaThrowdown/services/videoManager/controllers"
	"thaThrowdown/services/videoManager/repositories"
	"thaThrowdown/services/videoManager/routes"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

const (
	DEFAULT_DGRAPH_CONNECTION = "localhost:9080"
	LOGGER_FILENAME           = "/var/log/videoManager/videoManager.log"
)

func main() {
	//logFile := middle.RollingLog(LOGGER_FILENAME)
	//log.SetOutput(&logFile)

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Static("public", "/home")

	connStr := infrastructure.GetEnv("THATHROW_DGRAPH_CONNECTION", DEFAULT_DGRAPH_CONNECTION)

	// dgraph
	dgraph.Open(connStr)
	defer dgraph.Close()

	g := e.Group("/videos")
	routes := routes.NewRouter(g)

	videoController := controllers.NewVideoController(repositories.NewVideoRepository())
	genreController := controllers.NewGenreController(repositories.NewGenreRepository())

	routes.SetVideoRoutes(videoController)
	routes.SetGenreRoutes(genreController)

	log.Println("VideoManager - Listening on port 7050 ...")
	e.Start(":7050")
}
