package health

import (
	"github.com/mike-kimani/whitepointinventory/pkg/jsonresponses"
	"net/http"
)

func HandlerHealth(w http.ResponseWriter, r *http.Request) {
	jsonresponses.RespondWithJSON(w, 200, struct{}{})
}
