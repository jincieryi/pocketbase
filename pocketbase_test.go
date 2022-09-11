package pocketbase

import (
	"os"
	"testing"
)

func TestNew(t *testing.T) {
	// copy os.Args
	originalArgs := []string{}
	copy(originalArgs, os.Args)
	defer func() {
		// restore os.Args
		copy(os.Args, originalArgs)
	}()

	// change os.Args
	os.Args = os.Args[0:1]
	os.Args = append(
		os.Args,
		"--dir=test_dir",
		"--encryptionEnv=test_encryption_env",
		"--debug=true",
		"--mysqlDsnEnv=test_dsn_env",
	)

	app := New()

	if app == nil {
		t.Fatal("Expected initialized PocketBase instance, got nil")
	}

	if app.RootCmd == nil {
		t.Fatal("Expected RootCmd to be initialized, got nil")
	}

	if app.appWrapper == nil {
		t.Fatal("Expected appWrapper to be initialized, got nil")
	}

	if app.DataDir() != "test_dir" {
		t.Fatalf("Expected app.DataDir() %q, got %q", "test_dir", app.DataDir())
	}

	if app.EncryptionEnv() != "test_encryption_env" {
		t.Fatalf("Expected app.EncryptionEnv() test_encryption_env, got %q", app.EncryptionEnv())
	}

	if app.IsDebug() != true {
		t.Fatal("Expected app.IsDebug() true, got false")
	}

	if app.MysqlDsnEnv() != "test_dsn_env" {
		t.Fatalf("Expected app.MysqlDsnEnv() test_dsn_env, got %q", app.MysqlDsnEnv())
	}
}

func TestNewWithConfig(t *testing.T) {
	app := NewWithConfig(Config{
		DefaultDebug:         true,
		DefaultDataDir:       "test_dir",
		DefaultEncryptionEnv: "test_encryption_env",
		HideStartBanner:      true,
		DefaultMysqlDsnEnv:   "test_dsn_env",
	})

	if app == nil {
		t.Fatal("Expected initialized PocketBase instance, got nil")
	}

	if app.RootCmd == nil {
		t.Fatal("Expected RootCmd to be initialized, got nil")
	}

	if app.appWrapper == nil {
		t.Fatal("Expected appWrapper to be initialized, got nil")
	}

	if app.hideStartBanner != true {
		t.Fatal("Expected app.hideStartBanner to be true, got false")
	}

	if app.DataDir() != "test_dir" {
		t.Fatalf("Expected app.DataDir() %q, got %q", "test_dir", app.DataDir())
	}

	if app.EncryptionEnv() != "test_encryption_env" {
		t.Fatalf("Expected app.EncryptionEnv() test_encryption_env, got %q", app.EncryptionEnv())
	}

	if app.MysqlDsnEnv() != "test_dsn_env" {
		t.Fatalf("Expected app.MysqlDsnEnv() test_dsn_env, got %q", app.MysqlDsnEnv())
	}

	if app.IsDebug() != true {
		t.Fatal("Expected app.IsDebug() true, got false")
	}

}

func TestNewWithConfigAndFlags(t *testing.T) {
	// copy os.Args
	originalArgs := []string{}
	copy(originalArgs, os.Args)
	defer func() {
		// restore os.Args
		copy(os.Args, originalArgs)
	}()

	// change os.Args
	os.Args = os.Args[0:1]
	os.Args = append(
		os.Args,
		"--dir=test_dir_flag",
		"--encryptionEnv=test_encryption_env_flag",
		"--debug=false",
		"--mysqlDsnEnv=test_dsn_env",
	)

	app := NewWithConfig(Config{
		DefaultDebug:         true,
		DefaultDataDir:       "test_dir",
		DefaultEncryptionEnv: "test_encryption_env",
		HideStartBanner:      true,
		DefaultMysqlDsnEnv:   "test_mysql_dsn_env",
	})

	if app == nil {
		t.Fatal("Expected initialized PocketBase instance, got nil")
	}

	if app.RootCmd == nil {
		t.Fatal("Expected RootCmd to be initialized, got nil")
	}

	if app.appWrapper == nil {
		t.Fatal("Expected appWrapper to be initialized, got nil")
	}

	if app.hideStartBanner != true {
		t.Fatal("Expected app.hideStartBanner to be true, got false")
	}

	if app.DataDir() != "test_dir_flag" {
		t.Fatalf("Expected app.DataDir() %q, got %q", "test_dir_flag", app.DataDir())
	}

	if app.EncryptionEnv() != "test_encryption_env_flag" {
		t.Fatalf("Expected app.DataDir() %q, got %q", "test_encryption_env_flag", app.EncryptionEnv())
	}

	if app.IsDebug() != false {
		t.Fatal("Expected app.IsDebug() false, got true")
	}

	if app.MysqlDsnEnv() != "test_dsn_env" {
		t.Fatalf("Expected app.MysqlDsnEnv() %q, got %q", "test_dsn_env", app.MysqlDsnEnv())
	}
}
