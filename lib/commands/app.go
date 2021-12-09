// Copyright 2021-present Anon. All rights reserved.
// Use of this source code is governed by Apache 2.0 license
// that can be found in the LICENSE file.

package commands

import (
	"github.com/anonistas/notya/lib/services"
	"github.com/anonistas/notya/pkg"
	"github.com/spf13/cobra"
)

var (
	// NotyaPath is the global notya folder path
	// Which would be filled after executing the application.
	NotyaPath string
)

// service, is the default service of all commands.
var service services.ServiceRepo

// version is current version of application.
const version = "v1.0.0"

// AppCommand is the root command of application and genesis of all sub-commands.
var appCommand = &cobra.Command{
	Use:     "notya",
	Version: version,
	Short:   "\n 📝 Take notes quickly and expeditiously from terminal.",
}

// initCommands sets all special commands to application command.
func initCommands() {
	initCreateCommand()
	initSetupCommand()
}

// RunApp sets all special commands, checks notya initialization status,
// and then executes main app command.
func ExecuteApp() {
	initCommands()

	// Generate notya path.
	notyaPath, err := pkg.NotyaPWD()
	if err != nil {
		pkg.Alert(pkg.ErrorL, err.Error())
	}

	NotyaPath = *notyaPath

	// Initialize new local service.
	service = services.NewLocalService(NotyaPath)

	// Check initialization status of notya,
	// Setup working directories, if it's not initialized before.
	setupErr := initializeIfNotExists(NotyaPath)
	if setupErr != nil {
		pkg.Alert(pkg.ErrorL, setupErr.Error())
		return
	}

	_ = appCommand.Execute()
}