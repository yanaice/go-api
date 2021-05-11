package derror

import (
	"errors"
	"fmt"
	"net/http"

	"go-starter-project/pkg/log"
)

type Derror struct {
	Err            error                  `json:"error_raw,omitempty"`
	ErrStr         string                 `json:"error,omitempty"`
	ErrCode        DerrorCode             `json:"error_code"`
	SourceFunc     string                 `json:"source_func"`
	SourceLine     int                    `json:"source_line"`
	Debug          string                 `json:"debug_message,omitempty"`
	DebugExtraInfo map[string]interface{} `json:"debug_extra_info,omitempty"`
	ErrLevel       DerrorLevel            `json:"error_level"`
	HTTPCode       int                    `json:"-"`
}

func (e Derror) Error() string {
	return e.Err.Error()
}

func (e *Derror) SetDebug(debug string) *Derror {
	e.Debug = debug
	return e
}

func (e *Derror) SetExtraInfo(key string, value interface{}) *Derror {
	e.DebugExtraInfo[key] = value
	return e
}

func (e *Derror) SetCode(code DerrorCode) *Derror {
	e.ErrCode = code
	switch code {
	case ErrCodeUnauthorized:
		return e.SetHTTPCode(http.StatusUnauthorized)
	case ErrCodeForbidden:
		return e.SetHTTPCode(http.StatusForbidden)
	case ErrCodeInputError:
		return e.SetHTTPCode(http.StatusBadRequest)
	case ErrCodeNotFound:
		return e.SetHTTPCode(http.StatusNotFound)
	}
	return e
}

func (e *Derror) SetHTTPCode(code int) *Derror {
	e.HTTPCode = code
	return e
}

func (e Derror) Log() Derror {
	message := e.Error()
	message += ": "
	message += e.Debug
	for k, v := range e.DebugExtraInfo {
		message += fmt.Sprintf(" [%s = %+v]", k, v)
	}
	switch e.ErrLevel {
	case ErrLevelWarn:
		log.Warn(message)
	case ErrLevelError:
		log.Error(message)
	case ErrLevelFatal:
		log.Fatal(message)
	}
	return e
}

func (e Derror) Strip() DerrorReply {
	return DerrorReply{
		ErrCode:  e.ErrCode,
		ErrLevel: e.ErrLevel,
	}
}

func ErrorDebug(err error, debug string) Derror {
	return E(err).SetDebug(debug).Log()
}

func Error(err error) Derror {
	return E(err).Log()
}

func buildError(err error, level DerrorLevel) *Derror {
	file, line := log.GetCaller()
	errStr := ""
	if err != nil {
		errStr = err.Error()
	} else {
		err = errors.New("")
	}
	derror := &Derror{
		Err:            err,
		ErrStr:         errStr,
		SourceFunc:     file,
		SourceLine:     line,
		ErrCode:        ErrCodeServerError,
		ErrLevel:       level,
		DebugExtraInfo: make(map[string]interface{}),
	}
	return derror
}

func E(err error) *Derror {
	return buildError(err, ErrLevelError)
}

func W(err error) *Derror {
	return buildError(err, ErrLevelWarn)
}

func F(err error) *Derror {
	return buildError(err, ErrLevelFatal)
}

func Ecode(code DerrorCode) *Derror {
	return E(nil).SetCode(code)
}

func Wcode(code DerrorCode) *Derror {
	return W(nil).SetCode(code)
}

func Fcode(code DerrorCode) *Derror {
	return F(nil).SetCode(code)
}

func Wrap(derr Derror) *Derror {
	return buildError(derr, derr.ErrLevel).SetCode(derr.ErrCode)
}

func WrapReply(derr DerrorReply) *Derror {
	return buildError(derr, derr.ErrLevel).SetCode(derr.ErrCode)
}
