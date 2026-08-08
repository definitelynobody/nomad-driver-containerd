package containerd

import (
	"net/http"
	_ "net/http/pprof"
	"os"

	log "github.com/hashicorp/go-hclog"
)

// PprofAddrEnvVar names the listen address for the driver's pprof endpoint,
// for example "127.0.0.1:6060". Unset means no listener.
const PprofAddrEnvVar = "CONTAINERD_DRIVER_PPROF_ADDR"

// servePprof exposes net/http/pprof when PprofAddrEnvVar is set. The driver
// runs as a go-plugin child whose stderr belongs to the Nomad agent, so a
// SIGQUIT stack dump is lost and killing the process destroys the state
// worth reading. This makes a stuck DestroyTask inspectable in place:
//
//	curl "http://127.0.0.1:6060/debug/pprof/goroutine?debug=2"
//
// The variable is read from the Nomad agent's own environment, since the
// agent is what launches this plugin.
func servePprof(logger log.Logger) {
	addr := os.Getenv(PprofAddrEnvVar)
	if addr == "" {
		return
	}

	logger.Info("serving pprof", "addr", addr)
	go func() {
		if err := http.ListenAndServe(addr, nil); err != nil {
			logger.Error("pprof listener stopped", "error", err)
		}
	}()
}
