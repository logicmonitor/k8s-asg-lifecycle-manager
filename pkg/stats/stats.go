package stats

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

const statsPort = ":8080"

// Server stats and healthcheck server
type Server struct{}

// Start the statsserver
func (s *Server) Start() error {
	h := http.NewServeMux()
	h.HandleFunc("/", s.healthCheck)

	srv := &http.Server{
		Addr:           statsPort,
		Handler:        h,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Infof("Stats server started on localhost:" + statsPort)
	log.Fatal(srv.ListenAndServe())
	return nil
}

func (s *Server) healthCheck(
	w http.ResponseWriter,
	req *http.Request,
) {
	w.WriteHeader(http.StatusOK)
}
