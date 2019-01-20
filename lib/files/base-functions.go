package files

import (
	"devlab/lib/errors"
	"io"
	"os"
	"path/filepath"
	"text/template"
)

// AbsolutePath ...
func AbsolutePath(relativePath string) (absolutePath string, err error) {
	dir, err := filepath.Abs(relativePath)
	if errors.CheckAndReturnIfError(err) {
		return
	}

	return dir, nil
}

// ReadTextFile ...
func ReadTextFile(path string) (resultString string, err error) {
	resultString = ""

	filepath, err := AbsolutePath(path)
	if errors.CheckAndReturnIfError(err) {
		return
	}

	file, err := os.Open(filepath)
	if errors.CheckAndReturnIfError(err) {
		return
	}
	defer file.Close()

	data := make([]byte, 64)
	for {
		n, err := file.Read(data)
		if err == io.EOF {
			break
		}
		resultString += string(data[:n])
	}

	err = nil
	return
}

// IsExists ...
func IsExists(path string) (bool, error) {
	filepath, err := AbsolutePath(path)
	if errors.CheckAndReturnIfError(err) {
		return false, err
	}

	_, err = os.Stat(filepath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// CreateDir ...
func CreateDir(path string) error {
	filepath, err := AbsolutePath(path)
	if errors.CheckAndReturnIfError(err) {
		return err
	}

	return os.MkdirAll(filepath, 0755)
}

// Copy ...
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

// Delete ...
func Delete(path string) (err error) {
	isFileExists, _ := IsExists(path)
	if !isFileExists {
		return
	}

	absolutePath, _ := AbsolutePath(path)

	return os.Remove(absolutePath)
}

// RenderTextTemplate renders text template
func RenderTextTemplate(src, dst string, params interface{}) (err error) {
	if isSrcExists, err := IsExists(src); !isSrcExists {
		return err
	}

	dstDir := filepath.Dir(dst)
	if isDstDirExists, _ := IsExists(dstDir); !isDstDirExists {
		CreateDir(dstDir)
	}

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	sourceTemplate, err := template.ParseFiles(src)
	if err != nil {
		return err
	}

	return sourceTemplate.Execute(out, params)
}
