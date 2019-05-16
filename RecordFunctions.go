package main

import (
	"log"
	"os/exec"
)

func captureScreen ()  {
	var sout = ":sout=#transcode{vcodec=h264,vb072}:standard{access=file,mux=mp4,dst=C:\\BCSVideos\\temp_bcs_recording.mp4}"
	cmd := exec.Command(`C:\Program Files\VideoLAN\VLC\vlc.exe`, "-I", "dummy","--one-instance", "screen://", ":screen-fps=15", sout)

	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func stopCapturing () {
	kill_cmd := exec.Command(`C:\Program Files\VideoLAN\VLC\vlc.exe`, "--one-instance", "vlc://quit")

	if err := kill_cmd.Run(); err != nil {
		log.Fatal(err)
	}
}
