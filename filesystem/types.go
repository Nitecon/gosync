package filesystem

import (
    "time"
)

type FileData struct {
    Filename    string
    Directory   string
    IsDir       bool
    Checksum    string
    Mtime       time.Time
    Perms       string
    HostUpdated string
    HostIPs     string
}