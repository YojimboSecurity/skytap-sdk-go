package skytap

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	exampleVMRequest = `{
		"template_id": "%d",
    		"vm_ids": [
        		"%d"
    	]
	}`
)

func TestCreateVM(t *testing.T) {
	exampleCreateVMResponse := string(readTestFile(t, "createVMResponse.json"))

	response := fmt.Sprintf(exampleCreateVMResponse, 123, 123, 456)
	request := fmt.Sprintf(exampleVMRequest, 42, 43)

	skytap, hs, handler := createClient(t)
	defer hs.Close()

	*handler = func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/configurations/123", req.URL.Path, "Bad path")
		assert.Equal(t, "PUT", req.Method, "Bad method")

		body, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err, "Bad request body")
		assert.JSONEq(t, request, string(body), "Bad request body")

		io.WriteString(rw, response)
	}
	opts := &CreateVMRequest{
		TemplateID: "42",
		VMID:       []string{"43"},
	}

	environment, err := skytap.VMs.Create(context.Background(), "123", opts)
	assert.Nil(t, err, "Bad API method")

	var environmentExpected Environment
	err = json.Unmarshal([]byte(response), &environmentExpected)
	assert.Equal(t, environmentExpected, *environment, "Bad environment")
}

func TestReadVM(t *testing.T) {
	exampleVMResponse := string(readTestFile(t, "VMResponse.json"))

	response := fmt.Sprintf(exampleVMResponse, 456)

	skytap, hs, handler := createClient(t)
	defer hs.Close()

	*handler = func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/v2/configurations/123/vms/456", req.URL.Path, "Bad path")
		assert.Equal(t, "GET", req.Method, "Bad method")

		io.WriteString(rw, response)
	}

	vm, err := skytap.VMs.Get(context.Background(), "123", "456")
	assert.Nil(t, err, "Bad API method")

	var vmExpected VM
	err = json.Unmarshal([]byte(response), &vmExpected)
	assert.Equal(t, vmExpected, *vm, "Bad VM")
}

func TestUpdateVM(t *testing.T) {
	exampleVMResponse := string(readTestFile(t, "VMResponse.json"))

	response := fmt.Sprintf(exampleVMResponse, 456)

	skytap, hs, handler := createClient(t)
	defer hs.Close()

	var vm VM
	json.Unmarshal([]byte(response), &vm)
	*vm.Name = "updated vm"

	bytes, err := json.Marshal(&vm)
	assert.Nil(t, err, "Bad vm")

	*handler = func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/v2/configurations/123/vms/456", req.URL.Path, "Bad path")
		assert.Equal(t, "PUT", req.Method, "Bad method")

		body, err := ioutil.ReadAll(req.Body)
		assert.Nil(t, err, "Bad request body")
		assert.JSONEq(t, `{"name": "updated vm"}`, string(body), "Bad request body")

		io.WriteString(rw, string(bytes))
	}

	opts := &UpdateVMRequest{
		Name: strToPtr(*vm.Name),
	}
	vmUpdate, err := skytap.VMs.Update(context.Background(), "123", "456", opts)
	assert.Nil(t, err, "Bad API method")

	assert.Equal(t, vm, *vmUpdate, "Bad vm")
}

func TestDeleteVM(t *testing.T) {
	skytap, hs, handler := createClient(t)
	defer hs.Close()

	*handler = func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/v2/configurations/123/vms/456", req.URL.Path, "Bad path")
		assert.Equal(t, "DELETE", req.Method, "Bad method")
	}

	err := skytap.VMs.Delete(context.Background(), "123", "456")
	assert.Nil(t, err, "Bad API method")
}

func TestListVMs(t *testing.T) {
	exampleVMResponse := string(readTestFile(t, "VMResponse.json"))

	response := fmt.Sprintf(exampleVMResponse, 456)

	skytap, hs, handler := createClient(t)
	defer hs.Close()

	*handler = func(rw http.ResponseWriter, req *http.Request) {
		assert.Equal(t, "/v2/configurations/123/vms", req.URL.Path, "Bad path")
		assert.Equal(t, "GET", req.Method, "Bad method")

		io.WriteString(rw, fmt.Sprintf(`[%+v]`, response))
	}

	result, err := skytap.VMs.List(context.Background(), "123")
	assert.Nil(t, err, "Bad API method")

	var found = false
	for _, vm := range result.Value {
		if *vm.Name == "test vm" {
			found = true
			break
		}
	}
	assert.True(t, found, "VM not found")
}

func readTestFile(t *testing.T, name string) []byte {
	path := filepath.Join("testdata", name) // relative path
	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	return bytes
}
