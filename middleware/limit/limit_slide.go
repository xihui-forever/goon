package limit

import (
	"time"

	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
)

func SlideWindow(cfg Config) bool {
	err := cfg.Storage.ZAdd(cfg.Key, time.Now().Unix())
	if err != nil {
		log.Errorf("err:%v", err)
		return true
	}

	var items []string
	for {
		items, err = cfg.Storage.ZRange(cfg.Key, int64(0), int64(0+1))
		if err != nil {
			log.Errorf("err:%v", err)
			return true
		}

		d := time.Now().Unix() - utils.ToInt64(items[0])
		if d > int64(cfg.Expiration) {
			cfg.Storage.ZRem(cfg.Key, items[0])
		} else {
			break
		}
	}

	var val int64
	val, err = cfg.Storage.ZLen(cfg.Key)
	if err != nil {
		log.Errorf("err:%v", err)
		return true
	}

	return val > cfg.Max
}
