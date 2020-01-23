package main

/*
REST webserver that starts and stops recording the screen.
Handles start and stop requests from Shelter webapp
Records screen using VLC Player
*/

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

type Confirmation struct {
	Name  string
	Reset string
}

type Configuration struct {
	RECORDING_SOFTWARE_PATH   string
	RECORDING_SOFTWARE_PARAMS string
	CERTIFICATE_PATH          string
	ROOTHOST                  string
	PORT                      string
	VIDEO_SAVE_PATH           string
	VIDEO_FORMAT              string
	TIMEOUT_IN_MINUTES        time.Duration
	SAVE_NAME_IN_VIDEO        bool
	DEBUG                     bool
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
var configuration = Configuration{}

// Insuring that we dont start recording again after stopping for timeout
var timeout_triggered = false

func main() {
	// Load configuration
	file, err := os.Open("conf.json.txt")
	if err != nil {
		log.Fatal("Error: ", err)
	}
	defer file.Close()
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&configuration)
	if err != nil {
		log.Fatal("Failed to decode json configuration file: ", err)
	}

	if _, err := os.Stat(configuration.VIDEO_SAVE_PATH); os.IsNotExist(err) {
		err := os.Mkdir(configuration.VIDEO_SAVE_PATH, 0777)
		if err != nil {
			log.Fatal("Failed to create directory: ", configuration.VIDEO_SAVE_PATH)
		}
	}

	router := mux.NewRouter()
	router.HandleFunc("/start", startRecording)
	router.HandleFunc("/stop", stopRecording)
	router.HandleFunc("/status", getStatus)

	// Cors handling is needed due to the request coming from a different origin

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowCredentials: true,
	})

	handler := c.Handler(router)

	hostname := configuration.ROOTHOST + ":" + configuration.PORT
	if err := http.ListenAndServeTLS(hostname, configuration.CERTIFICATE_PATH+"bcomesafe.crt", configuration.CERTIFICATE_PATH+"bcomesafe.key", handler); err != nil {
		log.Fatal(err)
	}
}

func startRecording(w http.ResponseWriter, r *http.Request) {
	if recording == RECORDING || recording == WAITING || timeout_triggered == true {
		return
	} else {
		recording = RECORDING
		go captureScreen(configuration.RECORDING_SOFTWARE_PATH, configuration.RECORDING_SOFTWARE_PARAMS, configuration.VIDEO_FORMAT, configuration.VIDEO_SAVE_PATH, configuration.DEBUG)
		fmt.Println("Capturing Screen")
	}

	fmt.Println("Started Timer")

	// The timeout timer determines when the recording should automatically stop

	timeout_timer = time.AfterFunc(configuration.TIMEOUT_IN_MINUTES*time.Minute, stopTimer)

	w.Write([]byte("Capturing Screen"))

}

func stopRecording(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Timeout: " + strconv.FormatBool(timeout_triggered))
	fmt.Println("Recording: " + strconv.Itoa(recording))

	if recording == INACTIVE && !timeout_triggered {
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
		log.Print(err)
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
	var filename string
	if len(request.Name) > 1 {
		fmt.Println("Saving file")
		var time_now = time.Now().Format("2006-01-02_15.04.05")
		var additional_filename_text = ""
		if configuration.SAVE_NAME_IN_VIDEO {
			additional_filename_text += "_" + request.Name
		}
		filename = fmt.Sprintf(configuration.VIDEO_SAVE_PATH+"%s%s.%s", time_now, additional_filename_text, configuration.VIDEO_FORMAT)

		fmt.Println(filename)
		time.Sleep(1 * time.Second)

		err := os.Rename(configuration.VIDEO_SAVE_PATH+"temp_bcs_recording."+configuration.VIDEO_FORMAT, filename)
		if err != nil {
			fmt.Println("WARNING: File not found or is in use by another process (Timeout function might have saved file already)")
		}

	} else {

		http.Error(w, "ERROR: Length of name not longer than 1", 400)
		time.Sleep(1 * time.Second)
		fmt.Println("Deleting file")
		os.Remove(configuration.VIDEO_SAVE_PATH + "temp_bcs_recording.mp4")
	}

	jsonData, _ := json.Marshal(map[string]string{"filename": filename})
	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)
}

func stopTimer() {

	fmt.Println("Stopped timer and stopping recording")

	stopCapturing()
	recording = WAITING
}

// Returns the current recording status
// INACTIVE - If not recording
// RECORDING - If recording
// WAITING - If stopped by timeout

func getStatus(w http.ResponseWriter, r *http.Request) {

	fmt.Println("Timeout: " + strconv.FormatBool(timeout_triggered))
	fmt.Println("Recording: " + strconv.Itoa(recording))

	var r_status = RecordingStatus{Status: recording}

	jsonData, err := json.Marshal(r_status)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jsonData)

	if recording == WAITING {
		timeout_triggered = true
		recording = INACTIVE
	}
}
