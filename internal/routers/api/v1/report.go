package v1

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
)

type Report struct {
}

func NewReport() Report {
	return Report{}
}

func (f Report) Get(c *gin.Context) {
	if _, err := os.Stat("/dev/attestation/report"); os.IsNotExist(err) {
		fmt.Println("Cannot find `/dev/attestation/report`; are you running under SGX?")
		os.Exit(1)
	}

	myTargetInfo, err := ioutil.ReadFile("/dev/attestation/my_target_info")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("/dev/attestation/target_info", myTargetInfo, 0644)
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile("/dev/attestation/user_report_data", make([]byte, 64), 0644)
	if err != nil {
		panic(err)
	}

	report, err := ioutil.ReadFile("/dev/attestation/report")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Generated SGX report with size = %d and the following fields:\n", len(report))
	fmt.Printf("  ATTRIBUTES.FLAGS: %s  [ Debug bit: %t ]\n", hex.EncodeToString(report[48:56]), report[48]&2 > 0)
	fmt.Printf("  ATTRIBUTES.XFRM:  %s\n", hex.EncodeToString(report[56:64]))
	fmt.Printf("  MRENCLAVE:        %s\n", hex.EncodeToString(report[64:96]))
	fmt.Printf("  MRSIGNER:         %s\n", hex.EncodeToString(report[128:160]))
	fmt.Printf("  ISVPRODID:        %s\n", hex.EncodeToString(report[256:258]))
	fmt.Printf("  ISVSVN:           %s\n", hex.EncodeToString(report[258:260]))
	fmt.Printf("  REPORTDATA:       %s\n", hex.EncodeToString(report[320:352]))
	fmt.Printf("                    %s\n", hex.EncodeToString(report[352:384]))
}
