package main

import (
	`context`
	`io`
	`log`
	`math/rand`
	`time`

	`google.golang.org/grpc`

	pb `github.com/michaljirman/hsm/proto`
)

func main() {
	rand.Seed(time.Now().Unix())

	// dail server
	conn, err := grpc.Dial(":50005", grpc.WithInsecure())
	if err != nil {
		log.Fatalf("can not connect with server %v", err)
	}

	// create stream
	client := pb.NewMPCClient(conn)
	stream, err := client.Test(context.Background())
	if err != nil {
		log.Fatalf("openn stream error %v", err)
	}

	var max int32
	ctx := stream.Context()
	done := make(chan bool)

	// first goroutine sends random increasing numbers to stream
	// and closes it after 10 iterations
	go func() {
		req := pb.Request{Id: 1}
		if err := stream.Send(&req); err != nil {
			log.Fatalf("can not send %v", err)
		}
		log.Printf("req id %d sent", req.Id)
		time.Sleep(time.Millisecond * 200)
		if err := stream.CloseSend(); err != nil {
			log.Println(err)
		}
	}()

	// second goroutine receives data from stream
	// and saves result in max variable
	//
	// if stream is finished it closes done channel
	go func() {
		for {
			resp, err := stream.Recv()
			if err == io.EOF {
				close(done)
				return
			}
			if err != nil {
				log.Fatalf("can not receive %v", err)
			}
			log.Printf("new resp result %q received", resp.Result)
		}
	}()

	// third goroutine closes done channel
	// if context is done
	go func() {
		<-ctx.Done()
		if err := ctx.Err(); err != nil {
			log.Println(err)
		}
		close(done)
	}()

	<-done
	log.Printf("finished with max=%d", max)
}

//func randomBytes() []byte {
//	result := make([]byte, 32)
//	rand.Read(result)
//	return result
//}

//const rsaSize = 2048

//func main() {
//	secretArg := flag.String("s", "hello", "secret value to send")
//	serverAddr := flag.String("a", "localhost:8080", "server address")
//	flag.Parse()
//
//	url := "https://" + *serverAddr
//
//
//	ctx, err := crypto11.Configure(&crypto11.Config{
//		Path:              "/opt/nfast/toolkits/pkcs11/libcknfast.so",
//		TokenSerial:       "6D30-03E0-D947",
//		LoginNotSupported: true,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	id := randomBytes()
//	_, err = ctx.GenerateRSAKeyPair(id, rsaSize)
//	if err != nil {
//		log.Fatal(err)
//	}
//	log.Println("RSA key pair was successfully generated")
//
//
//	// Get server certificate and its report. Skip TLS certificate verification because
//	// the certificate is self-signed and we will verify it using the report instead.
//	tlsConfig := &tls.Config{InsecureSkipVerify: true}
//	certBytes := httpGet(tlsConfig, url+"/cert")
//	//reportBytes := httpGet(tlsConfig, url+"/report")
//
//	//if err := verifyReport(reportBytes, certBytes, signer); err != nil {
//	//	panic(err)
//	//}
//
//	// Create a TLS config that uses the server certificate as root
//	// CA so that future connections to the server can be verified.
//	cert, _ := x509.ParseCertificate(certBytes)
//	tlsConfig = &tls.Config{RootCAs: x509.NewCertPool(), ServerName: "localhost"}
//	tlsConfig.RootCAs.AddCert(cert)
//
//	httpGet(tlsConfig, fmt.Sprintf("%s/secret?s=%s", url, *secretArg))
//	fmt.Println("Sent secret over TLS channel.")
//}

//func verifyReport(reportBytes, certBytes, signer []byte) error {
//	report, err := eclient.VerifyRemoteReport(reportBytes)
//	if err != nil {
//		return err
//	}
//
//	hash := sha256.Sum256(certBytes)
//	if !bytes.Equal(report.Data[:len(hash)], hash[:]) {
//		return errors.New("report data does not match the certificate's hash")
//	}
//
//	// You can either verify the UniqueID or the tuple (SignerID, ProductID, SecurityVersion, Debug).
//
//	if report.SecurityVersion < 2 {
//		return errors.New("invalid security version")
//	}
//	if binary.LittleEndian.Uint16(report.ProductID) != 1234 {
//		return errors.New("invalid product")
//	}
//	if !bytes.Equal(report.SignerID, signer) {
//		return errors.New("invalid signer")
//	}
//
//	// For production, you must also verify that report.Debug == false
//
//	return nil
//}

//func httpGet(tlsConfig *tls.Config, url string) []byte {
//	client := http.Client{Transport: &http.Transport{TLSClientConfig: tlsConfig}}
//	resp, err := client.Get(url)
//	if err != nil {
//		panic(err)
//	}
//	defer resp.Body.Close()
//	if resp.StatusCode != http.StatusOK {
//		panic(resp.Status)
//	}
//	body, err := ioutil.ReadAll(resp.Body)
//	if err != nil {
//		panic(err)
//	}
//	return body
//}