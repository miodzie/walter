package trigger

import (
	"regexp"
	"strings"
)

const WORD = `(?i)\b{WORD}\b`

type Trigger struct {
	Id    uint64
	Word  string
	Reply string
	Regex string
}

func (t *Trigger) Check(str string) bool {
	reg := strings.ReplaceAll(WORD, "{WORD}", t.Word)
	r, _ := regexp.Compile(reg)

	return r.Match([]byte(str))
}

type Triggers interface {
	All() ([]Trigger, error)
	Add(*Trigger) error
	Remove(id uint64)
}
