package proto

type Opcode uint8

const (
	OpError Opcode = iota
	OpStartup
	OpReady
	OpAuthenticate
	OpCredentials
	OpOptions
	OpSupported
	OpQuery
	OpResult
	OpPrepare
	OpExecute
	OpRegister
	OpEvent
	OpBatch
	OpAuthChallenge
	OpAuthResponse
	OpAuthSuccess
)

var opcodeNames = map[Opcode]string{
	OpError:         "ERROR",
	OpStartup:       "STARTUP",
	OpReady:         "READY",
	OpAuthenticate:  "AUTHENTICATE",
	OpCredentials:   "CREDENTIALS",
	OpOptions:       "OPTIONS",
	OpSupported:     "SUPPORTED",
	OpQuery:         "QUERY",
	OpResult:        "RESULT",
	OpPrepare:       "PREPARE",
	OpExecute:       "EXECUTE",
	OpRegister:      "REGISTER",
	OpEvent:         "EVENT",
	OpBatch:         "BATCH",
	OpAuthChallenge: "AUTH_CHALLENGE",
	OpAuthResponse:  "AUTH_RESPONSE",
	OpAuthSuccess:   "AUTH_SUCCESS",
}

func (oc Opcode) String() string {
	name, found := opcodeNames[oc]
	if !found {
		return "UNKNOWN"
	}
	return name
}
