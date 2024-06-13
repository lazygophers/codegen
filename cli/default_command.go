package cli

import (
	"fmt"
	"github.com/lazygophers/codegen/state"
	"github.com/spf13/cobra"
	"strings"
)

func initCommandDefault(command *cobra.Command) {
	helpCommand := &cobra.Command{
		Use:   "help [command]",
		Short: state.Localize(state.I18nTagCliHelpShort),
		Long: state.Localize(state.I18nTagCliHelpLong, map[string]any{
			"Name": command.Name(),
		}),
		Run: func(c *cobra.Command, args []string) {
			cmd, _, e := c.Root().Find(args)
			if cmd == nil || e != nil {
				c.Printf("Unknown help topic %#q\n", args)
				cobra.CheckErr(c.Root().Usage())
			} else {
				cmd.InitDefaultHelpFlag()    // make possible 'help' flag to be shown
				cmd.InitDefaultVersionFlag() // make possible 'version' flag to be shown
				cobra.CheckErr(cmd.Help())
			}
		},
	}

	helpCommand.ValidArgsFunction = func(c *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
		var completions []string
		cmd, _, e := c.Root().Find(args)
		if e != nil {
			return nil, cobra.ShellCompDirectiveNoFileComp
		}
		if cmd == nil {
			// Root help command.
			cmd = c.Root()
		}
		for _, subCmd := range cmd.Commands() {
			if subCmd.IsAvailableCommand() || subCmd == helpCommand {
				if strings.HasPrefix(subCmd.Name(), toComplete) {
					completions = append(completions, fmt.Sprintf("%s\t%s", subCmd.Name(), subCmd.Short))
				}
			}
		}
		return completions, cobra.ShellCompDirectiveNoFileComp
	}

	command.Flags().BoolP("help", "h", false, state.Localize(state.I18nTagCliHelpUsage, map[string]string{
		"Name": command.Name(),
	}))

	command.SetHelpCommand(helpCommand)
	command.SetUsageTemplate(fmt.Sprintf(`Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

%s:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}{{$cmds := .Commands}}{{if eq (len .Groups) 0}}

%s:{{range $cmds}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{else}}{{range $group := .Groups}}

{{.Title}}{{range $cmds}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{if not .AllChildCommandsHaveGroup}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

%s:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

%s:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

%s "{{.CommandPath}} [command] --help" %s.{{end}}
`,
		state.Localize(state.I18nTagCliHelpAliases),
		state.Localize(state.I18nTagCliHelpAvailableCommands),
		//state.Localize(state.I18nTagCliHelpAdditionalCommands),
		state.Localize(state.I18nTagCliHelpFlags),
		state.Localize(state.I18nTagCliHelpGlobalFlags),
		state.Localize(state.I18nTagCliHelpUse),
		state.Localize(state.I18nTagCliHelpGetInfo),
	))

	initCommandsDefaule(command.Commands())
}

func initCommandsDefaule(commands []*cobra.Command) {
	for _, command := range commands {
		initCommandDefault(command)
	}
}

func initCompletion() {
	const (
		compCmdName           = "completion"
		compCmdNoDescFlagName = "no-descriptions"
	)

	completionCmd := &cobra.Command{
		Use:   compCmdName,
		Short: state.Localize(state.I18nTagCliCompletionShort),
		Long: state.Localize(state.I18nTagCliCompletionLong, map[string]any{
			"RootName": rootCmd.Root().Name(),
		}),
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		Hidden:            rootCmd.CompletionOptions.HiddenDefaultCmd,
	}

	out := rootCmd.OutOrStdout()
	noDesc := rootCmd.CompletionOptions.DisableDescriptions
	bash := &cobra.Command{
		Use: "bash",
		Long: fmt.Sprintf(`Generate the autocompletion script for the bash shell.

This script depends on the 'bash-completion' package.
If it is not installed already, you can install it via your OS's package manager.

To load completions in your current shell session:

	source <(%[1]s completion bash)

To load completions for every new session, execute once:

#### Linux:

	%[1]s completion bash > /etc/bash_completion.d/%[1]s

#### macOS:

	%[1]s completion bash > $(brew --prefix)/etc/bash_completion.d/%[1]s

You will need to start a new shell for this setup to take effect.
`, rootCmd.Name()),
		Args:                  cobra.NoArgs,
		DisableFlagsInUseLine: true,
		ValidArgsFunction:     cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenBashCompletionV2(out, !noDesc)
		},
	}
	bash.Short = state.Localize(state.I18nTagCliCompletionSubcommandShort, map[string]any{
		"Command": bash.Name(),
	})
	bash.Flags().BoolVar(&noDesc, compCmdNoDescFlagName, false, state.Localize(state.I18nTagCliCompletionFlagsNoDescriptions))

	zsh := &cobra.Command{
		Use: "zsh",
		Long: fmt.Sprintf(`Generate the autocompletion script for the zsh shell.

	If shell completion is not already enabled in your environment you will need
	to enable it.  You can execute the following once:

		echo "autoload -U compinit; compinit" >> ~/.zshrc

	To load completions in your current shell session:

		source <(%[1]s completion zsh)

	To load completions for every new session, execute once:

	#### Linux:

		%[1]s completion zsh > "${fpath[1]}/_%[1]s"

	#### macOS:

		%[1]s completion zsh > $(brew --prefix)/share/zsh/site-functions/_%[1]s

	You will need to start a new shell for this setup to take effect.
	`, rootCmd.Name()),
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if noDesc {
				return rootCmd.GenZshCompletionNoDesc(out)
			}
			return rootCmd.GenZshCompletion(out)
		},
	}
	zsh.Short = state.Localize(state.I18nTagCliCompletionSubcommandShort, map[string]any{
		"Command": zsh.Name(),
	})
	zsh.Flags().BoolVar(&noDesc, compCmdNoDescFlagName, false, state.Localize(state.I18nTagCliCompletionFlagsNoDescriptions))

	fish := &cobra.Command{
		Use: "fish",
		Long: fmt.Sprintf(`Generate the autocompletion script for the fish shell.

	To load completions in your current shell session:

		%[1]s completion fish | source

	To load completions for every new session, execute once:

		%[1]s completion fish > ~/.config/fish/completions/%[1]s.fish

	You will need to start a new shell for this setup to take effect.
	`, rootCmd.Name()),
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			return rootCmd.GenFishCompletion(out, !noDesc)
		},
	}
	fish.Short = state.Localize(state.I18nTagCliCompletionSubcommandShort, map[string]any{
		"Command": fish.Name(),
	})
	fish.Flags().BoolVar(&noDesc, compCmdNoDescFlagName, false, state.Localize(state.I18nTagCliCompletionFlagsNoDescriptions))

	powershell := &cobra.Command{
		Use: "powershell",
		Long: fmt.Sprintf(`Generate the autocompletion script for powershell.

	To load completions in your current shell session:

		%[1]s completion powershell | Out-String | Invoke-Expression

	To load completions for every new session, add the output of the above command
	to your powershell profile.
	`, rootCmd.Name()),
		Args:              cobra.NoArgs,
		ValidArgsFunction: cobra.NoFileCompletions,
		RunE: func(cmd *cobra.Command, args []string) error {
			if noDesc {
				return rootCmd.GenPowerShellCompletion(out)
			}
			return rootCmd.GenPowerShellCompletionWithDesc(out)

		},
	}
	powershell.Short = state.Localize(state.I18nTagCliCompletionSubcommandShort, map[string]any{
		"Command": powershell.Name(),
	})
	powershell.Flags().BoolVar(&noDesc, compCmdNoDescFlagName, false, state.Localize(state.I18nTagCliCompletionFlagsNoDescriptions))

	completionCmd.AddCommand(bash, zsh, fish, powershell)

	rootCmd.AddCommand(completionCmd)
	rootCmd.InitDefaultCompletionCmd()
}

func initDefault() {
	initCompletion()

	initCommandDefault(rootCmd)
}
