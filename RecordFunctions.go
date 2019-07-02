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

func captureScreen (video_software_path string, video_software_params string, video_path string)  {
	var cliParams string
	if len(video_software_params) > 0 {
		cliParams += fmt.Sprintf(video_software_params, video_path)
	} else {
		switch runtime.GOOS {
		case "linux":
			cliParams += fmt.Sprintf("-y -loglevel error -video_size 1920x1080 -f x11grab -i :0.0 -f pulse -i 0 -f pulse -i default -filter_complex amerge -ac 2 -preset veryfast %stemp_bcs_recording.mp4", video_path)
			break
		case "windows":
			cliParams += fmt.Sprintf("-f dshow -i video=screen-capture-recorder %stemp_bcs_recording.mp4", video_path)
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
	log.Print("Pid is " + strconv.Itoa(pid))

	if err := cmd.Wait(); err != nil {
		if len(stderr.String()) > 0 {
			log.Fatalf("Video recording process has exited: %v", stderr.String())
		}
	}

}

func stopCapturing () {
	log.Print("Trying to kill " + strconv.Itoa(pid))
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
