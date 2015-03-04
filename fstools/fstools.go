package fstools

import (
	"crypto/md5"
	"fmt"
	"io"
	"log"
	"os"
    "math"
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
			Perms:    fmt.Sprintf("%#o",info.Mode().Perm()),
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
        item.Perms = fmt.Sprintf("%#o",fi.Mode().Perm())
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


func ComputeMd5(filepath string) (string, error){
    const filechunk = 8192
    file, err := os.Open(filepath)

    if err != nil {
        return "", err
    }

    defer file.Close()

    // calculate the file size
    info, _ := file.Stat()

    filesize := info.Size()

    blocks := uint64(math.Ceil(float64(filesize) / float64(filechunk)))

    hash := md5.New()

    for i := uint64(0); i < blocks; i++ {
        blocksize := int(math.Min(filechunk, float64(filesize-int64(i*filechunk))))
        buf := make([] byte, blocksize)

        file.Read(buf)
        io.WriteString(hash, string(buf))   // append into the hash
    }

    //fmt.Printf("%s checksum is %x\n",file.Name(), hash.Sum(nil))
    return fmt.Sprintf("%x",hash.Sum(nil)), nil
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatalln(msg, err)
	}
}
