package resources

import "github.com/silenteh/gantryos/utils"

// get the OS model version
func OsModel() string {
	detectedOs := detectOS()

	version := "unknown - lsb_release missing ?"

	switch detectedOs {
	case BSD:
		productName := utils.ExecCommand(true, "sw_vers", "-productName")       // osx architecture Ex: x86_64
		productVersion := utils.ExecCommand(true, "sw_vers", "-productVersion") // osx architecture Ex: x86_64
		productBuild := utils.ExecCommand(true, "sw_vers", "-buildVersion")     // osx architecture Ex: x86_64
		version = utils.ConcatenateString(productName, ": ", productVersion, " - build: ", productBuild)
		break
	case LINUX:
		productName := utils.ExecCommand(true, "lsb_release", "-ds")
		codeName := utils.ExecCommand(true, "lsb_release", "-cs")
		version = utils.ConcatenateString(productName, " (", codeName, ")")
		break
	default:
		break
	}
	return version

}
