package naming

import (
	"context"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/lafikl/consistent"
	"github.com/oklog/run"
	"github.com/peizhong/codeplay/pkg/cache"
	"github.com/peizhong/codeplay/pkg/logger"
	"github.com/peizhong/codeplay/pkg/util"
	"github.com/redis/go-redis/v9"
)

var (
	currentInstance  string
	lastInstanceList []string
	ringMu           sync.RWMutex
	ring             *consistent.Consistent
)

func RegisterService(g *run.Group, instance string) {
	currentInstance = instance
	heartbeat(context.Background(), time.Now().Unix())
	ctx, cancel := context.WithCancel(context.Background())
	g.Add(func() error {
		tick := time.NewTicker(time.Second * 10)
		for {
			select {
			case <-ctx.Done():
				return nil
			default:
			}
			select {
			case tm := <-tick.C:
				heartbeat(context.Background(), tm.Unix())
			case <-ctx.Done():
				return nil
			}
		}
	}, func(err error) {
		cancel()
	})
}

func heartbeat(ctx context.Context, timestamp int64) {
	// 记录当前实例心跳
	if err := cache.ShareCache.ZAdd(ctx, cache.CODEPLAY_INSTANCE, redis.Z{
		Score:  float64(timestamp),
		Member: currentInstance,
	}).Err(); err != nil {
		return
	}
	// 删除过期数据
	if err := cache.ShareCache.ZRemRangeByScore(ctx, cache.CODEPLAY_INSTANCE, "0", strconv.FormatInt(timestamp-30, 10)).Err(); err != nil {
		return
	}
	// 拉取最新数据
	list, err := cache.ShareCache.ZRange(ctx, cache.CODEPLAY_INSTANCE, 0, -1).Result()
	if err != nil {
		return
	}
	sort.Strings(list)
	if !util.SliceEqual(lastInstanceList, list) {
		logger.Sugar().Infow("rebuild ring", "last", lastInstanceList, "current", list)
		lastInstanceList = lastInstanceList[:0]
		lastInstanceList = append(lastInstanceList, list...)
		ringMu.Lock()
		defer ringMu.Unlock()
		ring = consistent.New()
		for _, instance := range list {
			ring.Add(instance)
		}
	}
}

func IsCmdbNodeHit(nodeIp string) bool {
	if nodeIp == "" {
		return false
	}
	ringMu.RLock()
	defer ringMu.RUnlock()
	host, _ := ring.Get(nodeIp)
	return host == currentInstance
}

func Clusters() []string {
	return append(make([]string, 0, len(lastInstanceList)), lastInstanceList...)
}
