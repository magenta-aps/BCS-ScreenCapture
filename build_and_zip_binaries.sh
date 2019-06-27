#!/bin/bash
set -euxo pipefail

mkdir -p bin/BCSRecorder_win64/lib
mkdir -p bin/BCSRecorder_win64/certs
mkdir -p bin/BCSRecorder_win386/lib
mkdir -p bin/BCSRecorder_win386/certs
GOOS=windows GOARCH=amd64 go build -o bin/BCSRecorder_win64/BCSRecorder.exe .
GOOS=windows GOARCH=386 go build -o bin/BCSRecorder_win386/BCSRecorder.exe .
go build -o bin/BCSRecorder .
cat << EOF >> bin/conf.json.txt
{
	"RECORDING_SOFTWARE_PATH": "lib/ffmpeg/bin/ffmpeg",
    "CERTIFICATE_PATH": "./certs/",
	"ROOTHOST": "loc.bcomesafe.com",
	"PORT": "3032",
	"VIDEO_SAVE_PATH": "C:/BCSVideos/",
	"TIMEOUT_IN_MINUTES": 60,
	"SAVE_NAME_IN_VIDEO": true
}
EOF
cp bin/conf.json.txt bin/BCSRecorder_win64
mv bin/conf.json.txt bin/BCSRecorder_win386
cp -r bin/ffmpeg bin/BCSRecorder_win64/lib
cp -r bin/ffmpeg bin/BCSRecorder_win386/lib
pushd bin/BCSRecorder_win64
zip -r ../bcsrecorder_win64_full.zip .
zip ../bcsrecorder_win64_minimal.zip *.*
popd
pushd bin/BCSRecorder_win386
zip -r ../bcsrecorder_win386_full.zip .
zip ../bcsrecorder_win386_minimal.zip *.*