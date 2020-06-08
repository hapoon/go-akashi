package akashi

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().StringVar(&loginCompanyCode, "company-code", "", "Login company code")
	rootCmd.PersistentFlags().StringVarP(&accessToken, "token", "t", "", "Access token")
}

var rootCmd = &cobra.Command{
	Use:   "akactl",
	Short: "akactl is a command line tool for AKASHI",
	Long: `A command line tool for AKASHI
			Complete documentation is available at ...`,
	Run: func(cmd *cobra.Command, args []string) {
		// Do Stuff Here
	},
}

// Execute is execution root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
