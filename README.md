# BCS Screen Capture
A small https server that handles incoming requests from the
BCS shelter and records video and audio from the screen.  

### Installation instructions
##### Standard installation (recommended)
Download and unzip the win64 or win386 zip file, and make
sure you add BComeSafe certificates to `certs` folder.

##### Minimal installation (use when updating)
Download and unzip the win64_minimal or win386_minimal zip
file and replace `BCSRecorder.exe` where it was originally
extracted. You can also replace `conf.json.txt` to restore
configuration defaults.  
  
Note that the minimal installation requires `ffmpeg` to be
present in `lib/ffmpeg` (relative to the installation
folder) or in the Windows `$PATH`. Please adjust the value
of `RECORDING_SOFTWARE_PATH` in `conf.json.txt` accordingly.

### Developer instructions
##### Packages needed
* `rs/cors` - install with `go get -u github.com/rs/cors`
* `gorilla/mux` - install with `go get -u github.com/gorilla/mux`

Rename `conf_sample.json.txt` to `conf.json.txt` and make
sure the values are correct.

To build & run the project on UNIX, run `go build -o bin/BCSRecorder . && ./BCSRecorder`  
To package a Windows binary for production on UNIX, run the
`build_and_zip_binaries.sh` script.