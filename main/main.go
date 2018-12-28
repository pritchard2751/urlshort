package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	//"path/filepath"
	"flag"
	"strings"

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

var pathsToUrls = UrlPaths{
	"/urlshort-godoc": "https://godoc.org/github.com/gophercises/urlshort",
	"/yaml-godoc":     "https://godoc.org/gopkg.in/yaml.v2",
	"/goo":            "https://google.com",
}

// Here we have implemented the ServeHTTP method on the UrlPaths object
// meaning that it satisfies the Handler interface
func (paths UrlPaths) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	longURL, ok := paths[path]
	if ok {
		http.Redirect(w, r, longURL, http.StatusFound)
	} else {
		msg := fmt.Sprintf("No such page: %s\n", r.URL)
		http.Error(w, msg, http.StatusNotFound)
	}
}

func getFlagOptions(dir string) ([]string, error) {
	// Iterate specified directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}
	// Construct flag help text based on file types in 'data' folder
	var options = []string{"Specify data source type, options:"}
	for _, file := range files {
		fn := file.Name()
		extindex := strings.Index(fn, ".") + 1
		options = append(options, fn[extindex:len(fn)])
	}
	return options, nil

}

func main() {

	// Create flag to allow user to specify the data source type
	options, err := getFlagOptions("../data")
	if err != nil {
		panic(err)
	}

	var flagvar string
	var help string = strings.Join(options, "|")
	flag.StringVar(&flagvar, "dst", "map", help)
	flag.Parse()

	// Default HTTP request router
	mux := defaultMux()
	// Build the MapHandler using the mux as the fallback
	mapHandler := urlshort.MapHandler(pathsToUrls, mux)
	// Set as default handler
	handler := mapHandler

	switch flagvar {
	case "yaml":
		yaml, err := ioutil.ReadFile("../data/redirect.yaml")
		if err != nil {
			panic(err)
		}
		yamlHandler, err := urlshort.YAMLHandler(yaml, mapHandler)
		if err != nil {
			panic(err)
		}
		handler = yamlHandler
	case "json":
		json, err := ioutil.ReadFile("../data/redirect.json")
		if err != nil {
			panic(err)
		}
		jsonHandler, err := urlshort.JSONHandler(json, mapHandler)
		if err != nil {
			panic(err)
		}
		handler = jsonHandler
	}

	fmt.Println("Starting the server on :8080")
	//http.ListenAndServe(":8080", mux)
	//http.ListenAndServe(":8080", pathsToUrls)
	http.ListenAndServe(":8080", handler)
}
