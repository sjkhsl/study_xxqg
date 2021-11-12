package lib

import (
	"fmt"
	"testing"
)

func TestName(t *testing.T) {
	core := Core{}
	core.Init()
	cookies, err := core.Login()
	if err != nil {
		return
	}
	score, err := GetUserScore(cookies)
	if err != nil {
		return
	}
	fmt.Println(score)
}
