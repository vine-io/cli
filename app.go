// Copyright 2020 The vine Authors
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

package cli

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

var (
	changeLogURL            = "https://github.com/lack-io/cli/blob/master/docs/CHANGELOG.md"
	appActionDeprecationURL = fmt.Sprintf("%s#deprecated-cli-app-action-sigature", changeLogURL)
	contactSysadmin         = "This is an error in the application. Please contact the distributor of this application if this is not you."
	errInvalidActionType    = NewExitError("ERROR invalid Action type. "+
		fmt.Sprintf("Must be `func(*Context)` or `func(*Context) error). %s", contactSysadmin)+
		fmt.Sprintf("See %s", appActionDeprecationURL), 2)
)

// App is the main structure of a cli application. It is recommended that
// an app be created with the cli.NewApp() function
type App struct {
	// The name of the program. Defaults to path.Base(os.Args[0])
	Name string
	// Full name of command for help, defaults to Name
	HelpName string
	// Description of the program
	Usage string
	// Text to override the USAGE section of help
	UsageText string
	// Description of the program argument format.
	ArgsUsage string
	// Version of the program
	Version string
	// Description of the program
	Description string
	// List of commands to execute
	Commands []*Command
	// List of flags to parse
	Flags []Flag
	// Boolean to enable bash completion commands
	EnableBashCompletion bool
	// Boolean to hide built-in help command
	HideHelp bool
	// Boolean to hide built-in version flag and the VERSION section of help
	HideVersion bool
	// categories contains the categorized commands and is populated on app startup
	categories CommandCategories
	// An action to execute when shell completion flag is set
	BashComplete BashCompletionFunc
	// An action to execute before any subcommands are run, but after the context is ready
	// If a non-nil error is returned, no subcommands are run
	Before BeforeFunc
	// An action to execute after an subcommands are run, but after the subcommand has finished
	// It is run even if Action panics
	After AfterFunc
	// The action to execute when no subcommands are specified
	Action ActionFunc
	// Execute this function if the proper command cannot be found
	CommandNotFound CommandNotFoundFunc
	// Execute this function if an usage error occurs
	OnUsageError OnUsageErrorFunc
	// Compilation date
	Compiled time.Time
	// List of all authors two contributed
	Authors []*Author
	// Copyright of the binary if any
	Copyright string
	// Writer writer to write output to
	Writer io.Writer
	// ErrWriter writes error output
	ErrWriter io.Writer
	// Execute this function to handle ExitErrors. If not provided, HandleExitCoder is provided to
	// function as a default, so this is optional.
	ExitErrHandler ExitErrHandleFunc
	// Other custom info
	Metadata map[string]interface{}
	// Carries a function which returns app specific into.
	ExtraInfo func() map[string]string
	// CustomAppHelpTemplate the text template for app help topic.
	// cli.go uses text/template to render templates. You can
	// render custom help text by setting this variable.
	CustomAppHelpTemplate string
	// Boolean to enable short-option handling so user can combine several
	// single-character bool arguments into one
	// i.e. foobar -o -v -> foobar -ov
	UseShortOptionHandling bool

	didSetup bool
}

// Tries to find out when this binary was compiled.
// Returns ths current time if it fails to find it.
func compileTime() time.Time {
	info, err := os.Stat(os.Args[0])
	if err != nil {
		return time.Now()
	}
	return info.ModTime()
}

// NewApp creates a new cli Application with some reasonable defaults for Name.
// Usage, Version and Action.
func NewApp() *App {
	return &App{
		Name:         filepath.Base(os.Args[0]),
		HelpName:     filepath.Base(os.Args[0]),
		Usage:        "A new cli application",
		UsageText:    "",
		BashComplete: DefaultAppComplete,
		Action:       helpCommand.Action,
		Compiled:     compileTime(),
		Writer:       os.Stdout,
	}
}

// Setup runs initialization code to ensure all data structures are ready for
// `Run` or inspection prior to `Run`.  It is internally called by `Run`, but
// will return early if setup has already happened.
func (a *App) Setup() {
	if a.didSetup {
		return
	}

	a.didSetup = true

	if a.Name == "" {
		a.Name = filepath.Base(os.Args[0])
	}

	if a.HelpName == "" {
		a.HelpName = filepath.Base(os.Args[0])
	}

	if a.Usage == "" {
		a.Usage = "A new cli application"
	}

	if a.Version == "" {
		a.HideVersion = true
	}

	if a.BashComplete == nil {
		a.BashComplete = DefaultAppComplete
	}

	if a.Action == nil {
		a.Action = helpCommand.Action
	}

	if a.Compiled == (time.Time{}) {
		a.Compiled = compileTime()
	}

	if a.Writer == nil {
		a.Writer = os.Stdout
	}

	var newCommands []*Command
	for _, c := range a.Commands {
		if c.HelpName == "" {
			c.HelpName = fmt.Sprintf("%s %s", a.HelpName, c.Name)
		}
		newCommands = append(newCommands, c)
	}
	a.Commands = newCommands

	if a.Command(helpCommand.Name) == nil && !a.HideHelp {
		a.appendCommand(helpCommand)
	}
}

func (a *App) handleExitCoder(ctx *Context, err error) {
	if a.ExitErrHandler != nil {
		a.ExitErrHandler(ctx, err)
	} else {
		HandleExitCoder(err)
	}
}

// Author represents someone who has contributed to a cli project.
type Author struct {
	Name  string // The Authors name
	Email string // The Authors email
}

// String makes Author comply to the Stringer interface, to allow an easy print in the templating process
func (a *Author) String() string {
	e := ""
	if a.Email != "" {
		e = " <" + a.Email + ">"
	}

	return fmt.Sprintf("%v%v", a.Name, e)
}

// HandleAction attempts to figure out which Action signature was used.  If
// it's an ActionFunc or a func with the legacy signature for Action, the func
// is run!
func HandleAction(action interface{}, context *Context) (err error) {
	switch a := action.(type) {
	case ActionFunc:
		return a(context)
	case func(*Context) error:
		return a(context)
	case func(*Context):
		a(context)
		return
	}

	return errInvalidActionType
}