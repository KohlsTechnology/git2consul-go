// +build e2e

package e2e

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
	"testing"
	"time"

	"github.com/hashicorp/consul/api"
)

func TestE2E(t *testing.T) {
	// Setup local Consul test server
	consulCmd := exec.Command("consul", "agent", "-dev")
	consulPipeReader, consulPipeWriter := io.Pipe()
	consulScanner := bufio.NewScanner(consulPipeReader)
	consulTee := io.MultiWriter(consulPipeWriter, os.Stderr) // dumping the output on Stderr can be useful when debugging this test

	consulCmd.Stdout = consulTee
	consulCmd.Stderr = consulTee
	if err := consulCmd.Start(); err != nil {
		t.Fatalf("failed to start local consul server: %s", err)
	}
	defer consulCmd.Process.Signal(syscall.SIGTERM)

	consulDone := make(chan error, 0)
	go func() {
		consulDone <- consulCmd.Wait()
	}()

	err := waitForString(consulScanner, 20*time.Second, "==> Consul agent running!")
	if err != nil {
		t.Fatal("initialization of consul server failed", err)
	}
	t.Log("initialization of consul server finished")

	client, err := api.NewClient(api.DefaultConfig())
	if err != nil {
		t.Fatal("failed to initialize consul api client", err)
	}
	t.Log("initilization of consul api finished")
	kv := client.KV()

	projectDir, err := GetRootProjectDir()
	if err != nil {
		t.Fatal("failed to get working project directory", err)
	}

	// Command line argument validations
	// no flag provided
	g2cCmd := exec.Command(projectDir + "/git2consul")
	expectedNoConfigMsg := fmt.Sprintf("%v error No configuration file provided", time.Now().Format("2006/01/02 15:04:05"))
	t.Log("Testig argument parsing...")
	t.Logf("Expected output: %s \n", expectedNoConfigMsg)
	err = executeCommand(g2cCmd, expectedNoConfigMsg)
	if err != nil {
		t.Fatal("git2consul run failed", err)
	} else {
		t.Log("git2consul ran successfully")
	}

	// version
	g2cCmd = exec.Command(projectDir+"/git2consul",
		"-version",
	)
	t.Logf("Expected output: %s \n", "git2consul, version")
	err = executeCommand(g2cCmd, "git2consul, version")
	if err != nil {
		t.Fatal("git2consul run failed", err)
	} else {
		t.Log("git2consul ran successfully")
	}

	// Integration tests
	g2cCmd = exec.Command(projectDir+"/git2consul",
		"-config",
		projectDir+"/pkg/e2e/data/create-config.json",
		"-debug",
		"-once")
	err = executeCommand(g2cCmd, "Terminating git2consul")
	if err != nil {
		t.Fatal("git2consul run failed", err)
	} else {
		t.Log("git2consul ran successfully")
	}

	pair, _, err := kv.Get("e2e/master/genre", nil)
	if err != nil {
		t.Fatal("failed to get pair e2e/master/genre", err)
	}
	if pair == nil {
		t.Fatal("did not find expected pair at e2e/master/genre")
	}
	if strings.TrimSpace(string(pair.Value)) != "classical" {
		t.Errorf("got %s want %s", string(pair.Value), "classical")
	} else {
		t.Log("confirmed e2e/master/genre was properly set")
	}

	// TODO: remove delete when issue #31 is resolved this is a workaround because if
	// the ref matches the current commit, the kv pairs will not be updated
	_, err = kv.Delete("e2e/master.ref", nil)
	if err != nil {
		t.Fatal("failed to delete git reference")
	}

	g2cCmd = exec.Command(projectDir+"/git2consul",
		"-config",
		projectDir+"/pkg/e2e/data/update-config.json",
		"-debug",
		"-once")
	err = executeCommand(g2cCmd, "Terminating git2consul")
	if err != nil {
		t.Fatal("git2consul run failed", err)
	} else {
		t.Log("git2consul ran successfully")
	}

	pair, _, err = kv.Get("e2e/master/artist", nil)
	if err != nil {
		t.Fatal("failed to get pair e2e/master/artist", err)
	}
	if pair == nil {
		t.Fatal("did not find expected pair at e2e/master/artist")
	}
	if strings.TrimSpace(string(pair.Value)) != "beethoven" {
		t.Errorf("got %s want %s", string(pair.Value), "beethoven")
	} else {
		t.Log("confirmed e2e/master/genre was properly updated")
	}

	// TODO: add stage to simulate delete of kv pair
}

func executeCommand(command *exec.Cmd, expectedLog string) (err error) {
	commandPipeReader, commandPipeWriter := io.Pipe()

	commandScanner := bufio.NewScanner(commandPipeReader)
	commandTee := io.MultiWriter(commandPipeWriter, os.Stderr) // dumping the output on Stderr can be useful when debugging this test

	command.Stdout = commandTee
	command.Stderr = commandTee
	if err := command.Start(); err != nil {
		return err
	}
	defer command.Process.Signal(syscall.SIGTERM)

	commandDone := make(chan error, 0)
	go func() {
		commandDone <- command.Wait()
	}()

	err = waitForString(commandScanner, 20*time.Second, expectedLog)
	if err != nil {
		return err
	}

	return nil
}

func waitForString(s *bufio.Scanner, timeout time.Duration, want string) error {
	done := make(chan error)
	go func() {
		for s.Scan() {
			if strings.Contains(s.Text(), want) {
				done <- nil
				return
			}
		}
		if s.Err() != nil {
			done <- s.Err()
		}
		done <- fmt.Errorf("process finished without printing expected string %q", want)
	}()
	select {
	case err := <-done:
		return err
	case <-time.After(timeout):
		return fmt.Errorf("wait for string %q timed out", want)
	}
}

// GetRootProjectDir returns path to root directory of project
func GetRootProjectDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for !strings.HasSuffix(wd, "git2consul-go") {
		if wd == "/" {
			return "", errors.New(`cannot find project directory, "/" reached`)
		}
		wd = filepath.Dir(wd)
	}
	return wd, nil
}
