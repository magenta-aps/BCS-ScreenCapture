package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
)

var pid int
var procStdin io.WriteCloser

func captureScreen(videoSoftwarePath string, videoSoftwareParams string, videoFormat string, videoPath string, debug bool) {
	var cliParams string
	if len(videoSoftwareParams) > 0 {
		cliParams += fmt.Sprintf(videoSoftwareParams, videoPath, videoFormat)
	} else {
		switch runtime.GOOS {
		case "linux":
			cliParams += fmt.Sprintf("-y -loglevel error -video_size 1920x1080 -f x11grab -i :0.0 -f pulse -i 0 -f pulse -i default -filter_complex amerge -ac 2 -preset veryfast %stemp_bcs_recording.%s", videoPath, videoFormat)
			break
		case "windows":
			cliParams += fmt.Sprintf("-y -f dshow -i video=screen-capture-recorder %stemp_bcs_recording.%s", videoPath, videoFormat)
			break
		default:
			log.Fatal("Unsupported OS: " + runtime.GOOS)
		}
	}
	cmd := exec.Command(videoSoftwarePath, strings.Split(cliParams, " ")...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	procStdin, _ = cmd.StdinPipe()
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

func stopCapturing() {
	io.WriteString(procStdin, "q")
}
