package dbutil

import "fmt"

type Scanner interface {
	Scan(dest ...interface{}) error
}

func GenerateSeqString(offset, count int) string {
	s := ""
	for i := 1 + offset; i <= count+offset; i++ {
		s += fmt.Sprintf(",$%d", i)
	}
	if len(s) > 0 {
		return s[1:]
	} else {
		return s
	}
}
