package cmd

import (
	"bufio"
	"fmt"
	"os"

	"git.sr.ht/~rottenfishbone/go-cook/internal/pkg/server"
	"git.sr.ht/~rottenfishbone/go-cook/pkg/config"
	"git.sr.ht/~rottenfishbone/go-cook/pkg/users"
	"github.com/spf13/cobra"
	"golang.org/x/term"
)

var noWebapp bool

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Hosts a local webserver to view/manage recipes.",
	Long:  ``,

	PreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
		config.EnsureUsersInit()
	},

	Run: func(cmd *cobra.Command, args []string) {
		server.Start(6969, noWebapp)
	},
}

var serverRegisterCmd = &cobra.Command{
	Use:   "register [username]",
	Short: "Prompts a user registration for the webapp.",
	Long:  ``,

	PreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
		config.EnsureUsersInit()
	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		var usr string
		if len(args) < 1 {
			fmt.Print("username: ")
			if _, err = fmt.Scanln(&usr); err != nil {
				os.Stderr.WriteString("Aborted.\n")
				os.Exit(1)
			}
		} else if len(args) == 1 {
			usr = args[0]
		} else {
			cmd.Help()
			os.Exit(1)
		}

		// All this just to hide input crossplatform -_-
		fmt.Print("password: ")
		fd := int(os.Stdin.Fd())
		var oldState *term.State
		if oldState, err = term.MakeRaw(fd); err != nil {
			panic(err)
		}

		t := term.NewTerminal(
			bufio.NewReadWriter(
				bufio.NewReader(os.Stdin),
				bufio.NewWriter(os.Stdout)), "")

		var pass string
		if pass, err = t.ReadPassword(""); err != nil {
			term.Restore(fd, oldState)
			os.Stderr.WriteString("Aborted.\n")
			os.Exit(1)
		}
		term.Restore(fd, oldState)
		fmt.Println("")

		if err = users.AddUser(usr, pass); err != nil {
			errMsg := fmt.Sprintf("Failed to save: %v\n", err)
			os.Stderr.WriteString(errMsg)
			os.Exit(1)
		}

		fmt.Printf(`"%v" added to users.%v`, usr, "\n")
	},
}

var serverDeregisterCmd = &cobra.Command{
	Use:   "deregister [username]",
	Short: "Removes user from the webapp.",
	Long:  ``,

	PreRun: func(cmd *cobra.Command, args []string) {
		initConfig()
		config.EnsureUsersInit()
	},

	Run: func(cmd *cobra.Command, args []string) {
		var err error

		var usr string
		if len(args) < 1 {
			fmt.Print("username: ")
			if _, err = fmt.Scanln(&usr); err != nil {
				os.Stderr.WriteString("Aborted.\n")
				os.Exit(1)
			}
		} else if len(args) == 1 {
			usr = args[0]
		} else {
			cmd.Help()
			os.Exit(1)
		}

		if err = users.RemoveUser(usr); err != nil {
			errMsg := fmt.Sprintf("Failed to save: %v\n", err)
			os.Stderr.WriteString(errMsg)
			os.Exit(1)
		}

		fmt.Printf(`"%v" removed from users.%v`, usr, "\n")
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.AddCommand(serverRegisterCmd)
	serverCmd.AddCommand(serverDeregisterCmd)

	serverCmd.Flags().BoolVarP(&noWebapp, "no-webapp", "", false, "Host the API server without the web app")
}
