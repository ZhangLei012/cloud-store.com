package util

import (
	"bytes"
	"os/exec"
)

//执行linux shell command
func ExecLinuxShell(s string) (string, error) {
	//函数返回一个io.Write类型的*Cmd
	cmd := exec.Command("/bin/bash", "-c", s)

	//通过bytes.Buffer将byte类型转化为string
	var result bytes.Buffer
	cmd.Stdout = &result

	//Run执行cmd命令，并阻塞直至完成
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	return result.String(), err
}
