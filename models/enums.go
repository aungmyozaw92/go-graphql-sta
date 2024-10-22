package models

import (
	"errors"
	"io"
	"strconv"
)


type Precision string

const (
	PrecisionZero  Precision = "0"
	PrecisionOne   Precision = "1"
	PrecisionTwo   Precision = "2"
	PrecisionThree Precision = "3"
	PrecisionFour  Precision = "4"
)

func (p Precision) MarshalGQL(w io.Writer) {
	w.Write([]byte(strconv.Quote(string(p))))
}

func (p *Precision) UnmarshalGQL(i interface{}) error {
	str, ok := i.(string)
	if !ok {
		return errors.New("precision must be string")
	}

	switch str {
	case "0":
		*p = PrecisionZero
	case "1":
		*p = PrecisionOne
	case "2":
		*p = PrecisionTwo
	case "3":
		*p = PrecisionThree
	case "4":
		*p = PrecisionFour
	default:
		return errors.New("invalid precision")
	}
	return nil
}