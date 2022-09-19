package sed

import (
	"testing"
)

func TestParseSed(t *testing.T) {
	result := ParseSed("s/potato/tomato/g")

	if result.Options != "g" {
		t.Fail()
	}
	if result.Command != "s" {
		t.Fail()
	}
	if len(result.Replacements) != 2 {
		t.Fail()
	}
}

func TestItReplacesTheFirstInstance(t *testing.T) {
	sed := ParseSed("s/tomato/potato")

	if sed.Replace("I love tomatoes.") != "I love potatoes." {
		t.Logf("Failed to love potatos!")
		t.Fail()
	}

	if sed.Replace("tomato tomato, same thing.") != "potato tomato, same thing." {
		t.Fail()
	}
}

func TestItReplacesAllInstances(t *testing.T) {
	sed := ParseSed("s/tomato/potato/g")

	result := sed.Replace("tomato tomato, same thing.")
	if result != "potato potato, same thing." {
		t.Log("it's essentially the same")
		t.Log(result)
		t.Fail()
	}
}
