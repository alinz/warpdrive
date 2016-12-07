package cli

import (
	"os/exec"
	"sync"
)

type task interface {
	RunIt() error
}

func taskRunner(tasks []task) (ok bool, errors []error) {
	wg := new(sync.WaitGroup)
	errors = make([]error, len(tasks))
	ok = true

	for idx, item := range tasks {
		wg.Add(1)
		go func(idx int, item task, wg *sync.WaitGroup) {
			err := item.RunIt()
			if err != nil {
				ok = false
				errors[idx] = err
			}
			wg.Done()
		}(idx, item, wg)
	}

	wg.Wait()

	return
}

type cmdTask struct {
	command string
	args    []string
}

func (c *cmdTask) RunIt() error {
	cmd := exec.Command(c.command, c.args...)
	_, err := cmd.Output()

	return err
}

func newCmdTask(command string, args ...string) task {
	return &cmdTask{command: command, args: args}
}
