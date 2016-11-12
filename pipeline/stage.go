package pipeline

import (
	"strings"

	"github.com/supereagle/jenkins-pipeline/api"
)

type StageGenerator interface {
	GenerateCompileStage() string
	GenerateUnitTestStage() string
	GenerateBuildStage() string
}

type MavenPiplineStageGenerator struct {
	ProjectConfig api.MavenProject
}

func (generator *MavenPiplineStageGenerator) GenerateCompileStage() string {
	project := generator.ProjectConfig

	stageTmpl := strings.NewReplacer("${maven.rootpom}", project.RootPom,
		"${mvn.options}", project.Options).Replace(MAVEN_COMPILE_STAGE)

	return stageTmpl
}

func (generator *MavenPiplineStageGenerator) GenerateUnitTestStage() string {
	project := generator.ProjectConfig
	stage := project.UnitTest

	stageTmpl := strings.NewReplacer("${maven.rootpom}", project.RootPom,
		"${mvn.options}", project.Options,
		"${test.report.path}", stage.TestReportPath).Replace(MAVEN_UNIT_TEST_STAGE)

	return stageTmpl
}

func (generator *MavenPiplineStageGenerator) GenerateBuildStage() string {
	project := generator.ProjectConfig

	stageTmpl := strings.NewReplacer("${maven.rootpom}", project.RootPom,
		"${mvn.options}", project.Options).Replace(MAVEN_BUILD_STAGE)

	return stageTmpl
}

type GradlePiplineStageGenerator struct {
	ProjectConfig api.GradleProject
}

func (generator *GradlePiplineStageGenerator) GenerateCompileStage() string {
	project := generator.ProjectConfig

	stageTmpl := strings.Replace(GRADLE_COMPILE_STAGE, "${gradle.gradleOpts}", project.Options, 1)

	return stageTmpl
}

func (generator *GradlePiplineStageGenerator) GenerateUnitTestStage() string {
	project := generator.ProjectConfig
	stage := project.UnitTest

	stageTmpl := strings.NewReplacer("${gradle.gradleOpts}", project.Options,
		"${test.report.path}", stage.TestReportPath).Replace(GRADLE_UNIT_TEST_STAGE)

	return stageTmpl
}

func (generator *GradlePiplineStageGenerator) GenerateBuildStage() string {
	project := generator.ProjectConfig

	stageTmpl := strings.Replace(GRADLE_BUILD_STAGE, "${gradle.gradleOpts}", project.Options, 1)

	return stageTmpl
}

type ScriptPiplineStageGenerator struct {
	ProjectConfig api.ScriptProject
	ProjectType   api.ProjectType
}

func (generator *ScriptPiplineStageGenerator) GenerateCompileStage() string {
	project := generator.ProjectConfig
	stage := project.Compile

	stageTmpl := SCRIPT_COMPILE_STAGE
	if generator.ProjectType == api.BATCH {
		stageTmpl = strings.Replace(stageTmpl, "sh", "bat", 1)
	}

	stageTmpl = strings.NewReplacer("${script.compile.command}", stage.Command).Replace(stageTmpl)

	return stageTmpl
}

func (generator *ScriptPiplineStageGenerator) GenerateUnitTestStage() string {
	project := generator.ProjectConfig
	stage := project.UnitTest

	stageTmpl := SCRIPT_UNIT_TEST_STAGE
	if generator.ProjectType == api.BATCH {
		stageTmpl = strings.Replace(stageTmpl, "sh", "bat", 1)
	}

	stageTmpl = strings.NewReplacer("${script.ut.command}", stage.Command,
		"${test.report.path}", stage.TestReportPath).Replace(stageTmpl)

	return stageTmpl
}

func (generator *ScriptPiplineStageGenerator) GenerateBuildStage() string {
	project := generator.ProjectConfig
	stage := project.Build

	stageTmpl := SCRIPT_BUILD_STAGE
	if generator.ProjectType == api.BATCH {
		stageTmpl = strings.Replace(stageTmpl, "sh", "bat", 1)
	}

	stageTmpl = strings.NewReplacer("${script.build.command}", stage.Command).Replace(stageTmpl)

	return stageTmpl
}
