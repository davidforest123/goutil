package gcmd

import (
	"fmt"
	"goutil/sys/gproc"
	"os"
	"os/exec"
	"strings"
	"time"
)

// sudo 下执行ExecWait或者ExecNoWait接口时会有问题

type (
	Cmder struct {
		cmd *exec.Cmd
	}
)

func StartCmder(screenPrint bool, cmd string, arg ...string) *Cmder {
	var cmder Cmder
	cmder.cmd = exec.Command(cmd, arg...)

	fmt.Println("command line:")
	fmt.Println(cmd, strings.Join(append([]string{}, arg...), " "))
	if screenPrint {
		cmder.cmd.Stdout = os.Stdout
		cmder.cmd.Stderr = os.Stderr
		if err := cmder.cmd.Start(); err != nil {
			fmt.Println(err.Error())
		}
	}
	return &cmder
}

func (c *Cmder) GetPid() int32 {
	return int32(c.cmd.Process.Pid)
}

func (c *Cmder) Wait() error {
	return c.cmd.Wait()
}

func (c *Cmder) WaitWithTimeout(timeout time.Duration) (execTimeout bool, execError error) {
	tmr := time.NewTimer(timeout)
	err := error(nil)
	chDone := make(chan struct{})

	go func() {
		err = c.cmd.Wait()
		chDone <- struct{}{}
	}()

	select {
	case <-chDone:
		return false, err
	case <-tmr.C:
		return true, nil
	}
}

func (c *Cmder) Kill() error {
	return gproc.Terminate(gproc.ProcId(c.GetPid()))
}

func (c *Cmder) GetResultOutput() (string, error) {
	out, err := c.cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	return string(out), nil
}

// will print screen
func ExecWaitPrintScreen(name string, arg ...string) error {
	var cmd *exec.Cmd

	cmd = exec.Command(name, arg...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}

// will not print screen
func ExecWaitReturn(name string, arg ...string) ([]byte, error) {
	return exec.Command(name, arg...).CombinedOutput()
}

// CommandExists checks if command exists.
func CommandExists(cmd string) bool {
	_, err := exec.LookPath(cmd)
	return err == nil
}

// RunScript run shell script just like in the terminal.
// This function waits until the end of the command execution before returning the result.
func RunScript(shellCommand string) ([]byte, error) {
	// implement 1
	// return exec.Command("sh", "-c", shellCommand).CombinedOutput()

	// implement 2
	args := strings.Fields(shellCommand)
	return exec.Command(args[0], args[1:]...).CombinedOutput()
}

// RunScripts run shell scripts just like in the terminal.
func RunScripts(scripts ...string) ([]byte, error) {
	var result []byte
	for _, script := range scripts {
		out, err := RunScript(script)
		if len(out) > 0 {
			result = append(result, out...)
		}
		if err != nil {
			return result, err
		}
	}

	return result, nil
}
