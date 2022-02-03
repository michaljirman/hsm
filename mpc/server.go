package main

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	`fmt`
	`io`
	`log`
	"math/big"
	`net`
	"time"

	pb "github.com/michaljirman/hsm/proto"

	"google.golang.org/grpc"
	//"github.com/edgelesssys/ego/enclave"
)

type server struct {
	pb.UnimplementedMPCServer
}

func (s server) Test(srv pb.MPC_TestServer) error {

	log.Println("start new server")
	ctx := srv.Context()

	for {

		// exit if context is done
		// or continue
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		// receive data from stream
		req, err := srv.Recv()
		if err == io.EOF {
			// return will close stream from server side
			log.Println("exit")
			return nil
		}
		if err != nil {
			log.Printf("receive error %v", err)
			continue
		}

		fmt.Printf("received req id=%d\n",req.Id)

		resp := pb.Response{Result: "hello"}
		if err := srv.Send(&resp); err != nil {
			log.Printf("send error %v", err)
		}
		log.Println("send new resp")
	}
}

func main() {
	//certDER, priv := createCertificate()
	//// parse the resulting certificate so we can use it again
	////cert, err := x509.ParseCertificate(certDER)
	////if err != nil {
	////	return
	////}
	//
	//tlsCfg := &tls.Config{
	//	ServerName: "localhost:50005",
	//	//ClientAuth:   tls.RequireAndVerifyClientCert,
	//	Certificates: []tls.Certificate{
	//		{
	//			Certificate: [][]byte{certDER},
	//			PrivateKey:  priv,
	//		},
	//	},
	//	//ClientCAs:    certPool,
	//}
	//creds := credentials.NewTLS(tlsCfg)

	// create listener
	lis, err := net.Listen("tcp", ":50005")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	// create grpc server
	//s := grpc.NewServer(grpc.Creds(creds))
	s := grpc.NewServer()
	pb.RegisterMPCServer(s, &server{})

	// and start...
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	//// Create certificate and a report that includes the certificate's hash.
	//cert, priv := createCertificate()
	////hash := sha256.Sum256(cert)
	////report, err := enclave.GetRemoteReport(hash[:])
	////if err != nil {
	////	fmt.Println(err)
	////}
	//
	//// Create HTTPS server.
	//
	//http.HandleFunc("/cert", func(w http.ResponseWriter, r *http.Request) { w.Write(cert) })
	////http.HandleFunc("/report", func(w http.ResponseWriter, r *http.Request) { w.Write(report) })
	//http.HandleFunc("/secret", func(w http.ResponseWriter, r *http.Request) {
	//	fmt.Printf("%v sent secret %v\n", r.RemoteAddr, r.URL.Query()["s"])
	//})
	//
	//tlsCfg := tls.Config{
	//	Certificates: []tls.Certificate{
	//		{
	//			Certificate: [][]byte{cert},
	//			PrivateKey:  priv,
	//		},
	//	},
	//}
	//
	//server := http.Server{Addr: "0.0.0.0:8080", TLSConfig: &tlsCfg}
	//
	//fmt.Println("listening ...")
	//err := server.ListenAndServeTLS("", "")
	//fmt.Println(err)
}

func createCertificate() ([]byte, crypto.PrivateKey) {
	template := &x509.Certificate{
		SerialNumber: &big.Int{},
		Subject:      pkix.Name{CommonName: "localhost"},
		NotAfter:     time.Now().Add(time.Hour),
		//DNSNames:     []string{"localhost"},
	}
	priv, _ := rsa.GenerateKey(rand.Reader, 2048)
	cert, _ := x509.CreateCertificate(rand.Reader, template, template, &priv.PublicKey, priv)
	return cert, priv
}