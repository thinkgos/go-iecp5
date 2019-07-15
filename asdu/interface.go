package asdu

type Connect interface {
	Params() *Params
	Send(a *ASDU) error
}
