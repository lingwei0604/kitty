package service

import (
	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	pb "git.yingzhongshare.com/mkt/kitty/proto"
	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"
)

func ProvideService(logger log.Logger, repo Repository, dmpServer pb.DmpServer, redisClient redis.UniversalClient, hookClient contract.HttpDoer) Service {
	return &service{logger: logger, repo: repo, dmpServer: dmpServer, redisClient: redisClient, hookClient: hookClient}
}
