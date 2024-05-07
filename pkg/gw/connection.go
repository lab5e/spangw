package gw

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"os"

	"github.com/lab5e/spangw/pkg/pb/gateway/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

// Create creates a new gateway process, connects to the Span service and launches the command
// processing. The handler implements the actual gateway
func Create(config Parameters, handler CommandHandler) (*GatewayProcess, error) {
	creds, err := loadCertificates(config.CertFile, config.KeyFile)
	if err != nil {
		return nil, err
	}

	cc, err := grpc.Dial(
		config.SpanEndpoint,
		grpc.WithTransportCredentials(creds))

	if err != nil {
		return nil, err
	}

	demoGW := gateway.NewUserGatewayServiceClient(cc)
	stream, err := demoGW.ControlStream(context.Background())
	if err != nil {
		return nil, err
	}

	return NewGatewayProcess(config.StateFile, stream, handler), nil
}

func loadCertificates(certFile, keyFile string) (credentials.TransportCredentials, error) {
	certs, err := os.ReadFile(certFile)
	if err != nil {
		return nil, err
	}
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(certs)

	cCert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, err
	}
	return credentials.NewTLS(&tls.Config{
		Certificates: []tls.Certificate{cCert},
		GetClientCertificate: func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			return &cCert, nil
		},
		RootCAs: certPool,
	}), nil
}
