package Engine


import (
	"fmt"
	"net"

)


func FindAvailablePort() (string, error){
	listener , err := net.Listen("tcp" , "127.0.0.1:0")

	if err != nil { 
		return "" , err;
	}

	defer listener.Close()

	addr := listener.Addr().(*net.TCPAddr)

	return fmt.Sprintf("%d", addr.Port) , nil;

}