package repository

import (
	"context"
	"database/sql"
	"math/rand"
	"strconv"
	"testing"
	"time"

	"github.com/lingwei0604/kitty/app/entity"
	"github.com/pkg/errors"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func TestGetFromWechat(t *testing.T) {
	if !useMysql {
		t.Skip("GetFromWechat tests must be run on mysql")
	}
	setUp(t)
	defer tearDown()
	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
	ctx := context.Background()
	u, err := userRepo.GetFromWechat(ctx, "", "foo", &entity.Device{Suuid: "bar", SMID: "baz"}, entity.User{UserName: "baz"})
	if err != nil {
		t.Fatal(err)
	}
	if u.WechatOpenId.String != "foo" {
		t.Fatalf("want foo, got %v", u.WechatOpenId)
	}
	if u.Devices[0].Suuid != "bar" {
		t.Fatalf("want bar, got %s", u.Devices[0].Suuid)
	}
	if u.CommonSUUID != "bar" {
		t.Fatalf("want bar, got %s", u.CommonSUUID)
	}
	if u.UserName != "baz" {
		t.Fatalf("want baz, got %s", u.UserName)
	}
	if !u.IsNew {
		t.Fatalf("user must be new")
	}
	u2, err := userRepo.GetFromWechat(ctx, "", "foo", &entity.Device{Suuid: "bar2", SMID: "baz2"}, entity.User{UserName: "baz2"})
	if err != nil {
		t.Fatal(err)
	}
	if u2.WechatOpenId.String != "foo" {
		t.Fatalf("want foo, got %v", u2.WechatOpenId)
	}
	if u2.Devices[0].Suuid != "bar2" {
		t.Fatalf("want bar2, got %s", u2.Devices[0].Suuid)
	}
	if u2.CommonSUUID != "bar2" {
		t.Fatalf("want bar2, got %s", u2.CommonSUUID)
	}
	if u2.CommonSMID != "baz2" {
		t.Fatalf("want baz2, got %s", u2.CommonSMID)
	}
	if u2.UserName != "baz" {
		t.Fatalf("want baz, got %s", u2.Devices[0].Suuid)
	}
	if u2.IsNew {
		t.Fatalf("user must be new")
	}

	u3, err := userRepo.GetFromWechat(ctx, "", "xxx", &entity.Device{Suuid: "bar2", SMID: "baz2"}, entity.User{UserName: "baz2"})
	if err != nil {
		t.Fatal(err)
	}
	if !u3.IsNew {
		t.Fatalf("user must be new")
	}
}

func TestGetFromMobile(t *testing.T) {
	if !useMysql {
		t.Skip("GetFromMobile tests must be run on mysql")
	}
	setUp(t)
	defer tearDown()
	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
	ctx := context.Background()
	u, err := userRepo.GetFromMobile(ctx, "", "110", &entity.Device{Suuid: "bar", SMID: "baz"})
	if err != nil {
		t.Fatal(err)
	}
	if u.Mobile.String != "110" {
		t.Fatalf("want 110, got %v", u.Mobile.String)
	}
	if u.Devices[0].Suuid != "bar" {
		t.Fatalf("want bar, got %s", u.Devices[0].Suuid)
	}
	if u.CommonSUUID != "bar" {
		t.Fatalf("want bar, got %s", u.CommonSUUID)
	}
	if !u.IsNew {
		t.Fatalf("user must be new")
	}
	u2, err := userRepo.GetFromMobile(ctx, "", "110", &entity.Device{Suuid: "bar2", SMID: "baz2"})
	if err != nil {
		t.Fatal(err)
	}
	if u2.Mobile.String != "110" {
		t.Fatalf("want foo, got %v", u2.Mobile)
	}
	if u2.Devices[0].Suuid != "bar2" {
		t.Fatalf("want bar2, got %s", u2.Devices[0].Suuid)
	}
	if u2.CommonSMID != "baz2" {
		t.Fatalf("want baz2, got %s", u2.CommonSMID)
	}
	if u2.CommonSUUID != "bar2" {
		t.Fatalf("want bar2, got %s", u2.CommonSUUID)
	}
	if u2.IsNew {
		t.Fatalf("user must not be new")
	}

	u3, err := userRepo.GetFromMobile(ctx, "", "119", &entity.Device{Suuid: "bar2", SMID: "baz2"})
	if err != nil {
		t.Fatal(err)
	}
	if u3.Mobile.String != "119" {
		t.Fatalf("want 119, got %v", u2.Mobile)
	}
	if !u3.IsNew {
		t.Fatalf("user must not be new")
	}
}

func TestGetFromDevice(t *testing.T) {
	setUp(t)
	defer tearDown()
	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
	ctx := context.Background()
	u, err := userRepo.GetFromDevice(ctx, "", "110", &entity.Device{Suuid: "bar", SMID: "baz"})
	if err != nil {
		t.Fatal(err)
	}
	if u.CommonSUUID != "110" {
		t.Fatalf("want 110, got %s", u.CommonSUUID)
	}
	if u.Devices[0].Suuid != "bar" {
		t.Fatalf("want bar, got %s", u.Devices[0].Suuid)
	}
	u2, err := userRepo.GetFromDevice(ctx, "", "110", &entity.Device{Suuid: "bar2"})
	if err != nil {
		t.Fatal(err)
	}
	if u2.CommonSUUID != "110" {
		t.Fatalf("want foo, got %s", u2.CommonSUUID)
	}
	if u2.CommonSMID != "baz" {
		t.Fatalf("want baz, got %s", u2.CommonSMID)
	}
	if u2.Devices[0].Suuid != "bar2" {
		t.Fatalf("want bar2, got %s", u2.Devices[0].Suuid)
	}
}

func TestGetSave(t *testing.T) {
	setUp(t)
	defer tearDown()
	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
	ctx := context.Background()
	user := entity.User{}
	user.ID = 50
	err := userRepo.Save(ctx, &user)
	if err != nil {
		t.Fatal(err)
	}
	u, err := userRepo.Get(ctx, 50)
	if err != nil {
		t.Fatal(err)
	}
	if u.ID != user.ID {
		t.Fatalf("want %d, go %d", user.ID, u.ID)
	}
}

func TestUserRepo_Delete(t *testing.T) {
	setUp(t)
	defer tearDown()
	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
	ctx := context.Background()
	user := entity.User{Model: gorm.Model{ID: uint(1)}, UserName: "hello"}
	_ = userRepo.Save(ctx, &user)

	err := userRepo.Delete(ctx, 1)
	assert.NoError(t, err)

	_, err = userRepo.Get(ctx, 1)
	assert.True(t, errors.Is(err, ErrRecordNotFound))
}

func TestUserRepo_GetAll(t *testing.T) {
	setUp(t)
	defer tearDown()

	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
	ctx := context.Background()
	for i := 1; i < 5; i++ {
		user := entity.User{Model: gorm.Model{ID: uint(i)}, UserName: "hello"}
		_ = userRepo.Save(ctx, &user)
	}
	users, err := userRepo.GetAll(ctx, clause.Where{Exprs: []clause.Expression{clause.IN{
		Column: "id",
		Values: []interface{}{1, 2, 3, 4},
	}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 4 {
		t.Fatal("there should be four users")
	}
	users, err = userRepo.GetAll(ctx, clause.Where{Exprs: []clause.Expression{clause.Like{
		Column: "user_name",
		Value:  "%ell%",
	}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 4 {
		t.Fatal("there should be four users")
	}
	users, err = userRepo.GetAll(ctx, clause.Where{Exprs: []clause.Expression{clause.Gt{
		Column: "created_at",
		Value:  time.Unix(500, 0),
	}}})
	if err != nil {
		t.Fatal(err)
	}
	if len(users) != 4 {
		t.Fatal("there should be four users")
	}
}

func TestUserRepo_Get(t *testing.T) {
	if !useMysql {
		t.Skip("count all tests must be run on mysql")
	}

	setUp(t)
	defer tearDown()

	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
	ctx := context.Background()

	// ADD SOME USERS WITH DEVICES
	for i := 1; i < 5; i++ {
		user := entity.User{Model: gorm.Model{ID: uint(i)}, UserName: "hello", PackageName: "com.foo.bar"}
		user.Devices = []entity.Device{{
			Imei: "1234567",
			Hash: cast.ToString(i),
		}}
		_ = userRepo.Save(ctx, &user)
	}

	users, err := userRepo.Get(ctx, 1)
	if err != nil {
		t.Fatalf("there should be no err, got %s", err)
	}
	if len(users.Devices) != 1 {
		t.Fatalf("there should be one devices, got %d", len(users.Devices))
	}
}

func TestUserRepo_GetByDevice(t *testing.T) {
	if !useMysql {
		t.Skip("count all tests must be run on mysql")
	}

	setUp(t)
	defer tearDown()

	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
	ctx := context.Background()

	// ADD SOME USERS WITH DEVICES
	for i := 1; i < 5; i++ {
		user := entity.User{Model: gorm.Model{ID: uint(i)}, UserName: "hello", PackageName: "com.foo.bar"}
		user.Devices = []entity.Device{{
			Imei: "1234567",
			Hash: cast.ToString(i),
		}}
		_ = userRepo.Save(ctx, &user)
	}

	// ADD SOME ADDITIONAL DEVICES. THOSE SHOULD NOT BE SELECTED
	for i := 1; i < 5; i++ {
		device := entity.Device{
			UserID: uint(i),
			Imei:   "1234567",
			Hash:   cast.ToString(i),
		}
		userRepo.db.Save(&device)
	}

	users, _ := userRepo.GetByDevice(ctx, "com.foo.bar", entity.Device{
		Imei: "1234567",
	})
	if len(users) != 4 {
		t.Fatalf("there should be four users, got %d", len(users))
	}
}

func TestGetRecentByDevice(t *testing.T) {
	if !useMysql {
		t.Skip("count all tests must be run on mysql")
	}

	setUp(t)
	defer tearDown()
	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
	ctx := context.Background()
	for i := 1; i < 5; i++ {
		user := entity.User{PackageName: "foo", CommonSUUID: strconv.Itoa(i)}
		user.AddNewDevice(&entity.Device{Oaid: "baz", Suuid: user.CommonSUUID, Idfa: "baz", AndroidId: "baz baz"})
		_ = userRepo.Save(ctx, &user)
	}

	cases := []struct {
		name      string
		req       *entity.Device
		wantSuuid string
		wantID    uint
		wantError error
	}{
		{name: "same at suuid and oaid", req: &entity.Device{Suuid: "2", Oaid: "baz"}, wantSuuid: "2", wantID: 4},
		{name: "oaid and android", req: &entity.Device{Idfa: "baz"}, wantSuuid: "2", wantID: 4},
		{name: "idfa", req: &entity.Device{Idfa: "baz"}, wantSuuid: "2", wantID: 4},
		{name: "no result", req: &entity.Device{Oaid: "foo"}, wantSuuid: "", wantID: 0, wantError: ErrRecordNotFound},
		{name: "only same at oaid", req: &entity.Device{Suuid: "5", Oaid: "baz"}, wantSuuid: "5", wantID: 4},
	}

	for _, cc := range cases {
		c := cc
		t.Run(c.name, func(t *testing.T) {
			repo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
			u, err := repo.GetRecentUserByDevice(context.Background(), "foo", c.req)
			if err != c.wantError {
				t.Fatalf("want error: %v, unexcepted error: %v", c.wantError, err)
			}
			if !IsErrRecordNotFound(err) {
				if assert.NotNil(t, u) {
					assert.Equal(t, c.wantSuuid, u.CommonSUUID)
					assert.Equal(t, c.wantID, u.ID)
				}
			}
		})
	}
}

func TestUserRepo_GetRecentIDBySUUID(t *testing.T) {
	setUp(t)
	defer tearDown()

	ctx := context.Background()
	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
	for i := 1; i < 5; i++ {
		user := entity.User{Model: gorm.Model{ID: uint(i)}, PackageName: "foo", CommonSUUID: strconv.Itoa(i)}
		_ = userRepo.Save(ctx, &user)
	}
	id, _, err := userRepo.GetIDAndCreatedAtBySUUID(ctx, "foo", "2")
	assert.NoError(t, err)
	assert.Equal(t, uint(2), id)
}

func TestUserRepo_Count(t *testing.T) {
	setUp(t)
	defer tearDown()

	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
	ctx := context.Background()
	for i := 1; i < 5; i++ {
		user := entity.User{Model: gorm.Model{ID: uint(i)}, UserName: "hello"}
		_ = userRepo.Save(ctx, &user)
	}
	count, err := userRepo.Count(ctx, clause.Where{Exprs: []clause.Expression{clause.IN{
		Column: "id",
		Values: []interface{}{1, 2, 3, 4},
	}}})
	if err != nil {
		t.Fatal(err)
	}
	if count != 4 {
		t.Fatalf("there should be four users, got %d", count)
	}
}

func TestUserRepo_CountAll(t *testing.T) {
	if !useMysql {
		t.Skip("count all tests must be run on mysql")
	}

	setUp(t)
	defer tearDown()

	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
	ctx := context.Background()
	for i := 1; i < 500; i++ {
		user := entity.User{Model: gorm.Model{ID: uint(i)}, UserName: "hello"}
		_ = userRepo.Save(ctx, &user)
	}
	count, err := userRepo.CountAll(ctx)
	if err != nil {
		t.Fatal(err)
	}
	if count != 499 {
		t.Fatalf("there should be 500 users, got %d", count)
	}
}

func TestUniqueConstraint(t *testing.T) {
	if !useMysql {
		t.Skip("unique constraints tests must be run on mysql")
	}
	setUp(t)
	defer tearDown()
	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
	ctx := context.Background()
	user := entity.User{
		PackageName: "abc",
		CommonSUUID: "a",
		Mobile:      sql.NullString{String: "110", Valid: true},
	}
	err := userRepo.Save(ctx, &user)
	if err != nil {
		t.Fatal(err)
	}
	user2 := entity.User{
		PackageName: "abc",
		CommonSUUID: "b",
		Mobile:      sql.NullString{String: "110", Valid: true},
	}
	err = userRepo.Save(ctx, &user2)
	if err == nil {
		t.Fatal("duplicate entry should be reported")
	}

	user3 := entity.User{
		PackageName:  "abc",
		CommonSUUID:  "c",
		WechatOpenId: sql.NullString{String: "110", Valid: true},
	}
	err = userRepo.Save(ctx, &user3)
	if err != nil {
		t.Fatal(err)
	}
	user4 := entity.User{
		PackageName:  "abc",
		CommonSUUID:  "d",
		WechatOpenId: sql.NullString{String: "110", Valid: true},
	}
	err = userRepo.Save(ctx, &user4)
	if err == nil {
		t.Fatal("save should fail")
	}
}

func TestUpdate(t *testing.T) {
	if !useMysql {
		t.Skip("unique constraints tests must be run on mysql")
	}
	setUp(t)
	defer tearDown()
	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
	ctx := context.Background()
	user := entity.User{
		PackageName:  "abc",
		CommonSUUID:  "a",
		Mobile:       sql.NullString{String: "110", Valid: true},
		WechatOpenId: sql.NullString{String: "110", Valid: true},
		TaobaoOpenId: sql.NullString{String: "110", Valid: true},
	}
	userRepo.Save(ctx, &user)

	user2 := entity.User{
		PackageName:  "abc",
		CommonSUUID:  "b",
		Mobile:       sql.NullString{Valid: false},
		WechatOpenId: sql.NullString{Valid: false},
		TaobaoOpenId: sql.NullString{Valid: false},
	}
	userRepo.Save(ctx, &user2)

	userRepo.Update(ctx, user2.ID, entity.User{
		Mobile: sql.NullString{String: "110", Valid: true},
	})

	{
		var target entity.User
		db.First(&target, user2.ID)

		if !target.Mobile.Valid {
			t.Fatalf("expect valid, got not valid")
		}

		var oldUser entity.User
		db.First(&oldUser, user.ID)

		if oldUser.Mobile.Valid {
			t.Fatalf("expect valid, got not valid")
		}
	}

	_, err := userRepo.Update(ctx, user2.ID, entity.User{
		WechatOpenId: sql.NullString{String: "110", Valid: true},
	})
	if err != nil {
		t.Fatal(err)
	}

	{
		var target entity.User
		db.First(&target, user2.ID)

		if !target.WechatOpenId.Valid {
			t.Fatalf("expect valid, got not valid")
		}

		var oldUser entity.User
		db.First(&oldUser, user.ID)

		if oldUser.WechatOpenId.Valid {
			t.Fatalf("expect not valid, got valid")
		}
	}

	userRepo.Update(ctx, user2.ID, entity.User{
		TaobaoOpenId: sql.NullString{String: "110", Valid: true},
	})

	{
		var target entity.User
		db.First(&target, user2.ID)

		if !target.TaobaoOpenId.Valid {
			t.Fatalf("expect valid, got not valid")
		}

		var oldUser entity.User
		db.First(&oldUser, user.ID)

		if oldUser.TaobaoOpenId.Valid {
			t.Fatalf("expect not valid, got valid")
		}
	}
}

func TestUser_UpdateCallback(t *testing.T) {
	if !useMysql {
		t.Skip("update callback tests must be run on mysql")
	}
	setUp(t)
	defer tearDown()
	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
	err := userRepo.UpdateCallback(context.Background(), 1, func(user *entity.User) error {
		return nil
	})
	assert.Error(t, err)
}

func TestUserRepo_GetByWechat(t *testing.T) {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	randSeq := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}
	if !useMysql {
		t.Skip("count all tests must be run on mysql")
	}

	setUp(t)
	defer tearDown()

	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
	ctx := context.Background()
	openid := randSeq(28)
	// ADD SOME USERS WITH DEVICES
	for i := 1; i < 5; i++ {

		user := entity.User{Model: gorm.Model{ID: uint(i)}, UserName: "hello", PackageName: "com.foo.bar", WechatOpenId: sql.NullString{String: openid}}
		user.Devices = []entity.Device{{
			Imei: "1234567",
			Hash: cast.ToString(i),
		}}
		_ = userRepo.Save(ctx, &user)
	}

	users, err := userRepo.GetByWechat(ctx, "com.foo.bar", openid)
	if err != nil {
		t.Fatalf("there should be no err, got %s", err)
	}
	if len(users.Devices) != 1 {
		t.Fatalf("there should be one devices, got %d", len(users.Devices))
	}
}

func TestUser_Exists(t *testing.T) {
	setUp(t)
	defer tearDown()
	userRepo := NewUserRepo(db, NewFileRepo(nil, nil), &mockID{})
	ok := userRepo.Exists(context.Background(), 1)
	assert.False(t, ok)
	user := user(1)
	userRepo.Save(context.Background(), &user)
	ok = userRepo.Exists(context.Background(), 1)
	assert.True(t, ok)
}
