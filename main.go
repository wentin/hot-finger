package main 

import (
	"os"
	"fmt"
	"log"
	"net/http"
	"time"
	"strconv"

    "github.com/aberke/hot-finger/touches"
)


const DEFAULT_PORT = "8080";

func main() {

	var port = os.Getenv("PORT")
	if len(port) == 0 {
		fmt.Println("$PORT not set -- defaulting to", DEFAULT_PORT)
		port = DEFAULT_PORT
	}
	fmt.Println("using port:", port)

	server := touches.NewServer()
	go server.Listen()
	
	http.HandleFunc("/static/", serveStatic)
	http.HandleFunc("/widget/", serveStatic)
	http.HandleFunc("/", serveHome)
	http.HandleFunc("/new-grid-id", getNewGridId)

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Error listening, %v", err)
    }
}

var lastNewGridId int64 = 0
func getNewGridId(w http.ResponseWriter, r *http.Request) {
	newGridId := time.Now().UnixNano()
	if (newGridId == lastNewGridId) {
		newGridId = newGridId + 1
	}
	lastNewGridId = newGridId
	log.Println("newGridId: ", newGridId)
	fmt.Fprintf(w, strconv.FormatInt(newGridId, 10))
}


func serveHome(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/")
}
func serveStatic(w http.ResponseWriter, r *http.Request) {
	var staticFileHandler = http.FileServer(http.Dir("./public/"))
	staticFileHandler.ServeHTTP(w, r)
}