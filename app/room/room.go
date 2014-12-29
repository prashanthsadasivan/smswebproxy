package room


var (
    conduits map[string]*Conduit
)

func init() {
    conduits = make(map[string]*Conduit)
}

type Conduit struct {
    Received chan SMSMessage
    RegId string
}

type SMSMessage struct {
    Num string
    Message string
}

func GetConduit(number string) *Conduit {
    return conduits[number]
}

func AddConduit(num string, conduit *Conduit) {
    conduits[num] = conduit
}
