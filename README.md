# SDLab2
Laboratorio 2 sistemas distribuidos

* Create project -> go mod init SDLab2
* Update .mod and .sum -> go mod tidy

# Go
- Actualizar paquete de dependencias en go
    * go mod tidy

# grpc proto
* Install 
    - go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
    - go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
    - export PATH="$PATH:$(go env GOPATH)/bin/"
- Generar archivos de compilación desde message.proto
    protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative proto/message.proto

#### Aparte
- crear modulo en go: 
    go mod init <modulename>
    go mod tidy     
- Actualizar paquete de dependencias en go
    * go mod tidy
- Arreglar problema del archivo de proto:
    export PATH="$PATH:$(go env GOPATH)/bin/"