package main

import (
	pb "SDLab2/proto"
	"context"
	"fmt"
	"log"
	"net"
	"os"
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

func createDataFile() {
	f, err := os.Create("DATA.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
}

func writeInDataFile(tipo_ string, id_ string, data_ string) {
	f, err := os.OpenFile("DATA.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	newLine := tipo_ + " " + id_ + " " + data_ + "\n"

	_, err = fmt.Fprintln(f, newLine)
	if err != nil {
		fmt.Println(err)
		return
	}
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

	createDataFile()

	var forever chan struct{}
	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	<-forever
}
