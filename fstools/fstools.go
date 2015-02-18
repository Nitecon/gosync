package fstools

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"syscall"
)

type FsItem struct {
	Filename string
	Path     string
	Checksum string
	IsDir    bool
	Mtime    int
	Uid      int
	Gid      int
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
			Mtime:    int(info.ModTime().Unix()),
			Uid:      int(info.Sys().(*syscall.Stat_t).Uid),
			Gid:      int(info.Sys().(*syscall.Stat_t).Gid),
			Perms:    info.Mode().Perm().String(),
		})
		return nil
	})
	if err != nil {
		checkErr(err, "Unable to read file info")
	}

	return fsItems
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
