package filesystem

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func Dir(filename string) string {
	dir, err := filepath.Abs(filepath.Dir(filename))
	if err != nil {
		log.Fatal(err)
	}
	return strings.Replace(dir, "\\", "/", -1) //将\替换成/
}

func UploadFile(content []byte, fileName string, dir string) (fullFileName string, err error) {
	if exist, err := PathExists(dir); !exist || err != nil {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return "", err
		}
	}

	fullFileName = strings.TrimRight(dir, "/") + "/" + fileName
	file, err := os.OpenFile(fullFileName, os.O_CREATE, os.ModePerm)
	if err != nil {
		return "", err
	}

	if _, err = file.Write(content); err != nil {
		return "", err
	}

	return fullFileName, nil
}
