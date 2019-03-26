package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"time"
)

type Confirmation struct {
	Name string
}

func main () {
	router := mux.NewRouter()
	router.HandleFunc("/start", StartRecording)
	router.HandleFunc("/stop", StopRecording)

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	if err := http.ListenAndServe(":3032", handler); err != nil {
		log.Fatal("Couldnt open http server")
	}
	fmt.Println("Listening on port 3032")
}

func StartRecording (w http.ResponseWriter, r *http.Request) {
	captureScreen()
	fmt.Println("Capturing Screen")

}

func StopRecording (w http.ResponseWriter, r *http.Request) {
	stopCapturing()
	fmt.Println("Stopped capturing the screen")

	var request Confirmation

	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
	}

	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}

	if len(request.Name) > 1 {
		fmt.Println("Saving file")
		var time_now = time.Now().Format("00-00-0000 00:00:00")
		var filename = fmt.Sprintf("C:\\Desktop\\%s", time_now)

		os.Rename("C:\\Desktop\\temp_bcs_recording.mp4", filename)

	} else {
		fmt.Println("Deleting file")
		os.Remove("C:\\Desktop\\temp_bcs_recording.mp4")
	}
}



