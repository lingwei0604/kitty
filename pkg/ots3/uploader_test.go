package ots3

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"testing"
)

var useS3 bool

func init() {
	flag.BoolVar(&useS3, "s3", false, "use s3 for testing")
}

func TestUploader(t *testing.T) {
	if !useS3 {
		t.Skip("must use s3 for this test")
	}
	url := "https://aecpm.alicdn.com/simba/img/TB1jFYch8FR4u4jSZFPSuunzFXa.jpg"
	orig, err := http.Get(url)
	if err != nil {
		t.Fatal(err)
	}
	defer orig.Body.Close()

	uploader := NewManager("BOOJ2EFZ25YXETGTWDMI", "rogGgzGVo3bwXGCmeVzjwOTFR5Nkz2fXcDyQrLGd", "http://minio.xg.tagtic.cn", "asia", "ad-material")
	nu, err := uploader.Upload(context.Background(), orig.Body)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(nu)
}
