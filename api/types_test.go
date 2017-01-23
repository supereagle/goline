package api_test

import (
	"reflect"
	"testing"

	"github.com/supereagle/goline/api"
	"github.com/supereagle/goline/utils/json"
)

func TestScriptProject(t *testing.T) {
	pipeline := api.Pipeline{
		Name:        "script-pipeline",
		ProjectType: api.SHELL,
		Project: api.ScriptProject{
			Compile: &api.ScriptCompile{
				Command: "Shell compile",
			},
			Build: &api.ScriptBuild{
				Command: "Shell build",
			},
		},
	}

	//Unmarshal pipeline object from string
	pipelineStr := `{
		"name": "script-pipeline",
		"type": "shell",
		"project": {
			"compile": {
				"command": "Shell compile"
			},
			"build": {
				"command": "Shell build"
			}
		}
	}`

	pl := api.Pipeline{}
	err := json.UnmarshalJsonStr2Obj(pipelineStr, &pl)
	if err != nil {
		t.Fatalf("Fail to unmarshal json string as %s", err.Error())
	}

	mpStr, err := json.Marshal2JsonStr(pl.Project)
	if err != nil {
		t.Fatalf("Fail to unmarshal json string as %s", err.Error())
	}

	sp := api.ScriptProject{}
	err = json.UnmarshalJsonStr2Obj(mpStr, &sp)
	if err != nil {
		t.Fatalf("Fail to unmarshal json string as %s", err.Error())
	}
	pl.Project = sp

	if !reflect.DeepEqual(pl, pipeline) {
		t.Fatal("The object unmarshaled from json string is not as desired")
	}
}

func TestMavenProject(t *testing.T) {
	pipeline := api.Pipeline{
		Name:        "maven-pipeline",
		ProjectType: api.MAVEN,
		Project: api.MavenProject{
			RootPom: "pom.xml",
			Options: "-s settings.xml",
		},
	}

	//Unmarshal pipeline object from string
	pipelineStr := `{
		"name": "maven-pipeline",
		"type": "maven",
		"project": {
			"root_pom": "pom.xml",
			"options": "-s settings.xml"
		}
	}`

	pl := api.Pipeline{}
	err := json.UnmarshalJsonStr2Obj(pipelineStr, &pl)
	if err != nil {
		t.Fatalf("Fail to unmarshal json string as %s", err.Error())
	}

	mpStr, err := json.Marshal2JsonStr(pl.Project)
	if err != nil {
		t.Fatalf("Fail to unmarshal json string as %s", err.Error())
	}

	mp := api.MavenProject{}
	err = json.UnmarshalJsonStr2Obj(mpStr, &mp)
	if err != nil {
		t.Fatalf("Fail to unmarshal json string as %s", err.Error())
	}
	pl.Project = mp

	if !reflect.DeepEqual(pl, pipeline) {
		t.Fatal("The object unmarshaled from json string is not as desired")
	}
}
