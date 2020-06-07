package akashi

import (
	"context"
	"fmt"
	"log"
	"os"

	"hapoon/go-akashi/pkg/akashi"

	"github.com/spf13/cobra"
)

func init() {
	staffCmd.Flags().StringVar(&loginCompanyCode, "company-code", "", "Login company code")
	staffCmd.Flags().StringVarP(&accessToken, "token", "t", "", "Access token")
	staffCmd.Flags().IntVar(&staffID, "staff-id", 0, "Staff ID")
	staffCmd.Flags().StringVar(&target, "target", "", "Target staff ID")
	staffCmd.Flags().IntVar(&page, "page", 0, "Page number")
	rootCmd.AddCommand(staffCmd)
}

var (
	loginCompanyCode string
	accessToken      string
	staffID          int
	target           string
	page             int
)

var staffCmd = &cobra.Command{
	Use:   "staff",
	Short: "Access staff API",
	Long:  `Access staff API`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		p := akashi.GetStaffParam{
			LoginCompanyCode: loginCompanyCode,
			Token:            accessToken,
			StaffID:          staffID,
			Target:           target,
			Page:             page,
		}
		res, err := akashi.GetStaff(ctx, p)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		fmt.Println("response:")
		fmt.Printf("%+v\n", res)
	},
}
