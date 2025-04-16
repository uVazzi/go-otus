package hw10programoptimization

import (
	"encoding/json"
	"errors"
	"io"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

var ErrIncorrectEmail = errors.New("incorrect email")

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	var user User
	result := make(DomainStat)
	decoder := json.NewDecoder(r)

	for decoder.More() {
		err := decoder.Decode(&user)
		if err != nil {
			return nil, err
		}

		if strings.HasSuffix(user.Email, "."+domain) {
			key := strings.SplitN(user.Email, "@", 2)
			if len(key) != 2 {
				return nil, ErrIncorrectEmail
			}

			result[strings.ToLower(key[1])]++
		}
	}

	return result, nil
}
