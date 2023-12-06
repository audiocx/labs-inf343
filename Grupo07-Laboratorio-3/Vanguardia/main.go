package main

import (
	"bufio"
	"context"
	messages "lab3/proto"
	"log"
	"os"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serverName string = "MonotonicReads"
var brokerAddress string = ""

func GetSoldados(sector string, base string) int {
	conn, err := grpc.Dial(brokerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("No se pudo conectar a %s: %s", brokerAddress, err)
	}
	defer conn.Close()

	c := messages.NewMessageServiceClient(conn)

	request := &messages.Info{
		Sector: sector,
		Base:   base,
	}

	responseBroker, err := c.GetSoldados(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	return int(responseBroker.Valor)
}

func main() {
	brokerAddress = os.Args[1]
	eSector := ""
	eBase := ""
	var input string
	scanner := bufio.NewScanner(os.Stdin)
	for {
		log.Printf("[%s] Ingrese el sector y la base que se desea consultar <sector> <base>: ", serverName)

		if scanner.Scan() {
			input = scanner.Text()
		}
		if input == "" {
			break
		}

		params := strings.Split(input, " ")

		eSector = params[0]
		eBase = params[1]

		cantidad_soldados := GetSoldados(eSector, eBase)
		log.Printf("[%s] Hay %v soldados en sector %s y base %s", serverName, cantidad_soldados, eSector, eBase)
	}

}
