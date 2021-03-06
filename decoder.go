package ali_mns

import (
	"encoding/xml"
	"io"
	"encoding/json"
)

type MNSDecoder interface {
	Decode(reader io.Reader, v interface{}) (err error)
}

type AliMNSDecoder struct {
}

func NewAliMNSDecoder() MNSDecoder {
	return &AliMNSDecoder{}
}

func (p *AliMNSDecoder) Decode(reader io.Reader, v interface{}) (err error) {
	decoder := xml.NewDecoder(reader)
	err = decoder.Decode(&v)

	return
}

func (p *AliMNSDecoder) DecodeJson(reader io.Reader, v interface{}) (err error) {
	decoder := json.NewDecoder(reader)
	err = decoder.Decode(&v)

	return
}