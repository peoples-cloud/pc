package rpcmsg

type Message struct {
	Msg string
}

type Swarm struct {
	Name     string
	Password string
}

type Program struct {
	Swarm string
	Hash  string
	Key   string
}

type Info struct {
	Swarm    string
	Password string
	Path     string
	Language string
	Hash     string
	Key      string // NOTE: remove key
}
