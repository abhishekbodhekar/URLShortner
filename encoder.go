package main

import "encoding/base64"

// This encodes the string passed to base64 and return it
func encodeToB64(data string) string {
	return base64.StdEncoding.EncodeToString([]byte(data))
}

//It decodes the string passed and returns it with okidiom
func decodeFromB64(data string) (string, bool) {
	decodedString, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", false
	}
	return string(decodedString), true

}
