package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

/*
 *  Helpers
 */
func exit(level int, format string, args ...any) {
	echo(format, args...)
	os.Exit(level)
}

func echo(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
	fmt.Fprintln(os.Stderr)
}

func shell(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func capture(command string, args ...string) (string, error) {
	cmd := exec.Command(command, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		echo("error trying to run command %q with args %q", command, strings.Join(args, " "))
		echo(string(output))
	}
	return string(output), err
}
