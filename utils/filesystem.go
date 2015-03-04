package utils

import (
    "gosync/config"
    "strings"
    "os"
    "io"
    "io/ioutil"
    //"log"
    //"bytes"
)

func GetBasePath(listener string) string {
    cfg := config.GetConfig()
    return cfg.Listeners[listener].BasePath
}

func GetBaseDir(listener string) string {
    cfg := config.GetConfig()
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
    //log.Printf("=====> Absolute Path: %s <=====", absPath)
    return absPath
}

func GetListenerFromDir(dir string) string {
    cfg := config.GetConfig()
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

    w, err := os.Create(path)
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