package server

import (
	"fmt"
	"net/http"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/supereagle/jenkins-pipeline/api"
	"github.com/supereagle/jenkins-pipeline/config"
	"github.com/supereagle/jenkins-pipeline/pipeline"
	httputil "github.com/supereagle/jenkins-pipeline/utils/http"
	"github.com/supereagle/jenkins-pipeline/utils/json"
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
	pipeline, err := parseBody(req)
	if err != nil {
		err = fmt.Errorf("Fail to parse the pipeline config as %s", err.Error())
		log.Errorln(err.Error())
		httputil.WriteResponse(resp, http.StatusInternalServerError, nil, err)
		return
	}
	log.Infof("Create Pipeline %s", pipeline.Name)

	err = server.pm.Create(pipeline)
	if err != nil {
		err = fmt.Errorf("Fail to create the pipeline %s as %s", pipeline.Name, err.Error())
		log.Errorln(err.Error())
		httputil.WriteResponse(resp, http.StatusInternalServerError, nil, err)
		return
	}

	httputil.WriteResponse(resp, http.StatusCreated, pipeline, nil)
}

func (server *Server) updatePipeline(resp http.ResponseWriter, req *http.Request) {
	plName := mux.Vars(req)["pipelinename"]

	pipeline, err := parseBody(req)
	if err != nil {
		err = fmt.Errorf("Fail to parse the pipeline config as %s", err.Error())
		log.Errorln(err.Error())
		httputil.WriteResponse(resp, http.StatusInternalServerError, nil, err)
		return
	}
	pipeline.Name = plName
	log.Infof("Update Pipeline %s", pipeline.Name)

	err = server.pm.Update(pipeline)
	if err != nil {
		err = fmt.Errorf("Fail to update the pipeline %s as %s", pipeline.Name, err.Error())
		log.Errorln(err.Error())
		httputil.WriteResponse(resp, http.StatusInternalServerError, nil, err)
		return
	}

	httputil.WriteResponse(resp, http.StatusOK, pipeline, nil)
}

func (server *Server) deletePipeline(resp http.ResponseWriter, req *http.Request) {
	plName := mux.Vars(req)["pipelinename"]
	log.Infof("Delete Pipeline %s", plName)

	err := server.pm.Delete(plName)
	if err != nil {
		err = fmt.Errorf("Fail to delete the pipeline %s as %s", plName, err.Error())
		log.Errorln(err.Error())
		httputil.WriteResponse(resp, http.StatusInternalServerError, nil, err)
		return
	}

	httputil.WriteResponse(resp, http.StatusOK, nil, nil)
}

func (server *Server) Start() error {
	log.Infoln("Start the server!")
	return http.ListenAndServe(":8080", server.router)
}

func parseBody(req *http.Request) (*api.Pipeline, error) {
	pipeline := &api.Pipeline{}
	err := json.Unmarshal2JsonObj(req.Body, pipeline)
	if err != nil {
		err = fmt.Errorf("Bad request. Can't parse the request body to a json object as %s", err.Error())
		return nil, err
	}

	projectStr, err := json.Marshal2JsonStr(pipeline.Project)
	if err != nil {
		return nil, err
	}

	switch pipeline.ProjectType {
	case api.SHELL, api.BATCH:
		project := api.ScriptProject{}
		err = json.UnmarshalJsonStr2Obj(projectStr, &project)
		if err != nil {
			return nil, err
		}
		pipeline.Project = project
	case api.MAVEN:
		project := api.MavenProject{}
		err = json.UnmarshalJsonStr2Obj(projectStr, &project)
		if err != nil {
			return nil, err
		}
		pipeline.Project = project
	default:
		return nil, fmt.Errorf("The project type %s is not supported", pipeline.ProjectType)
	}

	return pipeline, nil
}
