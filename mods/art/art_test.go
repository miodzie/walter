package art

import (
	"fmt"
	"testing"
)

func TestArt_gm(t *testing.T) {
	art := &Picture{Art: gm}
	fmt.Println("test")
	for !art.Completed() {
		t.Log(art.NextLine())
	}
	t.Fail()
}
