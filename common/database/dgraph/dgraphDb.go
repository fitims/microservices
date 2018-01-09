package dgraph

import (
	"log"

	"github.com/dgraph-io/dgraph/client"
	"google.golang.org/grpc"
	"github.com/dgraph-io/dgraph/protos/api"
)

var Client *client.Dgraph

var connection *grpc.ClientConn
var dgClient api.DgraphClient

// Open opens a separate connection for each url provided
func Open(url string) {
	conn, err := grpc.Dial(url, grpc.WithInsecure())
	if err != nil {
		log.Println("Error connecting to DGraph : ", err)
		panic("Error connecting to DGraph")
	}

	connection = conn
	dgClient = api.NewDgraphClient(conn)
	Client = client.NewDgraphClient(dgClient)
}

// Close will close all the connections
func Close() {
	connection.Close()
}
