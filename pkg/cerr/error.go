package cerr

import "encoding/json"

type Bag struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

func (b Bag) Error() string {
	data, err := json.Marshal(b)
	if err != nil {
		return "could not marshall error bag"
	}

	return string(data)
}

type Code int

// common response errors

const (
	BodyParserErrCode Code = 10001
	ProcessingErrCode Code = 10002
)

func BodyParser() Bag {
	return Bag{
		Code:    BodyParserErrCode,
		Message: "could not parse request body.",
	}
}

func Processing() Bag {
	return Bag{
		Code:    ProcessingErrCode,
		Message: "Error occurred when processing the request.",
	}
}
