package assert

import "strings"
import "testing"

func IntEqual(t *testing.T, a, b int) {
	if a != b {
		t.Errorf("%d != %d", a, b)
	}
}

func StringEqual(t *testing.T, a, b string) {
	if a != b {
		t.Errorf("%s != %s", a, b)
	}
}

func StringsEqual(t *testing.T, a, b []string) {
	s1 := strings.Join(a, ",")
	s2 := strings.Join(b, ",")
	StringEqual(t, s1, s2)
}
