// Copyright © 2021 Kaleido, Inc.
//
// SPDX-License-Identifier: Apache-2.0
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

package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hyperledger/firefly-cli/internal/constants"
	"github.com/hyperledger/firefly-cli/internal/docker"
	"github.com/hyperledger/firefly-cli/internal/stacks"
	"github.com/spf13/cobra"
)

var removeCmd = &cobra.Command{
	Use:     "remove <stack_name>",
	Aliases: []string{"rm"},
	Short:   "Completely remove a stack",
	Long: `Completely remove a stack

This command will completely delete a stack, including all of its data
and configuration.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := docker.CheckDockerConfig(); err != nil {
			return err
		}
		stackManager := stacks.NewStackManager(logger)
		if len(args) == 0 {
			return fmt.Errorf("no stack specified")
		}
		stackName := args[0]

		if exists, err := stacks.CheckExists(stackName); err != nil {
			return err
		} else if !exists {
			return fmt.Errorf("stack '%s' does not exist", stackName)
		}

		if !force {
			fmt.Println("WARNING: This will completely remove your stack and all of its data. Are you sure this is what you want to do?")
			if err := confirm(fmt.Sprintf("completely delete FireFly stack '%s'", stackName)); err != nil {
				cancel()
			}
		}

		if err := stackManager.LoadStack(stackName, verbose); err != nil {
			return err
		}
		fmt.Printf("deleting FireFly stack '%s'... ", stackName)
		if err := stackManager.StopStack(verbose); err != nil {
			return err
		}
		if err := stackManager.RemoveStack(verbose); err != nil {
			return err
		}
		os.RemoveAll(filepath.Join(constants.StacksDir, stackName))
		fmt.Println("done")
		return nil
	},
}

func init() {
	removeCmd.Flags().BoolVarP(&force, "force", "f", false, "Remove the stack without prompting for confirmation")
	rootCmd.AddCommand(removeCmd)
}
