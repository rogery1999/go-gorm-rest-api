package types

import "encoding/json"

type CustomError struct {
	Status uint
	Body   interface{}
}

func (ce *CustomError) Error() string {
	sError, err := json.Marshal(ce.Body)

	if err != nil {
		return err.Error()
	}

	return string(sError)
}
