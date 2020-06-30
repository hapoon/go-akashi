package akashi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// PostTokenReissueParam トークン再発行リクエストパラメータ
type PostTokenReissueParam struct {
	LoginCompanyCode string // AKASHI企業ID
	Token            string `json:"token"` // アクセストークン
}

// PostTokenReissueResponse トークン再発行レスポンス
type PostTokenReissueResponse struct {
	LoginCompanyCode string  `json:"login_company_code"` // ログイン企業ID
	StaffID          int     `json:"staff_id"`           // 従業員ID
	AgencyManagerID  int     `json:"agency_manager_id"`  // -
	Token            string  `json:"token"`              // アクセストークン
	ExpiredAt        *AkTime `json:"expired_at"`         // アクセストークンの有効期限
}

// PostTokenReissue トークン再発行
func PostTokenReissue(ctx context.Context, param PostTokenReissueParam) (PostTokenReissueResponse, error) {
	if param.LoginCompanyCode == "" {
		return PostTokenReissueResponse{}, errors.New("LoginCompanyCode must be set")
	}
	if param.Token == "" {
		return PostTokenReissueResponse{}, errors.New("Token must be set")
	}

	endpoint := fmt.Sprintf("/token/reissue/%s", param.LoginCompanyCode)

	cli := NewClient(param.LoginCompanyCode, param.Token)
	res, err := cli.Post(ctx, endpoint, param)
	if err != nil {
		return PostTokenReissueResponse{}, err
	}
	if res.StatusCode != http.StatusOK {
		return PostTokenReissueResponse{}, fmt.Errorf("Status code=%d", res.StatusCode)
	}

	var ptr struct {
		Success  bool                     `json:"success"`
		Response PostTokenReissueResponse `json:"response"`
		Errors   []Error                  `json:"errors"`
	}

	if err := json.NewDecoder(res.Body).Decode(&ptr); err != nil {
		return PostTokenReissueResponse{}, err
	}
	if !ptr.Success {
		return PostTokenReissueResponse{}, errors.New("Requesting Token Reissue API failed")
	}
	return ptr.Response, nil
}
