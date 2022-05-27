package update

import "testing"

func Test_versionCompare(t *testing.T) {
	println(versionCompare("v2.0.1", "v2.0.3"))
	println(versionCompare("v2.0.2", "v2.0.2"))
	println(versionCompare("v2.0.2", "v2.0.2-beta1"))
	println(versionCompare("v2.0.2-beta1", "v2.0.2-beta3"))
	println(versionCompare("v2.0.2-beta1", "v2.0.2-beta1"))
}

func Test_CheckUpdate(t *testing.T) {
	CheckUpdate("v1.0.22")
}
