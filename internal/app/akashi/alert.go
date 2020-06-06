package akashi

import (
	"context"
	"fmt"
	"hapoon/go-akashi/pkg/akashi"
	"log"
	"os"

	"github.com/spf13/cobra"
)

func init() {
	alertCmd.Flags().StringVar(&loginCompanyCode, "company-code", "", "Login company code")
	alertCmd.Flags().StringVarP(&accessToken, "token", "t", "", "Access token")
	rootCmd.AddCommand(alertCmd)
}

var alertCmd = &cobra.Command{
	Use:   "alert",
	Short: "Access alert API",
	Long:  "Access alert API",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		p := akashi.GetAlertParam{
			LoginCompanyCode: loginCompanyCode,
			Token:            accessToken,
		}
		res, err := akashi.GetAlerts(ctx, p)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		fmt.Println("response:")
		fmt.Printf("%+v\n", res)
	},
}
