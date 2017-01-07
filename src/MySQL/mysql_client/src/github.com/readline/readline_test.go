package readline

import (
	"strings"
	"testing"
)

func TestReadLine(t *testing.T) {
	s := "line1\nline2"
	r := strings.NewReader(s)
	e := ""

	ReadLine(r, func(line string) {
		e += line + "\n"
	})

	if e != "line1\nline2\n" {
		t.Fail()
	}
}

func TestReadLine2(t *testing.T) {
	s := "line1\r\nline2"
	r := strings.NewReader(s)
	e := ""

	ReadLine(r, func(line string) {
		e += line + "\n"
	})

	if e != "line1\nline2\n" {
		t.Fail()
	}
}

func TestReadLine3(t *testing.T) {
	s := "line1\r\nline2\n"
	r := strings.NewReader(s)
	e := ""

	ReadLine(r, func(line string) {
		e += line + "\n"
	})

	if e != "line1\nline2\n" {
		t.Fail()
	}
}

func TestReadLine4(t *testing.T) {
	s := "\n"
	r := strings.NewReader(s)
	e := ""

	ReadLine(r, func(line string) {
		e = line
	})

	if e != "" {
		t.Fail()
	}
}
