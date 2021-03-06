package tests

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/metiago/zbx1/common/helper"
	"github.com/metiago/zbx1/repository"
)

func init() {
	mountBackEndURL()
}

func TestAuthUserOK(t *testing.T) {

	u := repository.User{
		Username: "metiago",
		Password: "zero",
	}

	data, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}

	status := helper.PostHTTP(fmt.Sprintf("%s/%s", baseURL, authPathURL), "", data)

	expected := 200
	if status != expected {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}

func TestAuthUserThanFail(t *testing.T) {

	u := repository.User{
		Username: "metiago",
		Password: "Xh1345",
	}

	data, err := json.Marshal(u)
	if err != nil {
		t.Error(err)
	}

	status := helper.PostHTTP(fmt.Sprintf("%s/%s", baseURL, authPathURL), "", data)

	expected := 401
	if status != expected {
		t.Errorf("Expected is %d but was: %d", expected, status)
	}
}
