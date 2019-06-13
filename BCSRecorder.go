package main

/*
REST webserver that starts and stops recording the screen.
Handles start and stop requests from Shelter webapp
Records screen using VLC Player
*/

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/rs/cors"
	"log"
	"net/http"
	"os"
	"time"
	"strconv"
)

type Confirmation struct {
	Name string
	Reset string
}

var timeout_timer *time.Timer

// For handling timeouts after 60 minutes
const INACTIVE = 0
const RECORDING = 1
const WAITING = 2

type RecordingStatus struct {
	Status int
}

var recording = INACTIVE

// Insuring that we dont start recording again after stopping for timeout
var timeout_triggered = false

func main () {
	if _, err := os.Stat("C:\\BCSVideos"); os.IsNotExist(err) {
		os.Mkdir("C:\\BCSVideos", 0777)
	}

	router := mux.NewRouter()
	router.HandleFunc("/start", startRecording)
	router.HandleFunc("/stop", stopRecording)
	router.HandleFunc("/status", getStatus)

	// Cors handling is needed due to the request coming from a different origin

	c := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	if err := http.ListenAndServeTLS("loc.bcomesafe.com:3032","bcomesafe.crt", "bcomesafe.key", handler); err != nil {
		log.Fatal(err)
	}
}

func startRecording (w http.ResponseWriter, r *http.Request) {
	if recording == RECORDING || recording == WAITING || timeout_triggered == true {
		return
	} else {
		recording = RECORDING
		go captureScreen()
		fmt.Println("Capturing Screen")
	}

	fmt.Println("Started Timer")

	// The timeout timer determines when the recording should automatically stop

	timeout_timer = time.AfterFunc(60 * time.Minute, stopTimer)

	w.Write([]byte("Capturing Screen"))

}

func stopRecording (w http.ResponseWriter, r *http.Request) {
	fmt.Println("STOPPING RECORDING")
	fmt.Println("Timeout: " + strconv.FormatBool(timeout_triggered))
	fmt.Println("Recording: " + strconv.Itoa(recording))

	if (recording == INACTIVE && !timeout_triggered) {
		return
	}

	stopCapturing()

	fmt.Println("Stopping Timer")
	timeout_timer.Stop()

	recording = INACTIVE

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

	// A Reset parameter is sent in the request to determine if the stop call is from shelter reset or timeout
	// If the recording timed out the Reset parameter will be false

	if request.Reset == "true" {
		timeout_triggered = false
	}


	// If a name is supplied in the webapp the video is saved in C:\\BCSVideoes and named its timestamp
	// If not the video is deleted
	if len(request.Name) > 1 {
		fmt.Println("Saving file")
		var time_now = time.Now().Format("2006_01_02-15_04_05")
		var filename = fmt.Sprintf("C:\\BCSVideos\\%s.mp4", time_now)

		fmt.Println(filename)
		time.Sleep(1 * time.Second)

		err := os.Rename(`C:\BCSVideos\temp_bcs_recording.mp4`, filename)
		if err != nil {
			fmt.Println("WARNING: File not found or is in use by another process (Timeout function might have saved file already)")
		}

	} else {

		http.Error(w, "ERROR: Length of name not longer than 1", 400)
		time.Sleep(1 * time.Second)
		fmt.Println("Deleting file")
		os.Remove(`C:\BCSVideos\temp_bcs_recording.mp4`)
	}

	w.Write([]byte("Stopped capturing screen"))
}

func stopTimer () {

	fmt.Println("Stopped timer and stopping recording")

	stopCapturing()
	recording = WAITING
}

// Returns the current recording status 
// INACTIVE - If not recording
// RECORDING - If recording
// WAITING - If stopped by timeout

func getStatus (w http.ResponseWriter, r *http.Request) {

	fmt.Println("Timeout: " + strconv.FormatBool(timeout_triggered))
	fmt.Println("Recording: " + strconv.Itoa(recording))

	var r_status = RecordingStatus{Status: recording}

	jsonData, err := json.Marshal(r_status)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Posting status v2")

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	if recording == WAITING {
		timeout_triggered = true
		recording = INACTIVE
	}
}
