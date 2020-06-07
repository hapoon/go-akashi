package akashi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
)

// Alert アラート情報
type Alert struct {
	Month     string `json:"month"`      // アラートの発生した月度
	Date      string `json:"date"`       // アラートの発生した日付
	AlertType int    `json:"alert_type"` // アラート種別 1:打刻忘れ 2:欠勤疑い 3:休憩過少/超過 4:出勤打刻乖離 5:退勤打刻乖離 6:休日出勤 7:遅刻 8:早退 9:残業時間域値越え 10:無断出勤
}

// GetAlertParam アラート情報取得リクエストパラメータ
type GetAlertParam struct {
	LoginCompanyCode string // AKASHI企業ID
	Token            string // アクセストークン
}

// GetAlertResponse アラート情報取得レスポンス
type GetAlertResponse struct {
	LoginCompanyCode string  `json:"login_company_code"` // AKASHI企業ID
	StaffID          int     `json:"staff_id"`           // 従業員ID
	Count            int     `json:"count"`              // 取得されたアラートの件数
	Alerts           []Alert `json:"alerts"`             // 取得されたアラートの配列
}

// GetAlerts アラート情報取得
func GetAlerts(ctx context.Context, param GetAlertParam) (GetAlertResponse, error) {
	if param.LoginCompanyCode == "" {
		return GetAlertResponse{}, errors.New("LoginCompanyCode must be set")
	}
	if param.Token == "" {
		return GetAlertResponse{}, errors.New("Token must be set")
	}
	endpoint := fmt.Sprintf("/%s/alerts", param.LoginCompanyCode)
	uv := url.Values{}
	uv.Add("token", param.Token)
	q := uv.Encode()
	if q != "" {
		endpoint += "?" + q
	}

	cli := NewClient(param.LoginCompanyCode, param.Token)
	res, err := cli.Get(ctx, endpoint)
	if err != nil {
		return GetAlertResponse{}, err
	}
	if res.StatusCode != http.StatusOK {
		return GetAlertResponse{}, fmt.Errorf("Status code=%d", res.StatusCode)
	}
	var gar struct {
		Success  bool             `json:"success"`
		Response GetAlertResponse `json:"response"`
		Errors   []Error          `json:"errors"`
	}
	if err := json.NewDecoder(res.Body).Decode(&gar); err != nil {
		return GetAlertResponse{}, err
	}
	if !gar.Success {
		return GetAlertResponse{}, errors.New("Requesting Alert API failed")
	}
	return gar.Response, nil
}
