package models

import (
	"fmt"
	"testing"
)

func TestRegexIPv4Address(t *testing.T) {
	tables := []struct {
		x string
		n bool
	}{
		{"1.2.3.4", true},
		{"192.168.10.10", true},
		{"0.0.0.0", true},
		{"255.255.255.255", true},
		{"400.98.15.10", false},
		{"158.228.500.65", false},
		{"11.52.7", false},
		{"69", false},
		{"ip.address", false},
		{"4087:47e4::9e7d", false},
		{"60a1:69b2:9f2e:2f5c:cf32:01e2:4a2a:fb6d", false},
		{"::1234:5678", false},
	}

	for i, table := range tables {
		i := i
		table := table
		name := fmt.Sprintf("[%d] Testing %s", i, table.x)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match := reIPv4Address.MatchString(table.x)
			if match != table.n {
				t.Errorf("[%d] regex match on %s failed, got: %v, want: %v,", i, table.x, match, table.n)
			}
		})
	}
}

func TestRegexIPv6Address(t *testing.T) {
	tables := []struct {
		x string
		n bool
	}{
		{"::1", true},
		{"bc24:ca41:bf48:1b2b:dc22:6d41:46ff:1ec7", true},
		{"f887:3aa8:24e1:a707:bda0:3447:57b3:0086", true},
		{"1bb5:e87f:2ec5:1a9d:6768:5186:3ebb:95ca", true},
		{"2607:5f90:3f40:b948:8d14:d53a:ba83:4560", true},
		{"23d6:3310:0889:3640:d446:4888:72e0:cdb2", true},
		{"db58:6c89:19d2:fd00:5a15:8d82:6e67:f294", true},
		{"c2a4:4584:65ca:b7fb:3c:18b9:0d5b:6567", true},
		{"60a1:69b2:9f2e:2f5c:cf32:01e2:4a2a:fb6d", true},
		{"1e61:4949:28bd:a3d5:05d2:41f4:3504:920b", true},
		{"71ea:0ae2:6322:9e5e:b354:7d5d:0062:53b9", true},
		{"f977:e5cd:1088:4f3f:c191:ea73:378c:8241", true},
		{"d157:59cb:2224:914f:3ec5:3ee2:1159:684a", true},
		{"f63:ad21:24ba:81c8:e0a:4394:fee6:8148", true},
		{"4087:47e4::9e7d", true},
		{"::1234:5678", true},
		{"2001:db8::", true},
		{"400.98.15.10", false},
		{"158.228.500.65", false},
		{"11.52.7", false},
		{"69", false},
		{"ip.address", false},
		{"google.com", false},
	}

	for i, table := range tables {
		i := i
		table := table
		name := fmt.Sprintf("[%d] Testing %s", i, table.x)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match := reIPv6Address.MatchString(table.x)
			if match != table.n {
				t.Errorf("[%d] regex match on %s failed, got: %v, want: %v,", i, table.x, match, table.n)
			}
		})
	}
}

func TestRegexMXDomain(t *testing.T) {
	tables := []struct {
		x string
		n bool
	}{
		{"google.com.", false},
		{"asdf2.", false},
		{"xn--c1yn36f.", false},
		{"blog.xn--c1yn36f.", false},
		{"userGroups.example.com.", false},
		{"google.com", true},
		{"asdf2", true},
		{"xn--c1yn36f", true},
		{"blog.xn--c1yn36f", true},
		{"userGroups.example.com", true},
		{".xn--c1yn36f.", false},
		{"what?.", false},
		{"google", true},
		{"@", false},
		{"_ssh.tcp.host1", false},
		{"_ssh.tcp", false},
		{"_ssh.udp.xn--c1yn36f", false},
		{"_ssh._tcp.host1", false},
		{"_ssh._tcp", false},
		{"_ssh._udp.xn--c1yn36f", false},
		{"_xmpp-client._tcp", false},
		{"_xmpp-server._tcp", false},
		{"_xmpp-client._tcp.chat", false},
		{"_xmpp-server._tcp.a.long.sub.domain", false},
	}

	for i, table := range tables {
		i := i
		table := table
		name := fmt.Sprintf("[%d] Testing %s", i, table.x)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match := reMXDomain.MatchString(table.x)
			if match != table.n {
				t.Errorf("[%d] regex match on %s failed, got: %v, want: %v,", i, table.x, match, table.n)
			}
		})
	}
}

func TestRegexNSDomain(t *testing.T) {
	tables := []struct {
		x string
		n bool
	}{
		{"@", false},
		{"google.com", true},
		{"asdf2", true},
		{"foo.bar.baz", true},
		{"foo-bar-baz", true},
		{"xn--c1yn36f", true},
		{"xn--", false},
		{"--c1yn36f", false},
		{"c1yn36f-", false},
		{".xn--c1yn36f", false},
		{"what?", false},
		{"google", true},
		{"_ssh.tcp.host1", false},
		{"_ssh.tcp", false},
		{"_ssh.udp.xn--c1yn36f", false},
		{"_ssh._tcp.host1", false},
		{"_ssh._tcp", false},
		{"_ssh._udp.xn--c1yn36f", false},
		{"_xmpp-client._tcp", false},
		{"_xmpp-server._tcp", false},
		{"_xmpp-client._tcp.chat", false},
		{"_xmpp-server._tcp.a.long.sub.domain", false},
	}

	for i, table := range tables {
		i := i
		table := table
		name := fmt.Sprintf("[%d] Testing %s", i, table.x)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match := reNSDomain.MatchString(table.x)
			if match != table.n {
				t.Errorf("[%d] regex match on %s failed, got: %v, want: %v,", i, table.x, match, table.n)
			}
		})
	}
}

func TestRegexSRVDomain(t *testing.T) {
	tables := []struct {
		x string
		n bool
	}{
		{"@", false},
		{"google.com", false},
		{"asdf2", false},
		{"foo.bar.baz", false},
		{"foo-bar-baz", false},
		{"xn--c1yn36f", false},
		{"xn--", false},
		{"--c1yn36f", false},
		{"c1yn36f-", false},
		{".xn--c1yn36f", false},
		{"what?", false},
		{"google", false},
		{"_ssh.tcp.host1", false},
		{"_ssh.tcp", false},
		{"_ssh.udp.xn--c1yn36f", false},
		{"_ssh._tcp.host1", true},
		{"_ssh._tcp", true},
		{"_ssh._udp.xn--c1yn36f", true},
		{"_xmpp-client._tcp", true},
		{"_xmpp-server._tcp", true},
		{"_xmpp-client._tcp.chat", true},
		{"_xmpp-server._tcp.a.long.sub.domain", true},
	}

	for i, table := range tables {
		i := i
		table := table
		name := fmt.Sprintf("[%d] Testing %s", i, table.x)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match := reSRVDomain.MatchString(table.x)
			if match != table.n {
				t.Errorf("[%d] regex match on %s failed, got: %v, want: %v,", i, table.x, match, table.n)
			}
		})
	}
}

func TestRegexSubDomain(t *testing.T) {
	tables := []struct {
		x string
		n bool
	}{
		{"@", true},
		{"google.com", true},
		{"asdf2", true},
		{"foo.bar.baz", true},
		{"foo-bar-baz", true},
		{"xn--c1yn36f", true},
		{"xn--", false},
		{"--c1yn36f", false},
		{"c1yn36f-", false},
		{".xn--c1yn36f", false},
		{"what?", false},
		{"google", true},
		{"_ssh.tcp.host1", false},
		{"_ssh.tcp", false},
		{"_ssh.udp.xn--c1yn36f", false},
		{"_ssh._tcp.host1", false},
		{"_ssh._tcp", false},
		{"_ssh._udp.xn--c1yn36f", false},
		{"_xmpp-client._tcp", false},
		{"_xmpp-server._tcp", false},
		{"_xmpp-client._tcp.chat", false},
		{"_xmpp-server._tcp.a.long.sub.domain", false},
	}

	for i, table := range tables {
		i := i
		table := table
		name := fmt.Sprintf("[%d] Testing %s", i, table.x)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match := reSubDomain.MatchString(table.x)
			if match != table.n {
				t.Errorf("[%d] regex match on %s failed, got: %v, want: %v,", i, table.x, match, table.n)
			}
		})
	}
}

func TestRegexTopDomain(t *testing.T) {
	tables := []struct {
		x string
		n bool
	}{
		{"google.com.", true},
		{"asdf2.", true},
		{"xn--c1yn36f.", true},
		{"blog.xn--c1yn36f.", true},
		{"userGroups.example.com.", true},
		{".xn--c1yn36f.", false},
		{"what?.", false},
		{"google", false},
		{"@", false},
		{"_ssh.tcp.host1", false},
		{"_ssh.tcp", false},
		{"_ssh.udp.xn--c1yn36f", false},
		{"_ssh._tcp.host1", false},
		{"_ssh._tcp", false},
		{"_ssh._udp.xn--c1yn36f", false},
		{"_xmpp-client._tcp", false},
		{"_xmpp-server._tcp", false},
		{"_xmpp-client._tcp.chat", false},
		{"_xmpp-server._tcp.a.long.sub.domain", false},
	}

	for i, table := range tables {
		i := i
		table := table
		name := fmt.Sprintf("[%d] Testing %s", i, table.x)
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			match := reTopDomain.MatchString(table.x)
			if match != table.n {
				t.Errorf("[%d] regex match on %s failed, got: %v, want: %v,", i, table.x, match, table.n)
			}
		})
	}
}
