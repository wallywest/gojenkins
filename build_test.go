package gojenkins

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestIsBuildQueued(t *testing.T) {
	setupJenkins(t)

	jobs, _ := jenkins.GetAllJobs()
	if assert.NotNil(t, jobs) {
		for _, job := range jobs {
			fmt.Printf("job url is: %s\n", job.Raw.URL)
			fmt.Printf("job lastbuild is: %v\n", job.Raw.LastBuild)
			fmt.Printf("job inqueue is: %v\n", job.Raw.InQueue)
		}
	}
}

func TestCreateBuilds(t *testing.T) {
	setupJenkins(t)

	jobs, _ := jenkins.GetAllJobs()
	if assert.NotNil(t, jobs) {
		for _, item := range jobs {
			item.InvokeSimple(map[string]string{"param1": "param1"})
			item.Poll()
			isQueued, _ := item.IsQueued()
			assert.Equal(t, true, isQueued)

			time.Sleep(10 * time.Second)
			builds, _ := item.GetAllBuildIds()

			assert.True(t, (len(builds) > 0))
		}
	}
}

func TestParseBuildHistory(t *testing.T) {
	r, err := os.Open("_tests/build_history.txt")
	if err != nil {
		panic(err)
	}
	history := parseBuildHistory(r)
	assert.True(t, len(history) == 3)
}

func TestGetAllBuilds(t *testing.T) {
	builds, _ := jenkins.GetAllBuildIds("Job1_test")
	for _, b := range builds {
		build, _ := jenkins.GetBuild("Job1_test", b.Number)
		assert.Equal(t, "SUCCESS", build.GetResult())
	}
	assert.Equal(t, 1, len(builds))
}

func TestBuildMethods(t *testing.T) {
	job, _ := jenkins.GetJob("Job1_test")
	build, _ := job.GetLastBuild()
	params := build.GetParameters()
	assert.Equal(t, "params1", params[0].Name)
}
