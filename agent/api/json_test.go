package api

import (
	"encoding/json"
	"reflect"
	"testing"
)

type testTaskStatus struct {
	SomeStatus TaskStatus `json:"status"`
}

func TestUnmarshalTaskStatus(t *testing.T) {
	status := TaskStatusNone

	err := json.Unmarshal([]byte(`"RUNNING"`), &status)
	if err != nil {
		t.Error(err)
	}
	if status != TaskRunning {
		t.Error("RUNNING should unmarshal to RUNNING, not " + status.String())
	}

	var test testTaskStatus
	err = json.Unmarshal([]byte(`{"status":"STOPPED"}`), &test)
	if err != nil {
		t.Error(err)
	}
	if test.SomeStatus != TaskStopped {
		t.Error("STOPPED should unmarshal to STOPPED, not " + test.SomeStatus.String())
	}
}

type testContainerStatus struct {
	SomeStatus ContainerStatus `json:"status"`
}

func TestUnmarshalContainerStatus(t *testing.T) {
	status := ContainerStatusNone

	err := json.Unmarshal([]byte(`"RUNNING"`), &status)
	if err != nil {
		t.Error(err)
	}
	if status != ContainerRunning {
		t.Error("RUNNING should unmarshal to RUNNING, not " + status.String())
	}

	var test testContainerStatus
	err = json.Unmarshal([]byte(`{"status":"STOPPED"}`), &test)
	if err != nil {
		t.Error(err)
	}
	if test.SomeStatus != ContainerStopped {
		t.Error("STOPPED should unmarshal to STOPPED, not " + test.SomeStatus.String())
	}
}

type testContainerOverrides struct {
	SomeContainerOverrides ContainerOverrides `json:"overrides"`
}

type testContainerOverrideInput struct {
	SomeContainerOverrides string `json:"overrides"`
}

func TestUnmarshalContainerOverrides(t *testing.T) {
	overrides := &ContainerOverrides{}

	err := json.Unmarshal([]byte(`{"command": ["sh", "-c", "sleep 300"]}`), &overrides)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(overrides.Command, &[]string{"sh", "-c", "sleep 300"}) {
		t.Error("Unmarshalled wrong result", overrides.Command)
	}

	var overrides3 testContainerOverrides
	err = json.Unmarshal([]byte(`{"overrides":{"command":["sh", "-c", "sleep 15"]}}`), &overrides3)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(overrides3.SomeContainerOverrides.Command, &[]string{"sh", "-c", "sleep 15"}) {
		t.Error("unmarshalled wrong result", overrides3)
	}

	overrides2 := ContainerOverrides{Command: &[]string{"sh", "-c", "sleep 1"}}

	strOverrides, err := json.Marshal(overrides2)
	if err != nil {
		t.Error(err)
	}
	input := testContainerOverrideInput{SomeContainerOverrides: string(strOverrides)}
	strInput, err := json.Marshal(input)
	if err != nil {
		t.Error(err)
	}

	var overrides4 testContainerOverrides
	err = json.Unmarshal(strInput, &overrides4)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(overrides4.SomeContainerOverrides.Command, &[]string{"sh", "-c", "sleep 1"}) {
		t.Error("Unmarshalled wrong result", overrides4)
	}

	// Test that marshalling an unknown key fails hard
	var overrides5 testContainerOverrides
	err = json.Unmarshal([]byte(`{"overrides":{"command":["ash","-c","sleep 1"],"containerPlanet":"mars"}}`), &overrides5)
	if err == nil {
		t.Error("No error on unknown json field containerPlanet")
	}

	// Test the same for the string
	err = json.Unmarshal([]byte(`{"overrides":"{\"command\":[\"ash\",\"-c\",\"sleep 1\"],\"containerPlanet\":\"mars\"}"}`), &overrides5)
	if err == nil {
		t.Error("No error for unknown json field in string, containerPlanet")
	}

	// Now error cases
	err = json.Unmarshal([]byte(`{"overrides":"a string that can't be json unmarshalled }{"}`), &overrides5)
	if err == nil {
		t.Error("No error when unmarshalling an invalid json string")
	}

	err = json.Unmarshal([]byte(`{"overrides": ["must be a string or object"] }`), &overrides5)
	if err == nil {
		t.Error("No error when unmarshalling a really invalid json string")
	}
}
