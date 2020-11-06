package checksum

import (
	"errors"
)

var (
	ErrFileIsNil            = errors.New("file is nil")
	ErrChecksumWriteStop    = errors.New("checksum write stop")
	ErrChecksumWriteAllStop = errors.New("checksum write all stop")
)
