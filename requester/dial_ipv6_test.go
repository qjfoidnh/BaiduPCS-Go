package requester

import "testing"

func TestPickAddrPreferIPv6(t *testing.T) {
	v4First := []string{"27.221.82.38", "2409:8c3c:900:10b::78dc:a826"}
	v6First := []string{"2409:8c3c:900:10b::78dc:a826", "27.221.82.38"}
	onlyV4 := []string{"27.221.82.38", "36.110.192.115"}

	// 默认行为：保持系统解析顺序（多数 Linux 下 IPv4 在前）
	SetPreferIPv6(false)
	if got, err := pickAddr(v4First); err != nil || got != "27.221.82.38" {
		t.Fatalf("disabled: want first address, got %q err=%v", got, err)
	}

	// 开启后：存在 AAAA 时选 IPv6
	SetPreferIPv6(true)
	if got, err := pickAddr(v4First); err != nil || got != "2409:8c3c:900:10b::78dc:a826" {
		t.Fatalf("enabled: want IPv6, got %q err=%v", got, err)
	}
	if got, err := pickAddr(v6First); err != nil || got != "2409:8c3c:900:10b::78dc:a826" {
		t.Fatalf("enabled: want IPv6 regardless of order, got %q err=%v", got, err)
	}

	// 无 AAAA 时回退 IPv4，纯 IPv4 域名不受影响
	if got, err := pickAddr(onlyV4); err != nil || got != "27.221.82.38" {
		t.Fatalf("fallback: want IPv4, got %q err=%v", got, err)
	}

	// 空结果返回错误
	if _, err := pickAddr(nil); err == nil {
		t.Fatal("empty address list should return error")
	}

	SetPreferIPv6(false)
}
