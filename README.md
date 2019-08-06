# BCS Screen Capture
A small https server that handles incoming requests from the
BCS shelter and records video and audio from the screen.  

### Installation instructions
##### Standard installation (recommended)
Download and unzip the latest release, and make
sure you add BComeSafe certificates to the `certs` folder.

Run `Setup.Screen.Capturer.Recorder.v0.12.10.exe` from the
`lib` folder and follow the installation instructions.
 
Now run `BCSRecorder.exe` and verify that the installation
is working by opening `test_status.html` in a browser. If
everything is working, this should display a message
indicating that the software is ready to record. To test
recording, please open `test_recording.html` which will
trigger a 10 second video & audio recording. 

##### Minimal installation (use when updating)
Download and unzip the latest win64_minimal zip
file and replace `BCSRecorder.exe` where it was originally
extracted. You can also replace `conf.json.txt` to restore
configuration defaults.  
  

### Usage
Please add `BCSRecorder.exe` to your startup configuration
or ensure that it is running whenever the computer is
restarted.

### Developer instructions
##### Packages needed
* `rs/cors`
* `gorilla/mux`

Rename `conf_sample.json.txt` to `conf.json.txt` and make
sure the values are correct.

To build & run the project on UNIX, run
`go build -o bin/BCSRecorder . && bin/BCSRecorder`  
To package a Windows binary for production on UNIX, run the
`build_and_zip_binaries.sh` script.