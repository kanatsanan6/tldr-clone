package utils

import (
	"bufio"
	"io"
)

func Render(markdown io.Reader) (string, error) {
	var result string
	scanner := bufio.NewScanner(markdown)
	for scanner.Scan() {
		line := scanner.Text()
		result += line
	}
	return result, scanner.Err()
}
