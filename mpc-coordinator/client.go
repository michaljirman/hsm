package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"sync"
	"time"

	"github.com/google/uuid"
	"google.golang.org/grpc"

	pb "github.com/michaljirman/hsm/proto"
)

func randomBytes() []byte {
	result := make([]byte, 32)
	rand.Read(result)
	return result
}

const rsaSize = 2048

const (
	mpcSignerDefaultName           = "mpc-signer"
	mpcSignerStartPort             = 3000
	mpcSignerExecutableDefaultName = "./signer"
)

type mpcSigner struct {
	id      uuid.UUID
	name    string
	port    int
	process *os.Process

	client pb.MpcSignerClient
}

func (signer *mpcSigner) FullName() string {
	return fmt.Sprintf("%s-%s", signer.name, signer.id.String())
}

func (signer *mpcSigner) Start() error {
	// come out of package b and then go inside package a to run the executable file as
	cmd := exec.Command(mpcSignerExecutableDefaultName,
		fmt.Sprintf("-port=%d", signer.port),
		fmt.Sprintf("-id=%s", signer.id))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to run mpc signer process: %w", err)
	}
	signer.process = cmd.Process

	// dail mpc signer server
	conn, err := grpc.Dial(fmt.Sprintf(":%d", signer.port), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("failed to connect with server %v", err)
	}

	// create stream
	signer.client = pb.NewMpcSignerClient(conn)

	return nil
}

func (signer *mpcSigner) Stop() error {
	if signer.process != nil {
		if err := signer.process.Signal(os.Interrupt); err != nil {
			return fmt.Errorf("failed to stop mpc signer proces: %w", err)
		}
		if _, err := signer.process.Wait(); err != nil {
			return fmt.Errorf("failed to wait for child process to finish: %w", err)
		}
	}
	return nil
}

func (signer *mpcSigner) Test(id int) error {
	fmt.Println("started Test() code")

	stream, err := signer.client.Test(context.Background())
	if err != nil {
		log.Fatalf("openn stream error %v", err)
	}

	var max int32
	streamCtx := stream.Context()
	done := make(chan bool)

	// first goroutine sends random increasing numbers to stream
	// and closes it after 10 iterations
	go func() {
		req := pb.Request{Id: int32(id)}
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
				fmt.Println("closing done 1")
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
		<-streamCtx.Done()
		if err := streamCtx.Err(); err != nil {
			log.Println(err)
		}
		fmt.Println("closing done 2")
		if !IsClosed(done) {
			close(done)
		}
	}()

	<-done
	log.Printf("finished with max=%d", max)
	return nil
}

func IsClosed(ch <-chan bool) bool {
	select {
	case <-ch:
		return true
	default:
	}

	return false
}

//func (signer *mpcSigner) Test() error {
//	stream, err := signer.client.Test(context.Background())
//	if err != nil {
//		log.Fatalf("openn stream error %v", err)
//	}
//
//	var max int32
//	ctx := stream.Context()
//	done := make(chan bool)
//
//	// first goroutine sends random increasing numbers to stream
//	// and closes it after 10 iterations
//	go func() {
//		ctx, err := crypto11.Configure(&crypto11.Config{
//			Path:              "/opt/nfast/toolkits/pkcs11/libcknfast.so",
//			TokenSerial:       "6D30-03E0-D947",
//			LoginNotSupported: true,
//		})
//		if err != nil {
//			log.Fatal(err)
//		}
//
//		id := randomBytes()
//		_, err = ctx.GenerateRSAKeyPair(id, rsaSize)
//		if err != nil {
//			log.Fatal(err)
//		}
//		log.Println("RSA key pair was successfully generated")
//
//		req := pb.Request{Id: 1}
//		if err := stream.Send(&req); err != nil {
//			log.Fatalf("can not send %v", err)
//		}
//		log.Printf("req id %d sent", req.Id)
//		time.Sleep(time.Millisecond * 200)
//		if err := stream.CloseSend(); err != nil {
//			log.Println(err)
//		}
//	}()
//
//	// second goroutine receives data from stream
//	// and saves result in max variable
//	//
//	// if stream is finished it closes done channel
//	go func() {
//		for {
//			resp, err := stream.Recv()
//			if err == io.EOF {
//				close(done)
//				return
//			}
//			if err != nil {
//				log.Fatalf("can not receive %v", err)
//			}
//			log.Printf("new resp result %q received", resp.Result)
//		}
//	}()
//
//	// third goroutine closes done channel
//	// if context is done
//	go func() {
//		<-ctx.Done()
//		if err := ctx.Err(); err != nil {
//			log.Println(err)
//		}
//		close(done)
//	}()
//
//	<-done
//	log.Printf("finished with max=%d", max)
//	return nil
//}

type mpcSignerManager struct {
	nextAvailablePort int

	signersMap map[string]*mpcSigner
	mutex      sync.RWMutex
}

func NewMpcSignerManager() *mpcSignerManager {
	return &mpcSignerManager{
		nextAvailablePort: mpcSignerStartPort,
		signersMap:        map[string]*mpcSigner{},
		mutex:             sync.RWMutex{},
	}
}

func (manager *mpcSignerManager) allocatePort() int {
	port := manager.nextAvailablePort
	manager.nextAvailablePort++
	return port
}

func (manager *mpcSignerManager) next() *mpcSigner {
	manager.mutex.Lock()
	defer manager.mutex.Unlock()
	signer := &mpcSigner{
		id:   uuid.New(),
		port: manager.allocatePort(),
		name: mpcSignerDefaultName,
	}
	manager.signersMap[signer.id.String()] = signer
	return signer
}

var (
	mpcSignerNum *int
)

func init() {
	mpcSignerNum = flag.Int("signers-num", 5, "number of mpc signers")
}

func main() {
	flag.Parse()

	mpcSignerManager := NewMpcSignerManager()

	ctx, cancel := context.WithCancel(context.Background())
	for i := 0; i < *mpcSignerNum; i++ {
		signer := mpcSignerManager.next()
		if err := signer.Start(); err != nil {
			log.Printf("failed to start signer: %+v\n", err)
			continue
		}

		go func(ctx context.Context) {
			for {
				select {
				case <-ctx.Done():
					fmt.Println("server ping is stopping")
					return
				default:
					resp, err := signer.client.Ready(ctx, &pb.ReadyRequest{})
					if err != nil {
						fmt.Printf("failed to check if server is ready: %+v\n", err)
					}
					fmt.Println(fmt.Sprintf("mpc-signer server with ID: %s reports readiness status %t", resp.Id, resp.Status))
					time.Sleep(5 * time.Second)
				}
			}
		}(ctx)

		//if err := signer.Test(i); err != nil {
		//	log.Printf("failed to test signer: %+v\n", err)
		//}
	}

	time.Sleep(20 * time.Second)
	cancel()

	defer func() {
		for _, signer := range mpcSignerManager.signersMap {
			if err := signer.Stop(); err != nil {
				log.Println("failed to stop signer", err)
				return
			}
		}
	}()
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
