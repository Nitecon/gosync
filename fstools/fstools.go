package fstools

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

type FsItem struct {
	Filename string
	Path     string
	Checksum string
	IsDir    bool
	Mtime    time.Time
	Perms    string
}

func ListFilesInDir(doc_path string) []FsItem {
	searchDir := doc_path

	fsItems := []FsItem{}
	err := filepath.Walk(searchDir, func(doc_path string, info os.FileInfo, err error) error {
		//add to slice here
		fsItems = append(fsItems, FsItem{
			Filename: doc_path,
			IsDir:    info.IsDir(),
			Checksum: GetMd5Checksum(doc_path),
			Mtime:    info.ModTime().UTC(),
			Perms:    info.Mode().Perm().String(),
		})
		return nil
	})
	if err != nil {
		checkErr(err, "Unable to read file info")
	}

	return fsItems
}

func GetFileInfo(doc_path string) (FsItem, error){
    fi, err := os.Stat(doc_path)
    var item = FsItem{}
    if err != nil{
        return item, err
    }else{
        item.Filename = doc_path
        item.IsDir = fi.IsDir()
        item.Checksum = GetMd5Checksum(doc_path)
        item.Mtime = fi.ModTime().UTC()
        item.Perms = fi.Mode().Perm().String()
    }
    return item, nil
}

func GetMd5Checksum(filepath string) string {
	if b, err := ComputeMd5(filepath); err == nil {
		md5string := fmt.Sprintf("%x", b)
		return md5string
	} else {
		return "DirectoryMD5"
	}
}

func ComputeMd5(filePath string) ([]byte, error) {
	var result []byte
	file, err := os.Open(filePath)
	if err != nil {
		return result, err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return result, err
	}

	return hash.Sum(result), nil

}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
