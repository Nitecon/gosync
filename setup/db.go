package setup
import (
    log "github.com/cihub/seelog"
    "github.com/Nitecon/gosync/config"
    "github.com/Nitecon/gosync/database"
    "github.com/Nitecon/gosync/filesystem"
)

func DbVerify(){
    cfg := config.GetConfig()
    database.SetupDB()
    for lname, listener := range cfg.Listeners{
        log.Infof("Validating db for %s, base directory: %s", lname, listener.Directory)
        ignFile, err := filesystem.GetFileInfo(listener.Directory+"/.goignore")
        if err != nil{
            log.Criticalf("Ignore file has disappeared: %s", listener.Directory+"/.goignore")
        }
        if database.Exists(lname, ignFile){
            log.Debugf("Ignore exists in DB")
        }else{
            log.Debugf("Ignore DOES NOT EXIST in DB... adding...")
            updated:=database.Add(lname,ignFile)
            log.Debugf("Ignore file added to db: %b", updated)
        }
    }

}