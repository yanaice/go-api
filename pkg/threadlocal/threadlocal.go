package threadlocal

import "github.com/jtolds/gls"

var ctx = gls.NewContextManager()

var correlationID = gls.GenSym()

func SetCorrelationID(id string, cb func()) {
	ctx.SetValues(gls.Values{
		correlationID: id,
	}, cb)
}

func GetCorrelationID() string {
	value, ok := ctx.GetValue(correlationID)
	if ok {
		return value.(string)
	} else {
		return ""
	}
}

func Go(cb func()) {
	gls.Go(cb)
}
