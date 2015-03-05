package utils

import (
    "strings"
    "os"
    "io"
    "io/ioutil"
    "crypto/md5"
    "fmt"
    "math"
    "path/filepath"
)

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
        //log.Fatalln(msg, err)
    }
}


func GetBasePath(listener string) string {
    cfg := GetConfig()
    return cfg.Listeners[listener].BasePath
}

func GetBaseDir(listener string) string {
    cfg := GetConfig()
    return cfg.Listeners[listener].Directory
}

func GetRelativeBasePath(listener, local_path string) string {
    lPath := strings.TrimPrefix(local_path, GetBaseDir(listener))
    return GetBasePath(listener) + lPath
}

func GetRelativePath(listener, local_path string) string {
    return strings.TrimPrefix(local_path, GetBaseDir(listener))
}

func GetAbsPath(listener, db_path string) string {
    absPath := GetBaseDir(listener) + db_path
    //logs.WriteF("=====> Absolute Path: %s <=====", absPath)
    return absPath
}

func GetListenerFromDir(dir string) string {
    cfg := GetConfig()
    var listener = ""
    for lname, ldata := range cfg.Listeners{
        if ldata.Directory == dir{
            listener = lname
        }
    }
    return listener
}

// exists returns whether the given file or directory exists or not
func ItemExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func FileWrite(path string, r io.Reader, uid, gid int, perms string) (int64, error){

    w, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0664)
    if err != nil {
        if path == "" {
            w = os.Stdout
        } else {
            return 0, err
        }
    }
    defer w.Close()

    size, err := io.Copy(w, r)


    if err != nil {
        return 0, err
    }

    return size, err
}

func FileWriteBytes(path string, content []byte, overwrite bool, uid, gid int, perms string) (int64,error ){
    /*buf := new(bytes.Buffer)
    buf.ReadFrom(r)
    content := buf.Bytes()*/
    err := ioutil.WriteFile(path, content, 0644)
    if err != nil {
        return 0, err
    }

    return 0, err
}