package main

import (
	"fmt"
	pgrepo "forms/pg"
	"forms/web/cvsui"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
)

type App struct {
	handler http.Handler
	q       *pgrepo.Queries
}

func NewApp(q *pgrepo.Queries) *App {
	app := &App{
		q: q,
	}
	app.setRoutes()
	return app
}

func (a *App) setRoutes() {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /cvs", a.cvs())
	mux.HandleFunc("GET /cvs/{id}", a.cv())
	mux.HandleFunc("GET /cvs/new", a.newCV())
	mux.HandleFunc("POST /cvs", a.createCV())

	a.handler = mux
}

func (a *App) cvs() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		files, err := a.q.GetFiles(r.Context())
		if err != nil {
			fmt.Println(err)
			http.Error(w, "error listing files", http.StatusInternalServerError)
			return
		}

		cvsui.Index(files).Render(r.Context(), w)
	}
}

func (a *App) newCV() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		cvsui.New().Render(r.Context(), w)
	}
}

func (a *App) cv() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id, err := strconv.Atoi(r.PathValue("id"))
		if err != nil {
			fmt.Println("error parsing id")
			http.Redirect(w, r, "/cvs", http.StatusFound)
			return
		}

		fr, err := a.q.GetFile(r.Context(), int64(id))
		if err != nil {
			fmt.Println("error retrieving file")
			http.Redirect(w, r, "/cvs", http.StatusFound)
			return
		}

		w.Header().Set("Content-Length", strconv.Itoa(len(fr.File)))
		w.Header().Set("Content-Type", fr.Mimetype)
		w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fr.Name))
		w.Write(fr.File)
	}
}

func getFileMimeType(fh *multipart.FileHeader) (string, error) {
	file, err := fh.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Lê até 512 bytes do arquivo
	buffer := make([]byte, 512)
	n, err := file.Read(buffer)
	if err != nil && err != io.EOF {
		return "", err
	}

	// Detecta o tipo de conteúdo
	mimeType := http.DetectContentType(buffer[:n])
	return mimeType, nil
}

func (a *App) createCV() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		maxMemory := maxMemoryForRequest(r)
		err := r.ParseMultipartForm(maxMemory)
		if err != nil {
			http.Error(w, "error parsing body", http.StatusInternalServerError)
			return
		}

		defer r.MultipartForm.RemoveAll()

		if len(r.MultipartForm.File["files"]) == 0 {
			fmt.Println("no file uploaded")
			http.Redirect(w, r, "/cvs/new", http.StatusFound)
			return
		}

		fhs := r.MultipartForm.File["files"]
		for _, fh := range fhs {
			params, err := buildCreateFileParams(fh)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "error reading file", http.StatusInternalServerError)
				return
			}

			params.Mimetype, err = getFileMimeType(fh)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "error finding mimetype out", http.StatusInternalServerError)
				return
			}

			_, err = a.q.CreateFile(r.Context(), params)
			if err != nil {
				fmt.Println(err)
				http.Error(w, "error reading file", http.StatusInternalServerError)
				return
			}
		}
		http.Redirect(w, r, "/cvs", http.StatusFound)
	}
}

func buildCreateFileParams(fh *multipart.FileHeader) (pgrepo.CreateFileParams, error) {
	f, err := fh.Open()
	if err != nil {
		return pgrepo.CreateFileParams{}, fmt.Errorf("open file: %w", err)
	}

	content, err := io.ReadAll(f)
	if err != nil {
		return pgrepo.CreateFileParams{}, fmt.Errorf("read file: %w", err)
	}

	return pgrepo.CreateFileParams{
		Name: fh.Filename,
		File: content,
	}, nil
}

func (a *App) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	a.handler.ServeHTTP(w, r)
}

func maxMemoryForRequest(r *http.Request) int64 {
	var m int64 = 10 << 20
	stringContentLength := r.Header.Get("Content-Length")
	if strings.TrimSpace(stringContentLength) == "" {
		return m
	}

	var contentLength int64
	var err error
	contentLength, err = strconv.ParseInt(stringContentLength, 10, 64)
	if err != nil {
		return m
	}

	if m < contentLength {
		return m
	}

	return contentLength
}
