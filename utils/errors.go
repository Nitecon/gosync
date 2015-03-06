package utils

import (
    "runtime"
    "fmt"
)

func Check(err error, code int, msg string) bool{
    appErr := AppError{}
    if err != nil{
        appErr.Error = err
        appErr.Message = msg
        appErr.Code = code
        if code >= 500{
            var stack [4096]byte
            runtime.Stack(stack[:], false)
            appErr.Stack = string(stack[:])
            ErrorLn(4,getLogVerbose(appErr))
        }else{
            ErrorLn(3,getLogVerbose(appErr))
        }
        return true
    }
    return false
}

func CheckF(err error, code int, format string, a ...interface{}) bool{
    appErr := AppError{}
    if err != nil{
        appErr.Error = err
        appErr.Message = fmt.Sprintf(format, a)
        appErr.Code = code
        if code >= 500{
            var stack [4096]byte
            runtime.Stack(stack[:], false)
            appErr.Stack = string(stack[:])
            ErrorLn(4,getLogVerbose(appErr))
        }else{
            ErrorLn(3,getLogVerbose(appErr))
        }
        return true
    }
    return false
}

func getLogVerbose(appErr AppError) string{
    conf := GetConfig()
    var errString = ""
    if conf.ServerConfig.LogLevel == 0{
        errString = fmt.Sprintf("%s\n============\n%v\n==============\n%+v", appErr.Message,appErr.Error, appErr.Stack )
    }
    if conf.ServerConfig.LogLevel == 1{
        errString = fmt.Sprintf( "%s\n============\n%v", appErr.Message,appErr.Error )
    }
    if conf.ServerConfig.LogLevel > 1{
        errString = fmt.Sprintf("%s", appErr.Message )
    }
    return errString
}