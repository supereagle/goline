package server

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/supereagle/jenkins-pipeline/config"
	"github.com/supereagle/jenkins-pipeline/pipeline"
)

type Server struct {
	router *mux.Router
	pm     *pipeline.Manager
}

func NewServer(cfg *config.Config) (server *Server, err error) {
	pm, err := pipeline.NewPipelineManager(cfg)
	if err != nil {
		return nil, fmt.Errorf("Fail to create the server as %s", err.Error())
	}

	server = &Server{
		router: mux.NewRouter(),
		pm:     pm,
	}
	server.registerRoutes()
	return
}

func (server *Server) registerRoutes() {
	router := server.router
	router.Path("/pipelines").Methods("POST").HandlerFunc(server.createPipeline)
	router.Path("/pipelines/{pipelinename}").Methods("PUT").HandlerFunc(server.updatePipeline)
	router.Path("/pipelines/{pipelinename}").Methods("DELETE").HandlerFunc(server.deletePipeline)
}

func (server *Server) createPipeline(resp http.ResponseWriter, req *http.Request) {
	log.Infoln("Create Pipeline!")
}

func (server *Server) updatePipeline(resp http.ResponseWriter, req *http.Request) {
	log.Infoln("Update Pipeline!")
}

func (server *Server) deletePipeline(resp http.ResponseWriter, req *http.Request) {
	log.Infoln("Delete Pipeline!")
}

func (server *Server) Start() error {
	log.Infoln("Start the server!")
	return http.ListenAndServe(":8080", server.router)
}
