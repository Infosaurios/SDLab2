package main

import (
	pb "SDLab2/proto"
	"bufio"
	"context"
	"fmt"
	"os/signal"
	"os"
	"strconv"
	"strings"
	"syscall"

	//"net"
	"time"

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
	s=strings.TrimSpace(s)
	match:=false
	if (s=="1" || s=="2" || s=="3"){
		match=true
	}else{
		match=false
	}
	if !match {
		fmt.Println("<< Debe ingresar uno de los sgtes números: {1,2,3} >>")
	}
	return match
}

func categorySelectedByRebels(opt string) string {
	selectedOption := ""
	opt = strings.TrimSpace(opt)
	opt_, err := strconv.Atoi(opt)
	if err != nil {
		fmt.Println("error categorySelectedByRebels", err)
	}

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
	//fmt.Println("catSelected", catSelected, "serviceClient", serviceClient, "err", err)

	res, errDisp := serviceClient.ReceiveCategorySendDataToRebels(
		context.Background(),
		&pb.CategorySelected{
			Category: catSelected,
		})
	if errDisp != nil {
		panic("No se puede crear el mensaje en 'RebelsNameNode'" + err.Error())
	}
	fmt.Println("INFORMATION OF " + catSelected + " CATEGORY:")

	//fmt.Println(res.IdData)
	for i := range res.IdData {
		ss := strings.Split(res.IdData[i], ":")
		fmt.Println(ss[0] + " " + ss[1])
	}
	fmt.Println("")

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
		cancelChan := make(chan os.Signal, 1)
		signal.Notify(cancelChan,os.Interrupt,syscall.SIGTERM)
		
			optSelected := readUserData()
			if checkOptCorrectStructure(optSelected) {
				sendCategoryToNameNodeReceiveData(categorySelectedByRebels(optSelected), serviceCombine, err)
				//categorySelectedByRebels(optSelected)
			}
		
		
		go func() {
			<-cancelChan
			r,err:=serviceCombine.ReqInterruption(
					context.Background(),
					&pb.Interruption{
						Adv:"cierre",
					})
			if err != nil {
						
					}
			fmt.Println("Solicitando cierre de conexion...")
			time.Sleep(1 * time.Second)
			os.Exit(1)
			fmt.Println(r)
			
		}()
	}

}
