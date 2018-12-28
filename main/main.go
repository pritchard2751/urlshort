package main

import (
	"fmt"
	"net/http"
	"io/ioutil"
	"github.com/pritchard2751/urlshort"
)

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", muxHandler)
	return mux
}

func muxHandler(w http.ResponseWriter, r *http.Request) {
	msg := fmt.Sprintf("Default Handler: No such page: %s\n", r.URL)
	http.Error(w, msg, http.StatusNotFound)
}

type UrlPaths map[string]string

// Here we have implemented the ServeHTTP method on the UrlPaths object
// meaning that it satisfies the Handler interface
func (paths UrlPaths) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	longURL, ok := paths[path]
	if ok {
		http.Redirect(w, r, longURL, http.StatusFound)
	}else {
		msg := fmt.Sprintf("No such page: %s\n" , r.URL)
		http.Error(w, msg, http.StatusNotFound)
	}
}

func main() {

     pathsToUrls := UrlPaths{
		"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
		"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	}

	// Default HTTP request router
	mux := defaultMux()
	
	// Build the MapHandler using the mux as the fallback
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	handler := mapHandler

	// Build the YAMLHandler using the mapHandler as the
	// fallback
	yaml, err := ioutil.ReadFile("../data/redirect.yaml")
	if err != nil {
		fmt.Println(err)
	} else {
		yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
		if err != nil {
			panic(err)
		}
		handler = yamlHandler
	}

	fmt.Println("Starting the server on :8080")
	//http.ListenAndServe(":8080", mux)
	//http.ListenAndServe(":8080", pathsToUrls)
	http.ListenAndServe(":8080", handler)
}

