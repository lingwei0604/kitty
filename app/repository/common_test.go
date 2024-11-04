package repository

import (
	"testing"
)

func TestIsDuplicateErr(t *testing.T) {
	if !useMysql {
		t.Skip("TestIsDuplicateErr must be run on mysql")
	}

	setUp(t)
	defer tearDown()

	ok := IsDuplicateErr(nil)
	if ok {
		t.Error("IsDuplicateErr(nil) should be false")
	}
	user := user(1)

	err := db.Create(&user).Error
	if IsDuplicateErr(err) {
		t.Error("IsDuplicateErr(err) should be false")
	}

	err = db.Create(&user).Error
	if !IsDuplicateErr(err) {
		t.Error("IsDuplicateErr(err) should be true")
	}
}
