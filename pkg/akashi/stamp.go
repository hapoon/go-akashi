package akashi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

var (
	// DateFormat yyyymmddHHMMSS形式
	DateFormat = "20060102150405"
	// ReturnDateFormat yyyy/mm/dd HH:MM:SS形式
	ReturnDateFormat = "2006/01/02 15:04:05"
)

// Stamp 打刻データ
type Stamp struct {
	StampedAt  *AkTime        `json:"stamped_at"` // 打刻日時
	Type       StampType      `json:"type"`       // 打刻種別
	LocalTime  *AkTime        `json:"local_time"` // 打刻機のローカル打刻時刻
	Timezone   string         `json:"timezone"`   // 打刻機のタイムゾーン
	Attributes StampAttribute `json:"attributes"` // 実績参照結果
}

// StampType 打刻種別
type StampType int

const (
	// StampTypeUnknown 打刻種別:不明
	StampTypeUnknown StampType = 0
	// StampTypeGoToWork 打刻種別:出勤
	StampTypeGoToWork StampType = 11
	// StampTypeLeaveWork 打刻種別:退勤
	StampTypeLeaveWork StampType = 12
	// StampTypeGoStraight 打刻種別:直行
	StampTypeGoStraight StampType = 21
	// StampTypeBounce 打刻種別:直帰
	StampTypeBounce StampType = 22
	// StampTypeBreak 打刻種別:休憩入
	StampTypeBreak StampType = 31
	// StampTypeBreakReturn 打刻種別:休憩戻
	StampTypeBreakReturn StampType = 32
)

func (s StampType) String() string {
	switch s {
	case StampTypeGoToWork:
		return "出勤"
	case StampTypeLeaveWork:
		return "退勤"
	case StampTypeGoStraight:
		return "直行"
	case StampTypeBounce:
		return "直帰"
	case StampTypeBreak:
		return "休憩入"
	case StampTypeBreakReturn:
		return "休憩戻"
	default:
		return ""
	}
}

// StampAttribute 打刻実績参照結果
type StampAttribute struct {
	Method      int     `json:"method"`       // 打刻方法
	OrgID       int     `json:"org_id"`       // 組織ID
	WorkplaceID int     `json:"workplace_id"` // 勤務地ID
	Latitude    float32 `json:"latitude"`     // 緯度
	Longitude   float32 `json:"longitude"`    // 経度
	IP          string  `json:"ip"`           // 打刻機のIPアドレス
}

// AkTime 時間
type AkTime struct {
	time.Time
}

// UnmarshalJSON time.UUnmarshalJSONの拡張
func (a *AkTime) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		return nil
	}
	t, err := time.Parse(`"`+ReturnDateFormat+`"`, string(data))
	*a = AkTime{t}
	return err
}

// GetStampParam 打刻情報取得リクエストパラメータ
type GetStampParam struct {
	LoginCompanyCode string    // AKASHI企業ID
	Token            string    // アクセストークン
	StartDate        time.Time // 打刻取得期間の開始日時
	EndDate          time.Time // 打刻取得期間の終了日時
	StaffID          int       // 取得対象の従業員ID
}

// GetStampResponse 打刻情報取得レスポンス
type GetStampResponse struct {
	LoginCompanyCode string  `json:"login_company_code"` // ログイン企業ID
	StaffID          int     `json:"staff_id"`           // 従業員ID
	Count            int     `json:"count"`              // 従業員の打刻数
	Stamps           []Stamp `json:"stamps"`             // 打刻データの配列
}

// GetStamps 打刻情報取得
func GetStamps(ctx context.Context, param GetStampParam) (GetStampResponse, error) {
	if param.LoginCompanyCode == "" {
		return GetStampResponse{}, errors.New("LoginCompanyCode must be set")
	}
	if param.Token == "" {
		return GetStampResponse{}, errors.New("Token must be set")
	}
	if param.StartDate.IsZero() {
		return GetStampResponse{}, errors.New("StartDate must be set")
	}
	if param.EndDate.IsZero() {
		return GetStampResponse{}, errors.New("EndDate must be set")
	}

	endpoint := fmt.Sprintf("/%s/stamps", param.LoginCompanyCode)
	if param.StaffID != 0 {
		endpoint = fmt.Sprintf("%s/%d", endpoint, param.StaffID)
	}
	uv := url.Values{}
	uv.Add("token", param.Token)
	uv.Add("start_date", param.StartDate.Format(DateFormat))
	uv.Add("end_date", param.EndDate.Format(DateFormat))
	q := uv.Encode()
	if q != "" {
		endpoint += "?" + q
	}

	cli := NewClient(param.LoginCompanyCode, param.Token)
	res, err := cli.Get(ctx, endpoint)
	if err != nil {
		return GetStampResponse{}, err
	}
	if res.StatusCode != http.StatusOK {
		return GetStampResponse{}, fmt.Errorf("Status code=%d", res.StatusCode)
	}
	var gsr struct {
		Success  bool             `json:"success"`
		Response GetStampResponse `json:"response"`
		Errors   []Error          `json:"errors"`
	}

	if err := json.NewDecoder(res.Body).Decode(&gsr); err != nil {
		return GetStampResponse{}, err
	}
	if !gsr.Success {
		return GetStampResponse{}, errors.New("Requesting Stamps API failed")
	}
	return gsr.Response, nil
}

// PostStampParam 打刻リクエストパラメータ
type PostStampParam struct {
	LoginCompanyCode string    // AKASHI企業ID
	Token            string    `json:"token"`               // アクセストークン
	Type             StampType `json:"type,omitempty"`      // 打刻種別
	StampedAt        *AkTime   `json:"stampedAt,omitempty"` // クライアントでの打刻日時
	Timezone         string    `json:"timezone,omitempty"`  // クライアントでのタイムゾーン
}

// PostStampResponse 打刻レスポンス
type PostStampResponse struct {
	LoginCompanyCode string    `json:"login_company_code"` // AKASHI企業ID
	StaffID          int       `json:"staff_id"`           // 従業員ID
	Type             StampType `json:"type"`               // 打刻種別
	StampedAt        *AkTime   `json:"stampedAt"`          // サーバ側での打刻日時
}

// PostStamp 打刻
func PostStamp(ctx context.Context, param PostStampParam) (PostStampResponse, error) {
	if param.LoginCompanyCode == "" {
		return PostStampResponse{}, errors.New("LoginCompanyCode must be set")
	}
	if param.Token == "" {
		return PostStampResponse{}, errors.New("Token must be set")
	}

	endpoint := fmt.Sprintf("/%s/stamps", param.LoginCompanyCode)

	cli := NewClient(param.LoginCompanyCode, param.Token)
	res, err := cli.Post(ctx, endpoint, param)
	if err != nil {
		return PostStampResponse{}, err
	}
	if res.StatusCode != http.StatusOK {
		return PostStampResponse{}, fmt.Errorf("Status code=%d", res.StatusCode)
	}

	var psr struct {
		Success  bool              `json:"success"`
		Response PostStampResponse `json:"response"`
		Errors   []Error           `json:"errors"`
	}

	if err := json.NewDecoder(res.Body).Decode(&psr); err != nil {
		return PostStampResponse{}, err
	}
	if !psr.Success {
		return PostStampResponse{}, errors.New("Requesting Stamp API failed")
	}
	return psr.Response, nil
}
