package utils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCepValidateOk(t *testing.T) {
	cep := NewCEP("03132010")

	e := cep.Validate()

	assert.Nil(t, e)
}

func TestCepValidateShort(t *testing.T) {
	cep := NewCEP("0313201")

	e := cep.Validate()

	assert.Error(t, e)
}

func TestCepValidateInvalid(t *testing.T) {
	cep := NewCEP("0313201d")

	e := cep.Validate()

	assert.Error(t, e)
}
