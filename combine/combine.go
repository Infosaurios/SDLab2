package main

import (
	pb "SDLab2/proto"
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strings"

	//"net"
	//"time"

	"google.golang.org/grpc"
)

type InfoToUpload struct {
	type_ string
	id    string
	data  string
}

var (
	portNameNode = ":50051"
	hostNameNode = "localhost" //Host de nameNode
)

func readUserData() string {
	screenDialog := []string{
		"Enter the information to be sent to nameNode (Category:Id:Data):",
	}
	/*data to upload*/
	fmt.Println(screenDialog[0])
	reader := bufio.NewReader(os.Stdin)
	msgToUpload, _ := reader.ReadString('\n') // Leer hasta el separador de salto de línea
	return msgToUpload
}

func checkMsgCorrectStructure(s string) bool {
	match, _ := regexp.MatchString("([^:])+:([^:])+:([^:])+", s)
	if !match {
		fmt.Println("<< The message is not in the correct format >>")
	}
	return match
}

func MsgProcessing(msg string) InfoToUpload {
	s := strings.Split(msg, ":")
	s[0] = strings.TrimLeft(s[0], "\t")
	s[0] = strings.TrimRight(s[0], "\t")

	s[1] = strings.TrimLeft(s[1], "\t")
	s[1] = strings.TrimRight(s[1], "\t")

	s[2] = strings.TrimLeft(s[2], "\t")
	s[2] = strings.TrimRight(s[2], "\t\n")

	toUpload := InfoToUpload{
		type_: s[0],
		id:    s[1],
		data:  s[2],
	}
	return toUpload
}

// envian y reciben
func uploadMsg(toUpload_ InfoToUpload, serviceClient pb.MessageServiceClient, err error) {
	//send info to nameNode
	res, errDisp := serviceClient.CombineMsg(
		context.Background(),
		&pb.MessageUploadCombine{
			Type_: toUpload_.type_,
			Id:    toUpload_.id,
			Data:  toUpload_.data,
		})
	if errDisp != nil {
		panic("No se puede crear el mensaje " + err.Error())
	}
	fmt.Println(res)
}

func main() {
	/** synchronous connection with nameNode **/
	connS, err := grpc.Dial(hostNameNode+portNameNode, grpc.WithInsecure())
	if err != nil {
		panic("No se pudo conectar con el servidor" + err.Error())
	}
	serviceCombine := pb.NewMessageServiceClient(connS)

	/** main loop of the program*/
	for {
		msgToUpload := readUserData()
		if checkMsgCorrectStructure(msgToUpload) {
			uploadMsg(MsgProcessing(msgToUpload), serviceCombine, err)
		}
	}

}
