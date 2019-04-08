package api

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"

	"github.com/metiago/zbx1/common/request"

	"github.com/metiago/zbx1/repository"
	"github.com/mitchellh/mapstructure"
)

const (
	authHeader    string = "Authorization"
	maxSizeInByte int64  = 16000000
)

func fileUpload(w http.ResponseWriter, r *http.Request) {

	r.Body = http.MaxBytesReader(w, r.Body, maxSizeInByte)
	if err := r.ParseMultipartForm(maxSizeInByte); err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}

	file, handler, err := r.FormFile("file")
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}
	defer file.Close()

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		log.Println(err)
		request.Handle500(w, err)
		return
	}

	// FIXME REFACTOR IT, CREATE A HELPER FOR IT
	tokenString := r.Header.Get(authHeader)
	claims := jwt.MapClaims{}
	_, err = jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return verifyKey, nil
	})
	if err != nil {
		request.Handle500(w, err)
	}

	// FIXME REFACTOR IT, CREATE A SERVICE LAYER
	var u repository.User
	mapstructure.Decode(claims["uinf"], &u)
	var f repository.File
	f.Name = handler.Filename
	f.Ext = filepath.Ext(handler.Filename)
	f.Data = buf.Bytes()
	err = repository.FileUpload(u, f)
	if err != nil {
		log.Println(err)
		request.Handle500(w, err)
	}

	// FIXME CHANGE RESPONSE
	fmt.Fprintf(w, "%v", handler.Header)
}

func fileDownload(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["ID"])
	if err != nil {
		request.Handle500(w, err)
		return
	}

	f, err := repository.FileDownload(id)
	if err != nil {
		if err == sql.ErrNoRows {
			request.Handle404(w)
			return
		}
		request.Handle500(w, err)
		return
	}

	mime := http.DetectContentType(f.Data)

	contentLength := len(string(f.Data))

	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Disposition", "attachment; filename="+f.Name+"")
	w.Header().Set("Expires", "0")
	w.Header().Set("Content-Transfer-Encoding", "binary")
	w.Header().Set("Content-Length", strconv.Itoa(contentLength))
	w.Header().Set("Content-Control", "private, no-transform, no-store, must-revalidate")

	b := bytes.NewBuffer(f.Data)

	//stream the body to the client without fully loading it into memory
	io.Copy(w, b)
}

func fileFindAllByUsername(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	username := vars["username"]
	result, err := repository.FindaAllFilesByUsername(username)
	if err != nil {
		request.Handle500(w, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		request.Handle500(w, err)
		return
	}
}
