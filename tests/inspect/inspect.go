// Copyright 2015 CoreOS, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	globalFlagset = flag.NewFlagSet("inspect", flag.ExitOnError)
	globalFlags   = struct {
		PrintMsg  string
		PrintEnv  string
		CheckCwd  string
		ExitCode  int
		ReadFile  bool
		WriteFile bool
	}{}
)

func init() {
	globalFlagset.StringVar(&globalFlags.PrintMsg, "print-msg", "", "Print the message given as parameter")
	globalFlagset.StringVar(&globalFlags.CheckCwd, "check-cwd", "", "Check if the current working directory is the one specified")
	globalFlagset.StringVar(&globalFlags.PrintEnv, "print-env", "", "Print the specified environment variable")
	globalFlagset.IntVar(&globalFlags.ExitCode, "exit-code", 0, "Return this exit code")
	globalFlagset.BoolVar(&globalFlags.ReadFile, "read-file", false, "Print the content of the file $FILE")
	globalFlagset.BoolVar(&globalFlags.WriteFile, "write-file", false, "Write $CONTENT in the file $FILE")
}

func main() {
	globalFlagset.Parse(os.Args[1:])
	args := globalFlagset.Args()
	if len(args) > 0 {
		fmt.Fprintln(os.Stderr, "Wrong parameters\n")
		os.Exit(1)
	}

	if globalFlags.PrintMsg != "" {
		fmt.Fprintf(os.Stdout, "%s\n", globalFlags.PrintMsg)
	}

	if globalFlags.PrintEnv != "" {
		fmt.Fprintf(os.Stdout, "%s=%s\n", globalFlags.PrintEnv, os.Getenv(globalFlags.PrintEnv))
	}

	if globalFlags.WriteFile {
		fileName := os.Getenv("FILE")
		err := ioutil.WriteFile(fileName, []byte(os.Getenv("CONTENT")), 0600)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot write to file %q: %v\n", fileName, err)
			os.Exit(1)
			return
		}
	}

	if globalFlags.ReadFile {
		fileName := os.Getenv("FILE")
		dat, err := ioutil.ReadFile(fileName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot read file %q: %v\n", fileName, err)
			os.Exit(1)
			return
		}
		fmt.Print("<<<")
		fmt.Print(string(dat))
		fmt.Print(">>>\n")
	}

	if globalFlags.CheckCwd != "" {
		wd, err := os.Getwd()
		if err != nil {
			fmt.Fprintf(os.Stderr, "Cannot get working directory: %v\n", err)
			os.Exit(1)
		}
		if wd != globalFlags.CheckCwd {
			fmt.Fprintf(os.Stderr, "Working directory: %q. Expected: %q.\n", wd, globalFlags.CheckCwd)
			os.Exit(1)
		}
	}

	os.Exit(globalFlags.ExitCode)
}
