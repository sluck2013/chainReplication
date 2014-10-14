package reply

import (
    "fmt"
    "strings"
    "strconv"
)

type TransOutcome int

const (
    Processed TransOutcome = 1 + iota
    InconsistentWithHistory
    InsufficientFunds
)

var outcomes = [...]string {
    "Processed",
    "InconsistentWithHistory",
    "InsufficientFunds",
}

func (out TransOutcome) String() string {
    return outcomes[out - 1]
}

type Reply struct {
    RequestId string
    AccountNum string
    Outcome TransOutcome
    Balance float64
}

func (re Reply) Serialize() string {
    return re.RequestId + "|" + re.AccountNum + "|" + strconv.Itoa(int(re.Outcome)) + "|" + fmt.Sprintf("%.2f", re.Balance)
}

func (re *Reply) Unserialize(s string) {
    r := strings.Split(s, "|")
    re.RequestId = r[0]
    re.AccountNum = r[1]
    outcome, _ := strconv.Atoi(r[2])
    re.Outcome = TransOutcome(outcome)
    re.Balance, _ = strconv.ParseFloat(r[3], 64)
}
