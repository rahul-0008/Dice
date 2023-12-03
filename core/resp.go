package core

import (
	"bytes"
	"errors"
	"fmt"
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

func Decode(data []byte) ([]interface{}, error) {
	// incoming data will be of  bytes array from network
	if len(data) == 0 {
		return nil, errors.New("No data")
	}

	// now need to decode the value
	var index = 0
	var values []interface{} = make([]interface{}, 0)
	for index < len(data) {
		value, delta, err := DecodeOne(data[index:])
		if err != nil {
			return values, err
		}
		index = index + delta

		values = append(values, value)
	}

	return values, nil

}

func encodeString(v string) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(v), v))
}

func Encode(value interface{}, isSimple bool) []byte {
	switch v := value.(type) {
	case string:
		if isSimple {
			return []byte(fmt.Sprintf("+%s\r\n", v))
		} else {
			return encodeString(v)
		}
	case int, int8, int16, int32, int64:
		return []byte(fmt.Sprintf(":%d\r\n", v))
	case []string:
		var b []byte
		buf := bytes.NewBuffer(b)
		for _, token := range value.([]string) {
			buf.Write(encodeString(token))
		}
		return []byte(fmt.Sprintf("*%d\r\n%s", len(v), buf.Bytes()))
	case error:
		return []byte(fmt.Sprintf("-%s\r\n", v))
	default:
		return RESP_NIL
	}
}
