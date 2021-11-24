package models

import (
	"database/sql"
	"fmt"
	"testing"
)

func TestRecordValidateTypeA(t *testing.T) {
	tables := []struct {
		x Record
		n error
	}{
		{Record{Type: RecordTypeA}, errMissingName},
		{Record{Type: RecordTypeA, Value: "8.4.5.9", TTL: 300}, errMissingName},
		{Record{Type: RecordTypeA, Name: "example", TTL: 300}, errMissingIP},
		{Record{Type: RecordTypeA, Name: "@", Value: "1.2.3.4"}, errMissingTTL},
		{Record{Type: RecordTypeA, Name: "@", Value: "1.2.3.4", TTL: 300}, nil},
		{Record{Type: RecordTypeA, Name: "@", Value: "example.com", TTL: 300}, errInvalidIP},
		{Record{Type: RecordTypeA, Name: "@", Value: "10.2.1.400", TTL: 300}, errInvalidIP},
		{Record{Type: RecordTypeA, Name: "@", Value: "10.1.400.2", TTL: 300}, errInvalidIP},
		{Record{Type: RecordTypeA, Name: "@", Value: "10.400.1.2", TTL: 300}, errInvalidIP},
		{Record{Type: RecordTypeA, Name: "@", Value: "400.10.1.2", TTL: 300}, errInvalidIP},
		{Record{Type: RecordTypeA, Name: "test1", Value: "192.168.0.11", TTL: 300}, nil},
		{Record{Type: RecordTypeA, Name: "test1", Value: "test1.dev", TTL: 300}, errInvalidIP},
		{Record{Type: RecordTypeA, Name: "test1.", Value: "55.195.15.4", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeA, Name: ".test1", Value: "4.88.15.98", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeA, Name: ".test1", Value: "4.88.15.98", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeA, Name: "-test1", Value: "4.88.15.98", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeA, Name: "test1-", Value: "4.88.15.98", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeA, Name: "sub.test1", Value: "125.45.155.7", TTL: 300}, nil},
		{Record{Type: RecordTypeA, Name: "sub.test1", Value: "125.45.155.7", TTL: -300}, errInvalidTTL},
	}

	for i, table := range tables {
		i := i
		table := table
		name := fmt.Sprintf("[%d] Validate S record Name: %s, IP: %s TTL: %d", i, table.x.Name, table.x.Value, table.x.TTL)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := table.x.Validate()
			if err != table.n {
				t.Errorf("[%d] validation failed, got: %v, want: %v,", i, err, table.n)
			}
		})
	}
}

func TestRecordValidateTypeAAAA(t *testing.T) {
	tables := []struct {
		x Record
		n error
	}{
		{Record{Type: RecordTypeAAAA}, errMissingName},
		{Record{Type: RecordTypeAAAA, Value: "bc24:ca41:bf48:1b2b:dc22:6d41:46ff:1ec7", TTL: 300}, errMissingName},
		{Record{Type: RecordTypeAAAA, Name: "example", TTL: 300}, errMissingIP},
		{Record{Type: RecordTypeAAAA, Name: "@", Value: "4087:47e4::9e7d"}, errMissingTTL},
		{Record{Type: RecordTypeAAAA, Name: "@", Value: "4087:47e4::9e7d", TTL: 300}, nil},
		{Record{Type: RecordTypeAAAA, Name: "@", Value: "2001:db8::", TTL: 300}, nil},
		{Record{Type: RecordTypeAAAA, Name: "@", Value: "example.com", TTL: 300}, errInvalidIP},
		{Record{Type: RecordTypeAAAA, Name: "@", Value: "jpxy:6c89:19d2:fd00:5a15:8d82:6e67:f294", TTL: 300}, errInvalidIP},
		{Record{Type: RecordTypeAAAA, Name: "@", Value: "71ea:0ae2:6322:XXXX:b354:7d5d:0062:53b9", TTL: 300}, errInvalidIP},
		{Record{Type: RecordTypeAAAA, Name: "@", Value: "d157:59cb:2224:914f:3ec5:3ee2:1159:poly", TTL: 300}, errInvalidIP},
		{Record{Type: RecordTypeAAAA, Name: "@", Value: "400.10.1.2", TTL: 300}, errInvalidIP},
		{Record{Type: RecordTypeAAAA, Name: "test1", Value: "f887:3aa8:24e1:a707:bda0:3447:57b3:0086", TTL: 300}, nil},
		{Record{Type: RecordTypeAAAA, Name: "test1", Value: "test1.dev", TTL: 300}, errInvalidIP},
		{Record{Type: RecordTypeAAAA, Name: "test1.", Value: "1e61:4949:28bd:a3d5:05d2:41f4:3504:920b", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeAAAA, Name: ".test1", Value: "1e61:4949:28bd:a3d5:05d2:41f4:3504:920b", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeAAAA, Name: ".test1", Value: "1e61:4949:28bd:a3d5:05d2:41f4:3504:920b", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeAAAA, Name: "-test1", Value: "1e61:4949:28bd:a3d5:05d2:41f4:3504:920b", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeAAAA, Name: "test1-", Value: "1e61:4949:28bd:a3d5:05d2:41f4:3504:920b", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeAAAA, Name: "sub.test1", Value: "d157:59cb:2224:914f:3ec5:3ee2:1159:684a", TTL: 300}, nil},
		{Record{Type: RecordTypeAAAA, Name: "sub.test1", Value: "d157:59cb:2224:914f:3ec5:3ee2:1159:684a", TTL: -300}, errInvalidTTL},
	}

	for i, table := range tables {
		i := i
		table := table
		name := fmt.Sprintf("Validate AAAA record Name: %s, IP: %s TTL: %d", table.x.Name, table.x.Value, table.x.TTL)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := table.x.Validate()
			if err != table.n {
				t.Errorf("[%d] validation failed, got: %v, want: %v,", i, err, table.n)
			}
		})
	}
}

func TestRecordValidateTypeCNAME(t *testing.T) {
	tables := []struct {
		x Record
		n error
	}{
		{Record{Type: RecordTypeCNAME}, errMissingName},
		{Record{Type: RecordTypeCNAME, Value: "example.com.", TTL: 300}, errMissingName},
		{Record{Type: RecordTypeCNAME, Name: "example", TTL: 300}, errMissingHost},
		{Record{Type: RecordTypeCNAME, Name: "@", Value: "example.com."}, errMissingTTL},
		{Record{Type: RecordTypeCNAME, Name: "@", Value: "example.com.", TTL: 300}, nil},
		{Record{Type: RecordTypeCNAME, Name: "@", Value: "example.com", TTL: 300}, errInvalidHost},
		{Record{Type: RecordTypeCNAME, Name: "@", Value: "10.2.1.400", TTL: 300}, errInvalidHost},
		{Record{Type: RecordTypeCNAME, Name: "test1", Value: "asdf2.", TTL: 300}, nil},
		{Record{Type: RecordTypeCNAME, Name: "test1", Value: "test1.dev", TTL: 300}, errInvalidHost},
		{Record{Type: RecordTypeCNAME, Name: "test1.", Value: "example.com.", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeCNAME, Name: ".test1", Value: "example.com.", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeCNAME, Name: ".test1", Value: "example.com.", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeCNAME, Name: "-test1", Value: "example.com.", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeCNAME, Name: "test1-", Value: "example.com.", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeCNAME, Name: "sub.test1", Value: "userGroups.google.com.", TTL: 300}, nil},
		{Record{Type: RecordTypeCNAME, Name: "sub.test1", Value: "userGroups.google.com.", TTL: -300}, errInvalidTTL},
	}

	for i, table := range tables {
		i := i
		table := table
		name := fmt.Sprintf("Validate CNAME record Name: %s, Host: %s TTL: %d", table.x.Name, table.x.Value, table.x.TTL)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := table.x.Validate()
			if err != table.n {
				t.Errorf("[%d] validation failed, got: %v, want: %v,", i, err, table.n)
			}
		})
	}
}

func TestRecordValidateTypeMX(t *testing.T) {
	tables := []struct {
		x Record
		n error
	}{
		{Record{Type: RecordTypeMX}, errMissingName},
		{Record{Type: RecordTypeMX, Value: "example.com.", TTL: 300}, errMissingName},
		{Record{Type: RecordTypeMX, Name: "example", TTL: 300}, errMissingHost},
		{Record{Type: RecordTypeMX, Name: "@", Value: "example.com."}, errMissingTTL},
		{Record{Type: RecordTypeMX, Name: "@", Value: "example.com.", TTL: 300}, errMissingPriority},
		{Record{Type: RecordTypeMX, Name: "@", Value: "example.com.", Priority: sql.NullInt32{Int32: 10, Valid: true}, TTL: 300}, errInvalidHost},
		{Record{Type: RecordTypeMX, Name: "@", Value: "example.com", Priority: sql.NullInt32{Int32: 10, Valid: true}, TTL: 300}, nil},
		{Record{Type: RecordTypeMX, Name: "@", Value: "10.2.1.252", Priority: sql.NullInt32{Int32: 10, Valid: true}, TTL: 300}, errInvalidHost},
		{Record{Type: RecordTypeMX, Name: "test1", Value: "asdf2.", Priority: sql.NullInt32{Int32: 10, Valid: true}, TTL: 300}, errInvalidHost},
		{Record{Type: RecordTypeMX, Name: "test1", Value: "test1.dev", Priority: sql.NullInt32{Int32: 10, Valid: true}, TTL: 300}, nil},
		{Record{Type: RecordTypeMX, Name: "test1.", Value: "example.com.", Priority: sql.NullInt32{Int32: 10, Valid: true}, TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeMX, Name: ".test1", Value: "example.com.", Priority: sql.NullInt32{Int32: 10, Valid: true}, TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeMX, Name: ".test1", Value: "example.com.", Priority: sql.NullInt32{Int32: 10, Valid: true}, TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeMX, Name: "-test1", Value: "example.com.", Priority: sql.NullInt32{Int32: 10, Valid: true}, TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeMX, Name: "test1-", Value: "example.com.", Priority: sql.NullInt32{Int32: 10, Valid: true}, TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeMX, Name: "sub.test1", Value: "userGroups.google.com", Priority: sql.NullInt32{Int32: 10, Valid: true}, TTL: 300}, nil},
		{Record{Type: RecordTypeMX, Name: "sub.test1", Value: "userGroups.google.com", Priority: sql.NullInt32{Int32: 10, Valid: true}, TTL: -300}, errInvalidTTL},
		{Record{Type: RecordTypeMX, Name: "sub.test1", Value: "userGroups.google.com", Priority: sql.NullInt32{Int32: 10, Valid: true}, TTL: -300}, errInvalidTTL},
		{Record{Type: RecordTypeMX, Name: "sub.test1", Value: "userGroups.google.com", Priority: sql.NullInt32{Int32: -3, Valid: true}, TTL: 300}, errInvalidPriority},
	}

	for i, table := range tables {
		i := i
		table := table
		name := fmt.Sprintf("Validate MX record Name: %s, Host: %s TTL: %d", table.x.Name, table.x.Value, table.x.TTL)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := table.x.Validate()
			if err != table.n {
				t.Errorf("[%d] validation failed, got: %v, want: %v,", i, err, table.n)
			}
		})
	}
}

func TestRecordValidateTypeNS(t *testing.T) {
	tables := []struct {
		x Record
		n error
	}{
		{Record{Type: RecordTypeNS}, errMissingName},
		{Record{Type: RecordTypeNS, Value: "example.com.", TTL: 300}, errMissingName},
		{Record{Type: RecordTypeNS, Name: "example", TTL: 300}, errMissingHost},
		{Record{Type: RecordTypeNS, Name: "something.something", Value: "example.com."}, errMissingTTL},
		{Record{Type: RecordTypeNS, Name: "@", Value: "example.com.", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeNS, Name: "@", Value: "example.com", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeNS, Name: "@", Value: "10.2.1.400", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeNS, Name: "test1", Value: "asdf2.", TTL: 300}, nil},
		{Record{Type: RecordTypeNS, Name: "test1", Value: "test1.dev", TTL: 300}, errInvalidHost},
		{Record{Type: RecordTypeNS, Name: "test1.", Value: "example.com.", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeNS, Name: ".test1", Value: "example.com.", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeNS, Name: ".test1", Value: "example.com.", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeNS, Name: "-test1", Value: "example.com.", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeNS, Name: "test1-", Value: "example.com.", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeNS, Name: "sub.test1", Value: "userGroups.google.com.", TTL: 300}, nil},
		{Record{Type: RecordTypeNS, Name: "sub.test1", Value: "userGroups.google.com.", TTL: -300}, errInvalidTTL},
	}

	for i, table := range tables {
		i := i
		table := table
		name := fmt.Sprintf("Validate NS record Name: %s, Host: %s TTL: %d", table.x.Name, table.x.Value, table.x.TTL)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := table.x.Validate()
			if err != table.n {
				t.Errorf("[%d] validation failed, got: %v, want: %v,", i, err, table.n)
			}
		})
	}
}

func TestRecordValidateTypeSOA(t *testing.T) {
	tables := []struct {
		x Record
		n error
	}{
		{Record{Type: RecordTypeSOA}, errMissingName},
		{Record{Type: RecordTypeSOA, Value: "example.com.", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errMissingName},
		{Record{Type: RecordTypeSOA, Name: "example", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errMissingNS},
		{Record{Type: RecordTypeSOA, Name: "something.something", Value: "example.com.", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}}, errMissingTTL},
		{Record{Type: RecordTypeSOA, Name: "@", Value: "example.com.", Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errMissingMBox},
		{Record{Type: RecordTypeSOA, Name: "@", Value: "example.com.", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errMissingExpire},
		{Record{Type: RecordTypeSOA, Name: "@", Value: "example.com.", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errMissingRefresh},
		{Record{Type: RecordTypeSOA, Name: "@", Value: "example.com.", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errMissingRetry},
		{Record{Type: RecordTypeSOA, Name: "@", Value: "example.com.", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, nil},
		{Record{Type: RecordTypeSOA, Name: "@", Value: "example.com", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errInvalidNS},
		{Record{Type: RecordTypeSOA, Name: "@", Value: "10.2.1.400", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errInvalidNS},
		{Record{Type: RecordTypeSOA, Name: "test1", Value: "asdf2.", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, nil},
		{Record{Type: RecordTypeSOA, Name: "test1", Value: "test1.dev", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errInvalidNS},
		{Record{Type: RecordTypeSOA, Name: "test1.", Value: "example.com.", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeSOA, Name: ".test1", Value: "example.com.", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeSOA, Name: ".test1", Value: "example.com.", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeSOA, Name: "-test1", Value: "example.com.", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeSOA, Name: "test1-", Value: "example.com.", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeSOA, Name: "sub.test1", Value: "userGroups.google.com.", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, nil},
		{Record{Type: RecordTypeSOA, Name: "sub.test1", Value: "userGroups.google.com.", MBox: sql.NullString{String: "hostmaster@example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errInvalidMBox},
		{Record{Type: RecordTypeSOA, Name: "sub.test1", Value: "userGroups.google.com.", MBox: sql.NullString{String: "testy@mctest.com", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errInvalidMBox},
		{Record{Type: RecordTypeSOA, Name: "sub.test1", Value: "userGroups.google.com.", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: -300}, errInvalidTTL},
		{Record{Type: RecordTypeSOA, Name: "sub.test1", Value: "userGroups.google.com.", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: -66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errInvalidExpire},
		{Record{Type: RecordTypeSOA, Name: "sub.test1", Value: "userGroups.google.com.", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: -44, Valid: true}, Retry: sql.NullInt32{Int32: 44, Valid: true}, TTL: 300}, errInvalidRefresh},
		{Record{Type: RecordTypeSOA, Name: "sub.test1", Value: "userGroups.google.com.", MBox: sql.NullString{String: "hostmaster.example.com.", Valid: true}, Expire: sql.NullInt32{Int32: 66, Valid: true}, Refresh: sql.NullInt32{Int32: 44, Valid: true}, Retry: sql.NullInt32{Int32: -44, Valid: true}, TTL: 300}, errInvalidRetry},
	}

	for i, table := range tables {
		i := i
		table := table
		name := fmt.Sprintf("Validate SOA record Name: %s, NS: %s TTL: %d", table.x.Name, table.x.Value, table.x.TTL)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := table.x.Validate()
			if err != table.n {
				t.Errorf("[%d] validation failed, got: %v, want: %v,", i, err, table.n)
			}
		})
	}
}

func TestRecordValidateTypeSRV(t *testing.T) {
	tables := []struct {
		x Record
		n error
	}{
		{Record{Type: RecordTypeSRV}, errMissingName},
		{Record{Type: RecordTypeSRV, Value: "example.com.", Port: sql.NullInt32{Int32: 5000, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}, TTL: 300}, errMissingName},
		{Record{Type: RecordTypeSRV, Name: "_xmpp-client._tcp", Port: sql.NullInt32{Int32: 5000, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}, TTL: 300}, errMissingHost},
		{Record{Type: RecordTypeSRV, Name: "_xmpp-client._tcp", Value: "example.com.", Port: sql.NullInt32{Int32: 5000, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}}, errMissingTTL},
		{Record{Type: RecordTypeSRV, Name: "_xmpp-client._tcp", Value: "example.com.", Priority: sql.NullInt32{Int32: 5000, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}, TTL: 300}, errMissingPort},
		{Record{Type: RecordTypeSRV, Name: "_xmpp-client._tcp", Value: "example.com.", Port: sql.NullInt32{Int32: 5000, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}, TTL: 300}, errMissingPriority},
		{Record{Type: RecordTypeSRV, Name: "_xmpp-client._tcp", Value: "example.com.", Port: sql.NullInt32{Int32: 5000, Valid: true}, Priority: sql.NullInt32{Int32: 10, Valid: true}, TTL: 300}, errMissingWeight},
		{Record{Type: RecordTypeSRV, Name: "_xmpp-client._tcp", Value: "example.com", Port: sql.NullInt32{Int32: 5000, Valid: true}, Priority: sql.NullInt32{Int32: 10, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}, TTL: 300}, errInvalidHost},
		{Record{Type: RecordTypeSRV, Name: "_ssh._udp.xn--c1yn36f", Value: "test1.dev.", Port: sql.NullInt32{Int32: 5000, Valid: true}, Priority: sql.NullInt32{Int32: 10, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}, TTL: 300}, nil},
		{Record{Type: RecordTypeSRV, Name: "test1.", Value: "example.com.", Port: sql.NullInt32{Int32: 5000, Valid: true}, Priority: sql.NullInt32{Int32: 10, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}, TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeSRV, Name: ".test1", Value: "example.com.", Port: sql.NullInt32{Int32: 5000, Valid: true}, Priority: sql.NullInt32{Int32: 10, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}, TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeSRV, Name: "www.-test1", Value: "example.com.", Port: sql.NullInt32{Int32: 5000, Valid: true}, Priority: sql.NullInt32{Int32: 10, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}, TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeSRV, Name: "-test1", Value: "example.com.", Port: sql.NullInt32{Int32: 5000, Valid: true}, Priority: sql.NullInt32{Int32: 10, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}, TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeSRV, Name: "test1-", Value: "example.com.", Port: sql.NullInt32{Int32: 5000, Valid: true}, Priority: sql.NullInt32{Int32: 10, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}, TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeSRV, Name: "_xmpp-server._tcp.a.long.sub.domain", Value: "userGroups.google.com.", Port: sql.NullInt32{Int32: 5000, Valid: true}, Priority: sql.NullInt32{Int32: 10, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}, TTL: 300}, nil},
		{Record{Type: RecordTypeSRV, Name: "_xmpp-server._tcp.a.long.sub.domain", Value: "userGroups.google.com.", Port: sql.NullInt32{Int32: 5000, Valid: true}, Priority: sql.NullInt32{Int32: 10, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}, TTL: -300}, errInvalidTTL},
		{Record{Type: RecordTypeSRV, Name: "_xmpp-server._tcp.a.long.sub.domain", Value: "userGroups.google.com.", Port: sql.NullInt32{Int32: 5000, Valid: true}, Priority: sql.NullInt32{Int32: 10, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}, TTL: -300}, errInvalidTTL},
		{Record{Type: RecordTypeSRV, Name: "_xmpp-server._tcp.a.long.sub.domain", Value: "userGroups.google.com.", Port: sql.NullInt32{Int32: 5000, Valid: true}, Priority: sql.NullInt32{Int32: -3, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}, TTL: 300}, errInvalidPriority},
		{Record{Type: RecordTypeSRV, Name: "_xmpp-server._tcp.a.long.sub.domain", Value: "userGroups.google.com.", Port: sql.NullInt32{Int32: -1, Valid: true}, Priority: sql.NullInt32{Int32: 10, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}, TTL: 300}, errInvalidPort},
		{Record{Type: RecordTypeSRV, Name: "_xmpp-server._tcp.a.long.sub.domain", Value: "userGroups.google.com.", Port: sql.NullInt32{Int32: 100000, Valid: true}, Priority: sql.NullInt32{Int32: 10, Valid: true}, Weight: sql.NullInt32{Int32: 100, Valid: true}, TTL: 300}, errInvalidPort},
		{Record{Type: RecordTypeSRV, Name: "_xmpp-server._tcp.a.long.sub.domain", Value: "userGroups.google.com.", Port: sql.NullInt32{Int32: 5000, Valid: true}, Priority: sql.NullInt32{Int32: 10, Valid: true}, Weight: sql.NullInt32{Int32: -100, Valid: true}, TTL: 300}, errInvalidWeight},
	}

	for i, table := range tables {
		i := i
		table := table
		name := fmt.Sprintf("Validate SRV record Name: %s, Host: %s TTL: %d", table.x.Name, table.x.Value, table.x.TTL)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := table.x.Validate()
			if err != table.n {
				t.Errorf("[%d] validation failed, got: %v, want: %v,", i, err, table.n)
			}
		})
	}
}

func TestRecordValidateTypeTXT(t *testing.T) {
	tables := []struct {
		x Record
		n error
	}{
		{Record{Type: RecordTypeTXT}, errMissingName},
		{Record{Type: RecordTypeTXT, Value: "example.com.", TTL: 300}, errMissingName},
		{Record{Type: RecordTypeTXT, Name: "example", TTL: 300}, errMissingText},
		{Record{Type: RecordTypeTXT, Name: "@", Value: "example.com."}, errMissingTTL},
		{Record{Type: RecordTypeTXT, Name: "@", Value: "example.com.", TTL: 300}, nil},
		{Record{Type: RecordTypeTXT, Name: "@", Value: "example.com", TTL: 300}, nil},
		{Record{Type: RecordTypeTXT, Name: "@", Value: "10.2.1.400", TTL: 300}, nil},
		{Record{Type: RecordTypeTXT, Name: "test1", Value: "asdf2.", TTL: 300}, nil},
		{Record{Type: RecordTypeTXT, Name: "test1", Value: "test1.dev", TTL: 300}, nil},
		{Record{Type: RecordTypeTXT, Name: "test1", Value: "si4guthe7Baish4eeke1uax1ewooci8faip3bie0tio6ooPhafohP3ooshish2ahbah2Yi3su6choh5ja4einekelohngie5oosu2Am4phahyuyia2Cuz1eiPh8sheit0einguaT4Vuhohviengoox4Faiviose2od2chua0Cee8yahpiaCh9rahghiewee4as6phaecheipoofaiChihai2mait6aichiomohghohfayohyei6chaishech8oox", TTL: 300}, errLengthExceededText},
		{Record{Type: RecordTypeTXT, Name: "test1.", Value: "example.com.", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeTXT, Name: ".test1", Value: "example.com.", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeTXT, Name: ".test1", Value: "example.com.", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeTXT, Name: "-test1", Value: "example.com.", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeTXT, Name: "test1-", Value: "example.com.", TTL: 300}, errInvalidName},
		{Record{Type: RecordTypeTXT, Name: "sub.test1", Value: "userGroups.google.com.", TTL: 300}, nil},
		{Record{Type: RecordTypeTXT, Name: "sub.test1", Value: "userGroups.google.com.", TTL: -300}, errInvalidTTL},
	}

	for i, table := range tables {
		i := i
		table := table
		name := fmt.Sprintf("Validate TXT record Name: %s, Text: %s TTL: %d", table.x.Name, table.x.Value, table.x.TTL)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := table.x.Validate()
			if err != table.n {
				t.Errorf("[%d] validation failed, got: %v, want: %v,", i, err, table.n)
			}
		})
	}
}

func TestRecordValidateInvalidType(t *testing.T) {
	tables := []struct {
		x Record
		n error
	}{
		{Record{}, errMissingType},
		{Record{Type: "1234"}, errUnknownType},
		{Record{Type: "NOTATYPE"}, errUnknownType},
	}

	for i, table := range tables {
		i := i
		table := table
		name := fmt.Sprintf("Validate invalid type: '%s'", table.x.Type)
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			err := table.x.Validate()
			if err != table.n {
				t.Errorf("[%d] validation failed, got: %v, want: %v,", i, err, table.n)
			}
		})
	}
}
