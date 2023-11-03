package aid

func ErrorBadRequest(errorMessage string) JSON {
	return JSON{
		"errorCode":          "errors.com.epicgames.bad_request",
		"errorMessage":       errorMessage,
		"numericErrorCode":   1001,
		"originatingService": "fortnite",
		"intent":             "prod-live",
	}
}

var (
	ErrorNotFound = JSON{
		"errorCode":          "errors.com.epicgames.common.not_found",
		"errorMessage":       "Resource Not found",
		"numericErrorCode":   1004,
		"originatingService": "fortnite",
		"intent":             "prod-live",
	}

	ErrorInternalServer = map[string]interface{}{
		"errorCode":          "errors.com.epicgames.common.server_error",
		"errorMessage":       "Internal Server Error",
		"numericErrorCode":   1000,
		"originatingService": "fortnite",
		"intent":             "prod-live",
	}
)