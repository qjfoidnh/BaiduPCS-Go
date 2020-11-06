package netdisksign

import (
	"crypto/sha1"
	"encoding/hex"
	"github.com/iikira/BaiduPCS-Go/pcsutil/cachepool"
	"github.com/iikira/BaiduPCS-Go/pcsutil/converter"
	"strconv"
	"time"
)

type (
	LocateDownloadSign struct {
		Time   int64
		Rand   string
		DevUID string
	}
)

func NewLocateDownloadSign(uid uint64, bduss string) *LocateDownloadSign {
	return NewLocateDownloadSignWithTimeAndDevUID(time.Now().Unix(), DevUID(bduss), uid, bduss)
}

func NewLocateDownloadSignWithTimeAndDevUID(timeunix int64, devuid string, uid uint64, bduss string) *LocateDownloadSign {
	l := &LocateDownloadSign{
		Time:   timeunix,
		DevUID: devuid,
	}
	l.Sign(uid, bduss)
	return l
}

func (s *LocateDownloadSign) Sign(uid uint64, bduss string) {
	randSha1 := sha1.New()
	bdussSha1 := sha1.New()
	bdussSha1.Write(converter.ToBytes(bduss))
	sha1ResHex := cachepool.RawMallocByteSlice(40)
	hex.Encode(sha1ResHex, bdussSha1.Sum(nil))
	randSha1.Write(sha1ResHex)
	uidStr := strconv.FormatUint(uid, 10)
	randSha1.Write(converter.ToBytes(uidStr))
	randSha1.Write([]byte{'\x65', '\x62', '\x72', '\x63', '\x55', '\x59', '\x69', '\x75', '\x78', '\x61', '\x5a', '\x76', '\x32', '\x58', '\x47', '\x75', '\x37', '\x4b', '\x49', '\x59', '\x4b', '\x78', '\x55', '\x72', '\x71', '\x66', '\x6e', '\x4f', '\x66', '\x70', '\x44', '\x46'})
	timeStr := strconv.FormatInt(s.Time, 10)
	randSha1.Write(converter.ToBytes(timeStr))
	randSha1.Write(converter.ToBytes(s.DevUID))
	hex.Encode(sha1ResHex, randSha1.Sum(nil))
	s.Rand = converter.ToString(sha1ResHex)
}

func (s *LocateDownloadSign) URLParam() string {
	return "time=" + strconv.FormatInt(s.Time, 10) + "&rand=" + s.Rand + "&devuid=" + s.DevUID
}
