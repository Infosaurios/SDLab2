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

type connDN struct {
	sdn pb.MessageServiceClient
	e   error
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

	dataSendToRebels = []string{}
	DATA             = []string{} //contains all the info in DATA.txt, and store new info <category:id:dataNode>
	finishReadDATA   = false
)

// reciben y envían
// Receive data from Combine, write in DATA.txt & DATA
func (s *server) CombineMsg(ctx context.Context, msg *pb.MessageUploadCombine) (*pb.ConfirmationFromNameNode, error) {
	fmt.Println(msg)
	sdn := selectRandomDataNode()
	//Store new data in DATA.txt
	writeInDataFile(msg.Type_, msg.Id, sdn, msg.Data)
	//Store new data in DATA
	newData := msg.Type_ + ":" + msg.Id + ":" + sdn.name
	DATA = append(DATA, newData)
	return &pb.ConfirmationFromNameNode{ValidMsg: true}, nil
}

// This function receive the category selected by the rebels, and send to them all the info requested
func (s *server) ReceiveCategorySendDataToRebels(ctx context.Context, msg *pb.CategorySelected) (*pb.DataFromOneCategory, error) {
	//Send the category selected by rebels to some dataNode
	for {
		toDataNode(msg.Category)
		if finishReadDATA {
			break
		}
	}
	dataToUpload := dataSendToRebels
	finishReadDATA = false
	return &pb.DataFromOneCategory{IdData: dataToUpload}, nil
}

// downloadDataToArray
func downloadDATA() []string {
	var data []string

	f, err := os.Open("DATA.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		//fmt.Println(scanner.Text())
		data = append(data, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	return data
}

/*
	This function returns an array that contains all <id:dataNodeName>

from the data (<category:id:dataNode>) filtered by category
in order (older to newer)
*/
func filterByCategory(category string) []string {
	data := DATA
	var filtered []string
	var ss []string
	category_ := ""

	for i := range data {
		ss = strings.Split(data[i], ":")
		category_ = ss[0]
		if strings.Contains(category_, category) {
			filtered = append(filtered, ss[1]+":"+ss[2])
		}
	}

	return filtered
}

// Send id to dataNode and receive the data <id:data>
func sendIdToDataNodeReceiveData(id_ string, serviceClient pb.MessageServiceClient, err error) string {
	//res -> Receive all the data <id:data> from the nameNode that correspond
	res, errDisp := serviceClient.ReceiveIdSendDataToNameNode(
		context.Background(),
		&pb.IdSelected{
			Id: id_,
		})
	if errDisp != nil {
		panic("No se puede enviar la id hacia data node" + err.Error())
	}
	//fmt.Println(res)
	//dataSendToRebels = append(dataSendToRebels, res.String())
	return res.IdData
}

func toDataNode(category string) string {
	id_dataNodeName_arr := filterByCategory(category) //[<id:dataNode>]
	fmt.Println("id_dataNodeName_arr", id_dataNodeName_arr)

	for i := range id_dataNodeName_arr {

		fmt.Println("i", i, "id_dataNodeName_arr", id_dataNodeName_arr)

		ss := strings.Split(id_dataNodeName_arr[i], ":")
		id := ss[0]
		dtaNodeName := ss[1]
		dtaNode := dataNode{"", "", ""}

		if strings.Compare(dtaNodeName, "CREATOR") == 0 {
			dtaNode = dataNode{name: "CREATOR", port: portDataNodeCreator, host: hostDataNodeCreator}
		} else if strings.Compare(dtaNodeName, "GRUNT") == 0 {
			dtaNode = dataNode{"GRUNT", portDataNodeGrunt, hostDataNodeGrunt}
		} else if strings.Compare(dtaNodeName, "SYNTH") == 0 {
			dtaNode = dataNode{"SYNTH", portDataNodeSynth, hostDataNodeSynth}
		}
		//Connect with the dataNode and Send it the id
		connData := createConnWithDataNode(dtaNode)
		//Send id to dataNode and receive one string with the format <id:data>
		res := sendIdToDataNodeReceiveData(id, connData.sdn, connData.e)
		//Acumulate the data from each dataNode
		accumulator(res)
	}

	finishReadDATA = true

	return "Change this!!"
}

func selectRandomDataNode() dataNode {
	dn := []dataNode{
		{name: "CREATOR", port: portDataNodeCreator, host: hostDataNodeCreator},
		{name: "GRUNT", port: portDataNodeGrunt, host: hostDataNodeGrunt},
		{name: "SYNTH", port: portDataNodeSynth, host: hostDataNodeSynth},
	}
	// max := 3
	// min := 0
	// random := rand.Intn(max-min) + min
	// return dn[random]

	///Delete this!///
	return dn[0]
}

// func createDataFile() {
// 	f, err := os.Create("DATA.txt")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer f.Close()
// }

func writeInDataFile(tipo_ string, id_ string, dataNode_ dataNode, data_ string) {
	f, err := os.OpenFile("DATA.txt", os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}

	newLine := tipo_ + ":" + id_ + ":" + dataNode_.name

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

	var connData = createConnWithDataNode(dataNode_)
	uploadMsgToDataNode(toUpload_, connData.sdn, connData.e)
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

/*
This function accumulates all the strings <category:id:data> to send to the rebels
accumulator must wait for all data to be loaded before sending it to the rebels
*/

func accumulator(dataFromEachDataNode string) {
	dataSendToRebels = append(dataSendToRebels, dataFromEachDataNode)
	fmt.Println("dataSendToRebels", dataSendToRebels)
	//In this part, accumulator must wait for all data to be loaded before sending it to the rebels
	// how can i achieve that? ...
}

/******************Conexión cola síncrona (proto): send to dataNode******************/
func createConnWithDataNode(dtaNode dataNode) connDN {
	connS, err := grpc.Dial(dtaNode.host+dtaNode.port, grpc.WithInsecure())
	if err != nil {
		panic("No se pudo conectar con el servidor " + dtaNode.name + " " + err.Error())
	}
	serviceDataNode := pb.NewMessageServiceClient(connS)
	var connectData = connDN{sdn: serviceDataNode, e: err}
	return connectData
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

	DATA = downloadDATA()

	var forever chan struct{}
	fmt.Printf(" [*] Waiting for messages. To exit press CTRL+C\n")
	<-forever
}
