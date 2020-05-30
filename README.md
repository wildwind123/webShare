# webShare
Share and upload your files from web server

# Compiled version
compiled version can download from 
https://github.com/wildwind123/webShare/tree/master/compiled
# Run
# Linux
download compiled version
chmod -R 777 web_share 
./webShare --path /home --port 8006
open http://your_ip:8006

or just run ./webShare
default path will be file folder, default port will be 8000
# Windows
web_share.exe 
default path will be file folder, default port will be 8000
# Help 
--port - select port, default = 8000 \
--path - select full path, default = programm runned folder \
-h --help - Help

# Multiple Build

env GOOS=linux GOARCH=386 go build -o compiled/linux86 ; \
env GOOS=linux GOARCH=amd64 go build -o compiled/linux64 ; \
env GOOS=windows GOARCH=386 go build -o compiled/windows86 ; \
env GOOS=windows GOARCH=amd64 go build -o compiled/windows64 ; \
env GOOS=darwin GOARCH=386 go build -o compiled/mac86 ; \
env GOOS=darwin GOARCH=amd64 go build -o compiled/mac64; 
