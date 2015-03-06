package utils

import (
	"fmt"
	"os"
    "log"
)

func WriteLn(msg string) {
    cfg := GetConfig()
    if cfg.ServerConfig.LogLocation != "stdout" {
        f, err := os.OpenFile(cfg.ServerConfig.LogLocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
        if err != nil {
            log.Fatalf("[FATAL] Error opening log file (%s)", err.Error())
        }
        defer f.Close()
        log.SetOutput(f)
    }
    log.Printf("%s", msg)
}

func LogWriteF(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a)
    cfg := GetConfig()
	if cfg.ServerConfig.LogLocation != "stdout" {
		f, err := os.OpenFile(cfg.ServerConfig.LogLocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
		if err != nil {
			log.Fatalf("[FATAL] Error opening log file (%s)", err.Error())
		}
		defer f.Close()
		log.SetOutput(f)
	}
	log.Printf("%s",msg)
}

func ErrorLn(level int, msg string){
    cfg := GetConfig()
    if cfg.ServerConfig.LogLocation != "stdout" {
        f, err := os.OpenFile(cfg.ServerConfig.LogLocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
        if err != nil {
            log.Fatalf("[FATAL] Error opening log file (%s)", err.Error())
        }
        defer f.Close()
        log.SetOutput(f)
    }
    if level >= 4{
        log.Fatalf("%s%s", getLevel(level), msg)
    }else{
        log.Printf("%s%s", getLevel(level), msg)
    }
}
/*func ErrorF(level int, format string, a ...interface{}){
    msg := fmt.Sprintf(format, a)
    cfg := config.GetConfig()
    if cfg.ServerConfig.LogLocation != "stdout" {
        f, err := os.OpenFile(cfg.ServerConfig.LogLocation, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0664)
        if err != nil {
            log.Fatalf("[FATAL] Error opening log file (%s)", err.Error())
        }
        defer f.Close()
        log.SetOutput(f)
    }
    if level >= 4{
        log.Fatalf("%s%s", getLevel(level), msg)
    }else{
        log.Printf("%s%s", getLevel(level), msg)
    }
}*/

func getLevel(level int) string {
	switch level {
	case 4:
		return fmt.Sprintf("%s",  "[FATAL] ")
	case 3:
		return fmt.Sprintf("%s",  "[Error] ")
	case 2:
		return fmt.Sprintf("%s",  "[Warning] ")
	case 1:
		return "" // We return no formatting as this is the default
	case 0:
		return fmt.Sprintf("%s",  "[Debug] ")
	}
	return fmt.Sprintf("%s",  "[FATAL] ")
}
