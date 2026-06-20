package db

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDBName(t *testing.T) {
	conf := map[string]interface{}{
		"database": map[string]interface{}{
			"name": "kohrvid_test",
		},
	}

	assert.Equal(
		t,
		"kohrvid_test",
		DBName(conf),
		"Should return the correct database name if one is supplied",
	)
}

func TestDBNameWithOldFormat(t *testing.T) {
	conf := map[string]interface{}{"database_name": "kohrvid_test"}

	assert.Equal(
		t,
		"kohrvid_test",
		DBName(conf),
		"Should return the correct database name if one is supplied",
	)
}

func TestDBNameWhenMissing(t *testing.T) {
	conf := map[string]interface{}(nil)

	assert.Equal(
		t,
		"",
		DBName(conf),
		"Should return an empty string if no database name is supplied",
	)
}

func TestDBUser(t *testing.T) {
	conf := map[string]interface{}{
		"database": map[string]interface{}{
			"user": "kohrvid",
		},
	}

	assert.Equal(
		t,
		"kohrvid",
		DBUser(conf),
		"Should return the correct database user if one is supplied",
	)
}

func TestDBUserWithOldFormat(t *testing.T) {
	conf := map[string]interface{}{"database_user": "kohrvid"}

	assert.Equal(
		t,
		"kohrvid",
		DBUser(conf),
		"Should return the correct database user if one is supplied",
	)
}

func TestDBUserWhenMissing(t *testing.T) {
	conf := map[string]interface{}(nil)

	assert.Equal(
		t,
		"",
		DBUser(conf),
		"Should return an empty string if no database user is supplied",
	)
}

func TestHost(t *testing.T) {
	conf := map[string]interface{}{
		"database": map[string]interface{}{
			"host": "kohrvid.com",
		},
	}

	assert.Equal(
		t,
		"kohrvid.com",
		host(conf),
		"Should return the correct database host if one is supplied",
	)
}

func TestHostWithOldFormat(t *testing.T) {
	conf := map[string]interface{}{"database_host": "kohrvid.com"}

	assert.Equal(
		t,
		"kohrvid.com",
		host(conf),
		"Should return the correct database host if one is supplied",
	)
}

func TestMissingHost(t *testing.T) {
	conf := map[string]interface{}{}

	assert.Equal(
		t,
		"localhost",
		host(conf),
		"Should return localhost if no host is supplied",
	)
}

func TestPort(t *testing.T) {
	conf := map[string]interface{}{
		"database": map[string]interface{}{
			"port": 5433,
		},
	}

	assert.Equal(
		t,
		5433,
		port(conf),
		"Should return the correct port number if one is supplied",
	)
}

func TestPortWithOldFormat(t *testing.T) {
	conf := map[string]interface{}{"database_port": 5433}

	assert.Equal(
		t,
		5433,
		port(conf),
		"Should return the correct port number if one is supplied",
	)
}

func TestMissingPort(t *testing.T) {
	conf := map[string]interface{}{}

	assert.Equal(
		t,
		5432,
		port(conf),
		"Should return localport if no port is supplied",
	)
}

func TestSslMode(t *testing.T) {
	ssl_modes := []string{
		"verify-full", "verify-ca",
		"require", "prefer",
		"allow", "disable",
	}

	for _, mode := range ssl_modes {
		conf := map[string]interface{}{
			"database": map[string]interface{}{
				"ssl_mode": mode,
			},
		}

		assert.Equal(
			t,
			mode,
			sslMode(conf),
			"Should return the correct SSL mode if one is supplied",
		)
	}
}

func TestSslModeWithOldFormat(t *testing.T) {
	ssl_modes := []string{
		"verify-full", "verify-ca",
		"require", "prefer",
		"allow", "disable",
	}

	for _, mode := range ssl_modes {
		conf := map[string]interface{}{"ssl_mode": mode}

		assert.Equal(
			t,
			mode,
			sslMode(conf),
			"Should return the correct SSL mode if one is supplied",
		)
	}
}

func TestMissingSslMode(t *testing.T) {
	conf := map[string]interface{}{}

	assert.Equal(
		t,
		"disable",
		sslMode(conf),
		"Should return disable if no SSL mode is supplied",
	)
}

func TestInvalidSslMode(t *testing.T) {
	conf := map[string]interface{}{"ssl_mode": "required"}

	assert.Equal(
		t,
		"disable",
		sslMode(conf),
		"Should return disable if a mode with a typo is supplied",
	)
}
