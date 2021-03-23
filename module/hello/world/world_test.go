package world

import "testing"

func TestWorld(t *testing.T) {
    want := "world"
    if got := World(); got != want {
        t.Errorf("World() = %q, want %q", got, want)
    }
}

