package urlshort

import (
	"fmt"
	"net/http"
	"gopkg.in/yaml.v2"
)

func MapHandler(pathsToUrls map[string]string, fallback http.Handler) http.HandlerFunc {
	// As the return type is http.Handlerfunc, this function is coerced into 
	// satisfying the Handler interface
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

// A struct for unmarshalling the YAML
type YamlUM []struct {
    Path string `yaml:"path,omitempty"`
	Url  string `yaml:"url,omitempty"`
}

func parseYAML(yml []byte, urlPaths *YamlUM) error {
	err := yaml.Unmarshal(yml, urlPaths)
	return err
}

func buildMap(pathStruct YamlUM) map[string]string {
	pathMap := make(map[string]string)
	// Traverse the struct and add data to the map
	for _, pathdata := range pathStruct {
		fmt.Printf("Path data:\n%v\n", pathdata.Path)
		pathMap[pathdata.Path] = pathdata.Url
	}
	return pathMap
}

func YAMLHandler(yml []byte, fallback http.Handler) (http.HandlerFunc, error) {
	// Attempt to parse YAML
	var parsedYaml YamlUM
    err := parseYAML(yml, &parsedYaml)
	if err != nil {
		return nil, err
	}

	// YAML has been successfully parsed so build the map
	pathMap := buildMap(parsedYaml)
	return MapHandler(pathMap, fallback), nil
}


