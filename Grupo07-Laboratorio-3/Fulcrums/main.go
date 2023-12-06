package main

import (
	"bufio"
	"context"
	messages "lab3/proto"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

var serverName string = "default"
var servePort string = ""       // Puerto en el que este broker listenea
var addressFulcrum1 string = "" // Direccion del fulcrum1
var addressFulcrum2 string = "" // ...
var addressFulcrum3 string = "" // ...

// Iniciamos el mapa de reloj de vectores asociados a cada archivo, sector -> vectorclock
var fileVectorClock map[string]*messages.VectorClock = make(map[string]*messages.VectorClock)

type Server struct {
	messages.UnimplementedMessageServiceServer
}

func (s *Server) Command(ctx context.Context, msg *messages.Cmd) (*messages.VectorClock, error) {
	log.Printf("[%s] Received command from Informants: %s %s %s", serverName, msg.Cmd, msg.Sector, msg.Base)

	// Ejecutamos el comando
	executeCommand(msg.Cmd, msg.Sector, msg.Base, msg.NuevaBase, msg.Valor)

	// Agregamos al log
	appendLog(msg)

	// Extraemos el reloj de vector (ya actualizado) del sector asociado
	response := fileVectorClock[msg.Sector]

	log.Printf("[%s] Executed succesfully: Sector[%s] VectorClock[%v %v %v]", serverName, msg.Sector, response.X, response.Y, response.Z)

	// Enviamos la respuesta
	return response, nil
}

func (s *Server) Soldados(ctx context.Context, msg *messages.Info) (*messages.ValueInfo, error) {
	log.Printf("[%s] Received request from Broker: Sector[%s] Base[%s]", serverName, msg.Sector, msg.Base)

	vc, exists := fileVectorClock[msg.Sector]
	if !exists {
		log.Printf("[%s] Such sector does not exists here", serverName)
		return &messages.ValueInfo{Valor: -1}, nil
	}

	info := sectorInfo(msg.Sector)

	cantidad, exists := info[msg.Base]
	if !exists {
		log.Printf("[%s] Such base does not exists in %s here", serverName, msg.Sector)
		return &messages.ValueInfo{Valor: -1}, nil
	}

	response := &messages.ValueInfo{Valor: int32(cantidad), VC: vc}

	return response, nil
}

func (s *Server) RequestMerge(ctx context.Context, msg *messages.AskMerge) (*messages.Ack, error) {
	log.Printf("[%s] Executing requested merge...", serverName)
	executeMerge()
	log.Printf("[%s] Merge completed...", serverName)
	return &messages.Ack{Value: 1}, nil
}

func (s *Server) Merge(ctx context.Context, msg *messages.AskMerge) (*messages.AllInfo, error) {
	log.Printf("[%s] Received merge request from %s, sending logs...", serverName, msg.FulcrumName)

	logAndVectorClocksInfo := logAndVectorClocks()

	return logAndVectorClocksInfo, nil
}

func (s *Server) PropagateChanges(ctx context.Context, msg *messages.CommandsAndVectorClocks) (*messages.Ack, error) {
	log.Printf("[%s] Received changes, merging...", serverName)

	// Vaciamos todos los archivos
	for sector := range fileVectorClock {
		os.Create(serverName + "/" + sector + ".txt")
	}

	// Re-ejecutamos todos los comandos mergeados
	for i := 0; i < int(msg.CmdLen); i++ {
		cmd := msg.CmdList[i]
		executeCommand(cmd.Cmd, cmd.Sector, cmd.Base, cmd.NuevaBase, cmd.Valor)
	}

	// Actualizamos el vector de relojes
	for k, _ := range fileVectorClock {
		delete(fileVectorClock, k)
	}

	for k, v := range msg.SectorVectorClock {
		fileVectorClock[k] = v
	}

	log.Printf("[%s] Merge completed...", serverName)

	// Vaciamos el log
	//os.Create(serverName + "_log.txt")

	return &messages.Ack{}, nil
}

func executeMerge() {

	var addr1 string = ""
	var addr2 string = ""

	switch serverName {
	case "Fulcrum1":
		addr1 = addressFulcrum2
		addr2 = addressFulcrum3
	case "Fulcrum2":
		addr1 = addressFulcrum1
		addr2 = addressFulcrum3
	case "Fulcrum3":
		addr1 = addressFulcrum1
		addr2 = addressFulcrum2
	default:
	}

	conn1, err := grpc.Dial(addr1, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn1.Close()

	conn2, err := grpc.Dial(addr2, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn2.Close()

	c1 := messages.NewMessageServiceClient(conn1)
	c2 := messages.NewMessageServiceClient(conn2)

	request := &messages.AskMerge{FulcrumName: serverName}

	allinfo1, err1 := c1.Merge(context.Background(), request)
	if err1 != nil {
		log.Fatal(err1)
	}

	allinfo2, err2 := c2.Merge(context.Background(), request)
	if err2 != nil {
		log.Fatal(err2)
	}

	// Mapa de sectores a relojes actualizado
	var syncFileVectorClock map[string]*messages.VectorClock = make(map[string]*messages.VectorClock)

	// Primero se insertan los cambios propios
	for sector, vc := range fileVectorClock {
		syncFileVectorClock[sector] = vc
	}

	// Luego, dependiendo del Fulcrum orquestador, se insertan los cambios de los demas
	switch serverName {
	case "Fulcrum1":
		for sector, vc := range allinfo1.SectorVectorClock { // Fulcrum 2
			_, exists := syncFileVectorClock[sector]
			if exists { // Si la key ya existe, actualizo el valor donde corresponda
				syncFileVectorClock[sector].Y = vc.Y
			} else { // Si no existe, inserto los cambios del otro fulcrum
				syncFileVectorClock[sector] = vc
			}
		}

		for sector, vc := range allinfo2.SectorVectorClock { // Fulcrum 3
			_, exists := syncFileVectorClock[sector]
			if exists {
				syncFileVectorClock[sector].Z = vc.Z
			} else {
				syncFileVectorClock[sector] = vc
			}
		}
	case "Fulcrum2":
		for sector, vc := range allinfo1.SectorVectorClock { // Fulcrum 1
			_, exists := syncFileVectorClock[sector]
			if exists {
				syncFileVectorClock[sector].X = vc.X
			} else {
				syncFileVectorClock[sector] = vc
			}
		}

		for sector, vc := range allinfo2.SectorVectorClock { // Fulcrum 3
			_, exists := syncFileVectorClock[sector]
			if exists {
				syncFileVectorClock[sector].Z = vc.Z
			} else {
				syncFileVectorClock[sector] = vc
			}
		}
	case "Fulcrum3":
		for sector, vc := range allinfo1.SectorVectorClock { // Fulcrum 1
			_, exists := syncFileVectorClock[sector]
			if exists {
				syncFileVectorClock[sector].X = vc.X
			} else {
				syncFileVectorClock[sector] = vc
			}
		}

		for sector, vc := range allinfo2.SectorVectorClock { // Fulcrum 2
			_, exists := syncFileVectorClock[sector]
			if exists {
				syncFileVectorClock[sector].Y = vc.Y
			} else {
				syncFileVectorClock[sector] = vc
			}
		}
	default:
	}

	// Vaciamos todos los archivos
	for sector := range fileVectorClock {
		os.Create(serverName + "/" + sector + ".txt")
	}

	allinfoSelf := logAndVectorClocks()

	changes := &messages.CommandsAndVectorClocks{
		CmdList:           []*messages.Cmd{},
		CmdLen:            int32(0),
		SectorVectorClock: syncFileVectorClock,
	}

	i := 0 // self
	j := 0 // info1
	k := 0 // info2

	LogLenSelf := len(allinfoSelf.LogList) // i
	LogLen1 := len(allinfo1.LogList)       // j
	LogLen2 := len(allinfo2.LogList)       // k

	for i < LogLenSelf || j < LogLen1 || k < LogLen2 {
		if i < LogLenSelf && j < LogLen1 && k < LogLen2 { // Sobreviven los 3
			LogTsSelf := allinfoSelf.LogList[i].TimeStamp
			LogTs1 := allinfo1.LogList[j].TimeStamp
			LogTs2 := allinfo2.LogList[k].TimeStamp
			if LogTsSelf < LogTs1 && LogTsSelf < LogTs2 { // self menor que todos
				changes.CmdList = append(changes.CmdList, allinfoSelf.LogList[i].Cmd)
				i++
			} else if LogTs1 < LogTsSelf && LogTs1 < LogTs2 { // 1 menor que todos
				changes.CmdList = append(changes.CmdList, allinfo1.LogList[j].Cmd)
				j++
			} else { // 2 menor que todos
				changes.CmdList = append(changes.CmdList, allinfo2.LogList[k].Cmd)
				k++
			}
		} else if i < LogLenSelf && j < LogLen1 { // Sobrevive self y 1
			LogTsSelf := allinfoSelf.LogList[i].TimeStamp
			LogTs1 := allinfo1.LogList[j].TimeStamp
			if LogTsSelf < LogTs1 { // self menor que 1
				changes.CmdList = append(changes.CmdList, allinfoSelf.LogList[i].Cmd)
				i++
			} else { // 1 menor que self
				changes.CmdList = append(changes.CmdList, allinfo1.LogList[j].Cmd)
				j++
			}
		} else if i < LogLenSelf && k < LogLen2 { // Sobrevive self y 2
			LogTsSelf := allinfoSelf.LogList[i].TimeStamp
			LogTs2 := allinfo2.LogList[k].TimeStamp
			if LogTsSelf < LogTs2 { // self menor que 2
				changes.CmdList = append(changes.CmdList, allinfoSelf.LogList[i].Cmd)
				i++
			} else { // 2 menor que self
				changes.CmdList = append(changes.CmdList, allinfo2.LogList[k].Cmd)
				k++
			}
		} else if j < LogLen1 && k < LogLen2 { // Sobrevive 1 y 2
			LogTs1 := allinfo1.LogList[j].TimeStamp
			LogTs2 := allinfo2.LogList[k].TimeStamp
			if LogTs1 < LogTs2 { // 1 menor que 2
				changes.CmdList = append(changes.CmdList, allinfo1.LogList[j].Cmd)
				j++
			} else { // 2 menor que 1
				changes.CmdList = append(changes.CmdList, allinfo2.LogList[k].Cmd)
				k++
			}
		} else if i < LogLenSelf { // Sobrevive self
			changes.CmdList = append(changes.CmdList, allinfoSelf.LogList[i].Cmd)
			i++
		} else if j < LogLen1 { // Sobrevive 1
			changes.CmdList = append(changes.CmdList, allinfo1.LogList[j].Cmd)
			j++
		} else if k < LogLen2 { // Sobrevive 2
			changes.CmdList = append(changes.CmdList, allinfo2.LogList[k].Cmd)
			k++
		}

		changes.CmdLen++
	}

	_, err1 = c1.PropagateChanges(context.Background(), changes)
	if err1 != nil {
		log.Fatal(err1)
	}

	_, err2 = c2.PropagateChanges(context.Background(), changes)
	if err2 != nil {
		log.Fatal(err2)
	}

	// Vaciamos el log
	//os.Create(serverName + "_log.txt")

	// Re-ejecutamos todos los comandos mergeados
	for i := 0; i < int(changes.CmdLen); i++ {
		cmd := changes.CmdList[i]
		executeCommand(cmd.Cmd, cmd.Sector, cmd.Base, cmd.NuevaBase, cmd.Valor)
	}

	// Reemplazamos el mapa de reloj de vectores con el actualizado
	for k, _ := range fileVectorClock {
		delete(fileVectorClock, k)
	}

	for k, v := range syncFileVectorClock {
		fileVectorClock[k] = v
	}

}

func executeCommand(cmd string, sector string, base string, nuevaBase string, valor int32) {
	// Extraemos el valor del reloj de vector del sector asociado, y si acaso existe este reloj de vector
	_, exists := fileVectorClock[sector]

	// Si no existe, creamos su archivo vacio
	if !exists {
		os.Create(serverName + "/" + sector + ".txt")
	}

	// Seleccionamos el fulcrum en el que nos ubicamos
	switch serverName {
	case "Fulcrum1":
		if exists { // Si existe el reloj de vector, lo actualizamos (sumando 1)
			fileVectorClock[sector].X++
		} else { // Si no existe, lo creamos con valor 1 donde corresponda (segun fulcrum) y 0 en los demas
			fileVectorClock[sector] = &messages.VectorClock{X: 1, Y: 0, Z: 0}
		}
	case "Fulcrum2":
		if exists {
			fileVectorClock[sector].Y++
		} else {
			fileVectorClock[sector] = &messages.VectorClock{X: 0, Y: 1, Z: 0}
		}
	case "Fulcrum3":
		if exists {
			fileVectorClock[sector].Z++
		} else {
			fileVectorClock[sector] = &messages.VectorClock{X: 0, Y: 0, Z: 1}
		}
	default:
		break
	}

	switch cmd { // Elegimos segun el comando a ejecutar
	case "AgregarBase":
		appendData(serverName+"/"+sector+".txt", sector+" "+base+" "+strconv.Itoa(int(valor)))
	case "RenombrarBase":

		_, err := os.Stat(serverName + "/" + sector + ".txt")

		if os.IsNotExist(err) {
			appendData(serverName+"/"+sector+".txt", sector+" "+base+" "+strconv.Itoa(0))
		} else {
			info := sectorInfo(sector)
			cantidad := info[base]
			delete(info, base)
			info[nuevaBase] = cantidad

			os.Create(serverName + "/" + sector + ".txt")
			for k_base, v_cantidad := range info {
				appendData(serverName+"/"+sector+".txt", sector+" "+k_base+" "+strconv.Itoa(v_cantidad))
			}
		}
	case "ActualizarValor":

		_, err := os.Stat(serverName + "/" + sector + ".txt")

		if os.IsNotExist(err) {
			appendData(serverName+"/"+sector+".txt", sector+" "+base+" "+strconv.Itoa(int(valor)))
		} else {
			info := sectorInfo(sector)
			info[base] = int(valor)

			os.Create(serverName + "/" + sector + ".txt")
			for k_base, v_cantidad := range info {
				appendData(serverName+"/"+sector+".txt", sector+" "+k_base+" "+strconv.Itoa(v_cantidad))
			}
		}
	case "BorrarBase":
		info := sectorInfo(sector)
		delete(info, base)

		os.Create(serverName + "/" + sector + ".txt")
		for k_base, v_cantidad := range info {
			appendData(serverName+"/"+sector+".txt", sector+" "+k_base+" "+strconv.Itoa(v_cantidad))
		}
	default:
		break
	}
}

func sectorInfo(sector string) map[string]int {
	info := make(map[string]int)
	file, err := os.Open(serverName + "/" + sector + ".txt")
	if err != nil {
		log.Print(err)
		return info
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, " ")
		// <Sector>  <Base>  <Cantidad>
		// parts[0] parts[1] parts[2]
		cantidad, _ := strconv.Atoi(parts[2])
		info[parts[1]] = cantidad
	}

	return info
}

func logAndVectorClocks() *messages.AllInfo {
	allInfo := messages.AllInfo{
		LogList:           []*messages.Log{},
		LogLen:            int32(0),
		SectorVectorClock: fileVectorClock,
	}

	file, err := os.Open(serverName + "_log.txt")
	if err != nil {
		log.Fatal(err)
		return &allInfo
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		cmd := messages.Cmd{}

		line := scanner.Text()
		parts := strings.Split(line, " ")

		ts, _ := strconv.Atoi(parts[0]) // Timestamp en milisegundos
		cmd.Cmd = parts[1]              // Comando AgregarBase, RenombrarBase, ActualizarValor, BorrarBase
		cmd.Sector = parts[2]           // Nombre del sector
		cmd.Base = parts[3]             // Nombre de la base

		switch cmd.Cmd {
		case "AgregarBase":
			val, _ := strconv.Atoi(parts[4])
			cmd.Valor = int32(val)
		case "RenombrarBase":
			cmd.NuevaBase = parts[4]
		case "ActualizarValor":
			val, _ := strconv.Atoi(parts[4])
			cmd.Valor = int32(val)
		//case "BorrarBase": // Caso no necesario de ser usado
		default:
			// Hacer nada
		}

		allInfo.LogList = append(allInfo.LogList, &messages.Log{TimeStamp: int64(ts), Cmd: &cmd})
		allInfo.LogLen++
	}

	return &allInfo
}

func serve(port string) {
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

func appendLog(msg *messages.Cmd) {
	// Extraemos el tiempo actual en milisegundos
	ts := time.Now().UnixMilli()

	switch msg.Cmd { // Elegimos segun el comando a ejecutar
	case "AgregarBase":
		appendData(serverName+"_log.txt", strconv.Itoa(int(ts))+" "+msg.Cmd+" "+msg.Sector+" "+msg.Base+" "+strconv.Itoa(int(msg.Valor)))
	case "RenombrarBase":
		appendData(serverName+"_log.txt", strconv.Itoa(int(ts))+" "+msg.Cmd+" "+msg.Sector+" "+msg.Base+" "+msg.NuevaBase)
	case "ActualizarValor":
		appendData(serverName+"_log.txt", strconv.Itoa(int(ts))+" "+msg.Cmd+" "+msg.Sector+" "+msg.Base+" "+strconv.Itoa(int(msg.Valor)))
	case "BorrarBase":
		appendData(serverName+"_log.txt", strconv.Itoa(int(ts))+" "+msg.Cmd+" "+msg.Sector+" "+msg.Base)
	default:
		break
	}
}

func main() {
	// Recogemos los argumentos de la linea de comandos
	serverName = os.Args[1]
	servePort = os.Args[2]
	addressFulcrum1 = os.Args[3]
	addressFulcrum2 = os.Args[4]
	addressFulcrum3 = os.Args[5]

	// Creamos el archivo de log asociado al Fulcrum
	os.Create(serverName + "_log.txt")

	os.Mkdir("Fulcrum1", 0777)
	os.Mkdir("Fulcrum2", 0777)
	os.Mkdir("Fulcrum3", 0777)

	log.Printf("[%s] Initializing on port %s...", serverName, servePort)

	go serve(servePort)

	for {

	}
}
