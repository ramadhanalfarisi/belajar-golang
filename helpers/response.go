package helpers

func SuccessResponse(code int, message string, data interface{}) map[string]interface{} {
	var response map[string]interface{}
	response = map[string]interface{}{}
	if data != nil {
		response["code"] = code
		response["status"] = "success"
		response["data"] = data
		response["message"] = message
	} else {
		response["code"] = code
		response["status"] = "success"
		response["message"] = message
	}

	return response
}

func FailedResponse(code int, message string) map[string]interface{} {
	var response map[string]interface{}
	response = map[string]interface{}{}
	response["code"] = code
	response["status"] = "failed"
	response["message"] = message

	return response
}

func InvalidResponse(code int, message []string) map[string]interface{} {
	var response map[string]interface{}
	response = map[string]interface{}{}
	response["code"] = code
	response["status"] = "failed"
	response["message"] = message

	return response
}
