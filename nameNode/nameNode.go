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

func failOnError(err error, msg string) {
	if err != nil {
		fmt.Println(err)
	}
}

// Permite conexión cola síncrona proto
type server struct {
	pb.UnimplementedMessageServiceServer
}

type dataNode struct {
	name string
	port string
	host string
}

type InfoToUploadToDataNode struct {
	type_ string
	id    string
	data  string
}

var (
	portNameNode        = ":50051"
	portDataNodeCreator = ":50053"
	portDataNodeGrunt   = ":50055"
	portDataNodeSynth   = ":50057"
	portRebels          = ":50059"
	portCombine         = ":50061"

	hostNameNode        = "localhost"
	hostDataNodeCreator = "localhost"
	hostDataNodeGrunt   = "localhost"
	hostDataNodeSynth   = "localhost"
	hostRebels          = "localhost"
	hostCombine         = "localhost"
)

func (s *server) CombineMsg(ctx context.Context, msg *pb.MessageUploadCombine) (*pb.ConfirmationFromNameNode, error) {
	fmt.Println(msg)
	sdn := selectRandomDataNode()
	writeInDataFile(msg.Type_, msg.Id, sdn, msg.Data)
	return &pb.ConfirmationFromNameNode{ValidMsg: true}, nil
}

func selectRandomDataNode() dataNode {
	dn := []dataNode{
		{name: "creator", port: portDataNodeCreator, host: hostDataNodeCreator},
		{name: "grunt", port: portDataNodeGrunt, host: hostDataNodeGrunt},
		{name: "synth", port: portDataNodeSynth, host: hostDataNodeSynth},
	}
	// max := 3
	// min := 0
	// random := rand.Intn(max-min) + min
	// return dn[random]

	///Delete this!///
	return dn[0]
}

func createDataFile() {
	f, err := os.Create("DATA.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
}

func writeInDataFile(tipo_ string, id_ string, dataNode_ dataNode, data_ string) {
	f, err := os.OpenFile("DATA.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	newLine := tipo_ + " " + id_ + " " + dataNode_.name + "\n"

	_, err = fmt.Fprintln(f, newLine)
	if err != nil {
		fmt.Println(err)
		return
	}

	toUpload_ := InfoToUploadToDataNode{
		type_: tipo_,
		id:    id_,
		data:  data_,
	}

	createConnWithDataNode(dataNode_, toUpload_)
}

func uploadMsgToDataNode(toUpload_ InfoToUploadToDataNode, serviceClient pb.MessageServiceClient, err error) {
	//send info to nameNode
	res, errDisp := serviceClient.ToDataNodeMsg(
		context.Background(),
		&pb.MessageUploadToDataNode{
			Type_: toUpload_.type_,
			Id:    toUpload_.id,
			Data:  toUpload_.data,
		})
	if errDisp != nil {
		panic("No se puede crear el mensaje hacia data node" + err.Error())
	}
	fmt.Println(res)
}

/******************Conexión cola síncrona (proto): send to dataNode******************/
func createConnWithDataNode(dtaNode dataNode, toUploadtoDN_ InfoToUploadToDataNode) {
	connS, err := grpc.Dial(dtaNode.host+dtaNode.port, grpc.WithInsecure())
	if err != nil {
		panic("No se pudo conectar con el servidor " + dtaNode.name + " " + err.Error())
	}
	serviceDataNode := pb.NewMessageServiceClient(connS)

	uploadMsgToDataNode(toUploadtoDN_, serviceDataNode, err)

}

func main() {
	/******************Conexión cola síncrona (proto)******************/
	go func() {
		listener, err := net.Listen("tcp", portNameNode) //conexion sincrona
		if err != nil {
			panic("La conexion con nameNode no se pudo crear" + err.Error())
		}
		grpcServer := grpc.NewServer()
		pb.RegisterMessageServiceServer(grpcServer, &server{})
		if err = grpcServer.Serve(listener); err != nil {
			panic("El servidor nameNode no se pudo iniciar" + err.Error())
		}
	}()
	time.Sleep(1 * time.Second)

	createDataFile()

	var forever chan struct{}
	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	<-forever
}
