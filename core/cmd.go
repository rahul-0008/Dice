package core

type Rediscmd struct {
	Cmd  string
	Args []string
}

type RedisCmds []*Rediscmd
