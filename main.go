package main

import (
	"flag"
	"io"
	"log"
	"net"
	"net/http"
	"net/http/httputil"
	"os"
	"strings"
)

//Global variable of the last backend sent.  This is for round-robin load balancing
var VERSION string = "v1.0.7"
var lastHost int = 0
var logfile string = "/var/log/eps-conduit"

func main() {

	//set up logging
	var _, err = os.Stat(logfile)
	if os.IsNotExist(err) {
		var file, err = os.Create(logfile)
		checkError(err)
		defer file.Close()
	}
	f, err := os.OpenFile(logfile, os.O_WRONLY|os.O_APPEND, 0644)
	checkError(err)
	defer f.Close()
	log.SetOutput(f)

	//setup database
	startupres := startup()
	if startupres > 0 {
		log.Println("error running startup procedures")
	}

	var config = ReadConfig()

	//Handling user flags
	backend := flag.String("b", config.Backends, "target ips for backend servers")
	bind := flag.String("bind", config.Bind, "port to bind to")
	monitorbind := flag.String("monitorbind", config.MonitorBind, "port for monitor to bind to")
	mode := flag.String("mode", config.Mode, "Balancing Mode")
	certFile := flag.String("cert", config.Certfile, "cert")
	keyFile := flag.String("key", config.Keyfile, "key")
	flag.Parse()

	//Error out if backend is not set
	if *backend == "" {
		log.Printf("no backends chosen, use the -b flag.  ex:\n")
		log.Printf("eps-conduit -b 10.1.8.1,10.1.8.2")
		os.Exit(1)
	}

	//Remove whitespace from backends
	strings.Replace(*backend, " ", "", -1)
	//Throwing backends into an array
	backends := strings.Split(*backend, ",")
	totesHost := len(backends)

	//Function for handling the http requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		proxy := httpBalancer(backends, totesHost)
		proxy.ServeHTTP(w, r)
	})

	//tell the user what info the load balancer is using
	for _, v := range backends {
		log.Println("using " + v + " as a backend")
	}
	log.Println("listening on port " + *bind)
	log.Println("eps-conduit " + VERSION + " has started sucessfully")

	//Start monitor webserver
	monitorMux := http.NewServeMux()
	monitorMux.HandleFunc("/", monitorRoot)
	monitorMux.Handle("/resources/", http.StripPrefix("/resources/", http.FileServer(http.Dir("resources"))))
	go http.ListenAndServe(":"+*monitorbind, monitorMux)

	//Start the http(s) listener depending on user's selected mode
	if *mode == "http" {
		http.ListenAndServe(":"+*bind, nil)
	} else if *mode == "https" {
		http.ListenAndServeTLS(":"+*bind, *certFile, *keyFile, nil)
	} else if *mode == "tcp" {
		ln, err := net.Listen("tcp", ":"+*bind)
		if err != nil {
			log.Fatalf("binding failure:  %s", err)
		}
		for {
			conn, err := ln.Accept()
			if err != nil {
				log.Printf("failed to accept %s", err)
				continue
			}
			go tcpConnectionHandler(conn, backends, totesHost)
		}
	} else {
		log.Printf("unknown mode or mode not set")
		os.Exit(1)
	}
	log.Println("end of main function")
}

//Function for determining the next backend host and handing that to the user as a reverse proxy
func httpBalancer(backend []string, totesHost int) *httputil.ReverseProxy {
	nextHost := pickHost(totesHost, backend)
	director := func(req *http.Request) {
		req.URL.Scheme = "http"
		req.URL.Host = backend[nextHost]
	}
	insertRes := insertRoute(backend[nextHost])
	if insertRes > 0 {
		log.Println("error inserting route to db")
	}
	proxy := &httputil.ReverseProxy{Director: director}
	return proxy
}

func tcpConnectionHandler(us net.Conn, backend []string, totesHost int) {
	nextHost := pickHost(totesHost, backend)
	ds, err := net.Dial("tcp", backend[nextHost])
	if err != nil {
		us.Close()
		log.Printf("failed to dial %s: %s", backend[nextHost], err)
		insertError(backend[nextHost], err)
		return
	}

	insertRes := insertRoute(backend[nextHost])
	if insertRes > 0 {
		log.Println("error inserting route to db")
	}

	go copy(ds, us)
	go copy(us, ds)
}

//Function for determining what the next backend host should be
func pickHost(totesHost int, backend []string) int {
	totesHost = totesHost - 1
	for {
		if lastHost == totesHost {
			lastHost = 0
			host := strings.Split(backend[lastHost], ":")
			ip := host[0]
			port := host[1]
			if checkAvailability(ip, port, true) {
				return lastHost
			} else {
				continue
			}
		} else {
			lastHost++
			host := strings.Split(backend[lastHost], ":")
			ip := host[0]
			port := host[1]
			if checkAvailability(ip, port, true) {
				return lastHost
			} else {
				continue
			}
		}
	}
}

func copy(wc io.WriteCloser, r io.Reader) {
	defer wc.Close()
	io.Copy(wc, r)
}

func checkError(err error) {
	if err != nil {
		log.Println(err.Error)
	}
}
