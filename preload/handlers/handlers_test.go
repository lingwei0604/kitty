package handlers

import (
	"github.com/go-kit/kit/log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func Test_getCompressAndMd5Url(t *testing.T) {
	type args struct {
		preloadHost string
	}
	tests := []struct {
		name       string
		args       args
		wantGzUrl  string
		wantMd5Url string
	}{
		{
			name:       "标准preload host地址(不带index.html)",
			args:       args{preloadHost: "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx"},
			wantGzUrl:  "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/preload.tar.gz",
			wantMd5Url: "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/md5.txt",
		},
		{
			name:       "标准preload host地址(不带index.html)(多带一个斜杠)",
			args:       args{preloadHost: "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/"},
			wantGzUrl:  "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/preload.tar.gz",
			wantMd5Url: "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/md5.txt",
		},
		{
			name:       "标准preload host地址(不带index.html)(带有前端路由)",
			args:       args{preloadHost: "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx#/view1"},
			wantGzUrl:  "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/preload.tar.gz",
			wantMd5Url: "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/md5.txt",
		},
		{
			name:       "标准preload host地址(不带index.html)(多带一个斜杠)(带有前端路由)",
			args:       args{preloadHost: "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx#/view1"},
			wantGzUrl:  "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/preload.tar.gz",
			wantMd5Url: "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/md5.txt",
		},
		{
			name:       "标准preload host地址(带index.html)",
			args:       args{preloadHost: "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/index.html"},
			wantGzUrl:  "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/preload.tar.gz",
			wantMd5Url: "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/md5.txt",
		},
		{
			name:       "标准preload host地址(带index.html)(多带一个斜杠)",
			args:       args{preloadHost: "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/index.html/"},
			wantGzUrl:  "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/preload.tar.gz",
			wantMd5Url: "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/md5.txt",
		},
		{
			name:       "标准preload host地址(带index.html)(带有前端路由)",
			args:       args{preloadHost: "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/index.html#/view1"},
			wantGzUrl:  "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/preload.tar.gz",
			wantMd5Url: "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/md5.txt",
		},
		{
			name:       "标准preload host地址(带index.html)(多带一个斜杠)(多带一个斜杠)",
			args:       args{preloadHost: "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/index.html/#/view1"},
			wantGzUrl:  "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/preload.tar.gz",
			wantMd5Url: "https://ad-static-xg.tagtic.cn/static/recharge-dev/cdbx/md5.txt",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGzUrl, gotMd5Url := getCompressAndMd5Url(tt.args.preloadHost)
			if gotGzUrl != tt.wantGzUrl {
				t.Errorf("getCompressAndMd5Url() gotGzUrl = %v, want %v", gotGzUrl, tt.wantGzUrl)
			}
			if gotMd5Url != tt.wantMd5Url {
				t.Errorf("getCompressAndMd5Url() gotMd5Url = %v, want %v", gotMd5Url, tt.wantMd5Url)
			}
		})
	}
}

func Test_getMd5ByHttpUrl(t *testing.T) {
	testServer := httptest.NewServer(http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		flags := request.URL.Query()["md5-right"]

		if len(flags) > 0 {
			writer.Write([]byte("05A8EC7AB03506116217268F89287FE5"))
		}
		if len(flags) <= 0 {
			writer.Write([]byte("body"))
		}
	}))
	defer func() { testServer.Close() }()

	type args struct {
		md5Url string
		log    log.Logger
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "md5url，无法抵达404",
			args: args{
				md5Url: "https://www.baidxxxxsdajfjdsaklfjadslkfjlkdsaju.com",
				log:    log.NewLogfmtLogger(os.Stdout),
			},
			want: "",
		},
		{
			name: "md5url,错误md5值",
			args: args{
				md5Url: testServer.URL,
				log:    log.NewLogfmtLogger(os.Stdout),
			},
			want: "",
		},
		{
			name: "md5url,正确md5值",
			args: args{
				md5Url: testServer.URL + "?md5-right=true",
				log:    log.NewLogfmtLogger(os.Stdout),
			},
			want: "05A8EC7AB03506116217268F89287FE5",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMd5ByHttpUrl(tt.args.md5Url, tt.args.log); got != tt.want {
				t.Errorf("getMd5ByHttpUrl() = %v, want %v", got, tt.want)
			}
		})
	}
}
