package model

func AddWechatUser(user *WechatUser) error {
	_ = engine.Ping()
	_, err := engine.InsertOne(user)
	if err != nil {
		return err
	}
	return err
}

func WechatUserCount(openID string) int {
	_ = engine.Ping()
	count, err := engine.Where("open_id=?", openID).Count(new(WechatUser))
	if err != nil {
		return 0
	}
	return int(count)
}

func UpdateWechatUser(user *WechatUser) error {
	_ = engine.Ping()
	if WechatUserCount(user.OpenId) < 1 {
		err := AddWechatUser(user)
		if err != nil {
			return err
		}
	} else {
		_, err := engine.Where("open_id=?", user.OpenId).Update(user)
		if err != nil {
			return err
		}
	}
	return nil
}

func FindWechatUser(openID string) (*WechatUser, error) {
	_ = engine.Ping()
	w := new(WechatUser)
	_, err := engine.Where("open_id=?", openID).Get(w)
	if err != nil {
		return w, err
	}
	return w, err
}

func QueryWechatUser() ([]*WechatUser, error) {
	var (
		users []*WechatUser
	)
	_, err := engine.Table(new(WechatUser)).Query(&users)
	if err != nil {
		return nil, err
	}
	return users, err
}
