package config

import "errors"

const KeySoaNS = "soa-ns"

var ErrorNotDefined = errors.New("config value not defined")

type Config interface {
	Get(string) (string, error)
	Set(string, string) error
}

