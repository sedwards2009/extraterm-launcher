package wordcase

import "strings"

func KababCaseToCamelCaseMapKeys(sourceMap map[string]string) map[string]string {
	result := map[string]string{}
	for key, value := range sourceMap {
		result[KababCaseToCamelCase(key)] = value
	}
	return result
}

func KababCaseToCamelCase(word string) string {
	parts := strings.Split(word[2:], "-")

	result := ""
	for i, part := range parts {
		if part == "" {
			continue
		}
		if i == 0 {
			result = result + part
		} else {
			result = result + strings.ToUpper(part[0:1]) + part[1:]
		}
	}
	return result
}
