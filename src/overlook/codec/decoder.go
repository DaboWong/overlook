package codec

import (
	"overlook/ds"
)

type Decoder interface {
	Decode(b []byte, v ...interface{}) (*ds.Event, error)
}

type Encoder interface {
	Encode(id int32, v ...interface{}) ([]byte, error)
}
