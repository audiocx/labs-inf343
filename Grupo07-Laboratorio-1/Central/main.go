package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	messages "sd/proto"
	"strconv"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"
)

func MinMaxIter(filename string) (int, int, int) {
	byteData, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatal(err)
	}

	data := string(byteData)
	data = strings.Replace(data, "\n", "/", 1)
	partes := strings.Split(data, "/")

	min, err := strconv.Atoi(strings.Split(partes[0], "-")[0])
	if err != nil {
		log.Fatal(err)
	}

	max, err := strconv.Atoi(strings.Split(partes[0], "-")[1])
	if err != nil {
		log.Fatal(err)
	}

	iter, err := strconv.Atoi(partes[1])
	if err != nil {
		log.Fatal(err)
	}

	return min, max, iter
}

// Set to virtual machine ips
var G_SERVER_IP_PORT = [4]string{
	"dist025:50001",
	"dist026:50002",
	"dist027:50003",
	"dist028:50004",
}

var G_OUTPUT_FILE_ROUTE string = "/home/dist/tarea1/Grupo07-Laboratorio-1/Central/output.txt"

func aleatorio(menor, mayor int) int {
	return rand.Intn(mayor-menor) + menor
}

func sendKeys(ip string, keys int) {
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		log.Printf("No se pudo conectar a %s: %s", ip, err)
	}
	defer conn.Close()

	c := messages.NewMessageServiceClient(conn)

	msg := messages.Keys{
		Amount: int64(keys),
	}

	resp, err := c.SendAvailableKeys(context.Background(), &msg)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	resp.Ack = 0
}

func sendNonDeliveredKeys(ip string, ndk int) {
	conn, err := grpc.Dial(ip, grpc.WithInsecure())
	if err != nil {
		log.Printf("Could not connect to %s: %s", ip, err)
	}
	defer conn.Close()

	c := messages.NewMessageServiceClient(conn)

	msg := messages.Keys{
		Amount: int64(ndk),
	}

	resp, err := c.SendNonDelivered(context.Background(), &msg)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	resp.Ack = 0
}

func appendToOutputFile(data string) {

	f, err := os.OpenFile(G_OUTPUT_FILE_ROUTE, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(data + "\n")); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}
}

func main() {

	min, max, iter := MinMaxIter("/app/Central/parametros_de_inicio.txt")
	// Nos conectamos a la cola RabbitMQ
	conn_queue, err := amqp.Dial("amqp://guest:guest@dist026:5672")
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

	_ = q

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	msgs, err := ch.Consume(
		"ApplicationQueue",
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	n_keys := aleatorio(min, max)
	currentTime := time.Now().Format("15:04:05")
	appendToOutputFile(currentTime + " - " + strconv.Itoa(n_keys))

	if iter != -1 {
		log.Printf("[CENTRAL] 1/%v generation, %v keys generated, sending to regional servers", iter, n_keys)
	} else {
		log.Printf("[CENTRAL] 1/Inf generation, %v keys generated, sending to regional servers", n_keys)
	}

	// Enviar informacion de llaves disponibles
	go sendKeys(G_SERVER_IP_PORT[0], n_keys)
	go sendKeys(G_SERVER_IP_PORT[1], n_keys)
	go sendKeys(G_SERVER_IP_PORT[2], n_keys)
	go sendKeys(G_SERVER_IP_PORT[3], n_keys)

	var n_servers_requesting int = 4
	var curr_iter int = 1
	var counter int = 0

	log.Println("[CENTRAL] Initializing RMQ...")
	log.Println("[From RMQ] Waiting for applications...")
	noStop := make(chan bool)
	go func() {

		for d := range msgs {

			if curr_iter == iter {
				log.Printf("[CENTRAL] max number of iterations reached, closing...")
				ch.Close()
			}

			if iter != -1 {
				if counter%n_servers_requesting == 0 && counter != 0 {
					n_keys = aleatorio(min, max)
					log.Printf("[CENTRAL] %v/%v generation, %v keys generated", curr_iter+1, iter, n_keys)
					curr_iter += 1
					currentTime = time.Now().Format("15:04:05")
					appendToOutputFile(currentTime + " - " + strconv.Itoa(n_keys))
				}
			} else {
				if counter%n_servers_requesting == 0 && counter != 0 {
					n_keys = aleatorio(min, max)
					log.Printf("[CENTRAL] %v/Inf generation, %v keys generated", curr_iter+1, n_keys)
					curr_iter += 1
					currentTime = time.Now().Format("15:04:05")
					appendToOutputFile(currentTime + " - " + strconv.Itoa(n_keys))
				}
			}

			data := string(d.Body)
			parts := strings.Split(data, "-")
			asked, err := strconv.Atoi(parts[0])
			if err != nil {
				log.Fatal(err)
			}
			from := parts[1]
			log.Printf("[From RMQ]: %s asked for %v keys\n", from, asked)

			if asked > n_keys { // Si la cantidad de pedidos es mayor a las disponibles (no satisface completamente)
				log.Printf("[CENTRAL] %v registered on server %s, 0 keys remaining", n_keys, from)
				sendNonDeliveredKeys(from, asked-n_keys)
				appendToOutputFile(from + "-" + strconv.Itoa(asked) + "-" + strconv.Itoa(n_keys) + "-" + strconv.Itoa(asked-n_keys))
				n_keys = 0
			} else { // Si la cantidad de pedidos es igual o menor a las disponibles (satisface completamente)
				log.Printf("[CENTRAL] %v registered on server %s, %v keys remaining", asked, from, n_keys-asked)
				sendNonDeliveredKeys(from, 0)
				n_keys -= asked
				n_servers_requesting -= 1
				appendToOutputFile(from + "-" + strconv.Itoa(asked) + "-" + strconv.Itoa(asked) + "-0")
			}

			if n_servers_requesting == 0 {
				log.Printf("[CENTRAL] all regional servers satisfied, closing with %v keys remaining...", n_keys)
				ch.Close()
			}

			counter += 1
		}
	}()
	<-noStop

}
