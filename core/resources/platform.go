package resources

import (
	"github.com/silenteh/gantryos/utils"
	"regexp"
	"strings"
)

type lsb struct{}

type platform struct {
	OS      string
	Version string
	Family  string
}

// this one matches Fedora Rawhide, which is the development branch of Fedora.
// all fedora distros in development have (Rawhide) string
var matchRedHatishReg *regexp.Regexp = regexp.MustCompile(`(?i)Rawhide`)
var matchedRedHatishReg *regexp.Regexp = regexp.MustCompile(`(?i)((\d+) \(Rawhide\))`)
var unmatchedRedHatishReg *regexp.Regexp = regexp.MustCompile(`release ([\d\.]+)`)
var redHatPlatformReg *regexp.Regexp = regexp.MustCompile(`(?i)^Red Hat`)
var nonRedHatredHatPlatformReg *regexp.Regexp = regexp.MustCompile(`(?i)(\w+)`)
var ubuntuReg *regexp.Regexp = regexp.MustCompile(`(?i)Ubuntu`)
var linuxMintReg *regexp.Regexp = regexp.MustCompile(`(?i)LinuxMint`)
var redHatReg *regexp.Regexp = regexp.MustCompile(`(?i)RedHat`)
var amazonReg *regexp.Regexp = regexp.MustCompile(`(?i)Amazon`)
var xenServerReg *regexp.Regexp = regexp.MustCompile(`(?i)XenServer`)

func (lsb *lsb) Version() string {
	return utils.ExecCommand(true, "lsb_release", "-rs")
}

func (lsb *lsb) Id() string {
	return utils.ExecCommand(true, "lsb_release", "-vs")
}

func (lsb *lsb) CodeName() string {
	return utils.ExecCommand(true, "lsb_release", "-cs")
}

func (lsb *lsb) Description() string {
	return utils.ExecCommand(true, "lsb_release", "-ds")
}

func detectPlatform() platform {
	platform := detectSpecificPlatform()
	switch platform.OS {
	case "debian", "ubuntu", "linuxmint", "raspbian":
		platform.Family = "debian"
	case "fedora":
		platform.Family = "fedora"
	case "oracle", "centos", "redhat", "scientific", "enterpriseenterprise", "amazon", "xenserver", "cloudlinux", "ibm_powerkvm", "parallels":
		platform.Family = "rhel"
	case "suse":
		platform.Family = "suse"
	case "gentoo":
		platform.Family = "gentoo"
	case "slackware":
		platform.Family = "slackware"
	case "arch":
		platform.Family = "arch"
	case "exherbo":
		platform.Family = "exherbo"
	default:
		platform.Family = "unknown"
	}

	return platform
}

func detectSpecificPlatform() platform {

	platform := platform{"Unknown", "Unknown", "Unknown"}

	if utils.FileExists("/etc/oracle-release") {
		content := utils.ReadFileToString("/etc/oracle-release")
		contentChomp := utils.Chomp(content, "\r\n")
		platform.OS = "oracle"
		platform.Version = getRedhatishVersion(contentChomp)
		return platform
	}
	if utils.FileExists("/etc/enterprise-release") {
		content := utils.ReadFileToString("/etc/enterprise-release")
		contentChomp := utils.Chomp(content, "\r\n")
		platform.OS = "oracle"
		platform.Version = getRedhatishVersion(contentChomp)
		return platform
	}
	if utils.FileExists("/etc/parallels-release") {
		content := utils.ReadFileToString("/etc/parallels-release")
		contentChomp := utils.Chomp(content, "\r\n")
		platform.OS = getRedhatishPlatform(contentChomp)
		platform.Version = getRedhatishVersion(contentChomp)
		return platform
	}
	if utils.FileExists("/etc/redhat-release") {
		content := utils.ReadFileToString("/etc/redhat-release")
		contentChomp := utils.Chomp(content, "\r\n")
		platform.OS = getRedhatishPlatform(contentChomp)
		platform.Version = getRedhatishVersion(contentChomp)
		return platform
	}
	if utils.FileExists("/etc/redhat-release") {
		content := utils.ReadFileToString("/etc/redhat-release")
		contentChomp := utils.Chomp(content, "\r\n")
		platform.OS = getRedhatishPlatform(contentChomp)
		platform.Version = getRedhatishVersion(contentChomp)
		return platform
	}
	if utils.FileExists("/etc/system-release") {
		content := utils.ReadFileToString("/etc/system-release")
		contentChomp := utils.Chomp(content, "\r\n")
		platform.OS = getRedhatishPlatform(contentChomp)
		platform.Version = getRedhatishVersion(contentChomp)
		return platform
	}
	if utils.FileExists("/etc/system-release") {
		content := utils.ReadFileToString("/etc/system-release")
		contentChomp := utils.Chomp(content, "\r\n")
		platform.OS = getRedhatishPlatform(contentChomp)
		platform.Version = getRedhatishVersion(contentChomp)
		return platform
	}

	if utils.FileExists("/etc/debian_version") {
		// Ubuntu and Debian both have /etc/debian_version
		// Ubuntu should always have a working lsb, debian does not by default
		lsb := lsb{}

		if ubuntuReg.MatchString(lsb.Id()) {
			platform.OS = "ubuntu"
			platform.Version = lsb.Version()
			return platform
		}
		if linuxMintReg.MatchString(lsb.Id()) {
			platform.OS = "linuxmint"
			platform.Version = lsb.Version()
			return platform
		}

		if utils.FileExists("/usr/bin/raspi-config") {
			platform.OS = "raspbian"
		} else {
			platform.OS = "debian"
		}

		content := utils.ReadFileToString("/etc/debian_version")
		contentChomp := utils.Chomp(content, "\r\n")
		platform.Version = contentChomp
		return platform
	}

	return platform
}

func getRedhatishVersion(content string) string {
	matched := matchRedHatishReg.MatchString(content)

	if matched {
		results := matchedRedHatishReg.FindStringSubmatch(content)
		return strings.ToLower(getSubContent(results, 1))
	}

	results := unmatchedRedHatishReg.FindStringSubmatch(content)
	return getSubContent(results, 1)

}

func getRedhatishPlatform(content string) string {
	matched := redHatPlatformReg.MatchString(content)
	if matched {
		return "redhat"
	}

	results := nonRedHatredHatPlatformReg.FindStringSubmatch(content)

	return strings.ToLower(getSubContent(results, 1))
}

// this function gets a specific index of the string array
// it checks the cases where the index may be > or < than the size
func getSubContent(data []string, pos int) string {

	if data == nil || len(data) == 0 {
		return ""
	}

	if len(data) >= pos+1 {
		return data[pos]
	}

	return ""

}
