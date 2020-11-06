package getip

import (
	"testing"
)

func TestGetIP(t *testing.T) {
	ipAddr, err := IPInfo(false)
	if err != nil {
		t.Errorf("err: %s\n", err)
		return
	}

	t.Logf("from ipify: %s\n", ipAddr)

	ipAddr, err = IPInfoFromNetease()
	if err != nil {
		t.Errorf("err: %s\n", err)
		return
	}

	t.Logf("from netease: %s\n", ipAddr)
}
