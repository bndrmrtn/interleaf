package extract

import (
	"io"
	"strings"
)

func cleanText(s string) string {
	s = strings.ReplaceAll(s, "\x00", "")
	s = strings.ReplaceAll(s, "\r\n", "\n")

	return strings.TrimSpace(s)
}

func mustClose(c io.Closer) {
	if err := c.Close(); err != nil {
		panic(c)
	}
}
