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
	"google.golang.org/grpc/credentials/insecure"

	messages "lab2/proto"
)

var sv_name string = "OMS"
var feed_port string = ":51000"
var req_port string = ":52000"

var DN1_ip_feed string = "dist027:53000"
var DN1_ip_req string = "dist027:54000"

var DN2_ip_feed string = "dist028:55000"
var DN2_ip_req string = "dist028:56000"

var id int = 0

const data = "DATA.txt"
const a_to_m = "abcdefghijklm"

type Server struct {
	messages.UnimplementedMessageServiceServer
}

func (s *Server) SendInfoState(ctx context.Context, msg *messages.InfoState) (*messages.Ack, error) {
	if msg.State == 0 {
		log.Printf("[%s] Received from %s: %s %s (DEAD)", sv_name, msg.From, msg.Info.Name, msg.Info.Lastname)
	} else {
		log.Printf("[%s] Received from %s: %s %s (INFECTED)", sv_name, msg.From, msg.Info.Name, msg.Info.Lastname)
	}

	DN := 2
	if strings.Contains(a_to_m, strings.ToLower(msg.Info.Lastname[0:1])) { // enviar a dn1
		DN = 1
	}

	// id - datanode - infectado/muerto
	id_aux := strconv.Itoa(id)
	texto := id_aux + " " + strconv.Itoa(DN) + " " + strconv.Itoa(int(msg.State))
	appendData(data, texto)

	if DN == 1 {
		SendNameLastnameID(DN1_ip_feed, msg.Info.Name, msg.Info.Lastname, id)
	} else {
		SendNameLastnameID(DN2_ip_feed, msg.Info.Name, msg.Info.Lastname, id)
	}

	id++
	return &messages.Ack{Ack: 1}, nil
}

func (s *Server) SendRequestDeadOrAlive(ctx context.Context, msg *messages.State) (*messages.Names, error) {

	if msg.Value == 0 {
		log.Printf("[%s] Received: request for DEAD people's names from ONU", sv_name)
	} else {
		log.Printf("[%s] Received: request for INFECTED people's names from ONU", sv_name)
	}

	resp := RequestNamesByID(DN1_ip_req, DN2_ip_req, int(msg.Value))

	return resp, nil
}

func feed_server(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	messages.RegisterMessageServiceServer(grpcServer, &Server{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func req_server(port string) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatal(err)
	}

	grpcServer := grpc.NewServer()

	messages.RegisterMessageServiceServer(grpcServer, &Server{})

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatal(err)
	}
}

func SendNameLastnameID(ip string, name string, lastname string, id int) {
	conn, err := grpc.Dial(ip, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Print(err)
	}
	defer conn.Close()

	c := messages.NewMessageServiceClient(conn)

	// Setea la informacion en el struct
	info := messages.Info{
		Name:     name,
		Lastname: lastname,
	}

	// Setea el struct del mensaje a enviar
	msg := messages.InfoID{
		Info: &info,
		Id:   int64(id),
	}

	// Envia el mensaje, responde con un ack
	resp, err := c.SendInfoID(context.Background(), &msg)

	log.Printf("[%s] Sent: %s %s (ID: %v) to %s", sv_name, name, lastname, id, ip)

	if err != nil {
		log.Fatalf("Error: %s", err)
	}

	resp.Ack = 0
}

func RequestNamesByID(ip_dn1 string, ip_dn2 string, state int) *messages.Names {
	// Inicializa la request de ids
	req_dn1 := messages.Ids{
		List: []int64{}, // Lista de ids a pedir
		Len:  int64(0),  // Longitud de la lista
	}

	req_dn2 := messages.Ids{
		List: []int64{}, // Lista de ids a pedir
		Len:  int64(0),  // Longitud de la lista
	}

	// Abre el archivo DATA.txt
	file, err := os.Open(data)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	// Escanea línea por línea y busca IDs con el estado solicitado
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		if len(parts) >= 3 {
			id, _ := strconv.Atoi(parts[0])
			dn, _ := strconv.Atoi(parts[1])
			estado, _ := strconv.Atoi(parts[2])
			if estado == state && dn == 1 {
				req_dn1.List = append(req_dn1.List, int64(id))
				req_dn1.Len++
			}
			if estado == state && dn == 2 {
				req_dn2.List = append(req_dn2.List, int64(id))
				req_dn2.Len++
			}
		}
	}

	log.Printf("[%s] Requesting: %v names to %s, %v names to %s...", sv_name, req_dn1.Len, DN1_ip_req, req_dn2.Len, DN2_ip_req)

	conn_dn1, err := grpc.Dial(ip_dn1, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn_dn1.Close()
	c1 := messages.NewMessageServiceClient(conn_dn1)

	conn_dn2, err := grpc.Dial(ip_dn2, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn_dn2.Close()
	c2 := messages.NewMessageServiceClient(conn_dn2)

	// Enviamos la request de nombres y los recibimos
	names_dn1, err1 := c1.SendRequestNamesByID(context.Background(), &req_dn1)
	if err1 != nil {
		log.Fatal(err1)
	}
	names_dn2, err2 := c2.SendRequestNamesByID(context.Background(), &req_dn2)
	if err2 != nil {
		log.Fatal(err2)
	}

	names := &messages.Names{
		List: []*messages.Info{},
		Len:  int64(0),
	}

	for i := 0; i < int(names_dn1.Len); i++ {
		names.List = append(names.List, names_dn1.List[i])
		names.Len++
	}

	for i := 0; i < int(names_dn2.Len); i++ {
		names.List = append(names.List, names_dn2.List[i])
		names.Len++
	}

	return names
}

func appendData(filename string, txt string) {
	f, err := os.OpenFile(filename,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(txt + "\n"); err != nil {
		log.Fatal(err)
	}
}

func main() {

	log.Printf("[%s] Starting...", sv_name)

	os.Create(data) // inicia el archivo vacio

	go feed_server(feed_port) // continente -> oms
	go req_server(req_port)   // onu -> oms

	for {

	}
}
