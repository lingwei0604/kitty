//go:build !wireinject
// +build !wireinject

package service

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	pb "git.yingzhongshare.com/mkt/kitty/proto"
	"git.yingzhongshare.com/mkt/kitty/rule/dto"
	"git.yingzhongshare.com/mkt/kitty/rule/entity"
	"git.yingzhongshare.com/mkt/kitty/rule/service/mocks"
	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/go-redis/redismock/v8"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockClient struct {
	redis.UniversalClient
}

type mockDmpServer struct {
}

type mockHttp func(req *http.Request) (*http.Response, error)

func (m mockHttp) Do(req *http.Request) (*http.Response, error) {
	return m(req)
}

func (m mockDmpServer) UserMore(ctx context.Context, req *pb.DmpReq) (*pb.DmpResp, error) {
	return &pb.DmpResp{
		AdClick:    100,
		AdComplete: 0,
		AdDisplay:  0,
		AdCtrDev:   0,
		Register:   "2020-01-01 00:00:00",
		Score:      0,
		ScoreTotal: 0,
		BlackType:  0,
		Ext:        "",
	}, nil
}

func TestService_CalculateRules(t *testing.T) {
	cases := []struct {
		text    string
		payload dto.Payload
		result  dto.Data
	}{
		{
			`
style: advanced
rule:
  - if: true
    then: 
      foo: bar
`,
			dto.Payload{},
			dto.Data{"foo": "bar"},
		},
		{
			`
style: advanced
rule:
  - if: false
    then: 
      foo: bar
`,
			dto.Payload{},
			dto.Data{},
		},
		{
			`
style: advanced
rule:
  - if: Imei == "456"
    then: 
      foo: bar
  - if: Imei == "123"
    then:
      foo: baz
`,
			dto.Payload{
				Imei: "123",
			},
			dto.Data{
				"foo": "baz",
			},
		},
		{
			`
style: advanced
rule:
- if: Imei == "456" && Oaid == "789"
  then: 
    foo: bar
- if: Imei == "123" && Oaid == "789"
  then:
    foo: baz
- if: Imei == "123"
  then:
    foo: quz
`,
			dto.Payload{
				Imei: "123",
				Oaid: "789",
			},
			dto.Data{
				"foo": "baz",
			},
		},
		{
			`
style: advanced
enrich: true
rule:
- if: DMP.AdClick > 10
  then: 
    foo: bar
- if: true
  then:
    foo: quz
`,
			dto.Payload{
				Imei: "123",
				Oaid: "789",
			},
			dto.Data{
				"foo": "bar",
			},
		},
		{
			`
style: advanced
enrich: true
rule:
- if: HoursAgo(DMP.Register) < 1
  then: 
    foo: foo
- if: HoursAgo(DMP.Register) > 100
  then: 
    foo: bar
- if: true
  then:
    foo: quz
`,
			dto.Payload{
				Imei: "123",
				Oaid: "789",
			},
			dto.Data{
				"foo": "bar",
			},
		},
		{
			`
style: advanced
rule:
- if: SIsMember("imei", Oaid)
  then:
    foo: bar
- if: SIsMember("imei", Imei)
  then:
    foo: foo
- if: true
  then:
    foo: quz
`,
			dto.Payload{
				Imei: "123",
				Oaid: "345",
			},
			dto.Data{
				"foo": "foo",
			},
		},
	}
	for _, c := range cases {
		cc := c
		t.Run("", func(t *testing.T) {
			repo := &mocks.Repository{}
			db, rmock := redismock.NewClientMock()
			rmock.ExpectSIsMember("imei", "123").SetVal(true)
			ser := ProvideService(log.NewNopLogger(), repo, mockDmpServer{}, db, http.DefaultClient)
			repo.On("GetCompiled", mock.Anything).Return(entity.NewRules(bytes.NewReader([]byte(cc.text)), log.NewNopLogger()))
			result, err := ser.CalculateRules(context.Background(), "", &cc.payload)
			if err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(result, cc.result) {
				t.Fatalf("want %#v, got %#v", cc.result, result)
			}
		})
	}
}

func TestService_CalculateMultipleRules(t *testing.T) {
	repo := &mocks.Repository{}
	db, rmock := redismock.NewClientMock()
	rmock.ExpectSIsMember("imei", "123").SetVal(true)
	ser := ProvideService(log.NewNopLogger(), repo, mockDmpServer{}, db, http.DefaultClient)
	repo.On("GetCompiled", mock.Anything).Return(func(ruleName string) entity.Ruler {
		v := map[string]string{
			"a": `
style: advanced
rule:
  - if: true
    then: 
      foo: bar`,
			"b": `
style: advanced
enrich: true
rule:
- if: Imei == '123'
  then: 
    foo: bar
- if: true
  then:
    foo: quz`,
		}
		if s, ok := v[ruleName]; ok {
			return entity.NewRules(bytes.NewReader([]byte(s)), log.NewNopLogger())
		}
		return nil
	})
	result, err := ser.CalculateMultipleRules(context.Background(), &dto.Payload{
		RuleNames: []string{"a", "b", "c"},
		Imei:      "123",
		Oaid:      "345",
	})
	if err != nil {
		t.Fatal(err)
	}
	t.Log(result)
}

func TestService_UpdateRules(t *testing.T) {
	repo := &mocks.Repository{}
	ser := NewService(log.NewNopLogger(), repo, mockClient{}, http.DefaultClient)
	repo.On("SetRaw", mock.Anything, mock.Anything, mock.Anything).Return(nil)
	repo.On("ValidateRules", mock.Anything, mock.Anything).Return(func(ruleName string, reader io.Reader) error {
		return entity.ValidateRules(reader)
	})
	err := ser.UpdateRules(context.Background(), "foo", []byte("invalid"), false)
	if err == nil {
		t.Fatal("err should not be null")
	}
	data := []byte(`
style: advanced
rule:
- if: true
  then:
    data: ok
`)
	err = ser.UpdateRules(context.Background(), "foo", data, false)
	if err != nil {
		t.Fatal(err)
	}
}

func TestService_UpdateRulesWithHooks(t *testing.T) {

	prepareService := func(bool2 *bool) *service {
		repo := &mocks.Repository{}
		ser := NewService(log.NewNopLogger(), repo, mockClient{}, mockHttp(func(req *http.Request) (*http.Response, error) {
			body := `{"code": 1, "msg": "failed!"}`
			resp := httptest.ResponseRecorder{
				Code: 200,
				Body: bytes.NewBufferString(body),
			}
			*bool2 = true
			return resp.Result(), nil
		}))
		repo.On("SetRaw", mock.Anything, mock.Anything, mock.Anything).Return(nil)
		repo.On("ValidateRules", mock.Anything, mock.Anything).Return(func(ruleName string, reader io.Reader) error {
			return entity.ValidateRules(reader)
		})
		return ser
	}

	t.Run("with onChange hook", func(t *testing.T) {
		var onChangeCalled bool
		ser := prepareService(&onChangeCalled)

		data := []byte(`
style: advanced
hooks:
  onChange: http://baid.com
rule:
- if: true
  then:
    data: ok
`)
		err := ser.UpdateRules(context.Background(), "foo", data, true)
		if err == nil {
			t.Fatal("err should be triggered")
		}
		if !onChangeCalled {
			t.Fatal("onChange hook should be called")
		}
	})

	t.Run("when dry run, preUpdate hook should not be called", func(t *testing.T) {
		var preUpdateCalled bool
		ser := prepareService(&preUpdateCalled)

		data := []byte(`
style: advanced
hooks:
  preUpdate: http://preUpdate
rule:
- if: true
  then:
    data: ok
`)
		err := ser.UpdateRules(context.Background(), "foo", data, true)
		if err != nil {
			t.Fatal("no error should occur")
		}
		if preUpdateCalled {
			t.Fatal("preUpdateCalled hook should not be called")
		}
	})

	t.Run("when not dry run, preUpdate hook should be called", func(t *testing.T) {
		var preUpdateCalled bool
		ser := prepareService(&preUpdateCalled)

		data := []byte(`
style: advanced
hooks:
  preUpdate: http://preUpdate
rule:
- if: true
  then:
    data: ok
`)
		err := ser.UpdateRules(context.Background(), "foo", data, false)
		if err == nil {
			t.Fatal("error should occur")
		}
		if !preUpdateCalled {
			t.Fatal("preUpdateCalled hook should be called")
		}
	})

	t.Run("when not dry run, postUpdate hook should be called", func(t *testing.T) {
		var postUpdateCalled bool
		ser := prepareService(&postUpdateCalled)

		data := []byte(`
style: advanced
hooks:
  postUpdate: http://preUpdate
rule:
- if: true
  then:
    data: ok
`)
		err := ser.UpdateRules(context.Background(), "foo", data, false)
		if err == nil {
			t.Fatal("error should occur")
		}
		if !postUpdateCalled {
			t.Fatal("postUpdateCalled hook should be called")
		}
	})
}

func TestService_Preflight(t *testing.T) {
	repo := &mocks.Repository{}
	ser := NewService(log.NewNopLogger(), repo, mockClient{}, http.DefaultClient)

	{
		repo.On("IsNewest", mock.Anything, mock.Anything, mock.Anything).Return(true, nil).Once()
		err := ser.Preflight(context.Background(), "foo", "fooo")
		if err != nil {
			t.Fatal("err should be null")
		}
	}

	{
		repo.On("IsNewest", mock.Anything, mock.Anything, mock.Anything).Return(false, nil).Once()
		err := ser.Preflight(context.Background(), "foo", "fooo")
		if err == nil {
			t.Fatal("err should not be null")
		}
	}
}

func TestService_GetRules(t *testing.T) {
	repo := &mocks.Repository{}
	ser := NewService(log.NewNopLogger(), repo, mockClient{}, http.DefaultClient)
	{
		repo.On("GetRaw", mock.Anything, mock.Anything).Return([]byte("foo"), nil).Once()
		byt, err := ser.GetRules(context.Background(), "foo")
		assert.Nil(t, err)
		assert.Equal(t, byt, []byte("foo"))
	}
}
