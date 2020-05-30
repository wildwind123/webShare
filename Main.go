package main

import (
	"./actions"
	"./html"
	"bytes"
	"fmt"
	//"github.com/gobuffalo/packr/v2"
	"html/template"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
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
var html_template bool
var browserSupportedFiles map[string]FileType

type Folder struct {
	FolderName string
	LinkFolder string
}

type HtmlValues struct {
	Header  string
	Email   string
	Files   []File
	Folders []Folder
	DirPath string
}

func init() {
	lastItemArgs := len(os.Args) - 1

	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		haveError = true
		fmt.Println(err.Error())
		return
	}
	rootPath = dir + "/"

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
		} else if args == "--template" {
			if os.Args[i+1] == "true" {
				html_template = true
			}
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
	setSupportedFiles()
}

type FileType struct {
	Extension []string
	Icon      string
}

func setSupportedFiles() {
	browserSupportedFiles = make(map[string]FileType)
	browserSupportedFiles["img"] = FileType{
		Extension: []string{".apng", ".bmp", ".gif", ".ico", ".cur", ".jpg", ".jpeg", ".jfif", ".pjpeg", ".pjp", ".png", ".svg", ".tif", ".tiff", ".webp"},
		Icon:      "img",
	}
	browserSupportedFiles["pdf"] = FileType{
		Extension: []string{".pdf"},
		Icon:      "pdf",
	}
	browserSupportedFiles["audio"] = FileType{
		Extension: []string{".aac", ".mp3", "wav", ".webm"},
		Icon:      "audio",
	}
	browserSupportedFiles["video"] = FileType{
		Extension: []string{".mp4", ".webm"},
		Icon:      "video",
	}
	browserSupportedFiles["txt"] = FileType{
		Extension: []string{".css", ".txt", ".php"},
		Icon:      "txt",
	}
}

func main() {
	if haveError {
		return
	}
	http.HandleFunc("/", handler)
	http.HandleFunc("/FileUpload", actions.FileUpload)
	http.HandleFunc("/file", actions.HandleClient)
	http.HandleFunc("/show", icoHandler)

	//protect from favicon request
	http.HandleFunc("/favicon.ico", actions.DoNothing)
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port+"", nil))

}
func icoHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Query().Get("file")
	http.ServeFile(w, r, filePath)
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
	IsDir    string
	Link     string
	LinkShow string
	Type     string
	Icon     string
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
				IsDir    string
				Link     string
				LinkShow string
				Type     string
				Icon     string
			}{Name: f.Name(), Size: 0, Modified: f.ModTime().Format(dateFormat), IsDir: "true",
				Link:     "/?path=" + allPath + f.Name(),
				LinkShow: "", Type: "directory", Icon: ""}
			dirIndex++

		} else {
			fname := f.Name()
			fileExt := filepath.Ext(fname)
			mapFiles[lastIndexFiles-fileIndex] = struct {
				Name     string
				Size     int64
				Modified string
				IsDir    string
				Link     string
				LinkShow string
				Type     string
				Icon     string
			}{Name: fname, Size: f.Size(), Modified: f.ModTime().Format(dateFormat), IsDir: "",
				Link:     "/file?file=" + allPath + fname + "&fileName=" + fname,
				LinkShow: "/show?file=" + allPath + fname, Type: fileExt, Icon: getIcon(fileExt)}
			fileIndex++
		}
	}

	w.Write([]byte(getRenderedHtml(mapFiles)))
}

func getIcon(fileExt string) string {
	fileExt = strings.ToLower(fileExt)
	for _, typeFiles := range browserSupportedFiles {
		for _, typeExtension := range typeFiles.Extension {
			if fileExt == typeExtension {
				return typeFiles.Icon
			}
		}
	}
	return ""
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
	folders = append(folders, Folder{"rootFolder", rootPath})

	for _, folderName := range folderNames {
		if folderName != "" && folderName != "/" {
			folderPath = rootPath + folderPath + folderName
			var reSlash = regexp.MustCompile(`[/\\]`)
			folderName = reSlash.ReplaceAllString(folderName, ``)

			folders = append(folders, Folder{folderName, folderPath})
		}
	}
	return folders
}

func getRenderedHtml(f map[int]File) string {
	if html_template {
		file, err := ioutil.ReadFile("assets/index.html")
		if err != nil {
			return err.Error()
		}

		htmlTemplate = string(file)

		//re := regexp.MustCompile(`\r?\n`)
		//input := re.ReplaceAllString(htmlTemplate, " ")
		//re = regexp.MustCompile(`"`)
		//input = re.ReplaceAllString(input, "'")
		//fmt.Println(input)
	}
	t := template.New("fieldname example")
	t, _ = t.Parse(htmlTemplate)
	var files []File

	for i := 0; i < len(f); i++ {
		files = append(files, File{f[i].Name, f[i].Size, f[i].Modified, f[i].IsDir,
			f[i].Link, f[i].LinkShow, f[i].Type, f[i].Icon})

	}

	p := HtmlValues{Header: "WebShare.", Email: "test@mail.ru", Files: files, Folders: getFolders(), DirPath: allPath}
	var tpl bytes.Buffer
	if err := t.Execute(&tpl, p); err != nil {
		fmt.Println("Error")
	}
	return tpl.String()
}

func printHelp() {
	fmt.Println("--port [port] - select port, default = 8000")
	fmt.Println("--path [fullPath]- select full path, default = programm runned folder")
	fmt.Println("--template [true] - if you want to use yourself template from assets/index.html")
	fmt.Println("-h --help - Help")
	fmt.Println("https://github.com/wildwind123/webShare")
}
