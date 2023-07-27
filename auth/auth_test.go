package auth_test

import (
	"fmt"
	"testing"

	"github.com/artificial-lua/example-account-go/auth"
)

const testEmail = "email1@email.email"
const testPassword = "password"
const testName = "name"
const testBirth = "2023-07-26"
const testGender = "M"

const testUpdatePassword = "password2"
const testUpdateSalt = "salt2"
const testUpdateName = "name2"
const testUpdateBirth = "2023-07-27"
const testUpdateGender = "F"

func TestCreateAcccount(t *testing.T) {
	account, err := auth.MakeAccountObject(testEmail, testPassword, "", "", testName, testBirth, testGender)
	if err != nil {
		t.Error(err)
	}
	user1, err := auth.CreateAcccount(account)

	if err != nil {
		t.Error(err)
	}

	fmt.Println(user1)
}

func TestReadAccount(t *testing.T) {
	account, err := auth.ReadAccountByEmail(testEmail)

	if err != nil {
		t.Error(err)
	}

	if account == nil {
		t.Error("account is nil")
	}

	fmt.Println(account.GetEmail())
	fmt.Println(account.GetHash())
	fmt.Println(account.GetSalt())
	fmt.Println(account.GetName())
	fmt.Println(account.GetBirth())
	fmt.Println(account.GetGender())
}

func TestUpdateAccount(t *testing.T) {
	updateAccount, err := auth.MakeAccountObject(
		testEmail,
		testUpdatePassword,
		"",
		testUpdateSalt,
		testUpdateName,
		testUpdateBirth,
		testUpdateGender,
	)
	if err != nil {
		t.Error(err)
	}
	updateAccount.CryptoPassword(updateAccount.GetSalt())

	result, err := auth.UpdateAccount(updateAccount)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)
}

func TestDeleteAccount(t *testing.T) {
	result, err := auth.DeleteAccountByEmail(testEmail)
	if err != nil {
		t.Error(err)
	}

	fmt.Println(result)
}
