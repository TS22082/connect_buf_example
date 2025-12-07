package internal

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func TestHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	projectId, exists := vars["id"]

	if !exists {
		http.Error(w, "Wrong input", http.StatusExpectationFailed)
		return
	}

	if err := json.NewEncoder(w).Encode(map[string]interface{}{
		"msg":       fmt.Sprintf("ProjectId: %s", projectId),
		"timestamp": time.Now().UTC(),
	}); err != nil {
		http.Error(w, "Can not send response", http.StatusExpectationFailed)
		return
	}
}
