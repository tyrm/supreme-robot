package tests

import (
	"github.com/tyrm/supreme-robot/kv"
	"testing"
)

// DoAddDomain tests the AddDomain function
func DoAddDomain(t *testing.T, client kv.DNS) {
	newDomain := "doadddomain."
	err := client.AddDomain(newDomain)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	domains, err := client.GetDomains()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	found := false
	for _, d := range *domains {
		if d == newDomain {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected domain %s, got: %v", newDomain, domains)
	}
}

// DoGetDomain tests the GetDomain function
func DoGetDomain(t *testing.T, client kv.DNS) {
	newDomain := "dogetdomain."
	err := client.AddDomain(newDomain)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	domains, err := client.GetDomains()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	found := false
	for _, d := range *domains {
		if d == newDomain {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected domain %s, got: %v", newDomain, domains)
	}
}

// DoRemoveDomain tests the AddDomain function
func DoRemoveDomain(t *testing.T, client kv.DNS) {
	newDomain := "doremovedomain."
	err := client.AddDomain(newDomain)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	domains, err := client.GetDomains()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	found := false
	for _, d := range *domains {
		if d == newDomain {
			found = true
			break
		}
	}

	if !found {
		t.Errorf("expected domain %s, got: %v", newDomain, domains)
	}

	err = client.RemoveDomain(newDomain)
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	domains2, err := client.GetDomains()
	if err != nil {
		t.Errorf("unexpected error, got: %s, want: nil.", err.Error())
		return
	}

	found2 := false
	for _, d := range *domains2 {
		if d == newDomain {
			found2 = true
			break
		}
	}

	if found2 {
		t.Errorf("unexpected domain %s, got: %v", newDomain, domains)
	}
}
