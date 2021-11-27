package filters

import "strings"

type header struct {
	key   string
	value string
}

func newHeader(s string) *header {
	parts := strings.Split(s, "=")
	if len(parts) != 2 {
		return nil
	}

	return &header{
		key:   parts[0],
		value: parts[1],
	}
}
