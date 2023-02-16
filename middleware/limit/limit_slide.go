package limit

import (
	"github.com/darabuchi/log"
	"time"
)

func SlideWindow(cfg Config) bool {
	err := cfg.Storage.AddItem(cfg.Key, time.Now().Unix())
	if err != nil {
		log.Errorf("err:%v", err)
		return true
	}

	var item int64
	item, err = cfg.Storage.GetNotValid(cfg.Key, cfg.Expiration)
	if err != nil {
		log.Errorf("err:%v", err)
		return true
	}

	cfg.Storage.DeleteItem(cfg.Key, item)

	var val int64
	val, err = cfg.Storage.LenItemList(cfg.Key)
	if err != nil {
		log.Errorf("err:%v", err)
		return true
	}

	return val > cfg.Max
}
