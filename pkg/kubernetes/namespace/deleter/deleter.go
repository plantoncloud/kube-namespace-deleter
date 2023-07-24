package deleter

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
)

type Deleter interface {
	RunCommand(cmd *exec.Cmd) error
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte) error
	RemoveFile(filename string) error
}

type RealDeleter struct{}

func (e RealDeleter) RunCommand(cmd *exec.Cmd) error {
	return cmd.Run()
}

func (e RealDeleter) ReadFile(filename string) ([]byte, error) {
	return ioutil.ReadFile(filename)
}

func (e RealDeleter) WriteFile(filename string, data []byte) error {
	return ioutil.WriteFile(filename, data, 0644)
}

func (e RealDeleter) RemoveFile(filename string) error {
	return os.Remove(filename)
}

func Delete(namespace string, executor Deleter) error {
	// Step 1: Dump the contents of the namespace in a temporary file called tmp.json
	cmd := exec.Command("kubectl", "get", "namespace", namespace, "-o", "json", ">", "tmp.json")
	err := executor.RunCommand(cmd)
	if err != nil {
		return err
	}

	// Step 2: Edit the temporary file to remove kubernetes from the finalizer array
	data, err := executor.ReadFile("tmp.json")
	if err != nil {
		return err
	}

	var ns map[string]interface{}
	err = json.Unmarshal(data, &ns)
	if err != nil {
		return err
	}

	spec, ok := ns["spec"].(map[string]interface{})
	if !ok {
		return err
	}

	finalizers, ok := spec["finalizers"].([]interface{})
	if !ok {
		return err
	}

	var newFinalizers []interface{}
	for _, finalizer := range finalizers {
		if finalizer != "kubernetes" {
			newFinalizers = append(newFinalizers, finalizer)
		}
	}
	spec["finalizers"] = newFinalizers

	newData, err := json.Marshal(ns)
	if err != nil {
		return err
	}

	err = executor.WriteFile("tmp.json", newData)
	if err != nil {
		return err
	}

	// Step 3: Call the Kubernetes API application/json against the /finalize endpoint for the namespace to update the JSON
	cmd = exec.Command("curl", "-k", "-H", "\"Content-Type: application/json\"", "-X", "PUT", "--data-binary", "@tmp.json", "http://127.0.0.1:8001/api/v1/namespaces/"+namespace+"/finalize")
	err = executor.RunCommand(cmd)
	if err != nil {
		return err
	}

	// Cleanup: Remove the temporary file
	err = executor.RemoveFile("tmp.json")
	if err != nil {
		return err
	}

	return nil
}
