package utils

import "strings"

func CleanAffiliateCodes(codes []string) []string {
	var cleanedCodes []string
	for _, code := range codes {
		trimmedCode := strings.TrimSpace(code)
		if trimmedCode != "" {
			cleanedCodes = append(cleanedCodes, trimmedCode)
		}
	}
	return cleanedCodes
}
