package filesystem

import (
    "os"
    "io/ioutil"
    "fmt"
    "path"
    "crypto/md5"
    "io"
    "math"
    "net"
    "strings"
    log "github.com/cihub/seelog"
)

func FileExists(path string) (exists bool, err error){
    fi, err := os.Stat(path)
    if err != nil{
        return
    }
    if fi.IsDir(){
        exists = false
        return
    }else{
        exists = true
        return
    }
    return
}

func IsDir(path string) (isdir bool, err error) {
    fi, err := os.Stat(path)
    if err != nil {
        return
    }

    // check if the source is indeed a directory or not
    if fi.IsDir() {
        isdir = true
        return
    }else{
        isdir = false
    }
    return
}

func FileWriteString(path, s string) (err error){
    b := []byte(s)
    err = ioutil.WriteFile(path, b, 0644)
    return
}

func GetFileInfo(doc_path string) (fi FileData, err error) {
    fitem, err := os.Stat(doc_path)
    if err != nil {
        return
    } else {
        hostname, _ := os.Hostname()
        fi.Filename = path.Base(doc_path)
        fi.Directory = path.Dir(doc_path)
        fi.IsDir = fitem.IsDir()
        fi.Checksum = GetMd5Checksum(doc_path)
        fi.Mtime = fitem.ModTime().UTC()
        fi.HostUpdated = hostname
        fi.HostIPs = GetLocalIp()
        fi.Perms = fmt.Sprintf("%#o", fitem.Mode().Perm())
    }
    return fi, nil
}

func GetMd5Checksum(filepath string) string {
    if b, err := computeMd5(filepath); err == nil {
        md5string := fmt.Sprintf("%x", b)
        return md5string
    } else {
        return "DirectoryMD5"
    }
}

func computeMd5(filepath string) (string, error) {
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
        buf := make([]byte, blocksize)

        file.Read(buf)
        io.WriteString(hash, string(buf)) // append into the hash
    }

    //fmt.Printf("%s checksum is %x\n",file.Name(), hash.Sum(nil))
    return fmt.Sprintf("%x", hash.Sum(nil)), nil
}

func GetLocalIp() string {
    addrs, err := net.InterfaceAddrs()
    if err != nil {
        log.Criticalf("No network interfaces found, %s", err.Error())
    }
    //log.Infof("Addresses: %+v", addrs)
    var ips []string
    for _, a := range addrs {
        if ipnet, ok := a.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
            if ipnet.IP.To4() != nil {
                ips = append(ips, ipnet.IP.String())
            }
        }
    }
    return strings.Join(ips, ",")
}
