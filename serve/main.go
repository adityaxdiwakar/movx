package main

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"net/http"

	"github.com/go-chi/chi/v5"
)

/* Response is a generic response object for all requests */
type Response[T any] struct {
	Payload    T   `json:"payload"`
	StatusCode int `json:"status_code"`
}

/* DirectoryResp is a response object for HTTP directory requests */
type DirectoryResp struct {
	Folders []string `json:"folders"`
	Files   []string `json:"files"`
}

/* DirectoryReq is the expected request format for `/directories` endpoint */
type DirectoryReq struct {
	Directory string `json:"directory"`
}

var err error

/* Payloads is generic type for what can be returned as payload */
type Payloads interface {
	string | DirectoryResp
}

/* writeHeaderAndBody helper type */
func writeHeaderAndBody[T Payloads](w http.ResponseWriter, body T, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response[T]{
		Payload:    body,
		StatusCode: status,
	})
}

func main() {
	r := chi.NewRouter()
	defer http.ListenAndServe(":3137", r)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		writeHeaderAndBody(w, "Welcome!", http.StatusOK)
	})

	r.Get("/directories", func(w http.ResponseWriter, r *http.Request) {
		/* parse request body into directory request format */
		var req DirectoryReq

		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeHeaderAndBody(w, "Malformed request", http.StatusBadRequest)
			return
		}

		if req.Directory == "" {
			writeHeaderAndBody(w, "Malformed request", http.StatusBadRequest)
			return
		}

		var files []fs.FileInfo
		if files, err = ioutil.ReadDir(req.Directory); err != nil {
			writeHeaderAndBody(w, "Malformed path provided", http.StatusBadRequest)
			return
		}

		/* split files and directories */
		var directoriesFound []string
		var filesFound []string
		for _, file := range files {
			if file.IsDir() {
				directoriesFound = append(directoriesFound, file.Name())
			} else {
				filesFound = append(filesFound, file.Name())
			}
		}

		/* return */
		writeHeaderAndBody(w, DirectoryResp{
			Folders: directoriesFound,
			Files:   filesFound,
		}, http.StatusOK)
	})

}
