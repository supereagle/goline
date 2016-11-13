# Jenkins-Pipeline API

- [Pipelines](#pipelines)
  - [Create](#create-pipeline)
  - [Update](#update-pipeline)
  - [Delete](#delete-pipeline)

## Pipelines

### Create Pipeline

#### POST /pipelines

#### Description

The POST route for the pipelines creates the Jenkins pipeline according to the configures from the request body.
In the request body, the configure of `project` is different for different type of project.
Some stage can be skipped in `stages`, and the skipped stages' configures are not needed in `project`.

Supported project types:
- [Maven](#maven-pipeline)
- [Gradle](#gradle-pipeline)
- [Shell/Batch](#script-pipeline)

#### Maven Pipeline

##### Example Request

```json
POST http://localhost:8080/pipelines  HTTP/1.1

{
	"name": "maven-pipeline",
	"node_label": "maven-slave",
	"jdk": "jdk1.7",
	"repo": {
		"repo_path": "https://github.com/supereagle/jenkins-pipeline.git",
		"branch": "master"
	},
	"type": "maven",
	"project": {
		"root_pom": "pom.xml",
		"options": "-Ptest",
		"unit_test": {
			"test_report_path": "target/surefire-reports"
		}
	},
	"stages": ["compile", "unit_test", "build"]
}
```

##### Example Response

```json
HTTP/1.1 201 Created
Content-Type: application/json

{
  "code": 201,
  "status": "Created",
  "json_object": {
    "name": "maven-pipeline",
    "node_label": "maven-slave",
    "jdk": "jdk1.7",
    "repo": {
      "repo_path": "https://github.com/supereagle/jenkins-pipeline.git",
      "branch": "master"
    },
    "type": "maven",
    "project": {
      "root_pom": "pom.xml",
      "options": "-Ptest",
      "unit_test": {
        "test_report_path": "target/surefire-reports"
      }
    },
    "stages": [
      "compile",
      "unit_test",
      "build"
    ]
  }
}
```

#### Gradle Pipeline

##### Example Request

```json
POST http://localhost:8080/pipelines  HTTP/1.1

{
	"name": "gradle-pipeline",
	"node_label": "gradle-slave",
	"jdk": "jdk1.8",
	"repo": {
		"repo_path": "https://github.com/supereagle/jenkins-pipeline.git",
		"branch": "dev"
	},
	"type": "gradle",
	"project": {
		"options": "-Pqa",
		"unit_test": {
			"test_report_path": "**/build/test-results/*.xml"
		}
	},
	"stages": ["unit_test", "build"]
}
```

##### Example Response

```json
HTTP/1.1 201 Created
Content-Type: application/json

{
  "code": 201,
  "status": "Created",
  "json_object": {
    "name": "gradle-pipeline",
    "node_label": "gradle-slave",
    "jdk": "jdk1.8",
    "repo": {
      "repo_path": "https://github.com/supereagle/jenkins-pipeline.git",
      "branch": "dev"
    },
    "type": "gradle",
    "project": {
      "options": "-Pqa",
      "unit_test": {
        "test_report_path": "**/build/test-results/*.xml"
      }
    },
    "stages": [
      "unit_test",
      "build"
    ]
  }
}
```

#### Script Pipeline

Script Pipeline includes both Shell script on Linux and Batch script on Windows. Their configures in request body are the same, only the `type` is different. Shell pipeline will use `sh` and batch pipeline will use `bat` to run scripts.

##### Example Request

```json
POST http://localhost:8080/pipelines  HTTP/1.1

{
	"name": "shell-pipeline",
	"node_label": "shell-slave",
	"jdk": "jdk1.7",
	"repo": {
		"repo_path": "https://github.com/supereagle/jenkins-pipeline.git",
		"branch": "master"
	},
	"type": "shell",
	"project": {
		"compile": {
			"command": "./compile.sh"
		},
		"unit_test": {
			"command": "./unit_test.sh",
			"test_report_path": "**/build/test-results/*.xml"
		},
		"build": {
			"command": "./build.sh"
		}
	},
	"stages": ["compile", "unit_test", "build"]
}
```

##### Example Response

```json
HTTP/1.1 201 Created
Content-Type: application/json

{
	"name": "shell-pipeline",
	"node_label": "shell-slave",
	"jdk": "jdk1.7",
	"repo": {
		"repo_path": "https://github.com/supereagle/jenkins-pipeline.git",
		"branch": "master"
	},
	"type": "shell",
	"project": {
		"compile": {
			"command": "./compile.sh"
		},
		"unit_test": {
			"command": "./unit_test.sh",
			"test_report_path": "**/build/test-results/*.xml"
		},
		"build": {
			"command": "./build.sh"
		}
	},
	"stages": ["compile", "unit_test", "build"]
}
```


#### PUT /pipelines/`:pipelinename`

#### Description

The PUT route for the pipelines updates the Jenkins pipeline specified in the REST path with the configures from the request body.
The request body is the same as that of creating pipeline, except that the `name` is NOT needed as it is specified in the REST path.

#### Example Request

**Changes**

* JDK: From JDK1.7 to JDK1.8
* Maven options: From `-Ptest` to `-Pstaging`
* Stages: Skip unit test stage

```json
PUT http://localhost:8080/pipelines/maven-pipeline  HTTP/1.1

{
	"node_label": "maven-slave",
	"jdk": "jdk1.8",
	"repo": {
		"repo_path": "https://github.com/supereagle/jenkins-pipeline.git",
		"branch": "master"
	},
	"type": "maven",
	"project": {
		"root_pom": "pom.xml",
		"options": "-Pstaging"
	},
	"stages": ["compile", "build"]
}
```

#### Example Response

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 200,
  "status": "OK",
  "json_object": {
    "name": "maven-pipeline",
    "node_label": "maven-slave",
    "jdk": "jdk1.8",
    "repo": {
      "repo_path": "https://github.com/supereagle/jenkins-pipeline.git",
      "branch": "master"
    },
    "type": "maven",
    "project": {
      "root_pom": "pom.xml",
      "options": "-Pstaging"
    },
    "stages": [
      "compile",
      "build"
    ]
  }
}
```

### Delete Pipeline

#### DELETE /pipelines/`:pipelinename`

#### Description

The DELETE route for the pipelines deletes the Jenkins pipeline specified in the REST path.

#### Example Request

```json
DELETE http://localhost:8080/pipelines/maven-pipeline  HTTP/1.1
```

#### Example Response

```json
HTTP/1.1 200 OK
Content-Type: application/json

{
  "code": 200,
  "status": "OK"
}
```