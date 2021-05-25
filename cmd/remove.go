/*
Copyright © 2021 Kaleido, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path"
	"strings"

	"github.com/kaleido-io/firefly-cli/internal/stacks"
	"github.com/nguyer/promptui"
	"github.com/spf13/cobra"
)

// removeCmd represents the remove command
var removeCmd = &cobra.Command{
	Use:   "remove <stack_name>",
	Short: "Completely remove a stack",
	Long: `Completely remove a stack

This command will completely delete a stack, including all of its data
and configuration. The stack must be stopped to run this command.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) == 0 {
			log.Fatal("No stack specified!")
		}
		stackName := args[0]

		if !stacks.CheckExists(stackName) {
			log.Fatalf("Stack '%s' does not exist!", stackName)
		}

		prompt := promptui.Prompt{
			Label:     fmt.Sprintf("Completely delete FireFly stack '%s'", stackName),
			IsConfirm: true,
		}

		fmt.Println("WARNING: This will completely remove your stack and all of its data. Are you sure this is what you want to do?")
		result, err := prompt.Run()

		if err != nil || strings.ToLower(result) != "y" {
			fmt.Printf("canceled")
			return nil
		} else {
			fmt.Printf("deleting FireFly stack '%s'... ", stackName)

			dockerCmd := exec.Command("docker", "compose", "rm", "-f")
			dockerCmd.Dir = path.Join(stacks.StacksDir, stackName)
			err := dockerCmd.Run()
			if err != nil {
				return fmt.Errorf("command finished with error: %v", err)
			}

			os.RemoveAll(path.Join(stacks.StacksDir, stackName))
			fmt.Println("done")
			return nil
		}
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
}
