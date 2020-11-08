package handlers

import (
	"io"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/specter25/microservices-in-go/products-images/files"

	"github.com/hashicorp/go-hclog"
)

// Files is a handler for reading and writing files
type Files struct {
	log   hclog.Logger
	store files.Storage
}

// NewFiles creates a new File handler
func NewFiles(s files.Storage, l hclog.Logger) *Files {
	return &Files{store: s, log: l}
}

//UploadRest something
func (f *Files) UploadRest(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	fn := vars["filename"]

	if id == "" || fn == "" {
		f.invalidURI(r.URL.String(), rw)
		return
	}
	f.saveFile(id, fn, rw, r.Body)

	f.log.Info("Handle POST", "id", id, "Filename", fn)
}

func (f *Files) invalidURI(uri string, rw http.ResponseWriter) {
	f.log.Error("Invalid path", "path", uri)
	http.Error(rw, "Invalid file path should be in the format: /[id]/[filepath]", http.StatusBadRequest)
}

//UploadMultipart something
func (f *Files) UploadMultipart(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024)
	if err != nil {
		http.Error(rw, "Expected Multipart form data", http.StatusBadRequest)
		f.log.Error("Bad request -- Expected Multipart form data")
		return
	}

	id, idErr := strconv.Atoi(r.FormValue("id"))
	f.log.Info("Process Form for id", id)

	if idErr != nil {
		http.Error(rw, "Expected interger id", http.StatusBadRequest)
		f.log.Error("Bad request -- Expected interger id")
		return
	}

	fi, mh, err := r.FormFile("file")
	if err != nil {
		http.Error(rw, "Bad request", http.StatusBadRequest)
		f.log.Error("Expected file")
		return
	}
	f.saveFile(r.FormValue("id"), mh.Filename, rw, fi)
}

// saveFile saves the contents of the request to a file
func (f *Files) saveFile(id, path string, rw http.ResponseWriter, r io.ReadCloser) {
	f.log.Info("Save file for product", "id", id, "path", path)

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r)
	if err != nil {
		f.log.Error("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}
