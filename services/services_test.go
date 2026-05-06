package services

import (
	"testing"
	"reflect"
)

func Test_Generate_Token(t *testing.T){
	got1, got2 := GenerateToken(123)
	if reflect.TypeOf(got1) != reflect.TypeOf("") {
		t.Errorf("Generate Token for 123, result not string %T and %T", got1, got2)
	}
}
