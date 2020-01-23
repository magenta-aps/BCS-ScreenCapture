#!/bin/bash
set -euxo pipefail

mkdir -p bin/BCSRecorder_win64/lib
mkdir -p bin/BCSRecorder_win64/certs
GOOS=windows GOARCH=amd64 go build -o bin/BCSRecorder_win64/BCSRecorder.exe .
go build -o bin/BCSRecorder .
cat << EOF > bin/conf.json.txt
{
  "RECORDING_SOFTWARE_PATH": "lib/ffmpeg/bin/ffmpeg.exe",
  "RECORDING_SOFTWARE_PARAMS": "-y -f dshow -i audio=virtual-audio-capturer:video=screen-capture-recorder -filter_complex amix=inputs=1 -vcodec libx264 -pix_fmt yuv420p -preset ultrafast -vsync vfr -acodec libmp3lame %stemp_bcs_recording.%s",
  "CERTIFICATE_PATH": "./certs/",
  "ROOTHOST": "loc.bcomesafe.com",
  "PORT": "3032",
  "VIDEO_SAVE_PATH": "C:/BCSVideos/",
  "VIDEO_FORMAT": "mp4",
  "TIMEOUT_IN_MINUTES": 60,
  "SAVE_NAME_IN_VIDEO": true,
  "DEBUG": false
}
EOF
cp bin/conf.json.txt bin/BCSRecorder_win64
cp test_recording.html bin/BCSRecorder_win64
cp -r bin/ffmpeg bin/BCSRecorder_win64/lib
cp bin/Setup.Screen.Capturer.Recorder.v0.12.10.exe bin/BCSRecorder_win64/lib
pushd bin/BCSRecorder_win64
zip -r ../bcsrecorder_win64_full.zip .
zip ../bcsrecorder_win64_minimal.zip *.*
popd
