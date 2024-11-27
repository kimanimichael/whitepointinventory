package httpauth

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"
)

/*
Get Password from http header
Authorization Password {Insert Password here}:Email {Insert Email here}
*/
func GetPasswordAndEmail(headers http.Header) (string, string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", "", errors.New("no authentication info found")
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

func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("no authentication info found in header")
	}
	vals := strings.Split(val, ":")
	if len(vals) != 2 {
		return "", errors.New("malformed authentication info:Wrong auth format")
	}
	valHeaderName := vals[0]
	if valHeaderName != "APIKey" {
		return "", errors.New("malformed authentication info:Wrong auth header - Bad header name")
	}
	valHeaderValue := vals[1]
	if valHeaderValue == "" {
		return "", errors.New("malformed authentication info:Wrong auth header - No APIKey")
	}
	return valHeaderValue, nil
}

func GetPasswordAndEmailFromBody(r *http.Request) (string, string, error) {
	type parameters struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	decoder := json.NewDecoder(r.Body)

	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		return "", "", err
	}

	if params.Email == "" || params.Password == "" {
		return "", "", errors.New("empty, email or password")
	}

	if !strings.Contains(params.Email, "@") {
		return "", "", errors.New("email malformed missing @")
	}
	if !strings.Contains(params.Email, ".") {
		return "", "", errors.New("email malformed missing a period")
	}

	return params.Email, params.Password, nil
}
