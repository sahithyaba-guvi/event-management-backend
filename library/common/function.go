package commonutils

import (
	commonResp "em_backend/responses/common"
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
