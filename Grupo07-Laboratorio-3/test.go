package main

import (
	"context"
	"log"
	"os"
	"time"

	messages "lab3/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var brokerAddress string = ""

func SendCommand(cmd string, sector string, base string, nuevaBase string, valor int) *messages.VectorClock {
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
		NuevaBase: nuevaBase,
		Valor:     int32(valor),
	}

	responseBroker, err := c.AskAddress(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	conn2, err := grpc.Dial(responseBroker.FulcrumAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("No se pudo conectar a %s: %s", brokerAddress, err)
	}
	defer conn2.Close()

	c2 := messages.NewMessageServiceClient(conn2)

	responseFulcrum, err := c2.Command(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

	return responseFulcrum
}

func mergeRequest(address string, fulcrumName string) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("No se pudo conectar a %s: %s", brokerAddress, err)
	}
	defer conn.Close()

	c := messages.NewMessageServiceClient(conn)

	request := &messages.AskMerge{FulcrumName: fulcrumName}

	_, err = c.RequestMerge(context.Background(), request)
	if err != nil {
		log.Fatal(err)
	}

}

func main() {
	brokerAddress = os.Args[1]

	SendCommand("AgregarBase", "SectorAlpha", "Base1", "", 300)
	time.Sleep(time.Millisecond * 500)
	SendCommand("AgregarBase", "SectorBeta", "BaseXD", "", 22)
	time.Sleep(time.Millisecond * 500)
	SendCommand("RenombrarBase", "SectorBeta", "BaseAAA", "", 0)
	time.Sleep(time.Millisecond * 500)
	SendCommand("ActualizarValor", "SectorAlpha", "Base1", "", 1)
	time.Sleep(time.Millisecond * 500)
	SendCommand("AgregarBase", "SectorGiga", "BaseGiga", "", 1)
	time.Sleep(time.Millisecond * 500)
	SendCommand("BorrarBase", "SectorGIga", "BaseGiga", "", 22)
	time.Sleep(time.Millisecond * 500)
	SendCommand("AgregarBase", "Sector1", "BaseS1", "", 22)
	time.Sleep(time.Millisecond * 500)
	SendCommand("AgregarBase", "Sector2", "BaseS2", "", 4)
	time.Sleep(time.Millisecond * 500)
	SendCommand("AgregarBase", "Sector3", "BaseS3", "", 22)
	time.Sleep(time.Millisecond * 500)
	SendCommand("AgregarBase", "Sector1", "Baseqwe", "", 22)
	time.Sleep(time.Millisecond * 500)
	SendCommand("AgregarBase", "Sector2", "BasQWE", "", 4)
	time.Sleep(time.Millisecond * 500)
	SendCommand("AgregarBase", "Sector3", "BaseEWQ", "", 22)
	time.Sleep(time.Millisecond * 500)

	mergeRequest(":5000", "Fulcrum1")
}
