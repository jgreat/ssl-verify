package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"github.com/urfave/cli"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var (
	version = "0.0.1"
)

func main() {

	app := cli.NewApp()
	app.Name = "ssl-verify"
	app.Usage = "Verify your ssl certificates/chain is valid."
	app.Action = run
	app.Version = version
	app.HideHelp = true
	app.ArgsUsage = ""
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "ca",
			Usage: "Path to CA cert file.",
		},
		cli.StringFlag{
			Name:  "cert",
			Usage: "Path to Certificate file. May include intermediate certificates.",
		},
		cli.StringFlag{
			Name:  "key",
			Usage: "Path to Private Key file",
		},
		cli.StringFlag{
			Name:  "hostname",
			Usage: "Certificate Common name for verification",
		},
		cli.StringFlag{
			Name:  "port",
			Usage: "Https local port",
			Value: "8443",
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func run(c *cli.Context) error {

	required := []string{
		"cert",
		"key",
		"hostname",
	}

	// cli package doesn't seem have a way to return the Usage for a flag :(
	for _, flag := range required {
		present := c.IsSet(flag)
		if !present {
			log.Fatalf("Missing Required Flag %s. See --help for details.", flag)
		}
	}

	cfg := &tls.Config{}

	if _, err := os.Stat(c.String("ca")); err == nil {
		caCert, err := ioutil.ReadFile(c.String("ca"))
		if err != nil {
			log.Fatal("caCert: ", err)
		}
		caCertPool := x509.NewCertPool()
		caCertPool.AppendCertsFromPEM(caCert)

		cfg = &tls.Config{
			RootCAs: caCertPool,
		}
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%s", c.String("port")),
		Handler: &handler{},
	}

	go httpsServer(srv, c.String("cert"), c.String("key"))

	log.Println("Server Started")

	time.Sleep(2 * time.Second)
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}
	tr := &http.Transport{
		TLSClientConfig: cfg,
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return dialer.DialContext(ctx, network, fmt.Sprintf("127.0.0.1:%s", c.String("port")))
		},
	}
	client := &http.Client{
		Transport: tr,
	}

	resp, err := client.Get(fmt.Sprintf("https://%s", c.String("hostname")))
	if err != nil {
		log.Fatal("Client Error: ", err)
	}

	log.Println(resp.Status)
	log.Println("Certificate Verified")
	return nil
}

type handler struct{}

func httpsServer(srv *http.Server, certFile string, keyFile string) {
	err := srv.ListenAndServeTLS(certFile, keyFile)
	if err != nil {
		log.Fatal("Server Error: ", err)
	}
}

func (h *handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("OK\n"))
}
