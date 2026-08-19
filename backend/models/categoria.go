package models

import (
	"errors"
	"regexp"
)

var corHexRegex = regexp.MustCompile(`^#[0-9A-Fa-f]{6}$`)

type Categoria struct {
	ID   int64  `json:"id"`
	Nome string `json:"nome"`
	Cor  string `json:"cor"`
}

func (c *Categoria) Validate() error {
	if !CorValida(c.Cor) {
		return errors.New("cor inválida, utilize o formato hexadecimal, ex: #FFAA00")
	}

	return nil
}

func CorValida(cor string) bool {
	return corHexRegex.MatchString(cor)
}
