package workflows

import (
	"os"
	"reshell/pkg/config"
	"testing"
)

func TestWorkflowSaveLoadDelete(t *testing.T) {
	tempHome, err := os.MkdirTemp("", "reshell-test-home-*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempHome)

	oldHome := os.Getenv("HOME")
	os.Setenv("HOME", tempHome)
	defer os.Setenv("HOME", oldHome)

	wfName := "test-wf"
	steps := []config.WorkflowStep{
		{Command: "echo 'Step 1'", Dir: "~", Comment: "First step"},
		{Command: "echo 'Step 2'", Dir: "~/projects", Comment: "Second step"},
	}

	err = AddOrUpdate(wfName, "This is a test workflow", steps)
	if err != nil {
		t.Fatalf("AddOrUpdate failed: %v", err)
	}

	wf, err := Get(wfName)
	if err != nil {
		t.Fatalf("Get failed: %v", err)
	}

	if wf.Name != wfName || len(wf.Steps) != 2 {
		t.Errorf("Workflow mismatch: %+v", wf)
	}

	if wf.Steps[0].Command != "echo 'Step 1'" {
		t.Errorf("Step 0 mismatch: %+v", wf.Steps[0])
	}

	err = Remove(wfName)
	if err != nil {
		t.Fatalf("Remove failed: %v", err)
	}

	_, err = Get(wfName)
	if err == nil {
		t.Errorf("Expected workflow to be deleted, but it was found")
	}
}
