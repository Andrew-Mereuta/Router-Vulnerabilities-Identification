package main

import (
    "fmt"
    "log"
    "os"
    "time"

    "github.com/gosnmp/gosnmp"
)

func main() {
    if len(os.Args) != 2 {
        log.Fatalf("Usage: %s <IP Address>", os.Args[0])
    }
    target := os.Args[1] // Get IP address from command line

    // SNMPv3 parameters (Replace these with your credentials)
    params := &gosnmp.GoSNMP{
        Target:    target,
        Port:      161,
        Version:   gosnmp.Version3,
        Timeout:   time.Duration(6) * time.Second,
        Logger:    gosnmp.NewLogger(log.New(os.Stdout, "", 0)),
        Retries:   3,
        MaxOids:   gosnmp.MaxOids,
        SecurityModel: gosnmp.UserSecurityModel,
        MsgFlags:      gosnmp.NoAuthNoPriv, 
        SecurityParameters: &gosnmp.UsmSecurityParameters{
            UserName: " ", 
        },
    }

    err := params.Connect()
    if err != nil {
        log.Fatalf("Error connecting to %s: %v", target, err)
    }
    defer params.Conn.Close()

    // oids := []string{"1.3.6.1.2.1.1.4.0", "1.3.6.1.2.1.1.7.0"}
    oids := []string{"1.3.6.1.2.1.1.1.0", "1.3.6.1.2.1.2.1.0"} // Example OIDs

    result, err := params.Get(oids)
    if err != nil {
        log.Fatalf("Error performing SNMP Get request: %v", err)
    }

    for _, variable := range result.Variables {
        fmt.Printf("%v = %v\n", variable.Name, variable.Value)
        
        // some defensive programming code taken from my teammates, thanks guys. 
        switch variable.Type {
            case gosnmp.OctetString:
                fmt.Printf("string: %s\n", string(variable.Value.([]byte)))
            default:
                fmt.Printf("number: %d\n", gosnmp.ToBigInt(variable.Value))
            }
    }
}
