package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadConfig_OK(t *testing.T) {
	t.Parallel()

	dir := t.TempDir()
	p := filepath.Join(dir, "config.yaml")
	err := os.WriteFile(p, []byte(
		"ip: 127.0.0.1\n"+
			"port: 8080\n"+
			"api:\n"+
			"  secret_key: test\n"+
			"  front_address: http://localhost:5173\n"+
			"  image_path: ./materials\n"+
			"db:\n"+
			"  host: localhost\n"+
			"  port: 5432\n"+
			"  user: u\n"+
			"  password: p\n"+
			"  dbname: db\n"+
			"usecase:\n"+
			"  default_message: ok\n"+
			"unknown_field: ignored\n",
	), 0o600)
	if err != nil {
		t.Fatalf("write temp config: %v", err)
	}

	cfg, err := LoadConfig(p)
	if err != nil {
		t.Fatalf("LoadConfig: %v", err)
	}
	if cfg.IP != "127.0.0.1" || cfg.Port != 8080 {
		t.Fatalf("unexpected addr: %#v", cfg)
	}
	if cfg.API.SecretKey != "test" {
		t.Fatalf("unexpected secret: %#v", cfg.API)
	}
	if cfg.DB.Port != 5432 || cfg.DB.DBname != "db" {
		t.Fatalf("unexpected db: %#v", cfg.DB)
	}
}

func TestLoadConfig_NoFile(t *testing.T) {
	t.Parallel()

	_, err := LoadConfig("no-such-file.yaml")
	if err == nil {
		t.Fatalf("expected error")
	}
}

/* old heredoc variant (kept for reference)
ip: 127.0.0.1
port: 8080
api:
  secret_key: test
  front_address: http://localhost:5173
		image_path: ./materials
		db:
		host: localhost
		port: 5432
		user: u
		password: p
		dbname: db
usecase:
  default_message: ok
unknown_field: ignored
*/
