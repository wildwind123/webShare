package actions

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

func FileUpload(w http.ResponseWriter, r *http.Request) {
	dirPath := r.URL.Query().Get("filePath")
	fmt.Println(dirPath)
	if err := r.ParseMultipartForm(5000 << 20); err != nil {
		w.Header().Add("Content-Type", "text/html")
		w.Write([]byte(err.Error() + "<br><a href=\"./\"> Back Home page </a>"))
		return
	}
	files := r.MultipartForm.File["myFile[]"]
	for _, f := range files {
		file, err := f.Open()
		if err != nil {
			w.Header().Add("Content-Type", "text/html")
			w.Write([]byte(err.Error() + "<br><a href=\"./\"> Back Home page </a>"))
			return
		}
		file.Close()

		tempFile, err := ioutil.TempFile(dirPath, "wb*_"+f.Filename)
		if err != nil {
			w.Header().Add("Content-Type", "text/html")
			w.Write([]byte(err.Error() + "<br><a href=\"./\"> Back Home page </a>"))
			return
		}
		defer tempFile.Close()
		fileBytes, err := ioutil.ReadAll(file)
		tempFile.Write(fileBytes)
	}

	if r.Method == "POST" {
		http.Redirect(w, r, r.Header.Get("Referer"), 302)
	}
}
