package files

import (
  "io"
  "os"
  "path/filepath"
  "devlab/lib/errors"
  "devlab/lib/yml"
  "devlab/lib/logger"
)

func AbsolutePath(relativePath string) (absolutePath string, err error) {
  dir, err := filepath.Abs(relativePath)
  if  errors.CheckAndReturnIfError(err) { return }

  return dir, nil
}

func ReadTextFile(path string) (resultString string, err error) {   
  resultString = ""

  filepath, err := AbsolutePath(path)
  if errors.CheckAndReturnIfError(err) { return }

  file, err := os.Open(filepath)
  if errors.CheckAndReturnIfError(err) { return }
  defer file.Close() 
     
  data := make([]byte, 64)     
  for {
    n, err := file.Read(data)
    if err == io.EOF { break }
    resultString += string(data[:n])
  }
  
  err = nil
  return
}

func IsExists(path string) (bool, error) {
  filepath, err := AbsolutePath(path)
  if errors.CheckAndReturnIfError(err) { return false, err}

  _, err = os.Stat(filepath)
  if err == nil { return true, nil }
  if os.IsNotExist(err) { return false, nil }
  return true, err
}

func CreateDir(path string) error {
  filepath, err := AbsolutePath(path)
  if errors.CheckAndReturnIfError(err) { return err }

  return os.MkdirAll(filepath, 0755)
}


func ReadMainConfig() (config map[string]string, err error) {
  _, err =  IsExists("/.config")
  if errors.CheckAndReturnIfError(err) { return make(map[string]string), err }

  pathToConfig, err := AbsolutePath(".config")
  
  configData, err := ReadTextFile(pathToConfig)
  if errors.CheckAndReturnIfError(err) { return  make(map[string]string), err }
  
  config, err = yml.ParseOneLevelYAML(configData)
  if errors.CheckAndReturnIfError(err) { return  make(map[string]string), err }
  
  return
}

func ReadContextConfig(relativePathToContextConfigFile string) (context map[string]map[string]map[string]string, err error) {
	contextData, err := ReadTextFile(relativePathToContextConfigFile)
	if errors.CheckAndReturnIfError(err) { return  make(map[string]map[string]map[string]string), err }
  
	context, err = yml.ParseThreeLevelYAML(contextData)
	if errors.CheckAndReturnIfError(err) { return  make(map[string]map[string]map[string]string), err }

	return
}

/**
* Writes string to the of file
*/
func WriteAppendFile(filenamePath string, text string) (result int, err error) { 
  absoluteFilenamePath, err := AbsolutePath(filenamePath)
  isFileExists, _ :=  IsExists(filenamePath)
  logger.Debug(" %s ", absoluteFilenamePath)

  if !isFileExists {
    logger.Debug("text:  %s ", text)
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

/**
* It adds spaces indent before string anf write this string to the end of file
*/
func WriteAppendFileWithIndent(filenamePath string, text string, indent int) (result int, err error) {
  return WriteAppendFile(filenamePath, indentInSpaces(indent) + text)
}

/**
* Writes tree stucture data (with 3 levels) to yaml file
*/
func WriteYaml(filenamePath string, data interface{}) (err error) {   
  return WriteYamlBranch(filenamePath, data, 0)
}

/**
* Writes one branch of tree stucture data to yaml file
*/
func WriteYamlBranch(filenamePath string, data interface{}, depth int) (err error) { 
  if data, ok := data.(map[string]map[string]map[string]string); ok {
    for key, value := range data {
      _, err = WriteAppendFileWithIndent(filenamePath, key + ": ", 2*depth)
      err = WriteYamlBranch(filenamePath, value, depth + 1)
    }
  }

  if data, ok := data.(map[string]map[string]string); ok {
    for key, value := range data {
      _, err = WriteAppendFileWithIndent(filenamePath, key + ": ", 2*depth)
      err = WriteYamlBranch(filenamePath, value, depth + 1)
    }
  }

  if data, ok := data.(map[string]string); ok {
    for key, value := range data {
      _, err = WriteAppendFileWithIndent(filenamePath, key + ": " + value, 2*depth)
    }
  }

  return
}

/**
* Returns indent with num spaces
*/
func indentInSpaces(indent int) (spacesIndent string) {
  spacesIndent = ""
  for i := 0; i < indent; i++ {
    spacesIndent += " "
  } 
  return spacesIndent
}
