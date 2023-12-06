package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	messages "sd/proto"
	"strconv"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

var interested int
var sv_name string = "Oceania"
var rmq_ip string = "amqp://guest:guest@dist026:5672/"
var self_ip string = "dist028:50004"
var self_port string = ":50004"

func InterestedUsers(filename string) int {
	byteData, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	data := string(byteData)
	data = strings.Replace(data, "\n", "", 1)

	users, err := strconv.Atoi(data)
	if err != nil {
		log.Fatal(err)
	}

	return users
}

func randomUsers(users int) int {
	min := users/2 - users/5
	max := users/2 + users/5
	return rand.Intn(max-min) + min
}

type Server struct {
	messages.UnimplementedMessageServiceServer
}

func (s *Server) SendAvailableKeys(ctx context.Context, msg *messages.Keys) (*messages.Ack, error) {
	log.Printf("[FROM Central] Available keys: %v", msg.Amount)
	if interested > 0 {
		log.Printf("[%s] %v interested users, asking %v keys to Central", sv_name, interested, interested)
		SendRequiredKeys(interested, rmq_ip, self_ip)
	}
	return &messages.Ack{Ack: 1}, nil
}

func (s *Server) SendNonDelivered(ctx context.Context, msg *messages.Keys) (*messages.Ack, error) {
	asked := interested
	interested = int(msg.Amount)
	log.Printf("[%s] %v users registered", sv_name, asked-interested)
	if interested > 0 {
		log.Printf("[FROM Central] Could not satisfy %v keys of %v asked", msg.Amount, asked)
		log.Printf("[%s] %v interested users remaining, asking %v keys to Central", sv_name, interested, interested)
		SendRequiredKeys(interested, rmq_ip, self_ip)
	} else {
		log.Printf("[%s] %v interested users remaining, skipping...", sv_name, interested)
	}
	return &messages.Ack{Ack: 1}, nil
}

func SendRequiredKeys(n_keys int, ip string, from string) {
	conn_queue, err := amqp.Dial(ip)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}
	defer conn_queue.Close()

	ch, err := conn_queue.Channel()
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	q, err := ch.QueueDeclare(
		"ApplicationQueue",
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	_ = q

	amqpContext, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err = ch.PublishWithContext(
		amqpContext,
		"",
		"ApplicationQueue",
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(strconv.Itoa(n_keys) + "-" + from),
		},
	)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	conn_queue.Close()
}

func conn_server(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("Failed to listen on port %s: %v", port, err)
	}

	grpcServer := grpc.NewServer()

	messages.RegisterMessageServiceServer(grpcServer, &Server{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to serve gRPC server over port %s: %v", port, err)
	}
}

func main() {

	users := InterestedUsers("/app/Europa/parametros_de_inicio.txt")

	interested = randomUsers(users)

	log.Printf("[%s] Server initialized, waiting for keys...\n", sv_name)

	go conn_server(self_port)

	for {

	}
}
