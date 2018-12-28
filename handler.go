package urlshort

import (
	"fmt"
	"net/http"
	"encoding/json"

	"gopkg.in/yaml.v2"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	// As the return type is http.Handlerfunc, this function is coerced into 
	// satisfying the Handler interface
	fmt.Println("Running Map Handler")
	return func(w http.ResponseWriter, r *http.Request){
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

// A struct for unmarshaling the data
type PathStore []struct {
    Path string `json:"path",omitempty yaml:"path,omitempty"`
	Url  string `json:"url",omitempty yaml:"url,omitempty"`
}

func parseYAML(yml []byte, urlPaths *PathStore) error {
	fmt.Println("Attempting to parse YAML")
	err := yaml.Unmarshal(yml, urlPaths)
	return err
}

func parseJSON(jsonData []byte, urlPaths *PathStore) error {
	fmt.Println("Attempting to parse JSON")
	err := json.Unmarshal(jsonData, urlPaths)
    return err
}

func buildMap(pathStruct PathStore) map[string]string {
	fmt.Println("Attempting to build map")
	pathMap := make(map[string]string)
	// Traverse the struct and add data to the map
	for _, pathdata := range pathStruct {
		//fmt.Printf("Path data:\n%v\n", pathdata.Path)
		pathMap[pathdata.Path] = pathdata.Url
	}
	return pathMap
}

func JSONHandler(json []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsedJson PathStore
	err := parseJSON(json, &parsedJson)
	if err != nil {
		print(err)
		return nil, err
	}
	pathMap := buildMap(parsedJson)
	return MapHandler(pathMap, fallback), nil
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	var parsedYaml PathStore
    err := parseYAML(yml, &parsedYaml)
	if err != nil {
		print(err)
		return nil, err
	}
	// YAML has been successfully parsed so build the map
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}


