package pipeline

import (
	"fmt"

	"github.com/bndr/gojenkins"
	"github.com/supereagle/jenkins-pipeline/api"
	"github.com/supereagle/jenkins-pipeline/config"
)

type Manager struct {
	Jenkins *gojenkins.Jenkins
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
		Jenkins: jenkins,
	}

	return

}

// Create Creates the pipeline according to the pipeline config
func (mgr *Manager) Create(pl *api.Pipeline) error {
	return fmt.Errorf("no implement!")
}
