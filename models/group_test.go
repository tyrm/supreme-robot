package models

import (
	"fmt"
	"github.com/google/uuid"
	"testing"
)

func TestGroupTitle(t *testing.T) {
	tables := []struct {
		x uuid.UUID
		n string
	}{
		{uuid.Must(uuid.Parse("0b1e0a88-37a4-48f0-8060-2814906fa9f7")), "DNS Admin"},
		{uuid.Must(uuid.Parse("71df8f2b-f293-4fde-93b1-e40dbe5c97ea")), "Super Admin"},
		{uuid.Must(uuid.Parse("fbc827a0-32db-4d71-b95e-632b414e7993")), "User Admin"},
		{uuid.Must(uuid.Parse("1d34fcbd-a027-40cc-91d6-bf42a7f2122a")), ""},
	}

	for i, table := range tables {
		i := i
		table := table
		name := fmt.Sprintf("Testing GroupTitle for %s", table.x)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			title := GroupTitle(table.x)

			if title != table.n {
				t.Errorf("[%d] function return for %s wrong, got: %v, want: %v,", i, table.x, title, table.n)
			}
		})
	}
}
