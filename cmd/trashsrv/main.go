package main

import (
	"archive/tar"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"mcquay.me/trash"
)

func init() {
	log.SetFlags(log.Lshortfile)
}

type fake struct {
	name string
	size int64
}

func (f fake) Name() string {
	return f.name
}

func (f fake) Size() int64 {
	return f.size
}

func (f fake) Mode() os.FileMode {
	return 0644
}

func (f fake) ModTime() time.Time {
	return time.Now()
}

func (f fake) IsDir() bool {
	return false
}

func (f fake) Sys() interface{} {
	return nil
}

func tarball(w http.ResponseWriter, req *http.Request) {
	if err := req.ParseForm(); err != nil {
		http.Error(w, fmt.Sprintf("error parsing form: %v", err), http.StatusBadRequest)
		return
	}

	size, ok := req.Form["size"]
	if !ok {
		http.Error(w, "must provide size parameter", http.StatusBadRequest)
		return
	}
	if len(size) != 1 {
		http.Error(w, fmt.Sprintf("only specify one size (%d specified)", len(size)), http.StatusBadRequest)
		return
	}
	sz, err := strconv.ParseFloat(size[0], 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("error parsing size: %v", err), http.StatusBadRequest)
		return
	}
	sz = sz * 1024 * 1024 * 1024
	_, name := filepath.Split(req.URL.Path)
	if name == "" {
		name = "file.tar"
	}
	w.Header().Set("Content-Disposition", fmt.Sprintf("inline; filename=%q", name))
	w.Header().Set("content-type", "application/x-tar")
	lr := io.LimitReader(trash.Reader, int64(sz))
	t := tar.NewWriter(w)
	hdr, err := tar.FileInfoHeader(fake{"foo", int64(sz)}, "")
	if err != nil {
		http.Error(w, fmt.Sprintf("error creating tar header: %v", err), http.StatusBadRequest)
		return
	}
	hdr.Name = "trash.dat"
	if err := t.WriteHeader(hdr); err != nil {
		http.Error(w, fmt.Sprintf("error writing tar header: %v", err), http.StatusBadRequest)
		return
	}
	if _, err := io.Copy(t, lr); err != nil {
		log.Printf("problem writing tar to http body: %v", err)
		return
	}
	if err := t.Close(); err != nil {
		log.Printf("problem closing tar: %v", err)
		return
	}
}

var port = flag.Int64("-port", 9292, "port")

func main() {
	flag.Parse()
	http.HandleFunc("/", tarball)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", *port), nil); err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
