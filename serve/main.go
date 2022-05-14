package main

import (
	"encoding/json"
	"fmt"
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

type DirectoryReq struct {
	Directory string `json:"directory"`
}

var err error

type Payloads interface {
	string | DirectoryResp
}

func writeHeaderAndBody[T Payloads](w http.ResponseWriter, body T, status int) {
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(Response[T]{
		Payload:    body,
		StatusCode: status,
	})
}

func main() {
	r := chi.NewRouter()

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

		var directoriesFound []string
		var filesFound []string
		for _, file := range files {
			fmt.Println(file.Name())
			if file.IsDir() {
				directoriesFound = append(directoriesFound, file.Name())
			} else {
				filesFound = append(filesFound, file.Name())
			}
		}

		writeHeaderAndBody(w, DirectoryResp{
			Folders: directoriesFound,
			Files:   filesFound,
		}, http.StatusOK)
	})

	http.ListenAndServe(":3137", r)
}
