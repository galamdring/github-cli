package root

import (
	"testing"

	"github.com/cli/cli/v2/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_applyAccountFlag(t *testing.T) {
	tests := []struct {
		name      string
		cfgStubs  func(*testing.T, *config.AuthConfig)
		account   string
		wantErrContains string
	}{
		{
			name:    "account not in users list",
			account: "unknown-user",
			cfgStubs: func(t *testing.T, authCfg *config.AuthConfig) {
				_, err := authCfg.Login("github.com", "existing-user", "token123", "https", false)
				require.NoError(t, err)
			},
			wantErrContains: `account "unknown-user" is not authenticated with github.com`,
		},
		{
			name:    "error message includes gh auth login guidance",
			account: "ghost",
			cfgStubs: func(t *testing.T, authCfg *config.AuthConfig) {
				_, err := authCfg.Login("github.com", "monalisa", "token456", "https", false)
				require.NoError(t, err)
			},
			wantErrContains: "gh auth login --hostname github.com",
		},
		{
			name:    "valid account with no error",
			account: "monalisa",
			cfgStubs: func(t *testing.T, authCfg *config.AuthConfig) {
				_, err := authCfg.Login("github.com", "monalisa", "monalisa-token", "https", false)
				require.NoError(t, err)
			},
		},
		{
			name:    "selects correct account when multiple users exist",
			account: "hubot",
			cfgStubs: func(t *testing.T, authCfg *config.AuthConfig) {
				_, err := authCfg.Login("github.com", "monalisa", "monalisa-token", "https", false)
				require.NoError(t, err)
				_, err = authCfg.Login("github.com", "hubot", "hubot-token", "https", false)
				require.NoError(t, err)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, _ := config.NewIsolatedTestConfig(t, "")

			if tt.cfgStubs != nil {
				authCfg, ok := cfg.Authentication().(*config.AuthConfig)
				require.True(t, ok, "expected *config.AuthConfig from Authentication()")
				tt.cfgStubs(t, authCfg)
			}

			err := applyAccountFlag(cfg, tt.account)

			if tt.wantErrContains == "" {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.wantErrContains)
			}
		})
	}
}
