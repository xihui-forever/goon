package admin

import (
	"github.com/darabuchi/utils/db"
	"github.com/spf13/viper"
	"github.com/xihui-forever/goon/config"
	"github.com/xihui-forever/goon/types"
	"testing"
)

func TestAddAdmin(t *testing.T) {
	config.Load()

	err := db.Connect(db.Config{
		Dsn:      viper.GetString(config.DbDsn),
		Database: db.MySql,
	},
		&types.ModelAdmin{}, // db.AutoMigrate(&types.ModelAdmin{})
	)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	a, err := AddAdmin("admin", "123456")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	t.Log(a)
}

func TestChangedPwd(t *testing.T) {
	err := db.Connect(db.Config{
		Dsn:      viper.GetString(config.DbDsn),
		Database: db.MySql,
	},
		&types.ModelAdmin{},
	)
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}

	err = ChangePassword("admin", "123456", "654321")
	if err != nil {
		t.Errorf("err:%v", err)
		return
	}
}
