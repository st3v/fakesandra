package v3

import "github.com/st3v/fakesandra/cql"

const (
	opError opcode = iota
	opStartup
	opReady
	opAuthenticate
	_ // DEPRECATED
	opOptions
	opSupported
	opQuery
	opResult
	opPrepare
	opExecute
	opRegister
	opEvent
	opBatch
	opAuthChallenge
	opAuthResponse
	opAuthSuccess
)

type router struct {
	handlers map[opcode]cql.Handler
}
