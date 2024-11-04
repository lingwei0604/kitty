package module

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
	//mw "gitee.com/tagtic/go-middleware/http/middleware"
)

type logServer struct {
	t *testing.T
}

func (l logServer) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	buf, _ := ioutil.ReadAll(request.Body)
	l.t.Log(string(buf))
}

func TestSecurity(t *testing.T) {
	t.Skip("debug only")

	req := httptest.NewRequest("POST", "/", bytes.NewBuffer([]byte(`A5tOuYu6SltPvia/Xc+7QPw2D3yFDu8sWaF7clBr9ltwLxoT2tklwnkxMdW1zHXMLOZGKyWRbJfLQZ0d5wRHr0QO14OnSLfFnKEbfNi1gYTpftBtzFkGRaax6EQh6nDge3U+xE1Tmyn59jXFQijbdI5gy04ERgkk5fYjZs+am1pxDXDgeMUHLOs6k5kor+8vcEksAYI5tYtorYYanI+vbeeSRvA2BUp+rcu7V8MDiDzT6cW3h3xaZ2uCeAkIsx1dmaSy5olU4D5MYLZ7nmLf6s+khUkoJWrj/PZQYGKKyn6gqwA01JQAZAxtg0i80ox3mPiLw7hpeP4r3cw32gB3svjFUsxb8eN2rY6l6445LRNoUeLHCyKt+vcEEdeB2weZP5bn+OcfhF1DaTU6TsQtXz3W1SJM5l9mOrs8dtVBmH591IbhUuN1QyWElbRFV1CgSmXIDGT8vFiQrgheYCtvDg==`)))
	req.Header.Set("X-AuthToken", "101:HXP7eThn4x4XnOQUSCUJMv4RqjsMznDYbJuh3j2LRHpKjPHAf2Yj57MBkMLkpgLTbPj8PafXOtr2RG2miGeX3+0rOf6zG91+rvSGiGeAFfBszggei8/6k+qiMnZmrN0iRuwa7KRg3Ow3S8nyvOf4uEbffM3TEbFKs99Wh67Cv1o=")

	handler := logServer{t}
	rc := httptest.NewRecorder()
	handler.ServeHTTP(rc, req)

	data, _ := ioutil.ReadAll(rc.Result().Body)
	t.Log(string(data))

}
