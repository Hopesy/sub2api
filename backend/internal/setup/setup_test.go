package setup

import (
	"os"
	"strings"
	"testing"
)

func TestDecideAdminBootstrap(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		totalUsers int64
		adminUsers int64
		should     bool
		reason     string
	}{
		{
			name:       "empty database should create admin",
			totalUsers: 0,
			adminUsers: 0,
			should:     true,
			reason:     adminBootstrapReasonEmptyDatabase,
		},
		{
			name:       "admin exists should skip",
			totalUsers: 10,
			adminUsers: 1,
			should:     false,
			reason:     adminBootstrapReasonAdminExists,
		},
		{
			name:       "users exist without admin should skip",
			totalUsers: 5,
			adminUsers: 0,
			should:     false,
			reason:     adminBootstrapReasonUsersExistWithoutAdmin,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := decideAdminBootstrap(tc.totalUsers, tc.adminUsers)
			if got.shouldCreate != tc.should {
				t.Fatalf("shouldCreate=%v, want %v", got.shouldCreate, tc.should)
			}
			if got.reason != tc.reason {
				t.Fatalf("reason=%q, want %q", got.reason, tc.reason)
			}
		})
	}
}

func TestSetupDefaultAdminConcurrency(t *testing.T) {
	t.Run("simple mode admin uses higher concurrency", func(t *testing.T) {
		t.Setenv("RUN_MODE", "simple")
		if got := setupDefaultAdminConcurrency(); got != simpleModeAdminConcurrency {
			t.Fatalf("setupDefaultAdminConcurrency()=%d, want %d", got, simpleModeAdminConcurrency)
		}
	})

	t.Run("standard mode keeps existing default", func(t *testing.T) {
		t.Setenv("RUN_MODE", "standard")
		if got := setupDefaultAdminConcurrency(); got != defaultUserConcurrency {
			t.Fatalf("setupDefaultAdminConcurrency()=%d, want %d", got, defaultUserConcurrency)
		}
	})
}

func TestNormalizeRedisHost(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		raw     string
		want    string
		wantErr string
	}{
		{name: "plain host", raw: "stable-warthog-101481.upstash.io", want: "stable-warthog-101481.upstash.io"},
		{name: "https url", raw: "https://stable-warthog-101481.upstash.io", want: "stable-warthog-101481.upstash.io"},
		{name: "rediss url with port", raw: "rediss://stable-warthog-101481.upstash.io:6379", want: "stable-warthog-101481.upstash.io"},
		{name: "host with port", raw: "stable-warthog-101481.upstash.io:6379", want: "stable-warthog-101481.upstash.io"},
		{name: "ipv6", raw: "[2001:db8::1]", want: "2001:db8::1"},
		{name: "path not allowed", raw: "https://stable-warthog-101481.upstash.io/foo", wantErr: "should not contain path"},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := normalizeRedisHost(tc.raw)
			if tc.wantErr != "" {
				if err == nil || !strings.Contains(err.Error(), tc.wantErr) {
					t.Fatalf("normalizeRedisHost() error = %v, want contains %q", err, tc.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("normalizeRedisHost() error = %v", err)
			}
			if got != tc.want {
				t.Fatalf("normalizeRedisHost() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestNormalizeRedisConfig(t *testing.T) {
	t.Parallel()

	cfg := &RedisConfig{Host: "https://stable-warthog-101481.upstash.io", Port: 6379}
	if err := normalizeRedisConfig(cfg); err != nil {
		t.Fatalf("normalizeRedisConfig() error = %v", err)
	}
	if cfg.Host != "stable-warthog-101481.upstash.io" {
		t.Fatalf("cfg.Host = %q, want normalized host", cfg.Host)
	}
}

func TestNormalizeRedisHostRejectsAmbiguousColonHost(t *testing.T) {
	t.Parallel()

	_, err := normalizeRedisHost("https://stable-warthog-101481.upstash.io:6379:extra")
	if err == nil {
		t.Fatal("normalizeRedisHost() error = nil, want error")
	}
}

func TestChooseMaintenanceDatabaseReturnsTargetWhenAvailable(t *testing.T) {
	t.Parallel()

	cfg := &DatabaseConfig{Host: "127.0.0.1", Port: 1, User: "postgres", Password: "postgres", DBName: "sub2api", SSLMode: "disable"}
	_, err := chooseMaintenanceDatabase(cfg)
	if err == nil {
		t.Fatal("chooseMaintenanceDatabase() error = nil, want error")
	}
	if !strings.Contains(err.Error(), "maintenance database") {
		t.Fatalf("chooseMaintenanceDatabase() error = %v, want maintenance database context", err)
	}
}

func TestWriteConfigFileKeepsDefaultUserConcurrency(t *testing.T) {
	t.Setenv("RUN_MODE", "simple")
	t.Setenv("DATA_DIR", t.TempDir())

	if err := writeConfigFile(&SetupConfig{}); err != nil {
		t.Fatalf("writeConfigFile() error = %v", err)
	}

	data, err := os.ReadFile(GetConfigFilePath())
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}

	if !strings.Contains(string(data), "user_concurrency: 5") {
		t.Fatalf("config missing default user concurrency, got:\n%s", string(data))
	}
}
