package hw10programoptimization

import (
	"bufio"
	"errors"
	"io"
	"strings"

	"github.com/valyala/fastjson" //nolint: depguard
)

type User struct {
	ID       int
	Name     string
	Username string
	Email    string
	Phone    string
	Password string
	Address  string
}

type DomainStat map[string]int

func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	return countDomains(r, domain)
}

func countDomains(r io.Reader, domain string) (DomainStat, error) {
	result := make(DomainStat)

	scanner := bufio.NewScanner(r)

	for scanner.Scan() {
		email := fastjson.GetString(scanner.Bytes(), "Email")

		if email == "" {
			return nil, errors.New("no email field")
		}

		if strings.Contains(email, domain) {
			result[strings.ToLower(strings.SplitN(email, "@", 2)[1])]++
		}
	}

	if err := scanner.Err(); err != nil {
		return result, err
	}
	return result, nil
}
