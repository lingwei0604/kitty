package svc

import (
	"context"

	"git.yingzhongshare.com/mkt/kitty/pkg/contract"
	"git.yingzhongshare.com/mkt/kitty/pkg/kkafka"
	pb "git.yingzhongshare.com/mkt/kitty/proto"
	"github.com/go-kit/kit/endpoint"
	"github.com/segmentio/kafka-go"
)

func DecodeBindAdRequest(ctx context.Context, msg *kafka.Message) (interface{}, error) {
	var UserBindAd pb.UserBindAdRequest
	err := UserBindAd.Unmarshal(msg.Value)
	if err != nil {
		return nil, err
	}
	return &UserBindAd, nil
}

func provideBindAdSubscriber(endpoint endpoint.Endpoint, options ...kkafka.SubscriberOption) kkafka.Handler {
	return kkafka.NewSubscriber(
		endpoint,
		DecodeBindAdRequest,
		options...,
	)
}

func MakeKafkaServer(endpoints Endpoints, factory *kkafka.KafkaFactory, conf contract.ConfigReader, options ...kkafka.SubscriberOption) kkafka.Server {
	group := conf.String("kafka.groupId")

	bindAd := provideBindAdSubscriber(endpoints.BindAdEndpoint, options...)

	return kkafka.NewMux(
		factory.MakeKafkaServer(conf.String("kafka.bindAd"), bindAd, kkafka.WithGroup(group)),
	)
}
