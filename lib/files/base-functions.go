package files

import (
  "io"
  "os"
  "text/template"
  "path/filepath"
  "devlab/lib/errors"
)

/**
*/
func AbsolutePath(relativePath string) (absolutePath string, err error) {
  dir, err := filepath.Abs(relativePath)
  if  errors.CheckAndReturnIfError(err) { return }

  return dir, nil
}

/**
*/
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

/**
*/
func IsExists(path string) (bool, error) {
  filepath, err := AbsolutePath(path)
  if errors.CheckAndReturnIfError(err) { return false, err}

  _, err = os.Stat(filepath)
  if err == nil { return true, nil }
  if os.IsNotExist(err) { return false, nil }
  return true, err
}

/**
*/
func CreateDir(path string) error {
  filepath, err := AbsolutePath(path)
  if errors.CheckAndReturnIfError(err) { return err }

  return os.MkdirAll(filepath, 0755)
}

/* Copy the src file to dst. Any existing file will be overwritten and will not
   copy file attributes.
*/ 
func Copy(src, dst string) error {
  in, err := os.Open(src)
  if err != nil {
      return err
  }
  defer in.Close()

  out, err := os.Create(dst)
  if err != nil {
      return err
  }
  defer out.Close()

  _, err = io.Copy(out, in)
  if err != nil {
      return err
  }
  return out.Close()
}

/**
*
*/
func Delete(path string) (err error) {
  isFileExists, _ :=  IsExists(path)
  if !isFileExists { return }

  absolutePath, _ := AbsolutePath(path)
  
  return os.Remove(absolutePath)
}

/** Render text template   
*/
func RenderTextTemplate(src, dst string, params interface{}) error { 
  out, err := os.Create(dst)
  if err != nil {
      return err
  }
  defer out.Close()

  sourceTemplate, err := template.ParseFiles(src)
  if err != nil { return err }
  
  return sourceTemplate.Execute(out, params)  
}


