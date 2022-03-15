package pkg

import (
	"encoding/json"
	"io/ioutil"
	admiss "k8s.io/api/admission/v1beta1"
	"testing"
)

func TestUnmarshalJson(t *testing.T) {
	bytes, err := ioutil.ReadFile("../deploy/TestData")
	if err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(string(bytes))
	req := &admiss.AdmissionReview{}
	if err := json.Unmarshal(bytes, req); err != nil {
		t.Error(err.Error())
		return
	}
	t.Log(*req)
}

func TestMap(t *testing.T) {
	maptest := make(map[string]string)
	maptest["11"] = "1212"
}
