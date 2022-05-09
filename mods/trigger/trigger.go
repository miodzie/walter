package trigger

const WORD = `(?i)\b{{.Word}}\b`

type Trigger struct {
	Id    uint64
	Word  string
	Reply string
	Regex string
}

type Triggers interface {
	All() ([]Trigger, error)
	Add(*Trigger) error
	Remove(id uint64)
}
