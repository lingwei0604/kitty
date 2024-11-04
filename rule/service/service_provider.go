package service

import (
	"github.com/go-kit/kit/log"
	"github.com/go-redis/redis/v8"
	"github.com/lingwei0604/kitty/pkg/contract"
	pb "github.com/lingwei0604/kitty/proto"
)

func ProvideService(logger log.Logger, repo Repository, dmpServer pb.DmpServer, redisClient redis.UniversalClient, hookClient contract.HttpDoer) Service {
	return &service{logger: logger, repo: repo, dmpServer: dmpServer, redisClient: redisClient, hookClient: hookClient}
}
