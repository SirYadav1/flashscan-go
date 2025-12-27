package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "flashscan-go",
	Short: "FlashScan - High Performance Network Scanner",
	Long:  "FlashScan - High Performance Network Scanner by SirYadav1",
}

var (
	globalFlagThreads      int
	globalFlagStatInterval float64
)

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	// ANSI Color codes for the template
	cReset := "\033[0m"
	cCyan := "\033[36m\033[1m"
	cGreen := "\033[32m"
	cYellow := "\033[33m"

	usageTemplate := `  ______ _           _      _____                 
 |  ____| |         | |    / ____|                
 | |__  | | __ _ ___| |__ | (___   ___ __ _ _ __  
 |  __| | |/ _` + "`" + ` / __| '_ \ \___ \ / __/ _` + "`" + ` | '_ \ 
 | |    | | (_| \__ \ | | |____) | (_| (_| | | | |
 |_|    |_|\__,_|___/_| |_|_____/ \___\__,_|_| |_| v2.0

` + cCyan + `Usage:` + cReset + `
  ` + cGreen + `{{.CommandPath}}` + cReset + ` [command] [flags]

` + cCyan + `Available Commands:` + cReset + `{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  ` + cGreen + `{{rpad .Name .NamePadding}}` + cReset + ` {{.Short}}{{end}}{{end}}

` + cCyan + `Flags:` + cReset + `
` + cYellow + `{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}` + cReset + `

{{if .HasAvailableSubCommands}}` + cCyan + `Global Flags:` + cReset + `
` + cYellow + `{{.PersistentFlags.FlagUsages | trimTrailingWhitespaces}}` + cReset + `

` + cCyan + `Use` + cReset + ` "` + cGreen + `{{.CommandPath}} [command] --help` + cReset + `" ` + cCyan + `for more information about a command.` + cReset + `{{end}}
`
	rootCmd.SetUsageTemplate(usageTemplate)

	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.PersistentFlags().IntVarP(&globalFlagThreads, "threads", "t", 64, "total threads to use")
	rootCmd.PersistentFlags().Float64Var(&globalFlagStatInterval, "stat-interval", 1.0, "stat interval in seconds")
}
