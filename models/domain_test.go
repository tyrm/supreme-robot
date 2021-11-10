package models

import "testing"

func TestDomain_Validate(t *testing.T) {
	tables := []struct {
		x string
		n bool
	}{
		{"google.com.", true},
		{"asdf2.", true},
		{"xn--c1yn36f.", true},
		{"blog.xn--c1yn36f.", true},
		{".xn--c1yn36f.", false},
		{"what?.", false},
		{"google", false},
		{"@", false},
	}

	for _, table := range tables {
		d := Domain{
			Domain: table.x,
		}

		valid := d.Validate()
		if valid != table.n {
			t.Errorf("regex match on %s failed, got: %v, want: %v,", table.x, valid, table.n)
		}
	}
}
