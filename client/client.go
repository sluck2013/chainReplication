package main

import (
    "common/config"
    "common/reply"
    "os"
    "fmt"
    "net"
    "strings"
    "time"
    "common/logger"
)

const MSGLEN = 1024 // max msg length
var configure config.Configure //configuration structure
var myLogger logger.Logger //log

/*
 * sendRequest send a request message reqMsg to head server
 * of bank with bankId. listenPort of the listening port of
 * the client, clientId is clientID of the client sending
 * the message
 */
func sendRequest(listenPort, reqMsg, bankId, clientId string) {
    destId := configure.GetHeadIdByBankId(bankId)
    strIP := configure.GetAddrByServerId(destId)
    destAddr := net.UDPAddr {
        Port : configure.GetPortByServerId(destId),
        IP: net.ParseIP(strIP),
    }
    conn, err := net.DialUDP("udp", nil, &destAddr)

    if err != nil {
        myLogger.Error.Println("UDP Connection failed. ClientID:", clientId, "Error Msg:", err)
        return
    }

    myLogger.Trace.Println("UDP Connection established. ClientID:", clientId)
    defer conn.Close()

    n, err := conn.Write([]byte(reqMsg + "|" + listenPort))
    fmt.Println("Request:",reqMsg)
    if err != nil {
        myLogger.Error.Println("Writing to UDP failed! ClientID:", clientId, "Error Msg:", err)
    }
    n = 1
    myLogger.Send.Println("#", myLogger.SendId, ": " +reqMsg + "|" + listenPort, "ClientID:", clientId)
    myLogger.SendId++
}


/*
 * startClient starts a client thread. ch is a channel used to
 * send request result to parent thread, clientId is the ClientID
 * of the client being started.
 */
func startClient(ch chan reply.Reply, closed chan bool, clientId string) {
    myLogger.Info.Println("Client thread started up. ClientID:", clientId)
    la := net.UDPAddr {
        Port : 0,
        IP : net.ParseIP("127.0.0.1"),
    }
    conn, err := net.ListenUDP("udp", &la)
    
    if err != nil {
        myLogger.Error.Println("UDP Server did not start. ClientID:", clientId, "Error msg:", err)
    }

    defer conn.Close()
    a := strings.Split(conn.LocalAddr().String(), ":")
    listenPort := a[len(a) - 1]

    reqIdx := 0
    reqNum := len(configure.Requests[clientId])

    for {
        if reqIdx < reqNum {
            req := configure.Requests[clientId][reqIdx]
            reqId := strings.Split(req, "|")[1]
            bankId := strings.Split(reqId, ".")[0]
            sendRequest(listenPort, req, bankId, clientId)
            reqIdx++
        } else {
            closed <- true
            conn.Close()
            myLogger.Info.Println("Connection closed! ClientID:", clientId)
            return
        }
        buf := make([]byte, MSGLEN)
        n, _, err := conn.ReadFromUDP(buf)
        re := reply.Reply{}

        if err != nil {
            myLogger.Error.Println("UDP read error. ClientID:", clientId, "Error msg:", err)
        }

        received := string(buf[:n])
        myLogger.Recv.Println("#", myLogger.RecvId, ":",  received, "ClientID:", clientId)
        myLogger.RecvId++
        re.Unserialize(received)
        ch <- re
        time.Sleep(500 * time.Millisecond)
    }
}

/* deprecated
func isAllTrue(mp *map[string] bool) bool {
    for _, v := range *mp {
        if !v {
            return false
        }
    }
    return true
}
*/

func main() {
    myLogger.Init("client.log")
    myLogger.Info.Println("Client process started up.")
    
    args := os.Args
    if len(args) != 2 {
        fmt.Println("ERROR: invalid argument")
        fmt.Println("Usage: ./server <configFile>")
        myLogger.Error.Println("invalid argument")
        return
    }
    
    configure.LoadConfig("config/" + os.Args[1])
    channel := make(chan reply.Reply)
    closed := make(chan bool)
    
    for k := range configure.Requests {
        go startClient(channel, closed, k)
    }

    closedClientNum := 0
    for {
        select {
            case r := <- channel:
                fmt.Println("")
                fmt.Println("RequestID:", r.RequestId)
                fmt.Println("Account#:", r.AccountNum)
                fmt.Println("Outcome:", r.Outcome.String())
                fmt.Println("Balance:", r.Balance)
            case c := <- closed:
                if c {
                    closedClientNum++
                }
                if closedClientNum == len(configure.Requests) {
                    return
                }
        }
    }
}
