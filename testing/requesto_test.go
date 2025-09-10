package testing

import (
	"testing"

	"github.com/Kaguya233qwq/requesto"
)

func TestClient_NewRequest(t *testing.T) {
	resp, err := requesto.Get("https://www.baidu.com")
	if err != nil {
		t.Fatalf("Failed to create request: %v", err)
	}
	if resp.StatusCode() != 200 {
		t.Errorf("Expected status code 200, got %d", resp.StatusCode())
	}
	text, err := resp.Text()
	if err != nil {
		t.Errorf("Failed to get response text: %v", err)
	} else {
		t.Logf("Response text: %s", text)
	}
}
