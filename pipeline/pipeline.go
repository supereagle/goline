package pipeline

import (
	"fmt"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/supereagle/jenkins-pipeline/api"
)

func generatePipelineJobConfig(pipeline *api.Pipeline, credentialId string) (jobCfg string, err error) {
	// Validate the pipeline config
	if ok := ValidatePipeline(pipeline); !ok {
		return "", fmt.Errorf("Pipeline config is not correct")
	}

	pipelineScriptTmpl, err := generatePipelineScriptTmpl(pipeline, credentialId)
	if err != nil {
		err = fmt.Errorf("Fail to generate the script template for pipeline: %s", pipeline.Name)
		fmt.Printf(err.Error())
		return
	}

	jobTmpl := strings.Replace(PIPELINE_JOB_TEMPLATE, "${pipeline.script}", pipelineScriptTmpl, 1)
	jobTmpl = strings.NewReplacer("${pipeline.perform.phases}", convertStagesToString(pipeline.Stages),
		"${project.branch}", pipeline.Repo.Branch).Replace(jobTmpl)

	// Generate the period trigger
	if pipeline.PeriodTrigger != nil && !pipeline.PeriodTrigger.Skipped {
		periodTriggerTmpl := strings.Replace(PIPELINE_TRIGGERS_TEMPLATE, "${period.trigger.strategy}",
			pipeline.PeriodTrigger.Strategy, 1)
		jobTmpl = strings.Replace(jobTmpl, "${pipeline.triggers}", periodTriggerTmpl, 1)
	} else {
		jobTmpl = strings.Replace(jobTmpl, "${pipeline.triggers}", "<triggers/>", 1)
	}

	jobCfg = jobTmpl
	return
}

func generatePipelineScriptTmpl(pipeline *api.Pipeline, credenitalId string) (pipelineScriptTmpl string, err error) {
	scriptTmpl := strings.NewReplacer("${pipeline.label.node}", pipeline.NodeLabel,
		"${jenkins.credentialId}", credenitalId,
		"${project.repoPath}", pipeline.Repo.RepoPath,
		"${jdk.version}", api.JDK_PATH[pipeline.Jdk]).Replace(PIPELINE_SCRIPT_TEMPLATE)

	// Judge the project type
	var stageGenerator StageGenerator
	pType := pipeline.ProjectType
	switch pType {
	case api.MAVEN:
		stageGenerator = &MavenPiplineStageGenerator{pipeline.Project.(api.MavenProject)}

		scriptTmpl += MAVEN_COMMAND_FUNCTION
	case api.GRADLE:
		stageGenerator = &GradlePiplineStageGenerator{pipeline.Project.(api.GradleProject)}
	case api.SHELL, api.BATCH:
		stageGenerator = &ScriptPiplineStageGenerator{
			ProjectConfig: pipeline.Project.(api.ScriptProject),
			ProjectType:   pType,
		}

		if pType == api.BATCH {
			fmt.Println(scriptTmpl)
			scriptTmpl = strings.Replace(scriptTmpl, " sh ", " bat ", -1)
		}
	default:
		err = fmt.Errorf("The project type %v is not supported", pType)
		return
	}

	stages := pipeline.Stages

	// Add the compile stage
	if containStage(stages, api.COMPILE) {
		scriptTmpl = strings.Replace(scriptTmpl, "${pipeline.script.stage.compile}", generatePipelineStageTmpl(api.COMPILE), 1)
		stageTmpl := stageGenerator.GenerateCompileStage()

		scriptTmpl += stageTmpl
	} else {
		scriptTmpl = strings.Replace(scriptTmpl, "${pipeline.script.stage.compile}", "// Skipped", 1)
	}

	// Add the unit test stage
	if containStage(stages, api.UT) {
		scriptTmpl = strings.Replace(scriptTmpl, "${pipeline.script.stage.unittest}", generatePipelineStageTmpl(api.UT), 1)
		stageTmpl := stageGenerator.GenerateUnitTestStage()

		scriptTmpl += stageTmpl
	} else {
		scriptTmpl = strings.Replace(scriptTmpl, "${pipeline.script.stage.unittest}", "// Skipped", 1)
	}

	// Add the build stage
	if containStage(stages, api.BUILD) {
		scriptTmpl = strings.Replace(scriptTmpl, "${pipeline.script.stage.build}", generatePipelineStageTmpl(api.BUILD), 1)
		stageTmpl := stageGenerator.GenerateBuildStage()

		scriptTmpl += stageTmpl
	} else {
		scriptTmpl = strings.Replace(scriptTmpl, "${pipeline.script.stage.build}", "// Skipped", 1)
	}

	pipelineScriptTmpl = scriptTmpl
	return
}

func generatePipelineStageTmpl(stage api.Stage) string {
	stageTmpl := strings.Replace(STAGE_TEMPLATE, "${pipeline.stage}", string(stage), 1)
	switch stage {
	case api.COMPILE:
		stageTmpl = strings.Replace(stageTmpl, "${pipeline.script.stage.function}", "compile()", 1)
	case api.UT:
		stageTmpl = strings.Replace(stageTmpl, "${pipeline.script.stage.function}", "unitTest()", 1)
	case api.BUILD:
		stageTmpl = strings.Replace(stageTmpl, "${pipeline.script.stage.function}", "build()", 1)
	case api.DEPLOY:
		stageTmpl = strings.Replace(stageTmpl, "${pipeline.script.stage.function}", "deploy()", 1)
	}

	return stageTmpl
}

func convertStagesToString(stages []api.Stage) (stageStr string) {
	stageArrays := []string{}

	for _, stage := range stages {
		stageArrays = append(stageArrays, string(stage))
	}

	stageStr = strings.Join(stageArrays, ",")

	return
}

func containStage(stages []api.Stage, desiredStage api.Stage) bool {
	for _, stage := range stages {
		if stage == desiredStage {
			return true
		}
	}

	return false
}

// validatePipeline Validates the pipeline config.
// Returns true if correct, or false if wrong.
func ValidatePipeline(pipeline *api.Pipeline) bool {
	// Check the JDK
	if _, ok := api.JDK_PATH[pipeline.Jdk]; !ok {
		log.Errorf("The jdk version %s is not supported", pipeline.Jdk)
		return false
	}

	// Check the period trigger
	if pipeline.PeriodTrigger != nil && !pipeline.PeriodTrigger.Skipped {
		// TODO (robin) Check strategy to follow the syntax of cron
		if strings.TrimSpace(pipeline.PeriodTrigger.Strategy) == "" {
			return false
		}
	}

	// Check the repo
	repo := pipeline.Repo
	if repo == nil {
		log.Errorln("The source code repo is not specified")
		return false
	}
	// TODO (robin) Check the repo path and branch pattern
	if len(repo.RepoPath) == 0 || len(repo.Branch) == 0 {
		log.Errorln("The source code repo path or branch is empty")
		return false
	}

	projectType := pipeline.ProjectType
	switch projectType {
	case api.SHELL, api.BATCH:
		project, ok := pipeline.Project.(api.ScriptProject)
		if !ok {
			log.Errorf("Project config is not compatiable with project type %s", projectType)
			return false
		}
		//TODO (robin) Check according to the pipeline.stages
		if project.Build == nil || project.Compile == nil {
			log.Errorln("The project stages' configs are empty")
			return false
		}
	case api.MAVEN:
		project, ok := pipeline.Project.(api.MavenProject)
		if !ok {
			log.Errorf("Project config is not compatiable with project type %s", projectType)
			return false
		}
		if len(project.RootPom) == 0 {
			log.Errorln("The maven root pom is not specified")
			return false
		}
	case api.GRADLE:
		_, ok := pipeline.Project.(api.GradleProject)
		if !ok {
			log.Errorf("Project config is not compatiable with project type %s", projectType)
			return false
		}
	default:
		log.Errorf("The project type %s is not supported", pipeline.ProjectType)
		return false
	}

	return true
}
