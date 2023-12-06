// ReadYourWrites.main.go
package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	messages "lab3/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serverName string = "MonotonicWrites"
var brokerAddress string

type Server struct {
	messages.UnimplementedMessageServiceServer
}

// Vars to keep modifiedSectors
var modifiedSectors map[string]*messages.VectorClock = make(map[string]*messages.VectorClock)
var fulcrumAddresses []string

func askAddres(cmd, sector, base, newBase string, value int32) string {
	conn, err := grpc.Dial(brokerAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("No se pudo conectar a %s: %s", brokerAddress, err)
	}
	defer conn.Close()

	c := messages.NewMessageServiceClient(conn)

	request := &messages.Cmd{
		Cmd:       cmd,
		Sector:    sector,
		Base:      base,
		NuevaBase: newBase,
		Valor:     value,
	}

	responseBroker, err := c.AskAddress(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	return responseBroker.FulcrumAddress
}

func executeCommand(cmd string, sector string, base string, newBase string, value int32, address string) *messages.VectorClock {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("No se pudo conectar a %s: %s", address, err)
	}
	defer conn.Close()

	c := messages.NewMessageServiceClient(conn)

	request := &messages.Cmd{
		Cmd:       cmd,
		Sector:    sector,
		Base:      base,
		NuevaBase: newBase,
		Valor:     value,
	}

	responseBroker, err := c.Command(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	return responseBroker
}

func RequestMerge() {
	conn, err := grpc.Dial(fulcrumAddresses[0], grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("No se pudo conectar a %s: %s", fulcrumAddresses[0], err)
	}
	defer conn.Close()

	c := messages.NewMessageServiceClient(conn)

	request := &messages.AskMerge{
		FulcrumName: "Fulcrum1",
	}

	responseBroker, err := c.RequestMerge(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(responseBroker)
}

func main() {

	// Inicializamos las direcciones de Fulcrum y broker
	fulcrumAddresses = append(fulcrumAddresses, os.Args[2], os.Args[3], os.Args[4])
	brokerAddress = os.Args[1]

	// Vars
	var input string

	var cmd string
	var sector string
	var base string
	newBase := ""
	value := int32(0)
	var tempValue int

	for {

		// Get input
		fmt.Println("Ingrese uno de los siguientes comandos:\n* AgregarBase <nombre_sector> <nombre_base> <valor>\n* RenombrarBase <nombre_sector> <nombre_base> <nuevo_nombre>\n* ActualizarValor <nombre_sector <nombre_base> <nuevo_valor>\n* BorrarBase <nombre_sector> <nombre_base>")
		scanner := bufio.NewScanner(os.Stdin)

		if scanner.Scan() {
			input = scanner.Text()
		}
		if input == "" {
			break
		}

		// Extract data
		params := strings.Split(input, " ")

		if len(params) < 3 {
			fmt.Println("Ninguno de los parametros puede estar vacio.")
			continue
		}

		cmd = params[0]
		sector = params[1]
		base = params[2]

		if cmd == "AgregarBase" {
			if len(params) == 4 {
				tempValue, _ = strconv.Atoi(params[3])
				value = int32(tempValue)
			} else {
				value = 0
			}
		} else if cmd == "RenombrarBase" {
			if len(params) == 4 {
				newBase = params[3]
			} else {
				fmt.Println("Ninguno de los parametros puede estar vacio.")
				continue
			}
		} else if cmd == "ActualizarValor" {
			if len(params) == 4 {
				tempValue, _ = strconv.Atoi(params[3])
				value = int32(tempValue)
			} else {
				fmt.Println("Ninguno de los parametros puede estar vacio.")
				continue
			}
		}

		// Merge and execute
		address := askAddres(cmd, sector, base, newBase, value)
		clock := executeCommand(cmd, sector, base, newBase, value, address)
		RequestMerge()
		modifiedSectors[sector] = clock
		RequestMerge()
	}
}
