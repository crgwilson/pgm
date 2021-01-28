package pg

import (
	"testing"
)

func TestPostgresConfig(t *testing.T) {
	t.Run("valid config struct", func(t *testing.T) {
		testConfig := PostgresConfig{
			Address:  "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			Database: "mydb",
			SslMode:  "verify-full",
		}

		got, err := testConfig.ConnectionString()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		want := "postgres://postgres:password@localhost:5432/mydb?sslmode=verify-full"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("valid config struct no sslmode", func(t *testing.T) {
		testConfig := PostgresConfig{
			Address:  "localhost",
			Port:     5432,
			User:     "postgres",
			Password: "password",
			Database: "mydb",
		}

		got, err := testConfig.ConnectionString()
		if err != nil {
			t.Errorf("expected no error, got %v", err)
		}

		want := "postgres://postgres:password@localhost:5432/mydb"

		if got != want {
			t.Errorf("got %q, want %q", got, want)
		}
	})

	t.Run("invalid config struct (no address)", func(t *testing.T) {
		testConfig := PostgresConfig{
			Port:     5432,
			User:     "postgres",
			Password: "password",
			Database: "mydb",
			SslMode:  "verify-full",
		}

		_, got := testConfig.ConnectionString()

		want := ErrMissingAddress

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})

	t.Run("invalid config struct (no port)", func(t *testing.T) {
		testConfig := PostgresConfig{
			Address:  "localhost",
			User:     "postgres",
			Password: "password",
			Database: "mydb",
			SslMode:  "verify-full",
		}

		_, got := testConfig.ConnectionString()
		want := ErrMissingPort

		if got != want {
			t.Errorf("got %v, want %v", got, want)
		}
	})
}
