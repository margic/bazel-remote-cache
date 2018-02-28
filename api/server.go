package api

import (
	"context"
	"io"
	"net/http"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/gorilla/mux"
	"github.com/margic/bazel-s3-cache/store"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Server the api server that serves the cache endpoints
// See https://docs.bazel.build/versions/master/remote-caching.html#http-caching-protocol
type Server struct {
	Handler http.Handler
	Store   store.Storer
}

// NewServer creates a new server
func NewServer(stor store.Storer) (*Server, error) {
	r := mux.NewRouter()
	s := &Server{
		Store: stor,
	}
	r.HandleFunc("/ac/{id}", s.Putac).Methods("Put")
	r.HandleFunc("/ac/{id}", s.Getac).Methods("Get")
	r.HandleFunc("/cas/{id}", s.Putcas).Methods("Put")
	r.HandleFunc("/cas/{id}", s.Getcas).Methods("Get")
	r.Handle("/metrics", promhttp.Handler())

	s.Handler = r
	return s, nil
}

// Putac  handles bazel cache put
func (s *Server) Putac(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := s.Store.Put(context.Background(), "ac", id, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}

// Putcas handles bazel cache put
func (s *Server) Putcas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	err := s.Store.Put(context.Background(), "cas", id, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.WriteHeader(http.StatusNoContent)
}

// Getac handles bazel cache get
func (s *Server) Getac(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, err := s.Get(w, id, "ac")
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Getcas handles bazel cache get
func (s *Server) Getcas(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]
	_, err := s.Get(w, id, "cas")
	if err != nil {
		if strings.Contains(err.Error(), "404") {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

// Get does the actual get
func (s *Server) Get(w io.Writer, id, path string) (int64, error) {
	// buf here is not going to work well for large objects will change this to a rolling buffer
	// as soon as write one
	// TODO ^^ THIS
	buf := aws.NewWriteAtBuffer(make([]byte, 2048))
	size, err := s.Store.Get(context.Background(), path, id, buf)
	if err != nil {
		return size, err
	}

	_, err = w.Write(buf.Bytes()[0:size])
	if err != nil {
		return 0, err
	}
	return size, nil
}
