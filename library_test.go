package go_tls_pkcs11

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"testing"

	"github.com/ThalesIgnite/crypto11"
)

var tokenCfg = &crypto11.Config{
	Path:        "/usr/lib/x86_64-linux-gnu/pkcs11/onepin-opensc-pkcs11.so",
	TokenSerial: "0503014814020212",
	Pin:         "123456",
}

func TestConnect(t *testing.T) {

	ctx, err := crypto11.Configure(tokenCfg)
	if err != nil {
		t.Fatal(err)
	}

	// cert, err := ctx.FindCertificate(nil, []byte("mainkey"), nil)
	// if err != nil {
	// 	t.Fatal(err)
	// } else {
	// 	t.Log(cert.Subject)
	// 	t.Log(hex.EncodeToString(cert.SubjectKeyId))
	// }

	cacert, err := ctx.FindCertificate(nil, []byte("CA"), nil)
	if err != nil {
		t.Fatal(err)
	} else {
		t.Log(cacert.Subject)
		t.Log(hex.EncodeToString(cacert.SubjectKeyId))
	}

	ca := x509.NewCertPool()
	ca.AddCert(cacert)

	certificates, err := ctx.FindAllPairedCertificates()
	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("total certificates: ", len(certificates))

	cert := certificates[0]

	log.Printf("%#v", cert)

	// keyp, err := ctx.FindKeyPair([]byte{1}, nil)
	// // key, err := ctx.FindKey([]byte{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0x42}, nil)
	// if err != nil {
	// 	t.Fatal(err)
	// } else {
	// 	t.Logf("%#+v", keyp)
	// }

	server := &http.Server{
		TLSConfig: &tls.Config{
			Certificates: []tls.Certificate{cert},
			RootCAs:      ca,
		},
		Addr: ":8000",
	}

	server.ListenAndServeTLS("", "")

	ctx.Close()
}
