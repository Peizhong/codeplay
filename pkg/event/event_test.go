package event

import (
	"context"
	"testing"

	"github.com/peizhong/codeplay/gen/mock_evaluator"
	"github.com/peizhong/codeplay/pkg/logger"
	"github.com/peizhong/codeplay/pkg/util"
	"github.com/segmentio/kafka-go"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestMain(m *testing.M) {
	logger.InitLogger()
	defer logger.Flush()

	if err := InitKafka("10.10.10.1:9092"); err != nil {
		panic(err)
	}
	defer Close()

	m.Run()
}

func TestPing(t *testing.T) {
	err := EnsureTopicExist("alert_event", 6)
	assert.NoError(t, err)
}

func TestCli(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	m := mock_evaluator.NewMockEvaluatorClient(ctrl)
	m.EXPECT().SayHello(gomock.Any(), gomock.Any()).Return(nil, nil)
}

func TestLogger(t *testing.T) {
	logger.Sugar().Infof("aaa")
	logger.GetWarnLogger().Printf("aaa")
	logger.GetWarnLogger().Println("aaa", "bb")
}

func TestSendMessage(t *testing.T) {
	writer := InitWriter([]string{"10.10.10.1:9092", "10.10.10.1:9192", "10.10.10.1:9292"}, "alert_event")
	for i := 0; i < 100; i++ {
		_ = writer.WriteMessages(context.Background(), kafka.Message{
			// Topic: topic,
			Key:   []byte(util.UUID()),
			Value: []byte("{}"),
		})
	}
	writer.Close()
}
