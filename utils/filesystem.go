package utils

import (
    "gosync/config"
    "strings"
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