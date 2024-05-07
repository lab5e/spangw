package gw

// Parameters holds the main command line parameters for the gateway interface
type Parameters struct {
	CertFile     string `kong:"help='Client Certificate',required,file"`
	KeyFile      string `kong:"help='Client key file',required,file"`
	StateFile    string `kong:"help='State file for gateway',default=''"`
	SpanEndpoint string `kong:"help='Endpoint for the Span service',default='gw.lab5e.com:6674'"`
}
