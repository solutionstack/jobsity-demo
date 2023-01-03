package utils

import (
	"encoding/json"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func UnmarshalRequestBody(r *http.Request, output interface{}) error {
	if r.Body == nil {
		return errors.New("invalid body in request")
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}

	err = r.Body.Close()
	if err != nil {
		return err
	}

	err = json.Unmarshal(body, &output)
	if err != nil {
		return err
	}

	return nil
}
