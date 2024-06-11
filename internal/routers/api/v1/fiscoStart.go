package v1

import (
	"fisco-sgx-go/global"
	"fmt"
	"log"
	"os"
	"os/exec"

	"github.com/gin-gonic/gin"
)

type Fisco struct {
}

func NewFisco() Fisco {
	return Fisco{}
}

func (f Fisco) Get(c *gin.Context) {

	// //启动fisco_sgx 脚本
	port := global.FiscoSetting.FiscoPort
	start := global.FiscoSetting.FiscoStartPath
	go func() {
		fmt.Printf("port is : %s,\nstart path is %s.\n", port, start)
		cmd := exec.Command(start+"fisco-bcos", "-c", start+"config.ini", "-g", start+"config.genesis")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
		}
	}()

	c.JSON(200, "string(output)")
	// cmd := exec.Command("bash", "-l")
	// output, err := cmd.Output()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(string(output))
	// c.JSON(200, string(output))
}

func (f Fisco) DELETE(c *gin.Context) {

	//停止fisco_sgx脚本
	port := global.FiscoSetting.FiscoPort
	stop := global.FiscoSetting.FiscoStopPath
	fmt.Printf("port is : %s,\nstop path is %s.\n", port, stop)
	cmd := exec.Command(stop)
	output, err := cmd.Output()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(string(output))
	c.JSON(200, string(output))
}
