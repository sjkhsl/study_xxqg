package utils

import (
	"os"
	"os/exec"

	"github.com/sjkhsl/study_xxqg/conf"
	"github.com/sjkhsl/study_xxqg/utils/update"
)

// Restart
/* @Description:
 */
func Restart() {
	cmd := exec.Command("./study_xxqg")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	go func() {
		cmd.Start()
		os.Exit(3)
	}()

}

func GetAbout() string {
	msg := "study_xxqg\n程序版本："
	msg += conf.GetVersion()
	msg += "\n" + update.CheckUpdate(conf.GetVersion())
	return msg
}
