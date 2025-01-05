package commonutils

import (
	"context"
	mongoSetup "em_backend/configs/mongo"
	commonModel "em_backend/models/common"
	commonResp "em_backend/responses/common"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
)

func CreateSuccessResponse(resp *commonResp.SuccessResponse) commonResp.SuccessResponse {
	if resp == nil {
		resp = &commonResp.SuccessResponse{} // Initialize a new SuccessResponse if nil
	}
	if resp.Status == "" {
		resp.Status = "success"
	}
	if resp.Message == "" {
		resp.Message = "Operation was successful."
	}
	resp.Access = true
	if resp.Data == nil {
		resp.Data = map[string]interface{}{} // Default empty map
	}
	return *resp
}

func CreateFailureResponse(resp *commonResp.FailureResponse) commonResp.FailureResponse {
	if resp == nil {
		resp = &commonResp.FailureResponse{} // Initialize a new FailureResponse if nil
	}
	if resp.Status == "" {
		resp.Status = "failure"
	}
	if resp.Message == "" {
		resp.Message = "Operation failed."
	}
	return *resp
}

func GetAdminList() commonModel.Admins {
	db, _, err := mongoSetup.ConnectMongo("admin_list")
	if err != nil {
		fmt.Println(err)
	}
	defer db.Client().Disconnect(context.TODO())
	mycol := db.Collection("admin_list")
	var admins commonModel.Admins
	errr := mycol.FindOne(context.TODO(), bson.D{}).Decode(&admins)
	if errr != nil {
		fmt.Println(errr)
	}
	return admins
}
func CheckAdmin(mail string) bool {
	var adminList = GetAdminList()
	flag := 0
	for _, admin := range adminList.Admin {
		if mail == admin {
			flag = 1
			break
		}
	}
	if flag == 0 {
		return false
	} else {
		return true
	}
}
