package models

import (
	"testing"

	v2 "github.com/2022AA/bytes-linked/backend/pkg/logging/v2"
	"github.com/stretchr/testify/require"
)

func TestGenerateUserSecretFile(t *testing.T) {
	v2.SetLevel(v2.DebugLevel)
	GenerateUserSecretFile("sadf")
}

func TestQuerySecretFile(t *testing.T) {
	testInitDB()
	username := "test01"
	secretFile, err := QuerySecretFile(username, SecretType_Private)
	require.NoError(t, err)
	t.Logf("%+v", secretFile)
}
