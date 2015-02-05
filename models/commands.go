package models

import "github.com/silenteh/gantryos/core/proto"

type command struct {
	Shell                bool // this needs to be set to TRUE in order to run a shell command
	ShellSources         *commandSources
	EnvironmentVariables environmentVariables
	Arguments            arguments
	User                 *user
}

type commandSources []*commandSource

type commandSource struct {
	Path       string // this should be a local path or a remote path (downloadable via HTTPS for example)
	Executable bool
	Extract    bool // not supported yet
}

func (cs *commandSource) toProtoBuf() *proto.CommandInfo_URI {

	cmdSrc := new(proto.CommandInfo_URI)
	cmdSrc.Value = &cs.Path
	cmdSrc.Extract = &cs.Extract
	cmdSrc.Executable = &cs.Executable

	return cmdSrc

}

func (css commandSources) toProtoBuf() []*proto.CommandInfo_URI {

	commandSourcesProto := make([]*proto.CommandInfo_URI, len(css))

	for index, res := range css {
		commandSourcesProto[index] = res.toProtoBuf()
	}

	return commandSourcesProto

}

func NewCommandSource(path string, executable, extract bool) *commandSource {
	return &commandSource{
		Path:       path,
		Executable: executable,
		Extract:    extract,
	}
}

func NewCommand(shell bool, source *commandSources, envs []*environmentVariable, args arguments, user *user) *command {
	return &command{
		Shell:                shell, // this needs to be set to TRUE in order to run a shell command
		ShellSources:         source,
		EnvironmentVariables: envs,
		Arguments:            args,
		User:                 user,
	}
}

func (c *command) toProtoBuf() *proto.CommandInfo {

	cmd := new(proto.CommandInfo)
	cmd.Shell = &c.Shell
	cmd.Uris = c.ShellSources.toProtoBuf()
	cmd.Environment = c.EnvironmentVariables.toProtoBuf()
	cmd.Arguments = c.Arguments
	cmd.User = c.User.toProtoBuf()

	return cmd
}
