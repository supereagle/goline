package pipeline

const (
	PIPELINE_JOB_TEMPLATE = `<?xml version='1.0' encoding='UTF-8'?>
<flow-definition plugin="workflow-job@2.3">
  <actions/>
  <description>Pipeline to compile, unit test, build and deploy.</description>
  <keepDependencies>false</keepDependencies>
  <properties>
    <jenkins.plugins.slack.SlackNotifier_-SlackJobProperty plugin="slack@1.8">
      <teamDomain></teamDomain>
      <token></token>
      <room></room>
      <startNotification>false</startNotification>
      <notifySuccess>false</notifySuccess>
      <notifyAborted>false</notifyAborted>
      <notifyNotBuilt>false</notifyNotBuilt>
      <notifyUnstable>false</notifyUnstable>
      <notifyFailure>false</notifyFailure>
      <notifyBackToNormal>false</notifyBackToNormal>
      <notifyRepeatedFailure>false</notifyRepeatedFailure>
      <includeTestSummary>false</includeTestSummary>
      <showCommitList>false</showCommitList>
      <includeCustomMessage>false</includeCustomMessage>
      <customMessage></customMessage>
    </jenkins.plugins.slack.SlackNotifier_-SlackJobProperty>
    <hudson.model.ParametersDefinitionProperty>
      <parameterDefinitions>
        <hudson.model.StringParameterDefinition>
          <name>branch</name>
          <description>The srouce code branch.</description>
          <defaultValue>${project.branch}</defaultValue>
        </hudson.model.StringParameterDefinition>
        <hudson.model.StringParameterDefinition>
          <name>performPhases</name>
          <description>The phase to be performed.</description>
          <defaultValue>${pipeline.perform.phases}</defaultValue>
        </hudson.model.StringParameterDefinition>
      </parameterDefinitions>
    </hudson.model.ParametersDefinitionProperty>
  </properties>
  <definition class="org.jenkinsci.plugins.workflow.cps.CpsFlowDefinition" plugin="workflow-cps@2.9">
    <script>${pipeline.script}</script>
    <sandbox>false</sandbox>
  </definition>
  <triggers/>
</flow-definition>`

	PIPELINE_SCRIPT_TEMPLATE = `
performPhases = "${performPhases}"

node("${pipeline.label.node}") {
	timestamps {
		catchError {
			timeout(time: 1, unit: 'HOURS') {	
				// Checkout the source code
				checkout([$class: 'GitSCM', branches: [[name: '${project.branch}']], userRemoteConfigs: [[credentialsId: '${jenkins.credentialId}', url: '${project.repoPath}']]])
				 sh "git checkout $branch"
				
				withEnv(["WORKSPACE=${pwd()}", "PATH+JAVA=${jdk.version}/bin", "JAVA_HOME=${jdk.version}"]) {
					// Compile Stage
					${pipeline.script.stage.compile}
					
					// Unit Test Stage
					${pipeline.script.stage.unittest}
					
					// Package Stage
					${pipeline.script.stage.build}
					
					// Deploy Stage
					${pipeline.script.stage.deploy}
				}
			}
		}
	}

    // Archive the workspace
    archiveArtifacts artifacts: '**/*', excludes: '**/*.war, **/*.tar.gz, **/*.tgz, **/*.zip'
}
	`

	STAGE_TEMPLATE = `if (performPhases.contains("${pipeline.stage}")) {
						${pipeline.script.stage.function}
					}`

	MAVEN_COMMAND_FUNCTION = `
def mvn(args) {
    sh "/opt/maven/latest/bin/mvn ${args}"
}
	`

	MAVEN_COMPILE_STAGE = `
def compile() {
    stage "Compile"
	
    mvn("-B -f ${maven.rootpom} clean install -e -U -DskipTests=true -Dfindbugs.skip=true ${mvn.options}")
}
	`

	MAVEN_UNIT_TEST_STAGE = `
def unitTest() {
    stage "Unit Test"
	
    mvn("-B -f ${maven.rootpom} clean org.jacoco:jacoco-maven-plugin:0.7.2.201409121644:prepare-agent test -Dfindbugs.skip=true ${mvn.options}")
	
	junit '**/${test.report.path}/TEST-*.xml'
}
	`

	MAVEN_BUILD_STAGE = `
def build() {
    stage "Build"
	
    mvn("-B -f ${maven.rootpom} clean package -e -U -DskipTests=true -Dfindbugs.skip=true ${mvn.options}")
}
	`

	GRADLE_COMPILE_STAGE = `
def compile() {
    stage "Compile"
	
	sh "gradle clean compile -x test -x check ${gradle.gradleOpts}"
}
	`

	GRADLE_UNIT_TEST_STAGE = `
def unitTest() {
    stage "Unit Test"
	
	sh "gradle clean test ${gradle.gradleOpts}"
	
	junit '${test.report.path}'
}
	`

	GRADLE_BUILD_STAGE = `
def build() {
    stage "Build"
	
	sh "gradle clean build ${gradle.gradleOpts} -x test"
}
	`

	SCRIPT_COMPILE_STAGE = `
def compile() {
    stage "Compile"
	
    sh '''${script.compile.command}'''
}
	`

	SCRIPT_UNIT_TEST_STAGE = `
def unitTest() {
    stage "Unit Test"
	
    sh '''${script.ut.command}'''
	
	junit '${test.report.path}'
}
	`

	SCRIPT_BUILD_STAGE = `
def build() {
    stage "Build"
	
    sh '''${script.build.command}'''
}
	`
)
