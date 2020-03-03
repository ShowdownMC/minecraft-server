package impl

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"minecraft-server/apis"
	"minecraft-server/apis/cmds"
	"minecraft-server/apis/data"
	"minecraft-server/apis/ents"
	"minecraft-server/apis/logs"
	"minecraft-server/apis/task"
	"minecraft-server/apis/util"

	"minecraft-server/impl/conn"
	"minecraft-server/impl/cons"
	"minecraft-server/impl/data/system"
	"minecraft-server/impl/data/values"
	"minecraft-server/impl/prot"

	apis_base "minecraft-server/apis/base"
	impl_base "minecraft-server/impl/base"
)

type server struct {
	message chan system.Message

	console *cons.Console
	logging *logs.Logging

	tasking *task.Tasking

	command *cmds.CommandManager

	network impl_base.Network
	packets impl_base.Packets
}

func NewServer(host string, port int) apis.Server {
	message := make(chan system.Message)

	console := cons.NewConsole(message)
	logging := logs.NewLogging("server", logs.EveryLevel...)

	tasking := task.NewTasking(values.MPT)

	join := make(chan impl_base.PlayerAndConnection)
	quit := make(chan impl_base.PlayerAndConnection)

	packets := prot.NewPackets(join, quit)
	network := conn.NewNetwork(host, port, packets, join, quit)

	command := cmds.NewCommandManager()

	return &server{
		message: message,

		console: console,

		logging: logging,
		tasking: tasking,

		command: command,

		packets: packets,
		network: network,
	}
}

func (s *server) Load() {
	apis.SetMinecraftServer(s)

	s.console.Load()
	s.command.Load()
	s.tasking.Load()
	s.network.Load()

	s.command.Register("stop", s.stopServerCommand)
	s.command.Register("time", func(sender ents.Sender, params []string) {
		var seconds int64 = 0

		if len(params) > 0 {
			param, err := strconv.Atoi(params[0])

			if err != nil {
				panic(err)
			}

			if param <= 0 {
				panic(fmt.Errorf("value must be a positive whole number. [1..]"))
			}

			seconds = int64(param)
		}

		sender.SendMessage(util.FormatTime(seconds))
	})

	go func() {
		for {
			// read input from console
			text := strings.Trim(<-s.console.IChannel, " ")
			if len(text) == 0 {
				continue
			}

			args := strings.Split(text, " ")
			if len(args) == 0 {
				continue
			}

			if command := s.command.Search(args[0]); command != nil {

				err := apis_base.Attempt(func() {
					(*command).Evaluate(s.console, args[1:])
				})

				if err != nil {
					s.logging.Fail(
						data.Red, "failed to evaluate ",
						data.DarkGray, "`",
						data.White, (*command).Name(),
						data.DarkGray, "`",
						data.Red, ": ", err.Error()[8:])
				}

				continue
			}

			s.console.SendMessage(text)
		}
	}()

	s.Wait()
}

func (s *server) Kill() {

	s.console.Kill()
	s.command.Kill()
	s.tasking.Kill()
	s.network.Kill()

	// push the stop message to the server exit channel
	s.message <- system.Make(system.STOP, "normal stop")
	close(s.message)

	s.logging.Info(data.DarkRed, "server stopped")
}

func (s *server) Wait() {
	// select over server commands channel
	select {
	case command := <-s.message:
		switch command.Command {
		// stop selecting when stop is received
		case system.STOP:
			return
		case system.FAIL:
			fmt.Printf("internal server error: %s\n", command.Message)
			return
		}
	}

	s.Wait()
}

func (s *server) Logging() *logs.Logging {
	return s.logging
}

func (s *server) Command() *cmds.CommandManager {
	return s.command
}

func (s *server) Tasking() *task.Tasking {
	return s.tasking
}

func (s *server) stopServerCommand(sender ents.Sender, params []string) {
	if _, ok := sender.(*cons.Console); !ok {
		s.logging.FailF("non console sender %s tried to stop the server", sender.Name())
		return
	}

	var after int64 = 0

	if len(params) > 0 {
		param, err := strconv.Atoi(params[0])

		if err != nil {
			panic(err)
		}

		if param <= 0 {
			panic(fmt.Errorf("value must be a positive whole number. [1..]"))
		}

		after = int64(param)
	}

	if after == 0 {

		s.Kill()

	} else {

		// inform future shutdown
		s.logging.Warn(data.Gold, "stopping server in ", data.Green, util.FormatTime(after))

		// schedule shutdown {after} seconds later
		s.tasking.AfterTime(after, time.Second, func(task *task.Task) {
			s.Kill()
		})

	}
}
