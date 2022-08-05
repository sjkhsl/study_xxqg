package utils

import (
	"os"
	"os/exec"
)

// Restart
/* @Description:
 */
func Restart() {
	cmd := exec.Command("./study_xxqg")
	go func() {
		cmd.Start()
		os.Exit(3)
	}()

}
