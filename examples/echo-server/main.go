package main

import (
	"log"
	"net"

	"github.com/ejuju/jus/pkg/jutp"
)

const welcomeMessage = `
"
Welcome to the echo server.
Say something and I will repeat it back to you!

"
write

*prompt ["?> " write read retrieve] define
[prompt true] repeat
`

func main() {
	log.Println("starting echo server on port 8080")
	jutp.Serve(&net.TCPAddr{Port: 8080}, func(rui *jutp.RemoteUI) {
		err := rui.Exec(welcomeMessage)
		if err != nil {
			log.Println(err)
			return
		}

		for {
			msg, err := rui.Read()
			if err != nil {
				log.Println(err)
				return
			}
			// msg = strings.ReplaceAll(msg, `"`, `\"`) // TODO: escape potential quote characters
			err = rui.Exec(`"Received: ` + string(msg) + `\n" write`)
			if err != nil {
				log.Println(err)
				return
			}
		}
	})
}
