package _tests

import (
	"github.com/transfer360/go-transfer360/search"
	"testing"
	"time"
)

func TestHelloName(t *testing.T) {

	sr := search.Request{
		VRM:       "test",
		DateTime:  time.Now().Format(time.RFC3339),
		Reference: "Your Reference",
	}

	err := sr.Validate()

	if err != nil {
		t.Fatal(err)
	}

}
