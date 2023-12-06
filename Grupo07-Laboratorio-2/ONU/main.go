package main

import (
	"bufio"
	"context"
	messages "lab2/proto"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var sv_name string = "ONU"
var OMS_ip_port string = "dist025:52000"

func RequestDeadOrAlive(ip string, state int) *messages.Names {
	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Print(err)
	}
	defer conn.Close()

	c := messages.NewMessageServiceClient(conn)

	// Setea la informacion en el struct
	req := messages.State{
		Value: int32(state),
	}

	// Envia el mensaje, responde con los nombres
	names, err := c.SendRequestDeadOrAlive(context.Background(), &req)

	if err != nil {
		log.Fatal(err)
	}

	return names
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		log.Printf("[%s] Waiting for requests (<0> DEAD, <1> INFECTED, <x> EXIT)...", sv_name)
		opt, err := reader.ReadString('\n')
		if err != nil {
			log.Fatal(err)
		}
		opt = strings.TrimSuffix(opt, "\n")
		log.Print(opt)

		if opt == "0" {
			log.Printf("[%s] Requesting the names of DEAD people to OMS...", sv_name)

			names := RequestDeadOrAlive(OMS_ip_port, 0)

			for i := 0; i < int(names.Len); i++ {
				log.Printf("[%s] DEAD: %s %s", sv_name, names.List[i].Name, names.List[i].Lastname)
			}
		} else if opt == "1" {
			log.Printf("[%s] Requesting the names of INFECTED people to OMS...", sv_name)

			names := RequestDeadOrAlive(OMS_ip_port, 1)

			for i := 0; i < int(names.Len); i++ {
				log.Printf("[%s] INFECTED: %s %s", sv_name, names.List[i].Name, names.List[i].Lastname)
			}
		} else if opt == "x" {
			log.Printf("[%s] Closing...", sv_name)
			break
		} else {
			log.Printf("[%s] Invalid option, skipping...", sv_name)
		}
	}

}
