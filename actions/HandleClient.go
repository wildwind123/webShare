package actions

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
)

func HandleClient(writer http.ResponseWriter, request *http.Request) {
	//First of check if Get is set in the URL
	FilePath := request.URL.Query().Get("file")
	Filename := request.URL.Query().Get("fileName")

	if FilePath == "" {
		//Get not set, send a 400 bad request
		http.Error(writer, "Get 'file' not specified in url.", 400)
		return
	} else if Filename == "" {
		http.Error(writer, "File name is empty", 400)
		return
	}

	fmt.Println("Client requests: " + FilePath)

	//Check if file exists and open
	Openfile, err := os.Open(FilePath)
	defer Openfile.Close() //Close after function return
	if err != nil {
		//File not found, send 404
		fmt.Println(err.Error())
		writer.Header().Add("Content-Type", "text/html")
		writer.Write([]byte(err.Error() + " File not found" + "<br><a href=\"./\"> Back Home page </a>"))
		return
	}

	//File is found, create and send the correct headers

	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	FileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	Openfile.Read(FileHeader)
	//Get content type of file
	FileContentType := http.DetectContentType(FileHeader)

	//Get the file size
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	writer.Header().Set("Content-Disposition", "attachment; filename="+Filename)
	writer.Header().Set("Content-Type", FileContentType)
	writer.Header().Set("Content-Length", FileSize)

	//Send the file
	//We read 512 bytes from the file already, so we reset the offset back to 0
	Openfile.Seek(0, 0)
	io.Copy(writer, Openfile) //'Copy' the file to the client
	return
}
