package logger

import (
    "log"
    "os"
)

type Mode int

const(
    STDIO = 1 + iota
    FILE
)

type Logger struct {
    LogMode Mode
    Trace *log.Logger
    Error *log.Logger
    Info *log.Logger
    Send *log.Logger
    Recv *log.Logger
    SendId int
    RecvId int
}


func (myLogger *Logger) Init(filename string) {
    myLogger.SendId = 1
    myLogger.RecvId = 1
    file, err := os.OpenFile("logs/" + filename, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0666)

    if err != nil {
        log.Fatalln("Failed to open log file", filename, ":", err)
    }

    flag := log.Ldate | log.Ltime
    myLogger.Trace = log.New(file, "TRACE: ", flag)
    myLogger.Error = log.New(file, "ERROR: ", flag)
    myLogger.Info = log.New(file, "INFO: ", flag)
    myLogger.Send = log.New(file, "MSG SENT: ", flag)
    myLogger.Recv = log.New(file, "MSG RECV: ", flag)
}
