package policing

import (
	"regexp"

	"github.com/miodzie/seras"
)

// If the string has more than 3 all caps words.
const allCapsRegex = `(\b[A-Z]+\s?\b){3,}`

func IsSpam(msg seras.Message) bool {
	r, _ := regexp.Compile(allCapsRegex)

	return r.MatchString(msg.Content)
}
