package utils

import "regexp"

func ConvertToRegexSlice(regexSlice []string) []*regexp.Regexp {
	var result []*regexp.Regexp

	for _, cf := range regexSlice {
		rgx, err := regexp.Compile(cf)
		if err != nil {
			// TODO logs
			continue
		}
		result = append(result, rgx)
	}

	return result
}

func RegexMatch(key string, regexSlice []*regexp.Regexp) bool {
	for _, r := range regexSlice {
		if r.MatchString(key) {
			return true
		}
	}
	return false
}
