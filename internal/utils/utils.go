package utils

import "os"

func ExpandEnv(s string) string {
	return os.Expand(s, func(key string) string {
		if val, ok := os.LookupEnv(key); ok {
			return val
		}

		return ""
	})
}
