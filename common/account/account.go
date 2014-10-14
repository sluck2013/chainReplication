package account

import (
    "common/reply"
    "time"
)

/* Operation types */
type OpType int
const (
    OpDeposit OpType = 1 + iota
    OpWithdraw
)

var optypes = [...]string {
    "deposit",
    "withdraw",
}

func (op OpType) String() string {
    return optypes[op - 1]
}

/* request types */
type reqType int
const (
    reqProcessed reqType = 1 + iota
    reqInconsistent
    reqNew
)

/* transaction record */
type Record struct {
    operation OpType
    amount float64
    balance float64
    updateTime time.Time
}

/* account */
type Account struct {
    Id string
    Balance float64
    UpdateRecords map[string]Record
}

/*
 * checkReqId checks if a request is a new request, or 
 * processed one, or inconsistent with history
 */
func (acc Account) checkReqId(reqId string, operation OpType, amount float64) (Record, reqType) {
    record, existed := acc.UpdateRecords[reqId]
    if existed {
        if record.operation == operation && record.amount == amount {
            //same request
            return record, reqProcessed
        } else {
            //inconsistent
            return record, reqInconsistent
        }
    } else {
        //new
        return record, reqNew
    }
}

func (acc Account) GetBalance(reqId string) reply.Reply {
    return reply.Reply {
        RequestId : reqId,
        AccountNum : acc.Id,
        Outcome : reply.Processed, //todo
        Balance : acc.Balance,
    }
}

func (acc *Account) Deposit(reqId string, amount float64) reply.Reply {
    record, rtype := acc.checkReqId(reqId, OpDeposit, amount)
    var ret reply.Reply
    switch rtype {
    case reqProcessed:
        ret = reply.Reply {
            RequestId : reqId,
            AccountNum : acc.Id,
            Outcome : reply.Processed,
            Balance : record.balance,
        }
    case reqInconsistent:
        ret = reply.Reply {
            RequestId : reqId,
            AccountNum : acc.Id,
            Outcome : reply.InconsistentWithHistory,
            Balance : record.balance,
        }
    case reqNew:
        acc.Balance += amount
        acc.UpdateRecords[reqId] = Record {
            operation : OpDeposit,
            amount : amount,
            balance : acc.Balance,
            updateTime : time.Now(),
        }
        ret = reply.Reply {
            RequestId : reqId,
            AccountNum : acc.Id,
            Outcome : reply.Processed,
            Balance : acc.Balance,
        }
    }
    return ret
}

func (acc *Account) Withdraw(reqId string, amount float64) reply.Reply {
    record, rtype := acc.checkReqId(reqId, OpWithdraw, amount)
    var ret reply.Reply

    switch rtype {
    case reqProcessed:
        ret = reply.Reply {
            RequestId : reqId,
            AccountNum : acc.Id,
            Outcome : reply.Processed,
            Balance : record.balance,
        }
    case reqInconsistent:
        ret = reply.Reply {
            RequestId : reqId,
            AccountNum : acc.Id,
            Outcome : reply.InconsistentWithHistory,
            Balance : record.balance,
        }
    case reqNew:
        if acc.Balance >= amount {
            acc.Balance -= amount
            acc.UpdateRecords[reqId] = Record {
                operation : OpWithdraw,
                amount : amount,
                balance : acc.Balance,
                updateTime : time.Now(),
            }
            ret = reply.Reply {
                RequestId : reqId,
                AccountNum : acc.Id,
                Outcome : reply.Processed,
                Balance : acc.Balance,
            }
        } else {
            ret = reply.Reply {
                RequestId : reqId,
                AccountNum : acc.Id,
                Outcome : reply.InsufficientFunds,
                Balance : record.balance,
            }
        }
    }
    return ret
}
