package main

import (
    "common/config"
    "fmt"
    "os"
    "net"
    "net/http"
    "strconv"
//    "io/ioutil"
    "common/bank"
    "common/account"
    "strings"
    "net/url"
    "common/logger"
)

const MSGLEN = 1024
var configure config.Configure
var myBank bank.Bank
var serverId string
var myLogger logger.Logger

/* On receiving deposit request from predecessor */
func depositHandler(w http.ResponseWriter, req *http.Request) {
    reqArr := strings.Split(req.RequestURI, "/")
    l := len(reqArr)
    clientAddr := reqArr[l - 1]
    s, _ := url.QueryUnescape(reqArr[l - 2])
    accArgs := strings.Split(s, "|")
    amount, _ := strconv.ParseFloat(accArgs[2], 64)
    r := myBank.Deposit(accArgs[0], accArgs[1], amount)
    myLogger.Trace.Println(amount, "was depsit on account", accArgs[1] + ". ReqID:", accArgs[0])

    rEncoded := []byte(r.Serialize())
    if configure.IsTail(serverId) {
        //response to Client
        sendResultToClient(rEncoded, clientAddr)
    } else {
        //send request to successor
        succId := configure.GetSuccessor(serverId)
        arg := accArgs[0] + "|" + accArgs[1] + "|" + accArgs[2]
        go sendHTTPRequest(succId, account.OpDeposit, arg, clientAddr)
    }
}

/* On receiving withdraw request from predecessor */
func withdrawHandler(w http.ResponseWriter, req *http.Request) {
    reqArr := strings.Split(req.RequestURI, "/")
    l := len(reqArr)
    clientAddr := reqArr[l - 1]
    s, _ := url.QueryUnescape(reqArr[l - 2])
    accArgs := strings.Split(s, "|")
    amount, _ := strconv.ParseFloat(accArgs[2], 64)
    r := myBank.Withdraw(accArgs[0], accArgs[1], amount)
    myLogger.Trace.Println("Trying to withdraw", amount, "from account", accArgs[1] + ". ReqID:", accArgs[0])

    rEncoded := []byte(r.Serialize())
    if configure.IsTail(serverId) {
        //response to Client
        sendResultToClient(rEncoded, clientAddr)
    } else {
        //send request to successor
        succId := configure.GetSuccessor(serverId)
        arg := accArgs[0] + "|" + accArgs[1] + "|" + accArgs[2]
        go sendHTTPRequest(succId, account.OpWithdraw, arg, clientAddr)
    }
}

/* sendResultToClient sends request result to client */
func sendResultToClient(encodedReply []byte, cliAddr string) {
    clientAddr, err := net.ResolveUDPAddr("udp", cliAddr)
    conn, err := net.DialUDP("udp", nil, clientAddr)

    if err != nil  {
        myLogger.Error.Println("Connection to client failed.", err)
        return
    }

    defer conn.Close()

    n, err := conn.Write(encodedReply)
    if err != nil {
        myLogger.Error.Println("Responding to client failed.", err)
        return
    }
    myLogger.Send.Println("#", myLogger.SendId, ":", string(encodedReply[:n]))
    myLogger.SendId++
}

/* start http server */
func startHTTPService(serverId string) {
    http.HandleFunc("/deposit/", depositHandler)
    http.HandleFunc("/withdraw/", withdrawHandler)
    port := configure.GetPortByServerId(serverId)
    err := http.ListenAndServe(":" + strconv.Itoa(port), nil)
    if err != nil {
        myLogger.Error.Println("Http Server did not start!", err)
    }
}

/* start udp server */
func startUDPService(serverId string) {
    localAddr := net.UDPAddr {
        Port : configure.GetPortByServerId(serverId),
        IP : net.ParseIP("127.0.0.1"),
    }

    conn, err := net.ListenUDP("udp", &localAddr)
    defer conn.Close()

    if err != nil {
        myLogger.Error.Println("UDP server did not start!", err)
        os.Exit(1)
    } else {
        myLogger.Trace.Println("UDP server launched!")
    }

    for {
        buf := make([]byte, MSGLEN)
        n, sourceAddr, err := conn.ReadFromUDP(buf)

        if err != nil {
            myLogger.Error.Println("UDP read error.", err)
            continue
        }
        recv := string(buf[:n])
        myLogger.Recv.Println("#", myLogger.RecvId, ":", recv)
        myLogger.RecvId++
        
        reqArray := strings.Split(string(buf[:n]), "|")
        l := len(reqArray)
        // reqArray[0] request function name
        // reqArray[1:l - 2] request arguments
        // reqArray[l - 1] client listening port
        clientAddr := *sourceAddr
        clientAddr.Port, _ = strconv.Atoi(reqArray[l - 1])

        switch reqArray[0] {
        case "getBalance":
            //reqArray[1] reqId, reqArray[2] accountNum
            r := myBank.GetBalance(reqArray[1], reqArray[2])
            myLogger.Trace.Println("Retrieving balance for account", reqArray[2]+". ReqID:", reqArray[1])
            rEncoded := []byte(r.Serialize())
            sendResultToClient(rEncoded, clientAddr.String())
            
        case "deposit":
            //reqArray[1] reqId, reqArray[2] accountNum
            //reqArray[3] amount
            amount, _ := strconv.ParseFloat(reqArray[3], 64)
            r := myBank.Deposit(reqArray[1], reqArray[2], amount)
            myLogger.Trace.Println(amount, "was depsit on account", reqArray[2] + ". ReqID:", reqArray[1])
            rEncoded := []byte(r.Serialize())
            if configure.IsTail(serverId) {
                sendResultToClient(rEncoded, clientAddr.String())
            } else {
                // send request to successor
                succId := configure.GetSuccessor(serverId)
                arg := reqArray[1] + "|" + reqArray[2] + "|" + reqArray[3]
                go sendHTTPRequest(succId, account.OpDeposit, arg, clientAddr.String())
            }
        case "withdraw":
            //reqArray[1] reqId, reqArray[2] accountNum
            //reqArray[3] amount
            amount, _ := strconv.ParseFloat(reqArray[3], 64)
            r := myBank.Withdraw(reqArray[1], reqArray[2], amount)
            myLogger.Trace.Println("Trying to withdraw", amount, "from account", reqArray[2] + ". ReqID:", reqArray[1])
            rEncoded := []byte(r.Serialize())
            if configure.IsTail(serverId) {
                sendResultToClient(rEncoded, clientAddr.String())
            } else {
                succId := configure.GetSuccessor(serverId)
                arg := reqArray[1] + "|" + reqArray[2] + "|" + reqArray[3]
                go sendHTTPRequest(succId, account.OpWithdraw, arg, clientAddr.String())
            }
        }
    }
}

/* send http request to destServer */
func sendHTTPRequest(destServerId string, operation account.OpType, argStr string, clientAddr string) {
    addr :=  configure.GetAddrByServerId(destServerId)
    port := configure.GetPortByServerId(destServerId)
    url := "http://" + addr + ":" + strconv.Itoa(port)
    url += "/" + operation.String() + "/" + argStr
    url += "/" + clientAddr
    //resp, err := http.Get(url)
    http.Get(url)
    myLogger.Send.Println("#", myLogger.SendId, ":", url)
    myLogger.SendId++

    /*if err != nil {
        fmt.Printf("ERROR: %v", err)
    }
    
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)

    if err != nil {
        fmt.Println("ERROR: ", err)
    }

    fmt.Println(body)
    */
}

func main() {
    args := os.Args
    if len(args) != 3 {
        fmt.Println("ERROR: invalid argument")
        fmt.Println("Usage: ./server <configFile> <serverID>")
        return
    }

    configure.LoadConfig("config/" + os.Args[1])
    serverId = os.Args[2]
    myLogger.Init("server" + serverId + ".log")

    bankId := configure.GetBankIdByServerId(serverId)
    bankInfo := configure.GetBankInfoById(bankId)
    myBank.SetBankInfo(&bankInfo)
    myLogger.Info.Println("Server", serverId, "started up.")
    loginfo := "IP: " + configure.Servers[serverId].IP + ", "
    loginfo += "Port: " + strconv.Itoa(configure.Servers[serverId].Port) + ", "
    loginfo += "BankId: " + bankId + ", "
    loginfo += "Startup Delay: " + strconv.Itoa(configure.Servers[serverId].StartDelay) + ", "
    loginfo += "Lifetime: " + configure.Servers[serverId].Lifetime
    myLogger.Info.Println(loginfo)

    // start service
    if configure.IsHead(serverId) {
        startUDPService(serverId)
    } else {
        if configure.IsTail(serverId) {
            go startUDPService(serverId)
            startHTTPService(serverId)
        } else {
            startHTTPService(serverId)
        }
    }
}
