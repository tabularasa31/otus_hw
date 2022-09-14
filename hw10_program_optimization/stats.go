package hw10programoptimization

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
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
	scanner := bufio.NewScanner(r)
	scanner.Split(bufio.ScanLines)

	user := new(User)
	result := make(DomainStat)
	suffix := "." + domain

	// read line by line
	for scanner.Scan() {
		if !strings.Contains(scanner.Text(), domain) {
			continue
		}

		if err := json.Unmarshal(scanner.Bytes(), &user); err != nil {
			return result, err
		}

		if strings.HasSuffix(user.Email, suffix) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
		user = &User{}
	}

	return result, nil
}
