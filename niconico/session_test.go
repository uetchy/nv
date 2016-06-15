package niconico

import (
	"os"
	"strings"
	"testing"
)

func TestGetSessionKey(t *testing.T) {
	email := os.Getenv("TEST_EMAIL")
	pass := os.Getenv("TEST_PASS")
	_, v := GetSessionKey(email, pass)
	if !strings.HasPrefix(v, "user_session=user_session_") {
		t.Error("Expected has prefix 'user_session=user_session_', got", v)
	}
}
