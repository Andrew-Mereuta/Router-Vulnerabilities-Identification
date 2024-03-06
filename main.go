package main

import (
	"fmt"
	"log"
	"time"
	"os"

	g "github.com/gosnmp/gosnmp"
)

func main() {
    	// build our own GoSNMP struct, rather than using g.Default
	// params := &g.GoSNMP{
	// 	Target:        "10.200.25.9",
	// 	Port:          161,
	// 	Version:       g.Version3,
	// 	SecurityModel: g.UserSecurityModel,
	// 	MsgFlags:      g.AuthPriv,
	// 	Timeout:       time.Duration(30) * time.Second,
	// 	SecurityParameters: &g.UsmSecurityParameters{UserName: "user",
	// 		AuthenticationProtocol:   g.SHA,
	// 		AuthenticationPassphrase: "password",
	// 		PrivacyProtocol:          g.DES,
	// 		PrivacyPassphrase:        "password",
	// 	},
	// }

	params := &g.GoSNMP{
		Target:        "10.200.25.9",
		Port:          161,
		Version:       g.Version3,
		Timeout:       time.Duration(2) * time.Second,
		Retries:       0,
		Logger:        g.NewLogger(log.New(os.Stdout, "", 0)),
		SecurityModel: g.UserSecurityModel,
		SecurityParameters: &g.UsmSecurityParameters{
			UserName: " ",
		},
	}

	err := params.Connect()
	if err != nil {
		log.Fatalf("Connect() err: %v", err)
	}
	defer params.Conn.Close()

	oids := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0"}
	result, err2 := params.Get(oids) // Get() accepts up to g.MAX_OIDS
	if err2 != nil {
		log.Fatalf("Get() err: %v", err2)
		fmt.Println(result.SecurityParameters.SafeString())
	}

	for i, variable := range result.Variables {
		fmt.Printf("%d: oid: %s ", i, variable.Name)

		// the Value of each variable returned by Get() implements
		// interface{}. You could do a type switch...
		switch variable.Type {
		case g.OctetString:
			fmt.Printf("string: %s\n", string(variable.Value.([]byte)))
		default:
			// ... or often you're just interested in numeric values.
			// ToBigInt() will return the Value as a BigInt, for plugging
			// into your calculations.
			fmt.Printf("number: %d\n", g.ToBigInt(variable.Value))
		}
	}
	
}