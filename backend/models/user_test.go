package models

import (
	"fmt"
	"testing"

	"github.com/2022AA/bytes-linked/backend/pkg/util"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/stretchr/testify/require"
)

func testInitDB() {
	uri := fmt.Sprintf("%s:%s@tcp(%s)/%s?charset=utf8&parseTime=True&loc=Local",
		"admin", "123456", "127.0.0.1", "filestore")
	fmt.Println(uri)
	dbSetUp(uri, "mysql")
}

func TestIsValid(t *testing.T) {
	testInitDB()
	code := "1234"
	ret, err := IsValid(code)
	t.Log(ret, err)
}

func TestCreateInviteCode(t *testing.T) {
	testInitDB()

	for i := 0; i < 1; i++ {
		code := util.GenInvite(uint64(i))
		err := CreateInviteCode(code)
		require.NoError(t, err)
		ret, err := QueryInviteCode(code)
		require.NoError(t, err)
		t.Logf("%+v", ret)
	}

}

func TestQueryInviteCode(t *testing.T) {
	testInitDB()
	code := "1234"
	ret, err := QueryInviteCode(code)
	require.NoError(t, err)
	t.Logf("%+v", ret)
}

func TestUserSignUp(t *testing.T) {
	testInitDB()
	t.Run("invalid code", func(t *testing.T) {
		userName := "test01"
		pwd := "123456"
		inviteCode := "12345"
		ok := UserSignup(userName, pwd, inviteCode, "", "")
		require.False(t, ok)
	})

	t.Run("insert user", func(t *testing.T) {
		userName := "test01"
		pwd := "123456"
		inviteCode := "1234"
		ok := UserSignup(userName, pwd, inviteCode, "", "")
		require.True(t, ok)
	})

	t.Run("insert multi user", func(t *testing.T) {
		userName := "test01"
		pwd := "123456"
		inviteCode := "1234"
		ok := UserSignup(userName, pwd, inviteCode, "", "")
		require.True(t, ok)

		userName = "test02"
		pwd = "22222"
		inviteCode = "12345"
		ok = UserSignup(userName, pwd, inviteCode, "", "")
		require.True(t, ok)
	})
}
