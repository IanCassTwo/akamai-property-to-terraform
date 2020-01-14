/*
 * Copyright 2018. Akamai Technologies, Inc
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	akamai "github.com/akamai/cli-common-golang"
	"github.com/urfave/cli"
)

var commandLocator akamai.CommandLocator = func() ([]cli.Command, error) {
	cmdFlags := []cli.Flag{
	}

	commands := []cli.Command{
		{
			Name:      "create",
			Usage:     "create",
			ArgsUsage: "<property name>",
			Action:    cmdCreate,
			Flags:     cmdFlags,
			BashComplete: akamai.DefaultAutoComplete,
		},
		{
			Name:   "list",
			Usage:  "List commands",
			Action: akamai.CmdList,
			BashComplete: akamai.DefaultAutoComplete,
		},
		cli.Command{
			Name:         "help",
			Description:  "Displays help information",
			ArgsUsage:    "[command] [sub-command]",
			Action:       akamai.CmdHelp,
			BashComplete: akamai.DefaultAutoComplete,
		},
	}

	return commands, nil
}
