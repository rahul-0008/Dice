package core

import (
	"errors"
)

// TODO : IMPLEMENT RESP PROROCOL :
// a . simple strings
//	b . bulk strings
// c . integer
// d . array
// e . error

func readLength(data []byte) (int, int) {

	pos, length := 0, 0

	for pos = range data {
		b := data[pos]
		if !(b >= '0' && b <= '9') {
			return length, pos + 2
		}
		length = length*10 + int(b-'0')
	}
	return 0, 0

}

func readSimpleString(data []byte) (string, int, error) {

	pos := 1
	for data[pos] != '\r' {
		pos++
	}
	return string(data[1:pos]), pos + 2, nil
}

func readBulkString(data []byte) (string, int, error) {
	pos := 1
	length, delta := readLength(data[pos:])
	pos += delta

	return string(data[pos:(pos + length)]), (pos + length) + 2, nil

}
func readError(data []byte) (string, int, error) {
	return readSimpleString(data)

}
func readInt64(data []byte) (int64, int, error) {
	pos := 1
	var value int64 = 0
	for data[pos] != '\r' {
		value = value*10 + int64(data[pos]-'0')
		pos++
	}
	return value, pos + 2, nil
}
func readArray(data []byte) (interface{}, int, error) {
	pos := 1
	length, delta := readLength(data[pos:])

	var elems []interface{} = make([]interface{}, length)
	pos += delta

	for i := range elems {
		elem, delta, err := DecodeOne(data[pos:])
		if err != nil {
			return nil, 0, err
		}
		elems[i] = elem
		pos += delta
	}
	return elems, pos, nil
}

func DecodeOne(data []byte) (interface{}, int, error) {
	if len(data) == 0 {
		return nil, 0, errors.New("No data")
	}
	switch data[0] {
	case '+':
		return readSimpleString(data)
	case '*':
		//array
		return readArray(data)
	case ':':
		// integer
		return readInt64(data)
	case '$':
		//bulk string
		return readBulkString(data)
	case '-':
		// error
		return readError(data)

	}
	return nil, 0, nil
}

func Decode(data []byte) (interface{}, error) {
	// incoming data will be of  bytes array from network
	if len(data) == 0 {
		return nil, errors.New("No data")
	}

	// now need to decode the value
	value, _, err := DecodeOne(data)

	return value, err

}
