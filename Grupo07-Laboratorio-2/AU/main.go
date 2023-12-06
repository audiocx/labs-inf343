package main

import (
	"bufio"
	"container/list"
	"context"
	"crypto/rand"
	messages "lab2/proto"
	"log"
	"math/big"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var sv_name string = "AU"
var OMS_ip_port string = "dist025:51000"
var data string = "/app/AU/names.txt"

func SendNameLastnameState(ip string, name string, lastname string, state int) {
	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("No se pudo conectar a %s: %s", ip, err)
	}
	defer conn.Close()

	c := messages.NewMessageServiceClient(conn)

	info := messages.Info{
		Name:     name,
		Lastname: lastname,
	}

	msg := messages.InfoState{
		Info:  &info,
		From:  sv_name,
		State: int32(state),
	}

	resp, err := c.SendInfoState(context.Background(), &msg)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	resp.Ack = 0
}

func readNames(filename string, names *list.List) {
	f, err := os.Open(filename)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		names.PushFront(scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

}

func randState() int {
	state := 0 // Muerto

	num, err := rand.Int(rand.Reader, big.NewInt(100))
	if err != nil {
		log.Fatal(err)
	}
	if num.Int64() < 55 {
		state = 1 // Vivo
	}

	return state
}

func main() {
	names := list.New()
	readNames(data, names)

	iter := names.Front()
	for i := 0; i < 5; i++ {
		info := iter.Value.(string)
		info = strings.Replace(info, "\n", "", 1)
		nameLastname := strings.Split(info, " ")

		state := randState()

		SendNameLastnameState(OMS_ip_port, nameLastname[0], nameLastname[1], state)

		if state == 0 {
			log.Printf("[%s] Sent: %s %s (DEAD)", sv_name, nameLastname[0], nameLastname[1])
		} else {
			log.Printf("[%s] Sent: %s %s (INFECTED)", sv_name, nameLastname[0], nameLastname[1])
		}

		iter = iter.Next()
		if iter == nil {
			break
		}
	}

	for i := 5; i < names.Len(); i++ {
		time.Sleep(3 * time.Second)

		info := iter.Value.(string)
		info = strings.Replace(info, "\n", "", 1)
		nameLastname := strings.Split(info, " ")

		state := randState()

		SendNameLastnameState(OMS_ip_port, nameLastname[0], nameLastname[1], state)

		if state == 0 {
			log.Printf("[%s] Sent: %s %s (DEAD)", sv_name, nameLastname[0], nameLastname[1])
		} else {
			log.Printf("[%s] Sent: %s %s (INFECTED)", sv_name, nameLastname[0], nameLastname[1])
		}

		iter = iter.Next()
		if iter == nil {
			break
		}
	}

}
