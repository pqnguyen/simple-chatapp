package utils

import "bytes"

func WrapperMessage(buf []byte) []byte {
	if ok := bytes.HasSuffix(buf, []byte{'\n'}); !ok {
		buf = append(buf, '\n')
	}
	return buf
}
