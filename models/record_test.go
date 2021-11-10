package models

import "testing"

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

	for _, table := range tables {
		err := table.x.Validate()

		if err != table.n {
			t.Errorf("validation failed, got: %v, want: %v,", err, table.n)
		}
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

	for _, table := range tables {
		err := table.x.Validate()

		if err != table.n {
			t.Errorf("validation failed, got: %v, want: %v,", err, table.n)
		}
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

	for _, table := range tables {
		err := table.x.Validate()

		if err != table.n {
			t.Errorf("validation failed, got: %v, want: %v,", err, table.n)
		}
	}
}
