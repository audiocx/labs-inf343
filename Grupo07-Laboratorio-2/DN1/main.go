package main

import (
	"bufio"
	"context"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"

	messages "lab2/proto"
)

var sv_name string = "DataNode1"
var feed_port string = ":53000"
var req_port string = ":54000"
var data string = "DATA.txt"

type Server struct {
	messages.UnimplementedMessageServiceServer
}

func (s *Server) SendInfoID(ctx context.Context, msg *messages.InfoID) (*messages.Ack, error) {
	log.Printf("[%s] Received from OMS: %s %s (ID: %v)", sv_name, msg.Info.Name, msg.Info.Lastname, msg.Id)

	// Implementar aca la recepcion de nombres, apellido e id y guardarlas en data.txt

	// CODIGO NO PROBADO
	archivo, err := os.OpenFile(data, os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer archivo.Close() // Asegúrate de cerrar el archivo al finalizar la función

	id_aux := strconv.Itoa(int(msg.Id))
	texto := id_aux + " " + msg.Info.Name + " " + msg.Info.Lastname + "\n"

	_, err = archivo.WriteString(texto)
	if err != nil {
		log.Fatal(err)
	}

	return &messages.Ack{Ack: 1}, nil
}

/*
funcion buscarPorID

Esta funcion busca por ID en DATA.txt y retorna el nombre y el apellido de la persona a quien corresponde el ID.

No retorna nada ya que imprime.
*/
func buscarPorID(idBuscado string) (string, string) {
	// Abrir el archivo DATA.txt para lectura
	file, err := os.Open(data)
	if err != nil {
		return "", ""
	}
	defer file.Close()

	// Leer línea por línea para encontrar el ID
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) >= 1 && parts[0] == idBuscado {
			// Encontramos el ID, retorna el nombre y el apellido
			if len(parts) >= 3 {
				return parts[1], parts[2]
			}
		}
	}

	// Si no se encontró el ID, retorna un error
	return "", ""
}

func (s *Server) SendRequestNamesByID(ctx context.Context, msg *messages.Ids) (*messages.Names, error) {
	log.Printf("[%s] %v IDs requested from OMS", sv_name, msg.Len)

	resp := messages.Names{
		List: []*messages.Info{},
		Len:  int64(0),
	}

	for i := 0; i < int(msg.Len); i++ {
		name, lastname := buscarPorID(strconv.Itoa(int(msg.List[i])))
		resp.List = append(resp.List, &messages.Info{Name: name, Lastname: lastname})
		resp.Len++
	}

	return &resp, nil
}

func feed_server(port string) {
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

func req_server(port string) {
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

	log.Printf("[%s] Starting...", sv_name)

	os.Create(data) // inicia el archivo vacio

	go feed_server(feed_port)
	go req_server(req_port)

	for {

	}
}
