package exec

import (
	"os"
	e "os/exec"
)

func Run(cmd string, args ...string) ([]byte, error) {
	return run(nil, cmd, args...)
}

func RunInDir(wd string, cmd string, args ...string) ([]byte, error) {
	env := make(map[string]string)
	env["WORKING_DIR"] = wd

	return run(env, cmd, args...)
}

func run(env map[string]string, cmd string, args ...string) ([]byte, error) {
	c := e.Command(cmd, args...)

	if val, ok := env["WORKING_DIR"]; ok {
		c.Dir = val
		delete(env, "WORKING_DIR")
	}

	c.Env = os.Environ()
	for k, v := range env {
		c.Env = append(c.Env, k+"="+v)
	}

	return c.CombinedOutput()
}
