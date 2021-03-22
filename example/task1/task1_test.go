package task1

import "testing"

func TestHello(t *testing.T) {
	if Hello() != "Hello, World!" {
		t.Error("Incorrect hello")
	}
}