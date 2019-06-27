# BCS Screen Capture
A small https server that handles incoming requests from the
BCS shelter and records video and audio from the screen.  

### Installation instructions
##### Prerequisites
* `ffmpeg` - the software package used to capture the screen
* bcomesafe.com wildcard certificates - needed to authenticate
server responses
* The executable file (available in the Releases tab)


### Developer instructions
##### Packages needed
* `rs/cors` - install with `go get -u github.com/rs/cors`
* `gorilla/mux` - install with `go get -u github.com/gorilla/mux`

Rename `conf_sample.json.txt` to `conf.json.txt` and make
sure the values are correct.

To build & run the project on UNIX, run `go build -o bin/BCSRecorder . && ./BCSRecorder`  
To package a Windows binary for production on UNIX, run:  
`GOOS=windows GOARCH=amd64 go build -o bin/BCSRecorder.exe .` for 64-bit  
and  
`GOOS=windows GOARCH=386 go build -o bin/BCSRecorder_32.exe .` for 32-bit  
