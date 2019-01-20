package files

import (
	"devlab/lib/errors"
	"devlab/lib/yml"
	"os"
	"path/filepath"
	"reflect"
)

// ReadMainConfig ...
func ReadMainConfig() (config map[string]map[string]string, err error) {
	_, err = IsExists("/.config")
	if errors.CheckAndReturnIfError(err) {
		return make(map[string]map[string]string), err
	}

	pathToConfig, err := AbsolutePath(".config")

	configData, err := ReadTextFile(pathToConfig)
	if errors.CheckAndReturnIfError(err) {
		return make(map[string]map[string]string), err
	}

	config, err = yml.ParseTwoLevelYAML(configData)
	if errors.CheckAndReturnIfError(err) {
		return make(map[string]map[string]string), err
	}

	return
}

// ReadContextConfig ...
func ReadContextConfig(relativePathToContextConfigFile string) (context map[string]map[string]map[string]string, err error) {
	contextData, err := ReadTextFile(relativePathToContextConfigFile)
	if errors.CheckAndReturnIfError(err) {
		return make(map[string]map[string]map[string]string), err
	}

	context, err = yml.ParseThreeLevelYAML(contextData)
	if errors.CheckAndReturnIfError(err) {
		return make(map[string]map[string]map[string]string), err
	}

	return
}

// ReadOneLevelYaml ...
func ReadOneLevelYaml(relativePathToYamlFile string) (data map[string]string, err error) {
	dataYAML, err := ReadTextFile(relativePathToYamlFile)
	if errors.CheckAndReturnIfError(err) {
		return make(map[string]string), err
	}

	data, err = yml.ParseOneLevelYAML(dataYAML)
	if errors.CheckAndReturnIfError(err) {
		return make(map[string]string), err
	}

	return
}

// ReadTwoLevelYaml ...
func ReadTwoLevelYaml(relativePathToYamlFile string) (data map[string]map[string]string, err error) {
	dataYAML, err := ReadTextFile(relativePathToYamlFile)
	if errors.CheckAndReturnIfError(err) {
		return make(map[string]map[string]string), err
	}

	data, err = yml.ParseTwoLevelYAML(dataYAML)
	if errors.CheckAndReturnIfError(err) {
		return make(map[string]map[string]string), err
	}

	return
}

// WriteAppendFile writes string to the of file
func WriteAppendFile(filenamePath string, text string) (result int, err error) {
	filenameDir := filepath.Dir(filenamePath)
	if isFilenameDirExists, _ := IsExists(filenameDir); !isFilenameDirExists {
		CreateDir(filenameDir)
	}

	absoluteFilenamePath, err := AbsolutePath(filenamePath)

	if isFileExists, _ := IsExists(filenamePath); !isFileExists {
		file, _ := os.Create(absoluteFilenamePath)
		defer file.Close()
	}

	file, _ := os.OpenFile(absoluteFilenamePath, os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return
	}
	defer file.Close()

	result, err = file.WriteString(text + "\n")

	return
}

// WriteYaml  writes tree stucture data (with 3 levels) to yaml file
func WriteYaml(filenamePath string, data interface{}) (err error) {
	return WriteYamlBranch(filenamePath, data, 0)
}

// WriteAppendFileWithIndent adds spaces indent before string anf write this string to the end of file
func WriteAppendFileWithIndent(filenamePath string, text string, indent int) (result int, err error) {
	return WriteAppendFile(filenamePath, indentInSpaces(indent)+text)
}

// WriteYamlBranch writes one branch of tree stucture data to yaml file
func WriteYamlBranch(filenamePath string, data interface{}, depth int) (err error) {
	if data, ok := data.(map[string]map[string]map[string]string); ok {
		for key, value := range data {
			_, err = WriteAppendFileWithIndent(filenamePath, key+": ", 2*depth)
			err = WriteYamlBranch(filenamePath, value, depth+1)
		}
	}

	if data, ok := data.(map[string]map[string]string); ok {
		for key, value := range data {
			_, err = WriteAppendFileWithIndent(filenamePath, key+": ", 2*depth)
			err = WriteYamlBranch(filenamePath, value, depth+1)
		}
	}

	if data, ok := data.(map[string]string); ok {
		for key, value := range data {
			_, err = WriteAppendFileWithIndent(filenamePath, key+": "+value, 2*depth)
		}
	}

	return
}

// indentInSpaces returns indent with num spaces
func indentInSpaces(indent int) (spacesIndent string) {
	spacesIndent = ""
	for i := 0; i < indent; i++ {
		spacesIndent += " "
	}
	return spacesIndent
}

// Pair ...
type Pair struct {
	key   string
	value *interface{}
}

// writeAppendYAMLFile ...
func (p Pair) writeAppendYAMLFile(filenamePath string, depth int) {
	v := reflect.ValueOf(*(p.value))
	switch v.Kind() {
	case reflect.String:
		WriteAppendFileWithIndent(filenamePath, p.key+": "+v.Elem().String(), depth)

	case reflect.Slice, reflect.Array:
		WriteAppendFileWithIndent(filenamePath, p.key+": ", depth)
		for i := 0; i < v.NumField(); i++ {
			item := v.Field(i)
			if reflect.ValueOf(item).Kind() == reflect.String {
				WriteAppendFileWithIndent(filenamePath, "- "+reflect.ValueOf(item).Elem().String(), depth+2)
			} else {
				WriteAppendFileWithIndent(filenamePath, " - the value couldn't be printed", depth+2)
			}
		}

	case reflect.Map:
		WriteAppendFileWithIndent(filenamePath, p.key+": ", depth)
		for _, key := range v.MapKeys() {
			// (Pair{key.Elem().String(), v.MapIndex(key).Pointer()}).writeAppendYAMLFile(filenamePath, depth + 2)
			// (Pair{key.Elem().String(), &"unknown value"}).writeAppendYAMLFile(filenamePath, depth + 2)
			WriteAppendFileWithIndent(filenamePath, key.Elem().String()+": unknown value", depth+2)
		}

	default:
		WriteAppendFileWithIndent(filenamePath, p.key+": the value couldn't be printed", depth)
	}
}

// YamlTree ...
type YamlTree []Pair

// WriteYaml ...
func (t YamlTree) WriteYaml(filenamePath string) {
	for _, pair := range t {
		pair.writeAppendYAMLFile(filenamePath, 0)
	}
}
