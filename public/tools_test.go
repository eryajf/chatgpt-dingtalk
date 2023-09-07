package public

import (
	"testing"

	"github.com/eryajf/chatgpt-dingtalk/config"
)

func TestCheckRequestWithCredentials_Pass_WithNilConfig(t *testing.T) {
	Config = &config.Configuration{
		Credentials: nil,
	}
	clientId, pass := CheckRequestWithCredentials("ts", "sg")
	if !pass {
		t.Errorf("pass should be true, but false")
		return
	}
	if len(clientId) > 0 {
		t.Errorf("client id should be empty")
		return
	}
}

func TestCheckRequestWithCredentials_Pass_WithEmptyConfig(t *testing.T) {
	Config = &config.Configuration{
		Credentials: []config.Credential{},
	}
	clientId, pass := CheckRequestWithCredentials("ts", "sg")
	if !pass {
		t.Errorf("pass should be true, but false")
		return
	}
	if len(clientId) > 0 {
		t.Errorf("client id should be empty")
		return
	}
}

func TestCheckRequestWithCredentials_Pass_WithValidConfig(t *testing.T) {
	Config = &config.Configuration{
		Credentials: []config.Credential{
			config.Credential{
				ClientID:     "client-id-for-test",
				ClientSecret: "client-secret-for-test",
			},
		},
	}
	clientId, pass := CheckRequestWithCredentials("1684493546276", "nwBJQmaBLv9+5/sSS/66jcFc1/kGY5wo38L88LOGfRU=")
	if !pass {
		t.Errorf("pass should be true, but false")
		return
	}
	if clientId != "client-id-for-test" {
		t.Errorf("client id should be \"%s\", but \"%s\"", "client-id-for-test", clientId)
		return
	}
}

func TestCheckRequestWithCredentials_Failed_WithInvalidConfig(t *testing.T) {
	Config = &config.Configuration{
		Credentials: []config.Credential{
			config.Credential{
				ClientID:     "client-id-for-test",
				ClientSecret: "invalid-client-secret-for-test",
			},
		},
	}
	clientId, pass := CheckRequestWithCredentials("1684493546276", "nwBJQmaBLv9+5/sSS/66jcFc1/kGY5wo38L88LOGfRU=")
	if pass {
		t.Errorf("pass should be false, but true")
		return
	}
	if clientId != "" {
		t.Errorf("client id should be empty")
		return
	}
}
