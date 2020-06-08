package akashi

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"hapoon/go-akashi/pkg/akashi"

	"github.com/spf13/cobra"
)

var (
	// GetStamp
	startDate string
	endDate   string
	// PostStamp
	stampType int
	stampedAt string
	timezone  string
)

func init() {
	stampGetCmd.Flags().StringVarP(&startDate, "start-date", "s", "", "Start date")
	stampGetCmd.Flags().StringVarP(&endDate, "end-date", "e", "", "End date")
	stampCmd.AddCommand(stampGetCmd)
	stampPostCmd.Flags().IntVar(&stampType, "type", 0, "Stamp type")
	//stampPostCmd.Flags().StringVar(&stampedAt, "stamped-at", "", "クライアントでの打刻日時")
	//stampPostCmd.Flags().StringVar(&timezone, "timezone", "", "タイムゾーン")
	stampCmd.AddCommand(stampPostCmd)
	rootCmd.AddCommand(stampCmd)
}

var stampCmd = &cobra.Command{
	Use:   "stamp",
	Short: "Access stamp API",
	Long:  `Access stamp API`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var stampGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get stamp information by stamp API",
	Long:  `Get stamp information by stamp API`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		ts, err := time.Parse(akashi.DateFormat, startDate)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		te, err := time.Parse(akashi.DateFormat, endDate)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		p := akashi.GetStampParam{
			LoginCompanyCode: loginCompanyCode,
			Token:            accessToken,
			StartDate:        ts,
			EndDate:          te,
		}
		res, err := akashi.GetStamps(ctx, p)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		fmt.Println("企業ID:", res.LoginCompanyCode)
		fmt.Println("従業員ID:", res.StaffID)
		fmt.Println("件数:", res.Count)
		for _, v := range res.Stamps {
			fmt.Println("------------------------------------")
			fmt.Println("打刻日時:", v.StampedAt)
			fmt.Println("打刻種別:", v.Type)
			fmt.Println("ローカル打刻時刻:", v.LocalTime)
			fmt.Println("タイムゾーン:", v.Timezone)
			fmt.Println("打刻方法:", v.Attributes.Method)
			fmt.Println("組織ID:", v.Attributes.OrgID)
			fmt.Println("勤務地ID:", v.Attributes.WorkplaceID)
			fmt.Println("緯度:", v.Attributes.Latitude)
			fmt.Println("経度:", v.Attributes.Longitude)
			fmt.Println("IPアドレス:", v.Attributes.IP)
		}
	},
}

var stampPostCmd = &cobra.Command{
	Use:   "post",
	Short: "Post stamp by stamp API",
	Long:  "Post stamp by stamp API",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		p := akashi.PostStampParam{
			LoginCompanyCode: loginCompanyCode,
			Token:            accessToken,
			Type:             akashi.StampType(stampType),
		}
		res, err := akashi.PostStamp(ctx, p)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		fmt.Println("企業ID:", res.LoginCompanyCode)
		fmt.Println("従業員ID:", res.StaffID)
		fmt.Println("打刻種別:", res.Type)
		fmt.Println("打刻日時:", res.StampedAt)
	},
}
