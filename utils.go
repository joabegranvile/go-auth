package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func getSecret(secretName string) string {
	path := fmt.Sprintf("/run/secrets/%s", secretName)
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatal("Erro ao ler segredo %s: %v", secretName, err)
	}
	return strings.TrimSpace(string(data))
}
