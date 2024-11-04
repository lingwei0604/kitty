package dto

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/schema"
	"github.com/pkg/errors"
)

type Decoder struct {
	decoder *schema.Decoder
}

func NewDecoder() *Decoder {
	decoder := schema.NewDecoder()
	decoder.IgnoreUnknownKeys(true)
	decoder.SetAliasTag("json")
	return &Decoder{decoder: decoder}
}

func (d *Decoder) Decode(payload *Payload, r *http.Request) error {
	if r.Method == "POST" {
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			return errors.Wrapf(err, "cannot read body of http request")
		}
		err = json.Unmarshal(buf, payload)
		if err != nil {
			return errors.Wrap(err, "cannot json unmarshal")
		}
		err = json.Unmarshal(buf, &payload.B)
		if err != nil {
			return errors.Wrap(err, "cannot json unmarshal")
		}
		return nil
	}
	query := r.URL.Query()
	err := d.decoder.Decode(payload, query)
	if err != nil {
		return errors.Wrap(err, "fails to decode")
	}
	// store extra queries here.
	payload.Q = query
	return nil
}
