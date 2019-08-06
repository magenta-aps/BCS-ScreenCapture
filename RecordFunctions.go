package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

var pid int

func captureScreen (video_software_path string, video_software_params string, video_format string, video_path string, debug bool)  {
	var cliParams string
	if len(video_software_params) > 0 {
		cliParams += fmt.Sprintf(video_software_params, video_path, video_format)
	} else {
		switch runtime.GOOS {
		case "linux":
			cliParams += fmt.Sprintf("-y -loglevel error -video_size 1920x1080 -f x11grab -i :0.0 -f pulse -i 0 -f pulse -i default -filter_complex amerge -ac 2 -preset veryfast %stemp_bcs_recording.%s", video_path, video_format)
			break
		case "windows":
			cliParams += fmt.Sprintf("-y -f dshow -i video=screen-capture-recorder %stemp_bcs_recording.%s", video_path, video_format)
			break
		default:
			log.Fatal("Unsupported OS: " + runtime.GOOS)
		}
	}
	cmd := exec.Command(video_software_path, strings.Split(cliParams, " ")...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Start()

	if err != nil {
		fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		log.Fatal("Exiting")
	}
	pid = cmd.Process.Pid
	if debug {
		log.Print("Recording screen, process id: " + strconv.Itoa(pid))
	}

	if err := cmd.Wait(); err != nil {
		if len(stderr.String()) > 0 {
			if runtime.GOOS != "windows" {
				log.Fatalf("Video recording process has exited: %v", stderr.String())
			} else if debug {
				log.Print(stderr.String())
			}
		}
	}

}

func stopCapturing () {
	process, _ := os.FindProcess(pid)
	var err error
	switch runtime.GOOS {
	case "windows":
		err = process.Kill()
		break
	default:
		err = process.Signal(os.Interrupt)
	}
	if err != nil {
		log.Print(err)
	}
}
