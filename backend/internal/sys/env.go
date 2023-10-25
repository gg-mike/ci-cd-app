package sys

import (
	"fmt"
	"os"
)

func GetEnvWithFallback(key, fallback string) string {
  if value, ok := os.LookupEnv(key); ok {
      return value
  }
  return fallback
}

func GetRequiredEnv(key string) (string, error) {
	if value, ok := os.LookupEnv(key); ok {
		return value, nil
	}
	return "", fmt.Errorf(`env variable "%s" is missing`, key)
}
