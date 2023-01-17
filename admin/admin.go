package admin

import (
	"encoding/base64"
	"github.com/darabuchi/log"
	"github.com/darabuchi/utils"
	"github.com/darabuchi/utils/db"
	"github.com/xihui-forever/goon/types"
	"gorm.io/gorm"
)

func AddAdmin(username string, pwd string) (*types.ModelAdmin, error) {
	a := types.ModelAdmin{
		Username: username,
		Password: Encrypt(pwd),
	}

	err := db.Create(&a).Error
	if err != nil {
		if types.IsUniqueErr(err) {
			return nil, ErrAdminExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}

	return &a, nil
}

func GetAdmin(username string) (*types.ModelAdmin, error) {
	var a types.ModelAdmin
	err := db.Where("username = ?", username).First(&a).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrAdminNotExist
		}
		log.Errorf("err:%v", err)
		return nil, err
	}
	return &a, nil
}

func ChangePassword(username string, oldPwd, newPwd string) error {
	a, err := GetAdmin(username)
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if a.Password != Encrypt(oldPwd) {
		return ErrPasswordWrong
	}

	res := db.Model(&types.ModelAdmin{}).Where("id = ?", a.Id).Update("password", Encrypt(newPwd))
	err = res.Error
	if err != nil {
		log.Errorf("err:%v", err)
		return err
	}

	if res.RowsAffected == 0 {
		return ErrPasswordChangeFailed
	}

	return nil
}

func Encrypt(pwd string) string {
	return base64.StdEncoding.EncodeToString([]byte(utils.HmacSha384("620dd0f8d3e5424f99548ed8b9a6a59f", pwd)))
}
