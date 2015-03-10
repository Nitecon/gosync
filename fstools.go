package gosync

import (
    "strings"
)

// Fetches the home directory of the listener from the config file
func GetListenerHomeDir(listener string) string{
    cfg := GetConfig()
    return cfg.Listeners[listener].Directory
}

// Fetches the base directory (upload path) from configuration file
func GetListenerUploadPath(listener string) string {
    cfg := GetConfig()
    return cfg.Listeners[listener].BasePath
}

// Gets the relative path based on the home directory, essentially stripping the listener,
// home directory from the absolute path
func GetRelativePath(listener, local_path string) string{
    return strings.TrimPrefix(local_path, GetListenerHomeDir(listener))
}

// Gets the relative path for the item by trimming the absolute path with the home dir,
// it then appends the base directory (Upload Path) to the front of the relative path.
func GetRelativeUploadPath(listener, local_path string) string{
    return GetListenerUploadPath(listener) + GetRelativePath(listener, local_path)
}