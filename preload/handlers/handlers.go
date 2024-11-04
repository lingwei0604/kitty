package handlers

import (
	"context"
	"fmt"
	pb "git.yingzhongshare.com/mkt/kitty/proto"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"io/ioutil"
	"net/http"
	"regexp"
	"strings"
	"sync"
)

// NewService returns a naïve, stateless implementation of Service.
func NewService(log log.Logger) pb.PreloadServer {
	return preloadService{log}
}

type preloadService struct {
	log log.Logger
}

// ListInfo 将所有的 preloadHost 参数进行处理获得 gzUrl和md5值
// 如果md5值为空或者错误时，gzUrl和md5值返回为空
func (p preloadService) ListInfo(ctx context.Context, in *pb.PreloadReq) (*pb.PreloadResp, error) {
	var resp pb.PreloadResp
	fmt.Println(in.PreloadHostList)

	err := in.Validate()
	if err != nil {
		resp.Code = 400
		resp.Msg = err.Error()
		return &resp, nil
	}

	preloadInfoList := p.getPreloadInfoList(in.PreloadHostList)

	resp = pb.PreloadResp{
		Code: 0,
		Msg:  "success",
		Data: preloadInfoList,
	}

	return &resp, nil
}

func (p preloadService) getPreloadInfoList(preloadHostList []string) []*pb.PreloadInfo {

	wg := sync.WaitGroup{}
	wg.Add(len(preloadHostList))

	var preloadInfoList = make([]*pb.PreloadInfo, len(preloadHostList))

	for i := 0; i < len(preloadHostList); i++ {
		go func(preloadHost string, i int) {
			defer wg.Done()

			gzUrl, md5Url := getCompressAndMd5Url(preloadHost)

			md5 := getMd5ByHttpUrl(md5Url, p.log)

			if md5 == "" {
				gzUrl = ""
			}

			preloadInfo := pb.PreloadInfo{
				Gzurl:  gzUrl,
				Md5:    strings.TrimSpace(md5),
				Weburl: preloadHost,
			}

			preloadInfoList[i] = &preloadInfo
		}(preloadHostList[i], i)
	}

	wg.Wait()

	return preloadInfoList

}

func getCompressAndMd5Url(preloadHost string) (gzUrl string, md5Url string) {
	// 去除 前端路由
	index := strings.Index(preloadHost, "#/")
	if index != -1 {
		preloadHost = preloadHost[0:index]
	}

	// 是否 [index/结尾]  和 [斜杠结尾] 和 [index.html/结尾]
	var endIndexHtml = strings.HasSuffix(preloadHost, "index.html")
	var endBackslash = strings.HasSuffix(preloadHost, "/")
	var endIndexHtmlAndBackslash = strings.HasSuffix(preloadHost, "index.html/")

	// [index.html/]  结尾
	if endIndexHtmlAndBackslash {
		md5Url = strings.Replace(preloadHost, "index.html/", "md5.txt", 1)
		gzUrl = strings.Replace(preloadHost, "index.html/", "preload.tar.gz", 1)
		return gzUrl, md5Url
	}

	// [index.html结尾]
	if endIndexHtml && !endBackslash {
		md5Url = strings.Replace(preloadHost, "index.html", "md5.txt", 1)
		gzUrl = strings.Replace(preloadHost, "index.html", "preload.tar.gz", 1)
		return gzUrl, md5Url
	}

	// [斜杠结尾]
	if !endIndexHtml && endBackslash {
		md5Url = preloadHost + "md5.txt"
		gzUrl = preloadHost + "preload.tar.gz"
		return gzUrl, md5Url
	}

	// 两者结尾都不是
	if !endBackslash && !endIndexHtml {
		md5Url = preloadHost + "/md5.txt"
		gzUrl = preloadHost + "/preload.tar.gz"
		return gzUrl, md5Url
	}

	return gzUrl, md5Url
}

func getMd5ByHttpUrl(md5Url string, log log.Logger) string {
	response, err := http.Get(md5Url)
	if err != nil {
		level.Error(log).Log(err)
		return ""
	}
	if response.StatusCode != 200 && response.StatusCode != 201 {
		level.Error(log).Log("请求md5错误", md5Url)
		return ""
	}
	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		level.Error(log).Log(err)
		return ""
	}

	content := string(responseData)
	content = strings.Trim(content, "\n\r")
	content = strings.Trim(content, "")

	matched, err := regexp.MatchString(`^[0-9a-fA-F]{32}$`, content)
	if err != nil || matched == false {
		level.Error(log).Log("md5url值返回md5值不正确", content, "md5Url", md5Url)
		return ""
	}

	return content
}
