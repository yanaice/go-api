package derror

import "encoding/json"

type DerrorReply struct {
	ErrCode  DerrorCode  `json:"error_code"`
	ErrLevel DerrorLevel `json:"error_level"`
}

func (e DerrorReply) Error() string {
	return string(e.ErrCode)
}

func UnmarshalDerrorReply(body []byte) (DerrorReply, error) {
	var derr DerrorReply
	err := json.Unmarshal(body, &derr)
	if err != nil {
		return DerrorReply{}, err
	}
	return derr, nil
}
