package app

import (
	"bytes"
	"encoding/json"
	v1 "k8s.io/api/authentication/v1"
	"strings"
)

// SanitizingFilter implements the LogFilter interface from klog with custom functions to detect and mask tokens.
type SanitizingFilter struct{}

// Filter is the filter function for non-formatting logging functions of klog.
func (sf *SanitizingFilter) Filter(args []interface{}) []interface{} {
	for i, v := range args {
		if strValue, ok := v.(string); ok {
			if strings.Contains(strValue, `"kind":"TokenReview"`) {
				args[i] = maskTokenInLog(strValue)
			}
		}
	}
	return args
}

// FilterF is the filter function for formatting logging functions of klog.
func (sf *SanitizingFilter) FilterF(format string, args []interface{}) (string, []interface{}) {
	for i, v := range args {
		if strValue, ok := v.(string); ok {
			if strings.Contains(strValue, `"kind":"TokenReview"`) {
				args[i] = maskTokenInLog(strValue)
			}
		}
	}
	return format, args
}

// FilterS is the filter function for structured logging functions of klog.
func (sf *SanitizingFilter) FilterS(msg string, keysAndValues []interface{}) (string, []interface{}) {
	for i, v := range keysAndValues {
		if strValue, ok := v.(string); ok {
			if strings.Contains(strValue, `"kind":"TokenReview"`) {
				keysAndValues[i] = maskTokenInLog(strValue)
			}
		}
	}
	return msg, keysAndValues
}

func maskTokenInLog(logStr string) string {
	var tokenReview v1.TokenReview
	if err := json.Unmarshal([]byte(logStr), &tokenReview); err != nil {
		return "<log content masked due to unmarshal failure>"
	}

	if tokenReview.Spec.Token != "" {
		tokenReview.Spec.Token = "<masked>"
	}

	var buf bytes.Buffer
	encoder := json.NewEncoder(&buf)
	encoder.SetEscapeHTML(false)
	if err := encoder.Encode(tokenReview); err != nil {
		return "<log content masked due to encoding failure>"
	}
	return buf.String()
}