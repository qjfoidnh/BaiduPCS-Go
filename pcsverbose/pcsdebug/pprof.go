package pcsdebug

import (
	"net/http"
	_ "net/http/pprof"
)

func StartPprofListen() {
	http.ListenAndServe("0.0.0.0:6060", nil)
}
