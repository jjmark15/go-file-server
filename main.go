package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"flag"
	"fmt"
	"log"
	"math/big"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// version of the application
const version = "1.0.0"

// custom response writer to capture status code for logging
type statusRecorder struct {
	http.ResponseWriter
	statusCode int
}

// writeheader captures the status code
func (r *statusRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

// simple logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// initialize recorder with default 200 status
		recorder := &statusRecorder{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
		}

		next.ServeHTTP(recorder, r)

		// format similar to python's http.server:
		// ::1 - - [08/Feb/2026 21:00:00] "GET / HTTP/1.1" 200 -
		log.Printf("%s - - [%s] \"%s %s %s\" %d %s",
			r.RemoteAddr,
			start.Format("02/Jan/2006 15:04:05"),
			r.Method,
			r.URL.Path,
			r.Proto,
			recorder.statusCode,
			time.Since(start),
		)
	})
}

// generate self-signed certificate in memory
func generateSelfSignedCert() (tls.Certificate, error) {
	// generate private key
	priv, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return tls.Certificate{}, err
	}

	// create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"Go Simple Server"},
		},
		NotBefore: time.Now(),
		NotAfter:  time.Now().Add(24 * time.Hour),

		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// create certificate
	derBytes, err := x509.CreateCertificate(rand.Reader, &template, &template, &priv.PublicKey, priv)
	if err != nil {
		return tls.Certificate{}, err
	}

	// encode certificate to pem
	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})

	// encode private key to pem
	b, err := x509.MarshalECPrivateKey(priv)
	if err != nil {
		return tls.Certificate{}, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: b})

	// create tls certificate
	tlsCert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return tls.Certificate{}, err
	}

	return tlsCert, nil
}

func main() {
	// setup flags
	port := flag.String("p", "8000", "port to serve on")
	host := flag.String("h", "0.0.0.0", "host interface to bind to")
	dir := flag.String("d", ".", "directory to serve")
	useTLS := flag.Bool("tls", false, "enable https with self-signed certificate")
	showVer := flag.Bool("v", false, "show version")

	// custom usage message
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of goserver:\n")
		fmt.Fprintf(os.Stderr, "  A simple http server similar to python -m http.server\n\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if *showVer {
		fmt.Printf("goserver version %s\n", version)
		return
	}

	// resolve absolute path for clarity in logs
	absDir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatalf("failed to resolve directory: %v", err)
	}

	// verify directory exists
	if _, err := os.Stat(absDir); os.IsNotExist(err) {
		log.Fatalf("directory does not exist: %s", absDir)
	}

	// setup file server
	fs := http.FileServer(http.Dir(absDir))
	handler := loggingMiddleware(fs)

	addr := net.JoinHostPort(*host, *port)
	server := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	// start server
	if *useTLS {
		// generate certs on the fly
		cert, err := generateSelfSignedCert()
		if err != nil {
			log.Fatalf("failed to generate self-signed certificate: %v", err)
		}

		server.TLSConfig = &tls.Config{
			Certificates: []tls.Certificate{cert},
		}

		fmt.Printf("Serving HTTPS on %s over %s ...\n", addr, absDir)
		log.Fatal(server.ListenAndServeTLS("", ""))
	} else {
		fmt.Printf("Serving HTTP on %s over %s ...\n", addr, absDir)
		log.Fatal(server.ListenAndServe())
	}
}
