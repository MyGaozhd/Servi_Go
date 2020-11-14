package t41

import (
	"encoding/json"
	"testing"
)

func Test41_0(t *testing.T) {

	em := new(Employee)

	job := new(JobInfo)
	skills := []string{"go", "java"}
	job.Skills = skills

	ba := BasicInfo{"servi", 31}

	em.BasicInfo = ba
	em.JobInfo = *job
	if v, err := json.Marshal(em); err != nil {
		t.Error(err)
	} else {
		t.Log(string(v))

		ne := new(Employee)
		if err := json.Unmarshal([]byte(v), ne); err != nil {
			t.Error(err)
		} else {
			t.Log(*ne)
		}
	}
}

type BasicInfo struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}

type JobInfo struct {
	Skills []string `json:"skills"`
}

type Employee struct {
	BasicInfo BasicInfo `json:"basic_info"`
	JobInfo   JobInfo   `json:"job_info"`
}
