package utils

import (
	"os"
	"strconv"
	"time"
)

func GetAddr() string {
	if addr, ok := os.LookupEnv("ADDR"); ok {
		return addr
	}

	return "127.0.0.1"
}

func GetPort() uint64 {
	if port, ok := os.LookupEnv("PORT"); ok {
		if parsed, err := strconv.ParseUint(port, 10, 64); err == nil {
			return parsed
		}
	}

	return 13000
}

func GetClientCount() uint64 {
	if client, ok := os.LookupEnv("CLIENT_COUNT"); ok {
		if parsed, err := strconv.ParseUint(client, 10, 64); err == nil {
			return parsed
		}
	}

	return 1
}

func GetRepeatCount() uint64 {
	if repeat, ok := os.LookupEnv("REPEAT_COUNT"); ok {
		if parsed, err := strconv.ParseUint(repeat, 10, 64); err == nil {
			return parsed
		}
	}

	return 1
}

func GetDelay() time.Duration {
	if delay, ok := os.LookupEnv("DELAY"); ok {
		if parsed, err := time.ParseDuration(delay); err == nil {
			return parsed
		}
	}

	return time.Millisecond * 100
}

func GetLoggingPreference() bool {
	if _log, ok := os.LookupEnv("LOG"); ok {
		if parsed, err := strconv.ParseBool(_log); err == nil {
			return parsed
		}
	}

	return true
}

func GetCapacity() uint64 {
	if capacity, ok := os.LookupEnv("CAPACITY"); ok {
		if parsed, err := strconv.ParseUint(capacity, 10, 64); err == nil {
			return parsed
		}
	}

	return 256
}
