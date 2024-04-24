package auth

import (
	"errors"
	"net/http"
	"strings"
)

/*
Get Password from http header
Authorization Password {Insert Password here}:Name {Insert Name here}
*/
func GetPasswordAndName(headers http.Header)(string, string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", "",  errors.New("no authentication info found")
	}

	vals := strings.Split(val, ":")
	if len(vals) != 2 {
		return "", "", errors.New("malformed authentication info")
	}

	valPassword := vals[0]
	valsPassword := strings.Split(valPassword, " ")
	if len(valsPassword) != 2 {
		return "", "", errors.New("malformed authentication info")
	}
	if valsPassword[0] != "Password" {
		return "", "", errors.New("malformed authentication info")
	}


	valName := vals[1]
	valsName := strings.Split(valName, " ")
	if len(valsName) < 2 {
		return "", "", errors.New("malformed authentication info")
	}
	if valsName[0] != "Name" {
		return "", "", errors.New("malformed authentication info")
	}
	
	
	return valsPassword[1], valsName[1], nil
}