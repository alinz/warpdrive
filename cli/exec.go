package cli

import (
	"fmt"
	"os/exec"
	"sync"
)

type task interface {
	RunIt(func(message string)) error
}

func taskRunner(tasks ...task) (ok bool, errors []error) {
	internal := make(chan struct{})
	wg := new(sync.WaitGroup)
	errors = make([]error, len(tasks))
	ok = true

	messageChann := make(chan string, 10)
	go func() {
		loop := true
		for loop {
			select {
			case message, ok := <-messageChann:
				if !ok {
					loop = false
					continue
				}
				fmt.Print(message)
			}
		}
		close(internal)
	}()

	log := func(message string) {
		messageChann <- message
	}

	for idx, item := range tasks {
		wg.Add(1)
		go func(idx int, item task, wg *sync.WaitGroup) {
			err := item.RunIt(log)
			if err != nil {
				ok = false
				errors[idx] = err
			}
			wg.Done()
		}(idx, item, wg)
	}

	wg.Wait()
	close(messageChann)

	<-internal

	return
}

type cmdTask struct {
	label   string
	command string
	args    []string
}

func (c *cmdTask) RunIt(log func(string)) error {
	log(fmt.Sprintf("%s started\n", c.label))

	cmd := exec.Command(c.command, c.args...)
	_, err := cmd.Output()

	log(fmt.Sprintf("%s finished\n", c.label))
	return err
}

func newCmdTask(label, command string, args ...string) task {
	return &cmdTask{label: label, command: command, args: args}
}
