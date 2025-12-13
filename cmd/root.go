package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "flashscan-go",
	Long: "FlashScan - High Performance Network Scanner by SirYadav1",
}

var (
	globalFlagThreads      int
	globalFlagStatInterval float64
)

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.SetHelpCommand(&cobra.Command{Hidden: true})
	rootCmd.PersistentFlags().IntVarP(&globalFlagThreads, "threads", "t", 64, "total threads to use")
	rootCmd.PersistentFlags().Float64Var(&globalFlagStatInterval, "stat-interval", 1.0, "stat interval in seconds")
}
