package config

const KeySoaNS = "soa-ns"
const KeyAllowSignups = "allow-signup"

var defaults = map[string]string{
	KeyAllowSignups: "false",
}

//var ErrorNotDefined = errors.New("config value not defined")

type Config interface {
	Get(string) (*string, error)
	MGet(*[]string) (*map[string]*string, error)
	MSet(*map[string]string) error
	Set(string, string) error
}

func Default(k string) *string {
	if d, ok := defaults[k]; ok {
		return &d
	}
	return nil
}
