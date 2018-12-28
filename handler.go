package urlshort

import (
	"encoding/json"
	"fmt"
	"net/http"

	"gopkg.in/yaml.v2"
)

// A struct for unmarshaling the YAML/JSON data
type PathStore []struct {
	Path string `json:"path",omitempty yaml:"path,omitempty"`
	Url  string `json:"url",omitempty yaml:"url,omitempty"`
}

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	// As the return type is http.Handlerfunc, this function is coerced into
	// satisfying the Handler interface
	fmt.Println("Map Handler")
	return func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path
		longURL, ok := pathsToUrls[path]
		if ok {
			http.Redirect(w, r, longURL, http.StatusFound)
		} else {
			// ServeHTTP method dispatches the request to the handler
			// defined for the default mux
			fallback.ServeHTTP(w, r)
		}
	}
}

func buildMap(pathStruct PathStore) map[string]string {
	fmt.Println("Building map")
	pathMap := make(map[string]string)
	// Traverse the struct and add data to the map
	for _, pathdata := range pathStruct {
		pathMap[pathdata.Path] = pathdata.Url
	}
	return pathMap
}

func JSONHandler(jsonData []byte, fallback http.Handler) (http.HandlerFunc, error) {
	fmt.Println("JSON Handler")
	var parsedJson PathStore
	err := json.Unmarshal(jsonData, &parsedJson)
	if err != nil {
		print(err)
		return nil, err
	}
	pathMap := buildMap(parsedJson)
	return MapHandler(pathMap, fallback), nil
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	fmt.Println("YAML Handler")
	var parsedYaml PathStore
	err := yaml.Unmarshal(yml, &parsedYaml)
	if err != nil {
		return nil, err
	}
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}
