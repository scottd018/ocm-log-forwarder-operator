/*
Copyright 2023.

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

package commands

import (
	"github.com/spf13/cobra"

	// common imports for subcommands
	cmdgenerate "github.com/scottd018/ocm-log-forwarder-operator/cmd/ocmlogctl/commands/generate"
	cmdinit "github.com/scottd018/ocm-log-forwarder-operator/cmd/ocmlogctl/commands/init"
	cmdversion "github.com/scottd018/ocm-log-forwarder-operator/cmd/ocmlogctl/commands/version"

	// specific imports for workloads
	generateapps "github.com/scottd018/ocm-log-forwarder-operator/cmd/ocmlogctl/commands/generate/apps"
	initapps "github.com/scottd018/ocm-log-forwarder-operator/cmd/ocmlogctl/commands/init/apps"
	versionapps "github.com/scottd018/ocm-log-forwarder-operator/cmd/ocmlogctl/commands/version/apps"
	//+kubebuilder:scaffold:operator-builder:subcommands:imports
)

// OcmlogctlCommand represents the base command when called without any subcommands.
type OcmlogctlCommand struct {
	*cobra.Command
}

// NewOcmlogctlCommand returns an instance of the OcmlogctlCommand.
func NewOcmlogctlCommand() *OcmlogctlCommand {
	c := &OcmlogctlCommand{
		Command: &cobra.Command{
			Use:   "ocmlogctl",
			Short: "Manage OCM Log Forwarder workload",
			Long:  "Manage OCM Log Forwarder workload",
		},
	}

	c.addSubCommands()

	return c
}

// Run represents the main entry point into the command
// This is called by main.main() to execute the root command.
func (c *OcmlogctlCommand) Run() {
	cobra.CheckErr(c.Execute())
}

func (c *OcmlogctlCommand) newInitSubCommand() {
	parentCommand := cmdinit.GetParent(c.Command)
	_ = parentCommand

	// add the init subcommands
	initapps.NewOCMLogForwarderSubCommand(parentCommand)
	//+kubebuilder:scaffold:operator-builder:subcommands:init
}

func (c *OcmlogctlCommand) newGenerateSubCommand() {
	parentCommand := cmdgenerate.GetParent(c.Command)
	_ = parentCommand

	// add the generate subcommands
	generateapps.NewOCMLogForwarderSubCommand(parentCommand)
	//+kubebuilder:scaffold:operator-builder:subcommands:generate
}

func (c *OcmlogctlCommand) newVersionSubCommand() {
	parentCommand := cmdversion.GetParent(c.Command)
	_ = parentCommand

	// add the version subcommands
	versionapps.NewOCMLogForwarderSubCommand(parentCommand)
	//+kubebuilder:scaffold:operator-builder:subcommands:version
}

// addSubCommands adds any additional subCommands to the root command.
func (c *OcmlogctlCommand) addSubCommands() {
	c.newInitSubCommand()
	c.newGenerateSubCommand()
	c.newVersionSubCommand()
}
