package utils

import (
    "gosync/config"
    "strings"
    "os"
    "io"
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
    return GetBaseDir(listener) + db_path
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