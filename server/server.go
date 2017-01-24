package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/supereagle/goline/api"
	"github.com/supereagle/goline/config"
	"github.com/supereagle/goline/pipeline"
	httputil "github.com/supereagle/goline/utils/http"
	jsonutil "github.com/supereagle/goline/utils/json"
)

const DefaultSwaggerPath = "./swagger.json"

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

	// register the pipeline handlers
	server.registerRoutes()

	// register the swagger handler
	server.registerSwaggerHandler()

	return
}

func (server *Server) registerRoutes() {
	router := server.router
	router.Path("/pipelines").Methods("POST").HandlerFunc(server.createPipeline)
	router.Path("/pipelines/{pipelinename}").Methods("PUT").HandlerFunc(server.updatePipeline)
	router.Path("/pipelines/{pipelinename}").Methods("DELETE").HandlerFunc(server.deletePipeline)
	router.Path("/pipelines/performance/{pipelinename}").Methods("PUT").HandlerFunc(server.performPipeline)
}

// createPipeline swagger:route POST /pipelines pipelines createPipeline
//
// Creates a pipeline.
//
// Responses:
//    default: genericErrorResponse
//        201: pipelineResponse
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

// updatePipeline swagger:route PUT /pipelines/{pipelinename} pipelines updatePipeline
//
// Updates the configure for a pipeline.
//
// Responses:
//    default: genericErrorResponse
//        200: pipelineResponse
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

// deletePipeline swagger:route DELETE /pipelines/{pipelinename} pipelines deletePipeline
//
// Deletes a pipeline.
//
// Responses:
//    default: genericErrorResponse
//        200: noObjectResponse
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

// performPipeline swagger:route PUT /pipelines/performance/{pipelinename} pipelines performPipeline
//
// Performs a pipeline.
//
// Responses:
//    default: genericErrorResponse
//        200: noObjectResponse
func (server *Server) performPipeline(resp http.ResponseWriter, req *http.Request) {
	defer req.Body.Close()
	plName := mux.Vars(req)["pipelinename"]
	params := &api.PerformParams{}
	err := jsonutil.Unmarshal2JsonObj(req.Body, params)
	if err != nil {
		err = fmt.Errorf("Bad request. Can't parse the request body to a json object as %s", err.Error())
		log.Errorln(err.Error())
		httputil.WriteResponse(resp, http.StatusInternalServerError, nil, err)
		return
	}
	log.Infof("Perform Pipeline %s with params %v", plName, params)

	err = server.pm.Perform(plName, params)
	if err != nil {
		err = fmt.Errorf("Fail to perform the pipeline %s as %s", plName, err.Error())
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

func (server *Server) registerSwaggerHandler() {
	server.router.HandleFunc("/swagger.json", func(resp http.ResponseWriter, req *http.Request) {
		path := strings.TrimSpace(req.URL.Query().Get("path"))
		if len(path) == 0 {
			path = DefaultSwaggerPath
		}

		// Read and response the swagger.json
		resp.Header().Set("Content-Type", "application/json")
		resp.Header().Set("Access-Control-Allow-Origin", "*")
		resp.Header().Set("Access-Control-Allow-Methods", "GET")

		// The query parameter path must end with swagger.json
		if !strings.HasSuffix(path, "swagger.json") {
			result, _ := json.Marshal(&map[string]string{"error": fmt.Sprintf("Path %s must end with swagger.json.", path)})
			resp.WriteHeader(400)
			resp.Write(result)
			return
		}

		result, err := ioutil.ReadFile(path)
		if err != nil {
			result, _ := json.Marshal(&map[string]string{"error": fmt.Sprintf("Fail to read file %s as %s", path, err.Error())})

			resp.WriteHeader(500)
			resp.Write(result)
			return
		}

		resp.WriteHeader(200)
		resp.Write(result)
	})
}

func parseBody(req *http.Request) (*api.Pipeline, error) {
	defer req.Body.Close()
	pipeline := &api.Pipeline{}
	err := jsonutil.Unmarshal2JsonObj(req.Body, pipeline)
	if err != nil {
		err = fmt.Errorf("Bad request. Can't parse the request body to a json object as %s", err.Error())
		return nil, err
	}

	projectStr, err := jsonutil.Marshal2JsonStr(pipeline.Project)
	if err != nil {
		return nil, err
	}

	projectType := pipeline.ProjectType
	switch projectType {
	case api.SHELL, api.BATCH:
		project := api.ScriptProject{}
		err = jsonutil.UnmarshalJsonStr2Obj(projectStr, &project)
		if err != nil {
			return nil, err
		}
		pipeline.Project = project
	case api.MAVEN:
		project := api.MavenProject{}
		err = jsonutil.UnmarshalJsonStr2Obj(projectStr, &project)
		if err != nil {
			return nil, err
		}
		pipeline.Project = project
	case api.GRADLE:
		project := api.GradleProject{}
		err = jsonutil.UnmarshalJsonStr2Obj(projectStr, &project)
		if err != nil {
			return nil, err
		}
		pipeline.Project = project
	default:
		return nil, fmt.Errorf("The project type %s is not supported", projectType)
	}

	return pipeline, nil
}
