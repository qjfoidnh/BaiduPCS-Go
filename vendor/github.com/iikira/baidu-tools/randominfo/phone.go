package randominfo

// SumIMEI 根据key计算出imei
func SumIMEI(key string) uint64 {
	var hash uint64 = 53202347234687234
	for k := range key {
		hash += (hash << 5) + uint64(key[k])
	}
	hash %= uint64(1e15)
	if hash < 1e14 {
		hash += 1e14
	}
	return hash
}

// GetPhoneModel 根据key, 从PhoneModelDataBase中取出手机型号
func GetPhoneModel(key string) string {
	if len(PhoneModelDataBase) <= 0 {
		return "S3"
	}
	var hash uint64 = 2134
	for k := range key {
		hash += (hash << 4) + uint64(key[k])
	}
	hash %= uint64(len(PhoneModelDataBase))
	return PhoneModelDataBase[int(hash)]
}
