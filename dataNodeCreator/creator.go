package main

import (
	pb "SDLab2/proto"
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"google.golang.org/grpc"
)

// func failOnError(err error, msg string) {
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// }

// Permite conexión cola síncrona proto
type server struct {
	pb.UnimplementedMessageServiceServer
}

type dataNode struct {
	name string
	port string
}

var (
	portNameNode        = ":50051"
	portDataNodeCreator = ":50053"
)

func (s *server) ToDataNodeMsg(ctx context.Context, msg *pb.MessageUploadToDataNode) (*pb.ConfirmationFromDataNode, error) {
	fmt.Println(msg)
	writeInDataFile(msg.Type_, msg.Id, msg.Data)
	return &pb.ConfirmationFromDataNode{ValidMsg: true}, nil
}

// func createDataFile() {
// 	f, err := os.Create("DATA.txt")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()
// }

func writeInDataFile(tipo_ string, id_ string, data_ string) {
	f, err := os.OpenFile("DATA.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	newLine := tipo_ + ":" + id_ + ":" + data_

	_, err = fmt.Fprintln(f, newLine)
	if err != nil {
		fmt.Println(err)
		return
	}
}

// Receive the msg's id from nameNode and send one string with the format <id:data>
func (s *server) ReceiveIdSendDataToNameNode(ctx context.Context, msg *pb.IdSelected) (*pb.InfoById, error) {
	//Crear funciones:
	//nameNode retorna un string <id:data>.
	// Para esto debe buscar en el archivo DATA.txt la fila con este id
	fmt.Println("id_==msg.Id (ReceiveIdSendDataToNameNode)", msg.Id)
	idData := dataById(msg.Id)
	fmt.Println("idData", idData)
	return &pb.InfoById{IdData: idData}, nil
}

// Search in the file DATA.txt, the row that contains the id. Return <id : data> of that row
func dataById(id string) string {
	idData := ""
	f, err := os.Open("DATA.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		//data = append(data, scanner.Text())
		ss := strings.Split(scanner.Text(), ":")
		id_ := ss[1]
		if strings.Compare(id_, id) == 0 {
			idData = ss[1] + ":" + ss[2] //<id : data>
		}
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return idData
}

func main() {
	/******************Conexión cola síncrona (proto): Listen ******************/
	go func() {
		listener, err := net.Listen("tcp", portDataNodeCreator) //conexion sincrona escucha
		if err != nil {
			panic("La conexion con dataNodeCreator no se pudo crear" + err.Error())
		}
		grpcServer := grpc.NewServer()
		pb.RegisterMessageServiceServer(grpcServer, &server{})
		if err = grpcServer.Serve(listener); err != nil {
			panic("El servidor dataNodeCreator no se pudo iniciar" + err.Error())
		}
	}()
	time.Sleep(1 * time.Second)

	//createDataFile() //delete?

	var forever chan struct{}
	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	<-forever
}
