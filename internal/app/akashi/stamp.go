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
	// 打刻情報取得
	stampGetCmd.Flags().StringVarP(&startDate, "start-date", "s", "", "Start date")
	stampGetCmd.Flags().StringVarP(&endDate, "end-date", "e", "", "End date")
	stampCmd.AddCommand(stampGetCmd)
	// 打刻
	stampCmd.AddCommand(stampTouchCmd)
	stampCmd.AddCommand(stampGoToWorkCmd)
	stampCmd.AddCommand(stampLeaveWorkCmd)
	stampCmd.AddCommand(stampBreakCmd)
	stampCmd.AddCommand(stampBreakReturnCmd)
	rootCmd.AddCommand(stampCmd)
}

var stampCmd = &cobra.Command{
	Use:   "stamp",
	Short: "Access stamp API",
	Long:  `Access stamp API`,
	Run: func(cmd *cobra.Command, args []string) {
		// このコマンド単体では動作しないのでヘルプを表示する
	},
}

var stampGetCmd = &cobra.Command{
	Use:   "get",
	Short: "打刻情報の取得",
	Long:  `打刻情報の取得`,
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.Println("args:", args)
		}
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

var stampTouchCmd = &cobra.Command{
	Use:   "touch",
	Short: "打刻",
	Long: `打刻
	状況に合わせた打刻を行います。
	未出勤時は出勤の打刻を、出勤時は退勤の打刻を、休憩中は休憩戻りの打刻を行います。
	`,
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		p := akashi.PostStampParam{
			LoginCompanyCode: loginCompanyCode,
			Token:            accessToken,
		}
		res, err := akashi.PostStamp(ctx, p)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		printStampPostResponse(res)
	},
}

var stampGoToWorkCmd = &cobra.Command{
	Use:   "work-in",
	Short: "出勤の打刻",
	Long:  "出勤の打刻",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		p := akashi.PostStampParam{
			LoginCompanyCode: loginCompanyCode,
			Token:            accessToken,
			Type:             akashi.StampTypeGoToWork,
		}
		res, err := akashi.PostStamp(ctx, p)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		printStampPostResponse(res)
	},
}

var stampLeaveWorkCmd = &cobra.Command{
	Use:   "work-out",
	Short: "退勤の打刻",
	Long:  "退勤の打刻",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		p := akashi.PostStampParam{
			LoginCompanyCode: loginCompanyCode,
			Token:            accessToken,
			Type:             akashi.StampTypeLeaveWork,
		}
		res, err := akashi.PostStamp(ctx, p)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		printStampPostResponse(res)
	},
}

var stampBreakCmd = &cobra.Command{
	Use:   "break-in",
	Short: "休憩入りの打刻",
	Long:  "休憩入りの打刻",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		p := akashi.PostStampParam{
			LoginCompanyCode: loginCompanyCode,
			Token:            accessToken,
			Type:             akashi.StampTypeBreak,
		}
		res, err := akashi.PostStamp(ctx, p)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		printStampPostResponse(res)
	},
}

var stampBreakReturnCmd = &cobra.Command{
	Use:   "break-out",
	Short: "休憩戻りの打刻",
	Long:  "休憩戻りの打刻",
	Run: func(cmd *cobra.Command, args []string) {
		ctx := context.Background()
		p := akashi.PostStampParam{
			LoginCompanyCode: loginCompanyCode,
			Token:            accessToken,
			Type:             akashi.StampTypeBreakReturn,
		}
		res, err := akashi.PostStamp(ctx, p)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		printStampPostResponse(res)
	},
}

func printStampPostResponse(res akashi.PostStampResponse) {
	fmt.Println("企業ID:", res.LoginCompanyCode)
	fmt.Println("従業員ID:", res.StaffID)
	fmt.Println("打刻種別:", res.Type)
	fmt.Println("打刻日時:", res.StampedAt)
}
