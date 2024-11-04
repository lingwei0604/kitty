//go:generate mockery --name=Repository

package service

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/go-redis/redis/v8"
	"github.com/lingwei0604/kitty/pkg/contract"
	"github.com/lingwei0604/kitty/pkg/kerr"
	pb "github.com/lingwei0604/kitty/proto"
	"github.com/lingwei0604/kitty/rule/dto"
	"github.com/lingwei0604/kitty/rule/entity"
	"github.com/lingwei0604/kitty/rule/msg"
	"github.com/pkg/errors"
)

var ErrDataHasChanged = errors.New(msg.ErrorRulesHasChanged)

type Service interface {
	CalculateRules(ctx context.Context, ruleName string, payload *dto.Payload) (dto.Data, error)
	CalculateMultipleRules(ctx context.Context, payload *dto.Payload) (dto.Data, error)
	GetRules(ctx context.Context, ruleName string) ([]byte, error)
	UpdateRules(ctx context.Context, ruleName string, content []byte, dryRun bool) error
	Preflight(ctx context.Context, ruleName string, hash string) error
}

type Repository interface {
	GetCompiled(ruleName string) entity.Ruler
	GetRaw(ctx context.Context, key string) (value []byte, e error)
	SetRaw(ctx context.Context, key string, value string) error
	IsNewest(ctx context.Context, key, value string) (bool, error)
	WatchConfigUpdate(ctx context.Context) error
	ValidateRules(ruleName string, reader io.Reader) error
}

type service struct {
	hookClient  contract.HttpDoer
	redisClient redis.UniversalClient
	dmpServer   pb.DmpServer
	logger      log.Logger
	repo        Repository
}

func NewService(logger log.Logger, repo Repository, redisClient redis.UniversalClient, hookClient contract.HttpDoer) *service {
	return &service{logger: logger, repo: repo, redisClient: redisClient, hookClient: hookClient}
}

func (r *service) findCompiled(packageName, ruleName string) (entity.Ruler, error) {

	var compiled entity.Ruler
	if packageName != "" {
		compiled = r.repo.GetCompiled(packageName + "-" + ruleName)
		if compiled != nil {
			return compiled, nil
		}

		parts := strings.Split(packageName, ".")
		codeName := parts[len(parts)-1]
		compiled = r.repo.GetCompiled(codeName + "-" + ruleName)
		if compiled != nil {
			return compiled, nil
		}
	}

	compiled = r.repo.GetCompiled(ruleName)
	if compiled != nil {
		return compiled, nil
	}

	return nil, fmt.Errorf("no suitable configuration found for %s", ruleName)
}

func (r *service) CalculateRules(ctx context.Context, ruleName string, payload *dto.Payload) (dto.Data, error) {
	rules, err := r.findCompiled(payload.PackageName, ruleName)
	if err != nil {
		return nil, err
	}
	if rules.ShouldEnrich() {
		resp, err := r.dmpServer.UserMore(ctx, &pb.DmpReq{
			UserId:      payload.UserId,
			PackageName: payload.PackageName,
			Suuid:       payload.Suuid,
			Channel:     payload.Channel,
		})
		if err != nil {
			level.Warn(r.logger).Log("err", errors.Wrap(err, "dmp server error"))
		}
		if resp == nil {
			resp = &pb.DmpResp{}
		}
		payload.DMP = dto.Dmp{DmpResp: *resp}
	}
	if payload.Redis == nil {
		payload.Redis = r.redisClient
	}
	if payload.Context == nil {
		payload.Context = ctx
	}
	return entity.Calculate(rules, payload)
}

func (r *service) CalculateMultipleRules(ctx context.Context, payload *dto.Payload) (dto.Data, error) {
	d := map[string]interface{}{}
	shouldEnrich := false
	if payload.Redis == nil {
		payload.Redis = r.redisClient
	}
	if payload.Context == nil {
		payload.Context = ctx
	}
	for _, one := range payload.RuleNames {
		if _, exist := d[one]; exist {
			continue
		}
		d[one] = map[string]string{}
		ruler, err := r.findCompiled(payload.PackageName, one)
		if err != nil {
			continue
		}
		if !shouldEnrich && ruler.ShouldEnrich() {
			shouldEnrich = true
			resp, err := r.dmpServer.UserMore(ctx, &pb.DmpReq{
				UserId:      payload.UserId,
				PackageName: payload.PackageName,
				Suuid:       payload.Suuid,
				Channel:     payload.Channel,
			})
			if err != nil {
				level.Warn(r.logger).Log("err", errors.Wrap(err, "dmp server error"))
			}
			if resp == nil {
				resp = &pb.DmpResp{}
			}
			payload.DMP = dto.Dmp{DmpResp: *resp}
		}
		data, err := entity.Calculate(ruler, payload)
		if err != nil {
			continue
		}
		d[one] = &data
	}
	return d, nil
}

func (r *service) GetRules(ctx context.Context, ruleName string) ([]byte, error) {
	return r.repo.GetRaw(ctx, ruleName)
}

// TODO: trigger hooks
func (r *service) UpdateRules(ctx context.Context, ruleName string, content []byte, dryRun bool) error {
	var (
		buf         bytes.Buffer
		err         error
		tee         io.Reader
		hookManager *HookManager
	)
	reader := bytes.NewReader(content)
	tee = io.TeeReader(reader, &buf)
	err = r.repo.ValidateRules(ruleName, tee)
	var invalid *entity.ErrInvalidRules
	if errors.As(err, &invalid) {
		return kerr.InvalidArgumentErr(invalid, msg.ErrorRules)
	}
	if err != nil {
		return err
	}
	if hookManager, err = NewHookManager(r.hookClient, buf.Bytes()); err != nil {
		return err
	}
	if err := hookManager.OnChange(); err != nil {
		return errors.Wrap(err, "an error was returned from onChange hook")
	}
	if dryRun {
		return nil
	}
	if err := hookManager.PreUpdate(); err != nil {
		return errors.Wrap(err, "an error was returned from preUpdate hook")
	}
	if err := r.repo.SetRaw(ctx, ruleName, buf.String()); err != nil {
		return err
	}
	if err := hookManager.PostUpdate(); err != nil {
		return errors.Wrap(err, "save succeeded, but an error was returned from postUpdate hook")
	}
	return err
}

func (r *service) Preflight(ctx context.Context, ruleName string, hash string) error {
	ok, err := r.repo.IsNewest(ctx, ruleName, hash)
	if err != nil {
		return errors.Wrap(err, msg.ErrorETCD)
	}
	if !ok {
		return ErrDataHasChanged
	}
	return nil
}
