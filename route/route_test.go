package route

import (
	"fmt"
	"testing"
)

func TestList(t *testing.T) {
	set, err := Routed()
	if err != nil {
		t.Error(err)
	}

	fmt.Sprint(set)
}
