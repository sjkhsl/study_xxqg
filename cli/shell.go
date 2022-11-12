package cli

import (
	"fmt"

	"github.com/desertbit/grumble"
	log "github.com/sirupsen/logrus"

	"github.com/johlanse/study_xxqg/conf"
	"github.com/johlanse/study_xxqg/lib"
	"github.com/johlanse/study_xxqg/model"
	"github.com/johlanse/study_xxqg/push"
	"github.com/johlanse/study_xxqg/utils"
)

func RunCli() {
	app := grumble.New(&grumble.Config{
		Name:        "StudyXXqg",
		Description: "a StudyXXqh app",
		Flags: func(f *grumble.Flags) {
			f.Bool("n", "now", false, "now")
			f.Bool("u", "update", false, "update")
			f.Bool("i", "init", false, "init")
			f.String("c", "config", "", "config")
		},
	})
	app.AddCommand(getUser())
	app.AddCommand(addUser())
	app.AddCommand(study())
	app.AddCommand(score())
	grumble.Main(app)
}

func study() *grumble.Command {
	return &grumble.Command{
		Name:     "study",
		Aliases:  []string{"study"},
		Help:     "study xxqg",
		LongHelp: "对选定序号的用户进行学习",
		Args: func(a *grumble.Args) {
			a.Int("index", "the index user")
		},
		Run: func(c *grumble.Context) error {
			index := c.Args.Int("index")
			users, err := model.Query()
			if err != nil {
				return err
			}
			if len(users) > index {
				_, _ = c.App.Println("the index not exist")
				return nil
			}
			user := users[index]
			go func() {
				config := conf.GetConfig()
				l := &lib.Core{Push: push.GetPush(conf.GetConfig()), ShowBrowser: config.ShowBrowser}
				l.Init()
				defer l.Quit()
				l.LearnArticle(user)
				l.LearnVideo(user)
				if config.Model == 2 {
					l.RespondDaily(user, "daily")
				} else if config.Model == 3 {
					l.RespondDaily(user, "daily")
					l.RespondDaily(user, "weekly")
					l.RespondDaily(user, "special")
				}
			}()
			return nil
		},
	}
}

func score() *grumble.Command {
	return &grumble.Command{
		Name:     "get score",
		Aliases:  []string{"score"},
		Help:     "get score",
		LongHelp: "查询用户积分",
		Run: func(c *grumble.Context) error {
			users, err := model.Query()
			if err != nil {
				return err
			}
			for _, user := range users {
				score, err := lib.GetUserScore(user.ToCookies())
				if err != nil {
					return err
				}
				_, _ = c.App.Println(user.Nick + "\n" + lib.FormatScore(score) + "\n\n")
			}

			return err
		},
	}
}

func addUser() *grumble.Command {
	return &grumble.Command{
		Name:     "add user",
		Aliases:  []string{"add"},
		Help:     "add a user",
		LongHelp: "添加一个用户",
		Run: func(c *grumble.Context) error {
			core := &lib.Core{
				Push: push.GetPush(conf.GetConfig()),
			}
			_, err := core.L(conf.GetConfig().Retry.Intervals, "")
			if err != nil {
				return err
			}
			return err
		},
	}
}

func getUser() *grumble.Command {
	return &grumble.Command{
		Name:      "getUser",
		Aliases:   []string{"user"},
		Help:      "get all user",
		LongHelp:  "input the user,can get all user",
		HelpGroup: "",
		Usage:     "获取用户列表",
		Run: func(c *grumble.Context) error {
			users, err := model.Query()
			if err != nil {
				log.Errorln(err.Error())
				return err
			}
			for i, user := range users {
				_, _ = c.App.Println(fmt.Printf("%d %v %v", i+1, user.Nick, utils.Stamp2Str(user.LoginTime)))
			}

			return nil
		},
		Completer: nil,
	}
}
