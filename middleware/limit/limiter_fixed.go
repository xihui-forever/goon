package limit

import "github.com/darabuchi/log"

// 固定窗口限频

func FixedWindow(cfg Config) bool {
	val, err := cfg.Storage.Inc(cfg.Key)
	if err != nil {
		log.Errorf("err:%v", err)
		return true
	}

	if val == 1 {
		err = cfg.Storage.Expire(cfg.Key, cfg.Expiration)
		if err != nil {
			log.Errorf("err:%v", err)
			return true
		}
	}

	return val > cfg.Max
}
