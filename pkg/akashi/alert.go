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
	Month     string    `json:"month"`      // アラートの発生した月度
	Date      string    `json:"date"`       // アラートの発生した日付
	AlertType AlertType `json:"alert_type"` // アラート種別
}

// AlertType アラート種別
type AlertType int

const (
	// AlertTypeForgetStamp 打刻忘れ
	AlertTypeForgetStamp AlertType = iota + 1
	// AlertTypeMayBeAbsent 欠勤疑い
	AlertTypeMayBeAbsent
	// AlertTypeProblemInBreak 休憩過少/超過
	AlertTypeProblemInBreak
	// AlertTypeDivergenceGoToWork 出勤打刻乖離
	AlertTypeDivergenceGoToWork
	// AlertTypeDivergenceLeaveWork 退勤打刻乖離
	AlertTypeDivergenceLeaveWork
	// AlertTypeWorkHoliday 休日出勤
	AlertTypeWorkHoliday
	// AlertTypeLateness 遅刻
	AlertTypeLateness
	// AlertTypeLeaveEarly 早退
	AlertTypeLeaveEarly
	// AlertTypeExceedThresholdOvertime 残業時間域値越え
	AlertTypeExceedThresholdOvertime
	// AlertTypeGoToWorkWithoutPermission 無断出勤
	AlertTypeGoToWorkWithoutPermission
)

func (a AlertType) String() string {
	switch a {
	case AlertTypeForgetStamp:
		return "打刻忘れ"
	case AlertTypeMayBeAbsent:
		return "欠勤疑い"
	case AlertTypeProblemInBreak:
		return "休憩過少/超過"
	case AlertTypeDivergenceGoToWork:
		return "出勤打刻乖離"
	case AlertTypeDivergenceLeaveWork:
		return "退勤打刻乖離"
	case AlertTypeWorkHoliday:
		return "休日出勤"
	case AlertTypeLateness:
		return "遅刻"
	case AlertTypeLeaveEarly:
		return "早退"
	case AlertTypeExceedThresholdOvertime:
		return "残業時間域値越え"
	case AlertTypeGoToWorkWithoutPermission:
		return "無断出勤"
	default:
		return ""
	}
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
