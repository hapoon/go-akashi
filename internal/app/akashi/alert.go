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
		fmt.Println("企業ID:", res.LoginCompanyCode)
		fmt.Println("従業員ID:", res.StaffID)
		fmt.Println("アラート件数:", res.Count, "件")
		for _, alert := range res.Alerts {
			fmt.Println(alert.Month, "月", alert.Date, "日", alert.AlertType)
		}
	},
}
