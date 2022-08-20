package model

import log "github.com/sirupsen/logrus"

func SearchAnswer(title string) string {
	initQuestionDb()
	if db1 == nil {
		return ""
	}
	var answer string
	row := db1.QueryRow("select answer from tiku where question like ?", title+"%")
	err := row.Scan(&answer)
	if err != nil {
		log.Errorln(err.Error())
		return ""
	}
	if answer == "" {
		row := db1.QueryRow("select answer from tikuNet where question like ?", title+"%")
		err := row.Scan(&answer)
		if err != nil {
			log.Errorln(err.Error())
			return ""
		}
	}
	log.Infoln("从数据库查询到答案：" + answer)
	return answer
}
