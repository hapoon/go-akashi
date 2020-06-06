package akashi

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

// Staff information struct
type Staff struct {
	ID                   int                `json:"staffId"`              // 従業員ID
	LastName             string             `json:"lastName"`             // 姓
	FirstName            string             `json:"firstName"`            // 名
	LastNameKana         string             `json:"lastNameKana"`         // カナ(姓)
	FirstNameKana        string             `json:"firstNameKana"`        // カナ(名)
	Organization         Organization       `json:"organization"`         // 組織(メイン)
	SubGroups            []Organization     `json:"subgroups"`            // 組織(サブグループ)の配列
	EmploymentCategory   EmploymentCategory `json:"employmentCategory"`   // 雇用区分
	Tag                  string             `json:"tag"`                  // タグ
	StaffNum             string             `json:"staffNum"`             // 従業員番号
	IDmNum               string             `json:"idmNum"`               // IDm番号
	CardTypeID           int                `json:"cardTypeId"`           // カード種別
	Remarks              string             `json:"remarks"`              // 備考
	PermissionGroup      PermissionGroup    `json:"permissionGroup"`      // 権限グループ
	ManagedOrganizations []Organization     `json:"managedOrganizations"` // 管理対象組織
}

// EmploymentCategory information struct
type EmploymentCategory struct {
	ID   int    `json:"employmentCategoryId"` // 雇用区分ID
	Name string `json:"Name"`                 // 雇用区分名称
}

// PermissionGroup information struct
type PermissionGroup struct {
	ID   int    `json:"permissionGroupId"` // 権限グループID
	Type int    `json:"permissionType"`    // 権限種別(1:企業管理者、2:一般管理者、3:従業員)
	Name string `json:"name"`              // 権限グループ名
}

// GetStaffParam struct
type GetStaffParam struct {
	LoginCompanyCode string // AKASHI企業ID
	Token            string // アクセストークン
	Target           string // 取得する従業員のトークン
	StaffID          int    // 取得対象の従業員ID
	Page             int    // 管理下にある従業員をすべて取得する場合のページ番号
}

// GetStaffResponse is response struct
type GetStaffResponse struct {
	LoginCompanyCode string  `json:"login_company_code"` // AKASHI企業ID
	Count            int     `json:"count"`              // 取得された従業員数
	TotalCount       int     `json:"total_count"`        // 取得することができる従業員数
	Staffs           []Staff `json:"staffs"`             // 取得した従業員情報の配列
}

// GetStaff fetch staff information
func GetStaff(ctx context.Context, param GetStaffParam) (GetStaffResponse, error) {
	log.Println("Get staff information")
	if param.LoginCompanyCode == "" {
		return GetStaffResponse{}, errors.New("LoginCompanyCode must be set")
	}
	if param.Token == "" {
		return GetStaffResponse{}, errors.New("Token must be set")
	}
	endpoint := fmt.Sprintf("/%s/staffs", param.LoginCompanyCode)
	if param.StaffID != 0 {
		log.Println("staff ID:", param.StaffID)
		endpoint += fmt.Sprintf("/%d", param.StaffID)
	}
	uv := url.Values{}
	uv.Add("token", param.Token)
	if param.Target != "" {
		uv.Add("target", param.Target)
	}
	if param.Page != 0 {
		uv.Add("page", strconv.FormatInt(int64(param.Page), 10))
	}
	q := uv.Encode()
	if q != "" {
		endpoint += "?" + q
	}

	cli := NewClient(param.LoginCompanyCode, param.Token)
	res, err := cli.Get(ctx, endpoint)
	if err != nil {
		return GetStaffResponse{}, err
	}
	if res.StatusCode != http.StatusOK {
		return GetStaffResponse{}, fmt.Errorf("Status code=%d", res.StatusCode)
	}
	var gsr struct {
		Success  bool             `json:"success"`
		Response GetStaffResponse `json:"response"`
		Errors   []Error          `json:"errors"`
	}
	if err := json.NewDecoder(res.Body).Decode(&gsr); err != nil {
		return GetStaffResponse{}, err
	}
	if !gsr.Success {
		return GetStaffResponse{}, errors.New("Requesting AKASHI API failed")
	}
	return gsr.Response, nil
}
