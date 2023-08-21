package lib

const (
	userInfoUrl            = "https://pc-api.xuexi.cn/open/api/user/info"
	userTotalscoreUrl      = "https://pc-proxy-api.xuexi.cn/delegate/score/get"
	userTodaytotalscoreUrl = "https://pc-proxy-api.xuexi.cn/delegate/score/today/query"
	userRatescoreUrl       = "https://pc-proxy-api.xuexi.cn/delegate/score/days/listScoreProgress?sence=score&deviceType=2"

	// pageSize=1000&pageNo=1
	querySpecialList = "https://pc-proxy-api.xuexi.cn/api/exam/service/paper/pc/list"
	queryWeekList    = "https://pc-proxy-api.xuexi.cn/api/exam/service/practice/pc/weekly/more"
)

// 专项答题JSON结构
type SpecialList struct {
	PageNo         int `json:"pageNo"`
	PageSize       int `json:"pageSize"`
	TotalPageCount int `json:"totalPageCount"`
	TotalCount     int `json:"totalCount"`
	List           []struct {
		TipScore    float64 `json:"tipScore"`
		EndDate     string  `json:"endDate"`
		Achievement struct {
			Score   int `json:"score"`
			Total   int `json:"total"`
			Correct int `json:"correct"`
		} `json:"achievement"`
		Year             int    `json:"year"`
		SeeSolution      bool   `json:"seeSolution"`
		Score            int    `json:"score"`
		ExamScoreId      int    `json:"examScoreId"`
		UsedTime         int    `json:"usedTime"`
		Overdue          bool   `json:"overdue"`
		Month            int    `json:"month"`
		Name             string `json:"name"`
		QuestionNum      int    `json:"questionNum"`
		AlreadyAnswerNum int    `json:"alreadyAnswerNum"`
		StartTime        string `json:"startTime"`
		Id               int    `json:"id"`
		ExamTime         int    `json:"examTime"`
		Forever          int    `json:"forever"`
		StartDate        string `json:"startDate"`
		Status           int    `json:"status"`
	} `json:"list"`
	PageNum int `json:"pageNum"`
}

type WeekList struct {
	PageNo         int `json:"pageNo"`
	PageSize       int `json:"pageSize"`
	TotalPageCount int `json:"totalPageCount"`
	TotalCount     int `json:"totalCount"`
	List           []struct {
		Month     string `json:"month"`
		Practices []struct {
			SeeSolution bool    `json:"seeSolution"`
			TipScore    float64 `json:"tipScore"`
			ExamScoreId int     `json:"examScoreId"`
			Overdue     bool    `json:"overdue"`
			Achievement struct {
				Total   int `json:"total"`
				Correct int `json:"correct"`
			} `json:"achievement"`
			Name               string `json:"name"`
			BeginYear          int    `json:"beginYear"`
			StartTime          string `json:"startTime"`
			Id                 int    `json:"id"`
			BeginMonth         int    `json:"beginMonth"`
			Status             int    `json:"status"`
			TipScoreReasonType int    `json:"tipScoreReasonType"`
		} `json:"practices"`
	} `json:"list"`
	PageNum int `json:"pageNum"`
}

/*

{"month":"5月","practices":[{"seeSolution":true,"tipScore":4,"examScoreId":-1,"overdue":false,"achievement":{"total":5,"correct":4},"name":"2022年5月第四周答题","beginYear":2022,"startTime":"2022-05-23 10:35:00","id":276,"beginMonth":5,"status":2,"tipScoreReasonType":0},{"seeSolution":true,"tipScore":4,"examScoreId":-1,"overdue":false,"achievement":{"total":5,"correct":4},"name":"2022年5月第三周答题","beginYear":2022,"startTime":"2022-05-16 09:00:00","id":275,"beginMonth":5,"status":2,"tipScoreReasonType":0},{"seeSolution":true,"tipScore":4,"examScoreId":-1,"overdue":false,"achievement":{"total":5,"correct":4},"name":"2022年5月第二周答题","beginYear":2022,"startTime":"2022-05-09 08:30:00","id":274,"beginMonth":5,"status":2,"tipScoreReasonType":0},{"seeSolution":true,"tipScore":5,"examScoreId":-1,"overdue":false,"achievement":{"total":5,"correct":5},"name":"2022年5月第一周答题","beginYear":2022,"startTime":"2022-05-02 08:30:00","id":273,"beginMonth":5,"status":2,"tipScoreReasonType":0}]}

*/
