package v1

import (
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-gonic/gin"
)

type Quote struct {
}

func NewQuote() Quote {
	return Quote{}
}

func (f Quote) Get(c *gin.Context) {
	if _, err := os.Stat("/dev/attestation/quote"); os.IsNotExist(err) {
		fmt.Println("Cannot find `/dev/attestation/quote`; are you running under SGX, with remote attestation enabled?")
		os.Exit(1)
	}

	attestationType, err := ioutil.ReadFile("/dev/attestation/attestation_type")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Detected attestation type: %s\n", string(attestationType))

	err = ioutil.WriteFile("/dev/attestation/user_report_data", make([]byte, 64), 0644)
	if err != nil {
		panic(err)
	}

	quote, err := ioutil.ReadFile("/dev/attestation/quote")
	if err != nil {
		panic(err)
	}

	fmt.Printf("Extracted SGX quote with size = %d and the following fields:\n", len(quote))
	fmt.Printf("  ATTRIBUTES.FLAGS: %s  [ Debug bit: %t ]\n", hex.EncodeToString(quote[96:104]), quote[96]&2 > 0)
	fmt.Printf("  ATTRIBUTES.XFRM:  %s\n", hex.EncodeToString(quote[104:112]))
	fmt.Printf("  MRENCLAVE:        %s\n", hex.EncodeToString(quote[112:144]))
	fmt.Printf("  MRSIGNER:         %s\n", hex.EncodeToString(quote[176:208]))
	fmt.Printf("  ISVPRODID:        %s\n", hex.EncodeToString(quote[304:306]))
	fmt.Printf("  ISVSVN:           %s\n", hex.EncodeToString(quote[306:308]))
	fmt.Printf("  REPORTDATA:       %s\n", hex.EncodeToString(quote[368:400]))
	fmt.Printf("                    %s\n", hex.EncodeToString(quote[400:432]))
}
