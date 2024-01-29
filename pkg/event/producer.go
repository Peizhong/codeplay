package event

import (
	"time"

	"github.com/peizhong/codeplay/pkg/logger"
	"github.com/segmentio/kafka-go"
)

type K struct {
	conn *kafka.Conn
}

var k *K

func InitKafka(addr string) error {
	conn, err := kafka.Dial("tcp", addr)
	if err != nil {
		return err
	}
	k = &K{
		conn: conn,
	}
	return nil
}

func Close() {
	if k != nil {
		k.conn.Close()
	}
}

func EnsureTopicExist(topic string, numPartitions int) error {
	partitions, err := k.conn.ReadPartitions(topic)
	if err != nil {
		return err
	}
	if len(partitions) != numPartitions {
		logger.Sugar().Warnw("not enough partitions, recreate topic", "expect", numPartitions, "actual", len(partitions))
		if err = k.conn.DeleteTopics(topic); err != nil {
			return err
		}
		topicConfigs := []kafka.TopicConfig{
			{
				Topic:             topic,
				NumPartitions:     numPartitions,
				ReplicationFactor: 1,
			},
		}
		if err = k.conn.CreateTopics(topicConfigs...); err != nil {
			return err
		}
		<-time.After(time.Millisecond * 500)
		partitions, err = k.conn.ReadPartitions(topic)
		if err != nil {
			return err
		}
	} else {
		logger.Sugar().Infow("topic already exist", "topic", topic)
	}
	for _, partition := range partitions {
		logger.Sugar().Infoln("partition", partition.ID, partition.Replicas[0].Host)
	}
	return nil
}

func InitWriter(broker []string, topic string) *kafka.Writer {
	w := kafka.NewWriter(kafka.WriterConfig{
		Topic:        topic,
		Brokers:      broker,
		RequiredAcks: int(kafka.RequireAll),
		Balancer:     &kafka.LeastBytes{},
		Async:        true,
		BatchSize:    2,
	})
	return w
}
