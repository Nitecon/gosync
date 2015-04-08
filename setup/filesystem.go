package setup

import (
    log "github.com/cihub/seelog"
    "github.com/Nitecon/gosync/config"
    "github.com/Nitecon/gosync/filesystem"
    "os"
)

func FsVerify() {
    cfg := config.GetConfig()
    log.Info("Verifying Filesystem directories...")
    for lname, listener := range cfg.Listeners {
        log.Infof("Current Listener: %s, Directory: %s", lname, listener.Directory)
        isdir, err := filesystem.IsDir(listener.Directory)
        if err != nil{
            log.Infof("Cannot stat directory %s, %s", listener.Directory, err.Error())
        }
        if isdir{
            // Now we check to see if the goignore file exists
            verifySyncIgnore(listener.Directory+"/.goignore")
        }else{
            err = os.Mkdir(listener.Directory,0775)
            if err != nil{
                log.Errorf("Cannot create non existent sync dir %s\n%s", listener.Directory, err.Error())
                os.Exit(1)
            }
            verifySyncIgnore(listener.Directory+"/.goignore")
        }
    }
    log.Info("Filesystem verify completed...")
}

func verifySyncIgnore(path string){
    _, err := filesystem.FileExists(path)
    if err != nil{
        log.Debugf("Error occurred validating %s\n%s", path, err.Error())
        var defIgnore = ".git/*\n*.pid\n*.log"
        log.Debugf("Creating default ignore file: %s", path)
        err := filesystem.FileWriteString(path, defIgnore)
        if err != nil{
            log.Infof("Cannot create a default .goignore file:\n%s\n%s",path, err.Error())
        }
    }else{
        log.Debugf("Using existing gosync ignore file: %s", path)
    }

}