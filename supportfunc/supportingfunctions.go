package supportfunc

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func GetRootPath(rootDir string) (string, error) {
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "", err
	}

	tmp := strings.Split(currentDir, "/")

	if tmp[len(tmp)-1] == rootDir {
		return currentDir, nil
	}

	var path string = ""
	for _, v := range tmp {
		path += v + "/"

		if v == rootDir {
			return path, nil
		}
	}

	return path, nil
}

func ReadFileJson(fpath, fname string) ([]byte, error) {
	var newResult []byte

	rootPath, err := GetRootPath("ISEMS-MRSICT")
	if err != nil {
		return newResult, err
	}

	fmt.Println("func 'readFileJson', path = ", path.Join(rootPath, fpath, fname))

	f, err := os.OpenFile(path.Join(rootPath, fpath, fname), os.O_RDONLY, os.ModePerm)
	if err != nil {
		return newResult, err
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	for sc.Scan() {
		newResult = append(newResult, sc.Bytes()...)
	}

	return newResult, nil
}

func GetWhitespace(num int) string {
	var str string

	if num == 0 {
		return str
	}

	for i := 0; i < num; i++ {
		str += "  "
	}

	return str
}
