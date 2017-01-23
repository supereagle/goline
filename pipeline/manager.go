package pipeline

import (
	"fmt"
	"net/http"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/bndr/gojenkins"
	"github.com/supereagle/goline/api"
	"github.com/supereagle/goline/config"
)

type Manager struct {
	Jenkins      *gojenkins.Jenkins
	credentialId string
}

func NewPipelineManager(cfg *config.Config) (mgr *Manager, err error) {
	if len(cfg.JenkinsServer) == 0 {
		return nil, fmt.Errorf("The Jenkins server url should not be empty")
	}

	// Create the Jenkins Instance
	jenkins, err := gojenkins.CreateJenkins(cfg.JenkinsServer, cfg.JenkinsUser, cfg.JenkinsPassword).Init()
	if err != nil {
		return nil, fmt.Errorf("Fail to create Jenkins Instance as %s", err.Error())
	}

	mgr = &Manager{
		Jenkins:      jenkins,
		credentialId: cfg.JenkinsCredentialId,
	}

	return

}

// Create Creates the pipeline according to the pipeline config
func (mgr *Manager) Create(pl *api.Pipeline) error {
	// Generate the pipeline job config
	jobCfg, err := generatePipelineJobConfig(pl, mgr.credentialId)
	if err != nil {
		err = fmt.Errorf("Fail to generate pipeline config as %s", err.Error())
		log.Errorln(err.Error())
		return err
	}

	// Create the pipeline job
	_, err = mgr.Jenkins.CreateJob(jobCfg, pl.Name)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	return nil
}

// Update Updates the pipeline according to the pipeline config
func (mgr *Manager) Update(pl *api.Pipeline) error {
	// Check the existence of the pipeline job
	job, err := mgr.getJob(pl.Name)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}

	// Generate the pipeline job config
	jobCfg, err := generatePipelineJobConfig(pl, mgr.credentialId)
	if err != nil {
		err = fmt.Errorf("Fail to generate pipeline config as %s", err.Error())
		log.Errorln(err.Error())
		return err
	}

	// Update the pipeline job
	err = job.UpdateConfig(jobCfg)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	return nil
}

// Delete Deletes the pipeline according to the pipeline name
func (mgr *Manager) Delete(plName string) error {
	// Check the existence of the pipeline job
	job, err := mgr.getJob(plName)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}

	// Delete the pipeline job
	_, err = job.Delete()
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	return nil
}

// Perform Performs the pipeline with the perform parameters
func (mgr *Manager) Perform(plName string, pParams *api.PerformParams) error {
	// Check the existence of the pipeline job
	job, err := mgr.getJob(plName)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}

	params := make(map[string]string)
	params["branch"] = pParams.Branch
	params["performPhases"] = pParams.PerformPhases

	// Invoke the pipeline job with params
	//_, err = job.Invoke(nil, false, params, "", "")
	_, err = job.InvokeSimple(params)
	if err != nil {
		log.Errorln(err.Error())
		return err
	}
	return nil
}

//getJob Gets the specified pipeline job, return error if not exists
func (mgr *Manager) getJob(plName string) (*gojenkins.Job, error) {
	job, err := mgr.Jenkins.GetJob(plName)
	if err != nil {
		if strings.Contains(err.Error(), string(http.StatusNotFound)) {
			err = fmt.Errorf("The pipeline %s does not exist", plName)
			return nil, err
		}
		err = fmt.Errorf("Fail to get the pipeline %s as %s", plName, err.Error())
		return nil, err
	}

	return job, nil
}
