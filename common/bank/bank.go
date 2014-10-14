package bank

import (
    "common/account"
    "common/config"
    "common/reply"
)

/* obtain max int */
const MaxUint = ^uint(0)
const MaxInt = MaxUint >> 1
/* in reality, it should be max int */
const MaxAccounts = 500

type Bank struct{
    AccMap map[string]int
    Accounts []account.Account
    Id string
    Name string
    ChainLen int
    ClientNum int
    ServerID []string
}

/*
 * checkAccountId checks whether account with accountId
 * exists, if not, create it
 */
func (myBank *Bank) checkAccountId(accountId string) bool {
    _, exist := myBank.AccMap[accountId]
    if exist {
        return true
    } 

    // if account not found, create it
    myBank.AccMap[accountId] = len(myBank.Accounts)
    myBank.Accounts = append(myBank.Accounts, account.Account{
        Id : accountId,
        Balance : 0, 
        UpdateRecords : make(map[string]account.Record),
    })

    return true
}

/*
 * initialize myBank with with ConfBank
 */
func (myBank *Bank) SetBankInfo(bankInfo *config.ConfBank) {
    myBank.Id = bankInfo.Id
    myBank.Name = bankInfo.Name
    myBank.ChainLen = bankInfo.ChainLen
    myBank.ClientNum = bankInfo.ClientNum
    myBank.ServerID = bankInfo.ServerID
    myBank.AccMap = make(map[string]int)
    myBank.Accounts = make([]account.Account, 0, MaxAccounts)
}

func (myBank *Bank) GetBalance(reqId string, accNum string) reply.Reply {
    myBank.checkAccountId(accNum)
    idx := myBank.AccMap[accNum]
    return myBank.Accounts[idx].GetBalance(reqId)
}

func (myBank *Bank) Deposit(reqId string, accNum string, amount float64) reply.Reply {
    myBank.checkAccountId(accNum)
    idx := myBank.AccMap[accNum]
    return myBank.Accounts[idx].Deposit(reqId, amount)
}

func (myBank *Bank) Withdraw(reqId string, accNum string, amount float64) reply.Reply {
    myBank.checkAccountId(accNum)
    idx := myBank.AccMap[accNum]
    return myBank.Accounts[idx].Withdraw(reqId, amount)
}
