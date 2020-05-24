package main

import (
	"./html"
	"bytes"
	"fmt"
	//"github.com/gobuffalo/packr/v2"
	"html/template"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

const dateFormat = "January 2, 2006, 15:04:05"

var i int
var rootPath string = ""
var path string
var port string = "8000"
var htmlTemplate string = html.HtmlTemplate
var allPath string
var help bool
var haveError bool

type Folder struct {
	FolderName string
	LinkFolder string
}

type HtmlValues struct {
	Header  string
	Email   string
	Files   []File
	Folders []Folder
}

func init() {
	// set args parameters
	lastItemArgs := len(os.Args) - 1
	for i, args := range os.Args {
		if args == "-h" || args == "--help" {
			help = true
		}

		if i == lastItemArgs {
			continue
		}
		if args == "--path" {
			// set path
			rootPath = os.Args[i+1] + "/"
		} else if args == "--port" {
			// set port
			port = os.Args[i+1]
		} else {
			dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
			if err != nil {
				haveError = true
				fmt.Println(err.Error())
				return
			}
			rootPath = dir + "/"
		}
	}
	var reSlash = regexp.MustCompile(`/[/]+`)
	rootPath = reSlash.ReplaceAllString(rootPath, `/`)
	var text string
	haveError, text = shouldStopServer()
	if haveError {
		fmt.Println(text)
		printHelp()
		return
	}
	printIpInterfaces()

	//get binaryFiles
	//box := packr.New("myBox","./assets")
	//html, err := box.FindString("index2.html")
	//htmlTemplate = html
	//if err != nil {
	//	fmt.Println(err)
	//}
}

func main() {
	if haveError {
		return
	}
	http.HandleFunc("/", handler)
	http.HandleFunc("/file", HandleClient)

	//protect from favicon request
	http.HandleFunc("/favicon.ico", doNothing)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port+"", nil))

}

func shouldStopServer() (bool, string) {
	if help {
		return true, ""
	}
	if rootPath != "./" && strings.HasPrefix(rootPath, ".") {
		return true, "Please select full path. Current selected path =" + rootPath
	} else if rootPath != "./" {
		if _, err := os.Stat(rootPath); os.IsNotExist(err) && rootPath != "" {
			return true, "Path not found. path=" + rootPath
		}
	}
	return false, ""
}

func doNothing(w http.ResponseWriter, r *http.Request) {}

func printIpInterfaces() {
	ifaces, _ := net.Interfaces()

	// handle err
	for _, i := range ifaces {
		addrs, _ := i.Addrs()
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			// process IP address
			fmt.Println(ip.String() + ":" + port)
		}
	}
	fmt.Println("Выберите физический адрес вашего устройства и введите адресс в браузер, устройство должно " +
		"находиться в одной сети с сервером")
}

type File struct {
	Name     string
	Size     int64
	Modified string
	IsDir    bool
	Link     string
	Type     string
}

func handler(w http.ResponseWriter, r *http.Request) {
	path = r.URL.Query().Get("path")
	allPath = getPath()
	mapFiles := make(map[int]File)

	files, err := ioutil.ReadDir(allPath)
	countFiles := len(files)
	lastIndexFiles := 0
	if countFiles > 0 {
		lastIndexFiles = len(files) - 1
	}

	if err != nil {
		http.Error(w, err.Error(), 404)
		fmt.Println(err.Error())
		return
	}

	dirIndex := 0
	fileIndex := 0
	for _, f := range files {
		if f.IsDir() {
			mapFiles[dirIndex] = struct {
				Name     string
				Size     int64
				Modified string
				IsDir    bool
				Link     string
				Type     string
			}{Name: f.Name(), Size: f.Size(), Modified: f.ModTime().Format(dateFormat), IsDir: true,
				Link: "/?path=" + allPath + f.Name(), Type: "directory"}
			dirIndex++

		} else {
			mapFiles[lastIndexFiles-fileIndex] = struct {
				Name     string
				Size     int64
				Modified string
				IsDir    bool
				Link     string
				Type     string
			}{Name: f.Name(), Size: f.Size(), Modified: f.ModTime().Format(dateFormat), IsDir: false,
				Link: "/file?file=" + allPath + f.Name() + "&fileName=" + f.Name(), Type: filepath.Ext(f.Name())}
			fileIndex++
		}
	}

	w.Write([]byte(getRenderedHtml(mapFiles)))
}

func getPath() string {
	fullPath := ""
	if rootPath == "" {
		rootPath = "./"
		fullPath = rootPath + "/" + path + "/"
	} else {
		path = strings.Replace(path, rootPath, "", -1)
		strings.HasPrefix(rootPath, ".")
	}
	fullPath = rootPath + "/" + path + "/"

	//clear, no extra characters

	var reSlash = regexp.MustCompile(`/[/]+`)
	clearPath1 := reSlash.ReplaceAllString(fullPath, `/`)
	clearPath1 = reSlash.ReplaceAllString(fullPath, `/`)
	clearPath2 := strings.Replace(clearPath1, "././", "./", -1)
	//validate from over folder
	var re = regexp.MustCompile(`/.[.]+/`)
	clearPath2 = re.ReplaceAllString(clearPath2, `/`)

	return clearPath2
}

func getFolders() []Folder {
	var folders []Folder
	var folderNames []string

	childFolder := strings.Replace(allPath, rootPath, "", -1)
	folderNames = strings.SplitAfter(childFolder, "/")
	folderPath := ""
	//append first root path
	folders = append(folders, Folder{"rootFolder/", rootPath})

	for _, folderName := range folderNames {
		if folderName != "" || folderName != "/" {
			folderPath = rootPath + folderPath + folderName

			folders = append(folders, Folder{folderName, folderPath})
		}
	}
	return folders
}

func getRenderedHtml(f map[int]File) string {

	t := template.New("fieldname example")
	t, _ = t.Parse(htmlTemplate)
	var files []File

	for i := 0; i < len(f); i++ {
		files = append(files, File{f[i].Name, f[i].Size, f[i].Modified, true,
			f[i].Link, f[i].Type})

	}

	p := HtmlValues{Header: "WebShare.", Email: "test@mail.ru", Files: files, Folders: getFolders()}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, p); err != nil {
		fmt.Println("Error")
	}
	return tpl.String()
}

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
		http.Error(writer, "File not found.", 404)
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

func printHelp() {
	fmt.Println("--port - select port, default = 8000")
	fmt.Println("--path - select full path, default = programm runned folder")
	fmt.Println("-h --help - Help")
}
