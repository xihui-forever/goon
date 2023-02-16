package limit

import (
	"github.com/darabuchi/log"
)

func SlideWindow(cfg Config) bool {
	val, err := cfg.Storage.Inc(cfg.Key)
	if err != nil {
		log.Errorf("err:%v", err)
		return true
	}

	val, err = cfg.Storage.DecBy(cfg.Key, int64(cfg.Expiration))
	if err != nil {
		log.Errorf("err:%v", err)
		return true
	}

	return val > cfg.Max
}
