package skytap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreateEnvironment(t *testing.T) {
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	requestCounter := 0

	*handler = func(rw http.ResponseWriter, req *http.Request) {

		// ignore user_data requests
		if req.RequestURI == "/v2/configurations/456/user_data.json" {
			_, err := io.WriteString(rw, `{"contents": ""}`)
			assert.NoError(t, err)
			return
		}

		log.Printf("Request: (%d)\n", requestCounter)
		if requestCounter == 0 {
			if req.URL.Path != "/configurations" {
				t.Error("Bad path")
			}
			if req.Method != "POST" {
				t.Error("Bad method")
			}
			body, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			assert.JSONEq(t, fmt.Sprintf(`{"template_id":%q, "project_id":%d, "description":"test environment"}`, "12345", 12345), string(body))
			_, err = io.WriteString(rw, `{"id": "456"}`)
			assert.NoError(t, err)
		} else if requestCounter == 1 {
			assert.Equal(t, "/v2/configurations/456", req.URL.Path, "Bad path")
			assert.Equal(t, http.MethodGet, req.Method, "Bad method")

			_, err := io.WriteString(rw, string(readTestFile(t, "exampleEnvironment.json")))
			assert.NoError(t, err)
		} else if requestCounter == 2 {
			assert.Equal(t, "/v2/configurations/456", req.URL.Path, "Bad path")
			assert.Equal(t, http.MethodGet, req.Method, "Bad method")

			_, err := io.WriteString(rw, string(readTestFile(t, "exampleEnvironment.json")))
			assert.NoError(t, err)
		} else if requestCounter == 3 {
			assert.Equal(t, "/v2/configurations/456", req.URL.Path, "Bad path")
			assert.Equal(t, http.MethodGet, req.Method, "Bad method")

			_, err := io.WriteString(rw, string(readTestFile(t, "exampleEnvironment.json")))
			assert.NoError(t, err)
		} else if requestCounter == 4 {
			if req.URL.Path != "/v2/configurations/456" {
				t.Error("Bad path")
			}
			if req.Method != "PUT" {
				t.Error("Bad method")
			}
			body, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			assert.JSONEq(t, `{"description": "test environment", "runstate":"running"}`, string(body))

			_, err = io.WriteString(rw, string(readTestFile(t, "exampleEnvironment.json")))
			assert.NoError(t, err)

		} else if requestCounter == 5 {
			assert.Equal(t, "/v2/configurations/456", req.URL.Path, "Bad path")
			assert.Equal(t, http.MethodGet, req.Method, "Bad method")

			var envRunning Environment
			err := json.Unmarshal(readTestFile(t, "exampleEnvironment.json"), &envRunning)
			assert.NoError(t, err)
			envRunning.Runstate = environmentRunStateToPtr(EnvironmentRunstateRunning)
			b, err := json.Marshal(&envRunning)
			assert.Nil(t, err)
			_, err = io.WriteString(rw, string(b))
			assert.NoError(t, err)
		}
		requestCounter++
	}

	opts := &CreateEnvironmentRequest{
		TemplateID:  strToPtr("12345"),
		ProjectID:   intToPtr(12345),
		Description: strToPtr("test environment"),
	}

	environment, err := skytap.Environments.Create(context.Background(), opts)

	assert.Nil(t, err)

	var environmentExpected Environment

	err = json.Unmarshal(readTestFile(t, "exampleEnvironment.json"), &environmentExpected)

	assert.Equal(t, environmentExpected, *environment)

	assert.Equal(t, 6, requestCounter)
}

func TestReadEnvironment(t *testing.T) {
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	*handler = func(rw http.ResponseWriter, req *http.Request) {
		// ignore user_data requests
		if req.RequestURI == "/v2/configurations/456/user_data.json" {
			return
		}
		if req.URL.Path != "/v2/configurations/456" {
			t.Error("Bad path")
		}
		if req.Method != "GET" {
			t.Error("Bad method")
		}
		_, err := io.WriteString(rw, string(readTestFile(t, "exampleEnvironment.json")))
		assert.NoError(t, err)
	}

	environment, err := skytap.Environments.Get(context.Background(), "456")

	assert.Nil(t, err)
	var environmentExpected Environment

	err = json.Unmarshal(readTestFile(t, "exampleEnvironment.json"), &environmentExpected)

	assert.Equal(t, environmentExpected, *environment)
}

func TestUpdateEnvironment(t *testing.T) {
	skytap, hs, handler := createClient(t)
	defer hs.Close()
	var environment Environment
	err := json.Unmarshal(readTestFile(t, "exampleEnvironment.json"), &environment)
	assert.NoError(t, err)
	*environment.Description = "updated environment"

	requestCounter := 0

	*handler = func(rw http.ResponseWriter, req *http.Request) {
		// ignore user_data requests
		if req.RequestURI == "/v2/configurations/456/user_data.json" {
			_, err = io.WriteString(rw, `{"contents": ""}`)
			assert.NoError(t, err)
			return
		}

		log.Printf("Request: (%d)\n", requestCounter)
		if requestCounter == 0 {
			assert.Equal(t, "/v2/configurations/456", req.URL.Path, "Bad path")
			assert.Equal(t, http.MethodGet, req.Method, "Bad method")

			_, err := io.WriteString(rw, string(readTestFile(t, "exampleEnvironment.json")))
			assert.NoError(t, err)
		} else if requestCounter == 1 {
			if req.URL.Path != "/v2/configurations/456" {
				t.Error("Bad path")
			}
			if req.Method != "PUT" {
				t.Error("Bad method")
			}
			body, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			assert.JSONEq(t, `{"description": "updated environment"}`, string(body))

			b, err := json.Marshal(&environment)
			assert.Nil(t, err)
			_, err = io.WriteString(rw, string(b))
			assert.NoError(t, err)
		} else if requestCounter == 2 {
			assert.Equal(t, "/v2/configurations/456", req.URL.Path, "Bad path")
			assert.Equal(t, http.MethodGet, req.Method, "Bad method")

			var envRunning Environment
			err := json.Unmarshal(readTestFile(t, "exampleEnvironment.json"), &envRunning)
			assert.NoError(t, err)
			envRunning.Description = strToPtr("updated environment")
			b, err := json.Marshal(&envRunning)
			assert.Nil(t, err)
			_, err = io.WriteString(rw, string(b))
		}
		requestCounter++
	}

	opts := &UpdateEnvironmentRequest{
		Description: strToPtr(*environment.Description),
	}

	environmentUpdate, err := skytap.Environments.Update(context.Background(), "456", opts)

	assert.Nil(t, err)
	assert.Equal(t, environment, *environmentUpdate)

	assert.Equal(t, 3, requestCounter)
}

func TestDeleteEnvironment(t *testing.T) {
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	requestCounter := 0

	*handler = func(rw http.ResponseWriter, req *http.Request) {
		// ignore user_data requests
		if req.RequestURI == "/v2/configurations/456/user_data.json" {
			_, err := io.WriteString(rw, `{"contents": ""}`)
			assert.NoError(t, err)
			return
		}
		log.Printf("Request: (%d)\n", requestCounter)
		if requestCounter == 0 {
			assert.Equal(t, "/v2/configurations/456", req.URL.Path, "Bad path")
			assert.Equal(t, http.MethodGet, req.Method, "Bad method")

			_, err := io.WriteString(rw, string(readTestFile(t, "exampleEnvironment.json")))
			assert.NoError(t, err)
		} else if requestCounter == 1 {
			if req.URL.Path != "/configurations/456" {
				t.Error("Bad path")
			}
			if req.Method != "DELETE" {
				t.Error("Bad method")
			}
		}
		requestCounter++
	}

	err := skytap.Environments.Delete(context.Background(), "456")
	assert.Nil(t, err)
	assert.Equal(t, 2, requestCounter)
}

func TestListEnvironments(t *testing.T) {
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	*handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path != "/v2/configurations" {
			t.Error("Bad path")
		}
		if req.Method != "GET" {
			t.Error("Bad method")
		}
		_, err := io.WriteString(rw, fmt.Sprintf(`[%+v]`, string(readTestFile(t, "exampleEnvironment.json"))))
		assert.NoError(t, err)
	}

	result, err := skytap.Environments.List(context.Background())

	assert.Nil(t, err)

	var found = false
	for _, environment := range result.Value {
		if *environment.Description == "test environment" {
			found = true
			break
		}
	}

	assert.True(t, found)
}

func TestCompareEnvironmentCreateTrue(t *testing.T) {
	exampleEnvironment := readTestFile(t, "exampleEnvironment.json")

	var environment Environment
	err := json.Unmarshal(exampleEnvironment, &environment)
	assert.NoError(t, err)
	opts := CreateEnvironmentRequest{
		TemplateID: strToPtr("12345"),
		ProjectID:  intToPtr(12345),
	}
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	*handler = func(rw http.ResponseWriter, req *http.Request) {
		_, err := io.WriteString(rw, string(exampleEnvironment))
		assert.NoError(t, err)
	}

	message, ok := opts.compareResponse(context.Background(), skytap, &environment, envRunStateNotBusy("123"))
	assert.True(t, ok)
	assert.Equal(t, "", message)
}

func TestCompareEnvironmentCreateFalse(t *testing.T) {
	exampleEnvironment := readTestFile(t, "exampleEnvironment.json")

	var environment Environment
	err := json.Unmarshal(exampleEnvironment, &environment)
	assert.NoError(t, err)
	opts := CreateEnvironmentRequest{
		TemplateID: strToPtr("12345"),
		ProjectID:  intToPtr(12345),
	}

	skytap, hs, handler := createClient(t)
	defer hs.Close()

	*handler = func(rw http.ResponseWriter, req *http.Request) {
		var envRunning Environment
		err := json.Unmarshal(exampleEnvironment, &envRunning)
		assert.NoError(t, err)
		envRunning.Runstate = environmentRunStateToPtr(EnvironmentRunstateBusy)
		b, err := json.Marshal(&envRunning)
		assert.Nil(t, err)
		_, err = io.WriteString(rw, string(b))

		assert.NoError(t, err)
	}
	message, ok := opts.compareResponse(context.Background(), skytap, &environment, envRunStateNotBusy("123"))
	assert.False(t, ok)
	assert.Equal(t, "environment not ready", message)
}

func TestCompareEnvironmentUpdateTrue(t *testing.T) {
	exampleEnvironment := readTestFile(t, "exampleEnvironment.json")

	var environment Environment
	err := json.Unmarshal(exampleEnvironment, &environment)
	assert.NoError(t, err)
	opts := UpdateEnvironmentRequest{
		Description: strToPtr(*environment.Description),
	}
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	*handler = func(rw http.ResponseWriter, req *http.Request) {
		_, err := io.WriteString(rw, string(exampleEnvironment))
		assert.NoError(t, err)
	}

	message, ok := opts.compareResponse(context.Background(), skytap, &environment, envRunStateNotBusy("123"))
	assert.True(t, ok)
	assert.Equal(t, "", message)
}

func TestCompareEnvironmentUpdateFalse(t *testing.T) {
	exampleEnvironment := readTestFile(t, "exampleEnvironment.json")

	var environment Environment
	err := json.Unmarshal(exampleEnvironment, &environment)
	assert.NoError(t, err)
	environment.Runstate = environmentRunStateToPtr(EnvironmentRunstateBusy)
	opts := UpdateEnvironmentRequest{
		Runstate: environmentRunStateToPtr(EnvironmentRunstateStopped),
	}
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	*handler = func(rw http.ResponseWriter, req *http.Request) {
		var envRunning Environment
		err := json.Unmarshal(exampleEnvironment, &envRunning)
		assert.NoError(t, err)
		envRunning.Runstate = environmentRunStateToPtr(EnvironmentRunstateRunning)
		b, err := json.Marshal(&envRunning)
		assert.Nil(t, err)
		_, err = io.WriteString(rw, string(b))

		assert.NoError(t, err)
	}
	message, ok := opts.compareResponse(context.Background(), skytap, &environment, envRunStateNotBusy("123"))
	assert.False(t, ok)
	assert.Equal(t, "environment not ready", message)
}

func TestConfirmNilRoutableAlwaysFalse(t *testing.T) {
	var environment Environment
	err := json.Unmarshal(readTestFile(t, "exampleEnvironment.json"), &environment)
	assert.NoError(t, err)
	environment.Routable = nil

	opts := UpdateEnvironmentRequest{
		Routable: boolToPtr(false),
	}
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	*handler = func(rw http.ResponseWriter, req *http.Request) {
		b, err := json.Marshal(&environment)
		assert.Nil(t, err)
		_, err = io.WriteString(rw, string(b))
		assert.NoError(t, err)
	}

	message, ok := opts.compareResponse(context.Background(), skytap, &environment, envRunStateNotBusy("123"))
	assert.True(t, ok)
	assert.Equal(t, "", message)
}

func TestConfirmCreateEnvironmentCreateTags(t *testing.T) {
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	environmentCreated := false
	environmentUpdated := false
	tagsCreated := false
	*handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/configurations" && req.Method == "POST" { // create
			_, err := io.WriteString(rw, `{"id": "456"}`)
			assert.NoError(t, err)
			environmentCreated = true
		}
		if req.URL.Path == "/v2/configurations/456" && req.Method == "GET" { // get
			io.WriteString(rw, string(readTestFile(t, "exampleEnvironment.json")))
		}
		if req.URL.Path == "/v2/configurations/456" && req.Method == "PUT" {
			io.WriteString(rw, string(readTestFile(t, "exampleEnvironment.json")))
			environmentUpdated = true
		}
		if req.URL.Path == "/v2/configurations/456/tags.json" && req.Method == "PUT" {
			tagsCreated = true
			body, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			assert.JSONEq(t, `[{"value": "foo"}, {"value": "bar"}]`, string(body))
		}
	}

	opts := &CreateEnvironmentRequest{
		TemplateID:  strToPtr("12345"),
		ProjectID:   intToPtr(12345),
		Description: strToPtr("test environment"),
		Tags: []*CreateTagRequest{
			{"foo"},
			{"bar"},
		},
	}
	_, err := skytap.Environments.Create(context.Background(), opts)

	assert.Nil(t, err)
	assert.True(t, environmentCreated)
	assert.True(t, environmentUpdated)
	assert.True(t, tagsCreated)
}

func TestEnvironmentAddTag(t *testing.T) {
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	tagsCreated := false
	*handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/v2/configurations/456/tags.json" && req.Method == "PUT" {
			tagsCreated = true
			body, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			assert.JSONEq(t, `[{"value": "foo"}, {"value": "bar"}]`, string(body))
		}

	}
	tags := []*CreateTagRequest{
		{"foo"},
		{"bar"},
	}

	err := skytap.Environments.CreateTags(context.Background(), "456", tags)
	assert.Nil(t, err)
	assert.True(t, tagsCreated)
}

func TestEnvironmentEmptyListOfTagHasNoEffect(t *testing.T) {
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	tagsCreated := false
	*handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/v2/configurations/456/tags.json" && req.Method == "PUT" {
			tagsCreated = true
		}
	}

	tags := make([]*CreateTagRequest, 0)
	err := skytap.Environments.CreateTags(context.Background(), "456", tags)

	assert.Nil(t, err)
	assert.False(t, tagsCreated)
}

func TestEnvironmentDeleteTag(t *testing.T) {
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	tagDeleted := false
	*handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/v2/configurations/456/tags/42.json" && req.Method == "DELETE" {
			tagDeleted = true
		}
	}

	err := skytap.Environments.DeleteTag(context.Background(), "456", "42")
	assert.Nil(t, err)
	assert.True(t, tagDeleted)
}

func TestConfirmCreateEnvironmentCreateUserData(t *testing.T) {
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	environmentCreated := false
	environmentUpdated := false
	environmentUserData := false
	*handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/configurations" && req.Method == "POST" { // create
			_, err := io.WriteString(rw, `{"id": "456"}`)
			assert.NoError(t, err)
			environmentCreated = true
		}
		if req.URL.Path == "/v2/configurations/456" && req.Method == "GET" { // get
			io.WriteString(rw, string(readTestFile(t, "exampleEnvironment.json")))
		}
		if req.URL.Path == "/v2/configurations/456" && req.Method == "PUT" {
			io.WriteString(rw, string(readTestFile(t, "exampleEnvironment.json")))
			environmentUpdated = true
		}
		if req.URL.Path == "/v2/configurations/456/user_data.json" && req.Method == "PUT" {
			environmentUserData = true
			body, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			assert.JSONEq(t, `{"contents": "echo \\proc\\cpu_info"}`, string(body))
		}

	}

	opts := &CreateEnvironmentRequest{
		TemplateID:  strToPtr("12345"),
		ProjectID:   intToPtr(12345),
		Description: strToPtr("test environment"),
		UserData:    strToPtr("echo \\proc\\cpu_info"),
	}
	_, err := skytap.Environments.Create(context.Background(), opts)

	assert.Nil(t, err)
	assert.True(t, environmentCreated)
	assert.True(t, environmentUpdated)
	assert.True(t, environmentUserData)
}

func TestGetEnvironmentReadUserData(t *testing.T) {
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	readUserData := false
	*handler = func(rw http.ResponseWriter, req *http.Request) {

		if req.URL.Path == "/v2/configurations/456/user_data.json" {
			readUserData = true
			assert.Equal(t, req.Method, "GET")
			_, err := io.WriteString(rw, `{"contents": "dataexample"}`)
			assert.NoError(t, err)
		} else {
			_, err := io.WriteString(rw, string(readTestFile(t, "exampleEnvironment.json")))
			assert.NoError(t, err)
		}
	}

	environment, err := skytap.Environments.Get(context.Background(), "456")
	assert.Nil(t, err)
	assert.True(t, readUserData)
	assert.Equal(t, "dataexample", *environment.UserData)
}

func TestConfirmUserDataUpdate(t *testing.T) {
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	environmentUserData := false
	*handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/v2/configurations/456/user_data.json" && req.Method == "PUT" {
			environmentUserData = true
			body, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			assert.JSONEq(t, `{"contents": "echo \\proc\\cpu_info"}`, string(body))
		}

	}
	err := skytap.Environments.UpdateUserData(context.Background(), "456", strToPtr("echo \\proc\\cpu_info"))
	assert.Nil(t, err)
	assert.True(t, environmentUserData)
}

func TestConfirmCreateEnvironmentCreateLabels(t *testing.T) {
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	environmentCreated := false
	environmentUpdated := false
	labelsCreated := false
	*handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/configurations" && req.Method == "POST" { // create
			_, err := io.WriteString(rw, `{"id": "456"}`)
			assert.NoError(t, err)
			environmentCreated = true
		}
		if req.URL.Path == "/v2/configurations/456" && req.Method == "GET" { // get
			io.WriteString(rw, string(readTestFile(t, "exampleEnvironment.json")))
		}
		if req.URL.Path == "/v2/configurations/456" && req.Method == "PUT" {
			io.WriteString(rw, string(readTestFile(t, "exampleEnvironment.json")))
			environmentUpdated = true
		}
		if req.URL.Path == "/v2/configurations/456/labels.json" && req.Method == "PUT" {
			labelsCreated = true
			body, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			assert.JSONEq(t, `[{"label_category": "foo", "value": "bar"},
									    {"label_category": "foozz", "value": "barzz"}]`, string(body))
		}
	}

	opts := &CreateEnvironmentRequest{
		TemplateID:  strToPtr("12345"),
		ProjectID:   intToPtr(12345),
		Description: strToPtr("test environment"),
		Labels: []*CreateLabelRequest{
			{Category: strToPtr("foo"), Value: strToPtr("bar")},
			{Category: strToPtr("foozz"), Value: strToPtr("barzz")},
		},
	}
	_, err := skytap.Environments.Create(context.Background(), opts)

	assert.Nil(t, err)
	assert.True(t, environmentCreated)
	assert.True(t, environmentUpdated)
	assert.True(t, labelsCreated)
}

func TestEnvironmentAddLabel(t *testing.T) {
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	labelsCreated := false
	*handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/v2/configurations/456/labels.json" && req.Method == "PUT" {
			labelsCreated = true
			body, err := ioutil.ReadAll(req.Body)
			assert.Nil(t, err)
			assert.JSONEq(t, `[{"label_category": "foo", "value": "bar"}]`, string(body))
		}

	}
	labels := []*CreateLabelRequest{
		{Category: strToPtr("foo"), Value: strToPtr("bar")},
	}
	err := skytap.Environments.CreateLabels(context.Background(), "456", labels)
	assert.Nil(t, err)
	assert.True(t, labelsCreated)
}

func TestEnvironmentDeleteLabel(t *testing.T) {
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	labelDeleted := false
	*handler = func(rw http.ResponseWriter, req *http.Request) {
		if req.URL.Path == "/v2/configurations/456/labels/42.json" && req.Method == "DELETE" {
			labelDeleted = true
		}
	}

	err := skytap.Environments.DeleteLabel(context.Background(), "456", "42")
	assert.Nil(t, err)
	assert.True(t, labelDeleted)
}
