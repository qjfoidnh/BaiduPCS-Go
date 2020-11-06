package prealloc

import (
	"fmt"
)

type (
	PreAllocError struct {
		ProcName string
		Err      error
	}
)

func (pe *PreAllocError) Error() string {
	if pe.Err == nil {
		return "<nil>"
	}
	return fmt.Sprintf("%s error: %s\n", pe.ProcName, pe.Err)
}
