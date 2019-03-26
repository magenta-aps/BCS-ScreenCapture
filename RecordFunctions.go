package main

import (
	"log"
	"os/exec"
)

func captureScreen ()  {
	var vlc_location = "%programfiles%/VideoLAN/VLC/vlc.exe"
	var sout = ":sout=#transcode{vcodec=h264,vb072}:standard{access=file,mux=mp4,dst=C:\\Desktop\\temp_bcs_recording.mp4}"
	cmd := exec.Command(vlc_location, "-I", "dummy","--one-instance", "screen://", ":screen-fps=15", sout)

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

}

func stopCapturing () {
	kill_cmd := exec.Command("vlc", "--one-instance", "vlc://quit")

	if err := kill_cmd.Start(); err != nil {
		log.Fatal(err)
	}
}
