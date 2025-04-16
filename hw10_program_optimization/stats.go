package hw10programoptimization

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

type User struct {
	Email string
}

type DomainStat map[string]int

var ErrIncorrectEmail = errors.New("incorrect email")

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	u, err := getUsers(r)
	if err != nil {
		return nil, fmt.Errorf("get users error: %w", err)
	}
	return countDomains(u, domain)
}

type users []User

func getUsers(r io.Reader) (result users, err error) {
	var user User
	decoder := json.NewDecoder(r)

	for decoder.More() {
		err = decoder.Decode(&user)
		if err != nil {
			return
		}
		result = append(result, user)
	}

	return
}

func countDomains(usersData users, domain string) (DomainStat, error) {
	result := make(DomainStat)
	for _, user := range usersData {
		if strings.HasSuffix(user.Email, domain) {
			key := strings.SplitN(user.Email, "@", 2)
			if len(key) != 2 {
				return nil, ErrIncorrectEmail
			}

			result[strings.ToLower(key[1])]++
		}
	}

	return result, nil
}
