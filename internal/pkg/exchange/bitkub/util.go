package bitkub

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/url"
	"strconv"
	"strings"
)

func genSign(apiSecret, payloadString string) string {
	mac := hmac.New(sha256.New, []byte(apiSecret))
	mac.Write([]byte(payloadString))
	return hex.EncodeToString(mac.Sum(nil))
}

func genQueryParam(baseURL string, queryParams map[string]interface{}) string {
	u, _ := url.Parse(baseURL)
	q := u.Query()
	for key, value := range queryParams {
		q.Set(key, fmt.Sprintf("%v", value))
	}
	u.RawQuery = q.Encode()
	return strings.Replace(u.String(), baseURL, "", 1)
}

func convertCustomToSeconds(durationStr string) (float64, error) {
	// Handle "D" for days manually
	if strings.HasSuffix(durationStr, "D") {
		// Remove the "D" suffix
		daysStr := strings.TrimSuffix(durationStr, "D")
		days, err := strconv.Atoi(daysStr)
		if err != nil {
			return 0, err
		}
		// Convert days to seconds (1 day = 86400 seconds)
		return float64(days * 86400), nil
	}

	// Otherwise, treat it as minutes
	minutes, err := strconv.Atoi(durationStr)
	if err != nil {
		return 0, err
	}
	// Convert minutes to seconds
	return float64(minutes * 60), nil
}

func mustMarshal(v interface{}) []byte {
	b, err := json.Marshal(v)
	if err != nil {
		panic(err)
	}
	return b
}
