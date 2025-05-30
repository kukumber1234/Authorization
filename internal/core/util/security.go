package util

import (
	"regexp"

	"github.com/microcosm-cc/bluemonday"
)

var xssPatterns = []*regexp.Regexp{
	regexp.MustCompile(`(?i)<script[^>]*>.*?</script>`),
	regexp.MustCompile(`(?i)on\w+\s*=`),   // onclick=, onload= и т.д.
	regexp.MustCompile(`(?i)javascript:`), // href="javascript:"
	regexp.MustCompile(`(?i)eval\(`),      // eval()
}

func isMalicious(input string) bool {
	for _, pattern := range xssPatterns {
		if pattern.MatchString(input) {
			return true
		}
	}
	return false
}

func CheckInput(message string) bool {
	p := bluemonday.UGCPolicy()
	cleaned := p.Sanitize(message)

	if message != cleaned || isMalicious(message) {
		return true
	}

	return false
}