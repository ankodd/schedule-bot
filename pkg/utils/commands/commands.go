package commands

import "strings"

func IsCommand(text string) bool {
	return strings.Count(text, "/") == 1 && text[0] == '/'
}
