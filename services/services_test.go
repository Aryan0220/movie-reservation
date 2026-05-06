package services

import (
	"reflect"
	"testing"
)

func Test_Generate_Token(t *testing.T){
	t.Setenv("JWT_SECRET", "test-secret")

	token, err := GenerateToken(123, "admin")
	if err != nil {
		t.Fatalf("GenerateToken returned error: %v", err)
	}
	if reflect.TypeOf(token) != reflect.TypeOf("") || token == "" {
		t.Errorf("GenerateToken returned invalid token value: %q", token)
	}
}
