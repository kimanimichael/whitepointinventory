package auth

import (
	"errors"
	"net/http"
	"strings"
)

/*
Get Password from http header
Authorization Password {Insert Password here}:Email {Insert Email here}
*/
func GetPasswordAndEmail(headers http.Header)(string, string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", "",  errors.New("no authentication info found")
	}

	vals := strings.Split(val, ":")
	if len(vals) != 2 {
		return "", "", errors.New("malformed authentication info:Wrong auth format")
	}

	valPassword := vals[0]
	valsPassword := strings.Split(valPassword, " ")
	if len(valsPassword) != 2 {
		return "", "", errors.New("malformed authentication info:Password info wrong format")
	}
	if valsPassword[0] != "Password" {
		return "", "", errors.New("malformed authentication info:Password phrase missing")
	}


	valName := vals[1]
	valsName := strings.Split(valName, " ")
	if len(valsName) < 2 {
		return "", "", errors.New("malformed authentication info:Email info wring format")
	}
	if valsName[0] != "Email" {
		return "", "", errors.New("malformed authentication info:Email phrase missing")
	}

	if !strings.Contains(valsName[1], "@") {
		return "", "", errors.New("malformed authentication info:Bad email entered")
	}
	
	
	return valsPassword[1], valsName[1], nil
}