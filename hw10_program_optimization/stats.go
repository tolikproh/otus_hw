package hw10programoptimization

import (
	"bufio"
	"io"
	"strings"

	jsoniter "github.com/json-iterator/go"
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
	scan := bufio.NewScanner(r)
	result := make(DomainStat)
	var user User

	for scan.Scan() {
		if !strings.Contains(scan.Text(), domain) {
			continue
		}
		if err := jsoniter.Unmarshal(scan.Bytes(), &user); err != nil {
			return nil, err
		}
		if strings.HasSuffix(user.Email, "."+domain) {
			result[strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])]++
		}
	}

	return result, nil
}
