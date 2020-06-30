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
	tokenCmd.AddCommand(tokenReissueCmd)
	rootCmd.AddCommand(tokenCmd)
}

var tokenCmd = &cobra.Command{
	Use:   "token",
	Short: "Access token API",
	Long:  "Access token API",
	Run: func(cmd *cobra.Command, args []string) {
		// このコマンド単体では動作しないのでヘルプを表示する
	},
}

var tokenReissueCmd = &cobra.Command{
	Use:   "reissue",
	Short: "アクセストークンの再発行",
	Long: `トークンにて認証した従業員のアクセストークンを再発行します。
再発行時に有効期限切れのトークンが存在する場合、自動的に削除されます。`,
	Run: func(cmd *cobra.Command, args []string) {
		if verbose {
			log.Println("args:", args)
		}
		ctx := context.Background()
		p := akashi.PostTokenReissueParam{
			LoginCompanyCode: loginCompanyCode,
			Token:            accessToken,
		}
		res, err := akashi.PostTokenReissue(ctx, p)
		if err != nil {
			log.Fatalln(err)
			os.Exit(1)
		}
		fmt.Println("企業ID:", res.LoginCompanyCode)
		fmt.Println("従業員ID:", res.StaffID)
		fmt.Println("agency_manager_id:", res.AgencyManagerID)
		fmt.Println("アクセストークン:", res.Token)
		fmt.Println("有効期限:", res.ExpiredAt)
	},
}
