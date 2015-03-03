package utils

import (
    "gosync/config"
    "strings"
    "os"
    "io"
    //"log"
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

// exists returns whether the given file or directory exists or not
func ItemExists(path string) (bool, error) {
    _, err := os.Stat(path)
    if err == nil { return true, nil }
    if os.IsNotExist(err) { return false, nil }
    return false, err
}

func FileWrite(path string, content io.Reader, overwrite bool, uid, gid int, perms string) (int64,error ){
    if _, err := os.Stat(path); err == nil {
        // We wipe the file as we need to replace with a new one
        err := os.Remove(path)
        if err != nil {
            return 0, err
        }

    }

    file, err := os.Create(path)

    if err != nil {
        return 0, err
    }
    defer file.Close()

    size, err := io.Copy(file, content)

    file.Chown(uid, gid)
    //file.Chmod(perms)

    return size, err
}