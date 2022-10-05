package art

import "strings"

var gm = `
 ▄██████▄    ▄▄▄▄███▄▄▄▄         ▄████████    ▄█    █▄       ▄████████     ███
 ███    ███ ▄██▀▀▀███▀▀▀██▄      ███    ███   ███    ███     ███    ███ ▀█████████▄
 ███    █▀  ███   ███   ███      ███    █▀    ███    ███     ███    ███    ▀███▀▀██
▄███        ███   ███   ███      ███         ▄███▄▄▄▄███▄▄   ███    ███     ███   ▀
 ███ ████▄  ███   ███   ███      ███        ▀▀███▀▀▀▀███▀  ▀███████████     ███
 ███    ███ ███   ███   ███      ███    █▄    ███    ███     ███    ███     ███
 ███    ███ ███   ███   ███      ███    ███   ███    ███     ███    ███     ███
 ████████▀   ▀█   ███   █▀       ████████▀    ███    █▀      ███    █▀     ▄████▀
`

type Picture struct {
	Art         string
	CurrentLine int
}

func (p *Picture) NextLine() string {
	if p.Completed() {
		return ""
	}
	lines := strings.Split(p.Art, "\n")
	defer func() { p.CurrentLine++ }()
	if p.CurrentLine == len(lines) {
		p.CurrentLine = 0
		return ""
	}
	return lines[p.CurrentLine]
}

func (p *Picture) Completed() bool {
	lines := strings.Split(p.Art, "\n")

	return p.CurrentLine == len(lines)
}
