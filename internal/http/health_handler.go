package httpapi

import (
	"github.com/go-chi/chi"
)
import "github.com/nelkinda/health-go"
import "github.com/nelkinda/health-go/checks/uptime"
import "github.com/nelkinda/health-go/checks/sysinfo"

func RegisterHealthHandlerRoutes(router chi.Router) {
	h := health.New(
		health.Health{
			Version:   "2.0.0",
			ReleaseID: "2.0.0-SNAPSHOT",
		},
		uptime.System(),
		uptime.Process(),
		sysinfo.Health(),
	)
	router.Get("/health", h.Handler)
}
