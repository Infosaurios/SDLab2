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

var (
	portNameNode = ":50051"
	hostNameNode = "localhost" //Host de nameNode
)

// Permite conexión cola síncrona proto
type server struct {
	pb.UnimplementedMessageServiceServer
}

func readUserData() string {

	screenDialog := []string{
		"Selecciona el tipo de información, ingresando el número de la opción:",
		"1. MILITAR",
		"2. FINANCIERA",
		"3. LOGISTICA",
		"Ingrese un número:",
	}
	for i := 0; i < len(screenDialog); i++ {
		fmt.Println(screenDialog[i])
	}

	reader := bufio.NewReader(os.Stdin)
	msgToUpload, _ := reader.ReadString('\n') // Leer hasta el separador de salto de línea
	return msgToUpload
}

func checkOptCorrectStructure(s string) bool {
	match, _ := regexp.MatchString("^[1-3]", s)
	if !match {
		fmt.Println("<< Debe ingresar uno de los sgtes números: {1,2,3} >>")
	}
	return match
}

func categorySelected(opt string) string {
	selectedOption := ""

	if strings.Compare(opt, "1") == 0 {
		selectedOption = "MILITAR"
	} else if strings.Compare(opt, "2") == 0 {
		selectedOption = "FINANCIERA"
	} else if strings.Compare(opt, "3") == 0 {
		selectedOption = "LOGISTICA"
	}
	return selectedOption
}

// Send category to nameNode
func sendCategoryToNameNodeReceiveData(catSelected string, serviceClient pb.MessageServiceClient, err error) {
	//res -> Receive all the data of the category selected from nameNode
	res, errDisp := serviceClient.ReceiveCategorySendDataToRebels(
		context.Background(),
		&pb.CategorySelected{
			Category: catSelected,
		})
	if errDisp != nil {
		panic("No se puede crear el mensaje en 'RebelsNameNode'" + err.Error())
	}
	fmt.Println(res.IdData, res)
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
		optSelected := readUserData()
		if checkOptCorrectStructure(optSelected) {
			sendCategoryToNameNodeReceiveData(categorySelected(optSelected), serviceCombine, err)
		}
	}

}
