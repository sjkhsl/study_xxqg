package model

import "database/sql"

type WechatUser struct {
	OpenID          string `json:"open_id"`
	Remark          string `json:"remark"`
	Status          int    `json:"status"`
	LastRequestTime int64  `json:"last_request_time"`
}

func AddWechatUser(user *WechatUser) error {
	ping()
	_, err := db.Exec(`insert into wechat_user (open_id, remark, status, last_request_time) VALUES (?,?,?,?)`, user.OpenID, user.Remark, user.Status, user.LastRequestTime)
	return err
}

func WechatUserCount(openID string) int {
	ping()
	var count int
	err := db.QueryRow("select count(*) from wechat_user where open_id = ?", openID).Scan(&count)
	if err != nil {
		return 0
	}
	return count
}

func UpdateWechatUser(user *WechatUser) error {
	ping()
	_, err := db.Exec(`update wechat_user set remark=?,status=?,last_request_time=? where open_id=?;`, user.Remark, user.Status, user.LastRequestTime, user.OpenID)
	return err
}

func FindWechatUser(openID string) (*WechatUser, error) {
	ping()
	w := new(WechatUser)
	err := db.QueryRow(`select * from wechat_user  where open_id=?;`, openID).Scan(&w.OpenID, &w.Remark, &w.Status, &w.LastRequestTime)
	if err == sql.ErrNoRows {
		return nil, err
	}
	return w, err
}

func QueryWechatByCondition(condition string) ([]*WechatUser, error) {
	var users []*WechatUser
	if condition == "" {
		condition = "1=1"
	}
	ping()
	results, err := db.Query("select * from wechat_user where ?", condition)
	if err != nil {
		return users, err
	}

	for results.Next() {
		w := new(WechatUser)
		err := results.Scan(&w.OpenID, &w.Remark, &w.Status, &w.LastRequestTime)
		if err != nil {
			_ = results.Close()
			return users, err
		}
		users = append(users, w)
	}
	return users, err
}

func DeleteWechatUser(openID string) error {
	ping()
	_, err := db.Exec(`delete from wechat_user where open_id=?;`, openID)
	return err
}
