package main

import (
	"fmt"
	"crypto/tls"
	"time"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"github.com/docker/swarmkit/remotes"
	"github.com/Diddern/gIntercept/pb/api"
	"crypto/x509"
	cfsigner "github.com/cloudflare/cfssl/signer"
	"github.com/opencontainers/go-digest"
	"crypto"
	"sync"
)

type RootCA struct {
	// Certs contains a bundle of self-signed, PEM encoded certificates for the Root CA to be used
	// as the root of trust.
	Certs []byte

	// Intermediates contains a bundle of PEM encoded intermediate CA certificates to append to any
	// issued TLS (leaf) certificates. The first one must have the same public key and subject as the
	// signing root certificate, and the rest must form a chain, each one certifying the one above it,
	// as per RFC5246 section 7.4.2.
	Intermediates []byte

	// Pool is the root pool used to validate TLS certificates
	Pool *x509.CertPool

	// Digest of the serialized bytes of the certificate(s)
	Digest digest.Digest

	// This signer will be nil if the node doesn't have the appropriate key material
	signer *LocalSigner
}

type LocalSigner struct {
	cfsigner.Signer

	// Key will only be used by the original manager to put the private
	// key-material in raft, no signing operations depend on it.
	Key []byte

	// Cert is one PEM encoded Certificate used as the signing CA.  It must correspond to the key.
	Cert []byte

	// just cached parsed values for validation, etc.
	parsedCert   *x509.Certificate
	cryptoSigner crypto.Signer
}

type Conn struct {
	*grpc.ClientConn
	isLocal bool
	remotes remotes.Remotes
	peer    api.Peer
}

type Broker struct {
	mu        sync.Mutex
	remotes   remotes.Remotes
	localConn *grpc.ClientConn
}

func main() {
	cert := getRemoteCA()
	fmt.Print(cert)
}

func getRemoteCA() ([]byte) {

	ctx := context.TODO()
	fmt.Println("Starting gIntercept")
	insecureCreds := credentials.NewTLS(&tls.Config{InsecureSkipVerify: true})

	dialOpts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecureCreds),
		grpc.WithTimeout(5 * time.Second),
		grpc.WithBackoffMaxDelay(5 * time.Second),
	}

	// Connect securely to GCD service
	conn, err := grpc.Dial("127.0.0.1:4242", dialOpts...)
	if err != nil {
		fmt.Errorf("failed to start gRPC connection: %v", err)
	}
	defer conn.Close()

	client := pb.NewCAClient(conn.ClientConn)
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	defer func() {
		conn.Close(err == nil)
	}()
	response, err := client.GetRootCACertificate(ctx, &api.GetRootCACertificateRequest{})
	if err != nil {
		fmt.Errorf("problems getting the respoonse %s", err)
	}

	return response.Certificate
	// NewRootCA will validate that the certificates are otherwise valid and create a RootCA object.
	// Since there is no key, the certificate expiry does not matter and will not be used.
	//return NewRootCA(response.Certificate, nil, nil, (2160 * time.Hour), nil)
}