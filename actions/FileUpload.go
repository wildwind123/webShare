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
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	files := r.MultipartForm.File["myFile[]"]
	for _, f := range files {
		file, err := f.Open()
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		file.Close()

		tempFile, err := ioutil.TempFile(dirPath, "wb*_"+f.Filename)
		if err != nil {
			http.Error(w, err.Error(), 400)
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
