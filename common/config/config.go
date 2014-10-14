package config

import (
    "encoding/json"
    "os"
    "fmt"
)

type ConfServer struct {
    Id string
    IP string
    Port int
    BankId string
    StartDelay int
    Lifetime string
}

type ConfBank struct {
    Id string
    Name string
    ChainLen int
    ClientNum int
    ServerID []string
}

type ConfRequest struct {
    ClientID string
    Request []string
}

type ConfigureRead struct {
    Servers []ConfServer
    Banks []ConfBank
    Requests []ConfRequest
}

type Configure struct {
    Servers map[string]ConfServer
    Banks map[string]ConfBank
    Requests map[string][]string
}


func (config *Configure) LoadConfig(filename string) {
    /* read configuration */
    file, _ := os.Open(filename)
    decoder := json.NewDecoder(file)
    confRead := ConfigureRead{}
    err := decoder.Decode(&confRead)
    if err != nil {
        fmt.Println("Error while reading config file: ", err)
    }

    /* package it into map */
    config.Servers = make(map[string]ConfServer)
    for i := range confRead.Servers {
        config.Servers[confRead.Servers[i].Id] = confRead.Servers[i]
    }

    config.Banks = make(map[string]ConfBank)
    for i := range confRead.Banks {
        config.Banks[confRead.Banks[i].Id] = confRead.Banks[i]
    }

    config.Requests = make(map[string][]string)
    for i := range confRead.Requests {
        config.Requests[confRead.Requests[i].ClientID] = confRead.Requests[i].Request
    }
}

func (config *Configure) GetAddrByServerId(serverId string) string {
    return config.Servers[serverId].IP
}
func (config *Configure) GetPortByServerId(serverId string) int {
    return config.Servers[serverId].Port
}

func (config *Configure) GetBankIdByServerId(serverId string) string {
    return config.Servers[serverId].BankId
}

func (config *Configure) GetHeadIdByBankId(bankId string) string {
    return config.Banks[bankId].ServerID[0]
}
func (config *Configure) IsTail(serverId string) bool {
    bankId := config.GetBankIdByServerId(serverId)
    serverList := config.Banks[bankId].ServerID
    return serverList[len(serverList) - 1] == serverId
}

func (config *Configure) IsHead(serverId string) bool {
    bankId := config.GetBankIdByServerId(serverId)
    serverList := config.Banks[bankId].ServerID
    return serverList[0] == serverId
}

func (config *Configure) GetSuccessor(serverId string) string {
    if config.IsTail(serverId) {
        return ""
    }
    bankId := config.GetBankIdByServerId(serverId)
    serverList := config.Banks[bankId].ServerID
    for i := range serverList {
        if serverList[i] == serverId {
            return serverList[i + 1]
        }
    }
    return ""
}

func (config *Configure) GetBankInfoById(bankId string) ConfBank {
    return config.Banks[bankId]
}
