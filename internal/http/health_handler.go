package httpapi

import (
	"github.com/go-chi/chi"
	"github.com/nelkinda/health-go"
	"github.com/nelkinda/health-go/checks/sysinfo"
	"github.com/nelkinda/health-go/checks/uptime"
)

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
