package api

type ProjectType string
type Stage string

const (
	// Project types
	MAVEN  ProjectType = "maven"
	GRADLE             = "gradle"
	SHELL              = "shell"
	BATCH              = "batch"

	// Pipeline stages
	COMPILE Stage = "compile"
	UT            = "unit_test"
	BUILD         = "build"
	DEPLOY        = "deploy"
)

var (
	JDK_PATH = map[string]string{
		"jdk1.6": "/usr/lib/jvm/java-1.6.0",
		"jdk1.7": "/usr/lib/jvm/java-1.7.0",
		"jdk1.8": "/usr/lib/jvm/java-1.8.0",
	}
)

type Pipeline struct {
	Name          string         `json:"name,omitemtpy"`
	NodeLabel     string         `json:"node_label,omitempty"`
	Jdk           string         `json:"jdk,omitempty"`
	Repo          *Repo          `json:"repo,omitempty"`
	PeriodTrigger *PeriodTrigger `json:"period_trigger,omitempty"`
	ProjectType   ProjectType    `json:"type,omitemtpy"`
	Project       interface{}    `json:"project,omitempty"`
	Stages        []Stage        `json:"stages,omitemtpy"`
}

type Repo struct {
	RepoPath string `json:"repo_path,omitempty"`
	Branch   string `json:"branch,omitempty"`
}

type PeriodTrigger struct {
	Skipped  bool   `json:"skipped"`
	Strategy string `json:"strategy,omitempty"`
}

type ScriptProject struct {
	Compile  *ScriptCompile  `json:"compile,omitemtpy"`
	UnitTest *ScriptUnitTest `json:"unit_test,omitempty"`
	Build    *ScriptBuild    `json:"build,omitemtpy"`
}

type ScriptCompile struct {
	Command string `json:"command,omitempty"`
}

type ScriptUnitTest struct {
	Command        string `json:"command,omitempty"`
	TestReportPath string `json:"test_report_path,omitempty"`
}

type ScriptBuild struct {
	Command string `json:"command,omitempty"`
}

type MavenProject struct {
	RootPom  string         `json:"root_pom,omitempty"`
	Options  string         `json:"options,omitempty"`
	UnitTest *MavenUnitTest `json:"unit_test,omitempty"`
}

type MavenUnitTest struct {
	TestReportPath string `json:"test_report_path,omitempty"`
}

type GradleProject struct {
	Options  string          `json:"options,omitempty"`
	UnitTest *GradleUnitTest `json:"unit_test,omitempty"`
}

type GradleUnitTest struct {
	TestReportPath string `json:"test_report_path,omitempty"`
}
