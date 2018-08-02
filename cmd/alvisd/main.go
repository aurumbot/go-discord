package main

import (
	"fmt"
	dsg "github.com/bwmarrin/discordgo"
	c "github.com/skilstak/go-colors"
	"github.com/takama/daemon"
	f "github.com/whitman-colm/go-discord"
	"github.com/whitman-colm/go-discord/cmd/handler"
	"github.com/whitman-colm/go-discord/dat"
	"net"
	"os"
	"os/signal"
	"syscall"
)

type Service struct {
	daemon.Daemon
}

const (
	name        = "godiscordd"
	description = "A go-discord based bot daemon"
)

var port string

func (service *Service) Manage() (string, error) {
	usage := "Usage: " + name + " install | remove | start | stop | status"

	if len(os.Args) > 1 {
		command := os.Args[1]
		switch command {
		case "install":
			return service.Install()
		case "remove":
			return service.Remove()
		case "start":
			return service.Start()
		case "stop":
			return service.Stop()
		case "status":
			return service.Status()
		default:
			return usage, nil
		}
	}

	runBot()

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, os.Kill, syscall.SIGTERM)

	listener, err := net.Listen("tcp", port)
	if err != nil {
		dat.Log.Println(err)
		return "Possibly was a problem with the port binding", err
	}

	listen := make(chan net.Conn, 100)
	go acceptConnection(listener, listen)

	for {
		select {
		case conn := <-listen:
			go handleClient(conn)
		case killSignal := <-interrupt:
			f.DG.Close()
			dat.Log.Println("Recived signal:", killSignal)
			dat.Log.Println("Stopping listening on ", listener.Addr())
			listener.Close()
			if killSignal == os.Interrupt {
				dat.Log.Println("Daemon was interrupted by system signal")
				return "Daemon was interrupted by system signa", nil
			}
			dat.Log.Println("Daemon was killed")
			return "Daemon was killed", nil
		}
	}
	return usage, nil
}

func acceptConnection(listener net.Listener, listen chan<- net.Conn) {
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		listen <- conn
	}
}

func handleClient(client net.Conn) {
	for {
		buf := make([]byte, 4096)
		numbytes, err := client.Read(buf)
		if numbytes == 0 || err != nil {
			return
		}
		client.Write(buf[:numbytes])
	}
}

func main() {
	srv, err := daemon.New(name, description, "")
	if err != nil {
		dat.Log.Println(err)
		os.Exit(1)
	}
	service := &Service{srv}
	status, err := service.Manage()
	if err != nil {
		dat.Log.Println(err)
		os.Exit(1)
	}
	fmt.Println(status)

	// And the main should end here. so what to do about this?
}

func runBot() {
	bot, err := dat.GetBotInfo()
	fmt.Println(c.B0 + "Reading bot prefrences file...")
	dat.Log.Println("Reading bot prefs file...")
	if err != nil {
		dat.Log.Fatalln(err.Error())
	} else {
		dat.Log.Println("Bot prefrences recived.")
	}
	f.MyBot = bot
	port = bot.Auth.Port

	dat.Log.Println("Creating bot session")
	dg, err := dsg.New("Bot " + bot.Auth.Token)
	if err != nil {
		dat.Log.Fatalln(err)
	} else {
		dat.Log.Println("Session successfully created.")
	}

	dg.AddHandler(cmd.MessageCreate)

	dat.Log.Println("Opening websocket to Discord")
	err = dg.Open()
	if err != nil {
		dat.Log.Fatalln(err.Error())
	} else {
		dat.Log.Println("Socket successfully opened.")
	}
	f.DG = dg

	f.DG.UpdateStatus(0, f.MyBot.Prefs.Playing)

	dat.Log.Println("Escape for bot called. The system is now closing cleanly")
}

/*func installWizard(service *Service) (string, error) {
	fmt.Println("You will be taken through the process of installing an instance of")
	fmt.Println("the go-discord bot. Please have a bot user ready from")
	fmt.Println("https://discordapp.com/developers/applications/ before progressing.")
	clid, err := input.Prompt("Please enter the client ID:\n> ")
	if err != nil {
		dat.Log.Println(err)
		return "", err
	}
	clsc, err := input.Prompt("Please enter the client secret:\n> ")
	if err != nil {
		dat.Log.Println(err)
		return "", err
	}
	cltk, err := input.Prompt("Please enter the token:\n> ")
	if err != nil {
		dat.Log.Println(err)
		return "", err
	}
	pref, err := input.Prompt("Please decide on a prefix. This is used to call the bot (leave blank to only have an @ mention trigger)\nPrefix: ")
	if err != nil {
		dat.Log.Println(err)
		return "", err
	}
}*/
