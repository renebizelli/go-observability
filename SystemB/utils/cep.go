package utils

import (
	"errors"
	"regexp"
	"strings"
)

type CEP struct {
	searchedCEP string
}

func NewCEP(searchedCEP string) *CEP {
	return &CEP{
		searchedCEP: searchedCEP,
	}
}

func (c *CEP) Validate() error {

	if len(c.searchedCEP) != 8 {
		return errors.New("invalid cep")
	}

	re := regexp.MustCompile("[0-9]+")
	c.searchedCEP = strings.Join(re.FindAllString(c.searchedCEP, -1)[:], "")

	if len(c.searchedCEP) != 8 {
		return errors.New("invalid cep")
	}

	return nil

}
