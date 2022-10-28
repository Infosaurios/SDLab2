package main

import (
	pb "SDLab2/proto"
	"bufio"
	"context"
	"fmt"
	"os"
	"regexp"
	"strconv"
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

func categorySelectedByRebels(opt string) string {
	selectedOption := ""
	opt = strings.ReplaceAll(opt, "\n", "")

	opt_, err := strconv.Atoi(opt)
	fmt.Println("error categorySelectedByRebels", err)

	if opt_ == 1 {
		selectedOption = "MILITAR"
	} else if opt_ == 2 {
		selectedOption = "FINANCIERA"
	} else if opt_ == 3 {
		selectedOption = "LOGISTICA"
	}
	return selectedOption
}

// Send category to nameNode
func sendCategoryToNameNodeReceiveData(catSelected string, serviceClient pb.MessageServiceClient, err error) {
	//res -> Receive all the data of the category selected from nameNode
	fmt.Println("catSelected", catSelected, "serviceClient", serviceClient, "err", err)

	res, errDisp := serviceClient.ReceiveCategorySendDataToRebels(
		context.Background(),
		&pb.CategorySelected{
			Category: catSelected,
		})
	if errDisp != nil {
		panic("No se puede crear el mensaje en 'RebelsNameNode'" + err.Error())
	}
	fmt.Println("INFORMATION OF " + catSelected + " CATEGORY:")
	//fmt.Println(res.IdData, res)
	fmt.Println(res.IdData)

	//cleanProtoList(serviceClient, err)

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
			sendCategoryToNameNodeReceiveData(categorySelectedByRebels(optSelected), serviceCombine, err)
			//categorySelectedByRebels(optSelected)
		}
	}

}
