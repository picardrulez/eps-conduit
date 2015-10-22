package main

import (
    "log"
    "fmt"
    "net/http"
    "strings"
    "strconv"
    "net/http/httputil"
    "os"
    "flag"
)

//Global variable of the last backend sent.  This is for round-robin load balancing
var lastHost int = 0 

func main() {

    //Handling user flags
    backend := flag.String("b", "", "target ips for backend servers")
    bind := flag.String("bind", "8000", "port to bind to")
    mode := flag.String("mode", "http", "Balancing Mode")
    certFile := flag.String("cert", "", "cert")
    keyFile := flag.String("key", "", "key")
    flag.Parse()

    //Error out if backend is not set
    if *backend == "" {
        fmt.Fprintf(os.Stderr, "no backends chosen, use the -b flag.  ex:\n") 
        fmt.Fprintf(os.Stderr, "eps-conduit -b 10.1.8.1,10.1.8.2")
        os.Exit(1)
    }

    //Throwing backends into an array
    backends := strings.Split(*backend, ",")
    totesHost := len(backends)

    //Function for handling the http requests
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        log.Println("lastHost is " + strconv.Itoa(lastHost))
        proxy := httpBalancer(backends, totesHost)
        proxy.ServeHTTP(w, r)
    })

    //tell the user what info the load balancer is using
    for _, v := range backends {
        log.Println("using " + v + " as a backend") 
    }
    log.Println("listening on port " + *bind)

    //Start the http(s) listener depending on user's selected mode
    if *mode == "http" {
        http.ListenAndServe(":" + *bind, nil)
    } else if *mode == "https" {
        http.ListenAndServeTLS(":" + *bind, *certFile, *keyFile, nil) 
    } else {
        fmt.Fprintf(os.Stderr, "unknown mode or mode not set")
        os.Exit(1)
    }
}

//Function for determining the next backend host and handing that to the user as a reverse proxy
func httpBalancer(backend []string, totesHost int) *httputil.ReverseProxy {
    log.Println("picking a host")
    nextHost := pickHost(totesHost)
    lastHost = nextHost
    log.Println("PickHost done, lastHost is " + strconv.Itoa(lastHost))
    log.Println("serving " + backend[nextHost])
    director := func(req *http.Request) {
        req.URL.Scheme = "http"
        req.URL.Host = backend[nextHost]
    }
    proxy := &httputil.ReverseProxy{Director: director}
    return proxy
}

//Function for determining what the next backend host should be
func pickHost(totesHost int) int {
    totesHost = totesHost - 1
    log.Println("totesHost is " + strconv.Itoa(totesHost))
    log.Println("lastHost is " + strconv.Itoa(lastHost))
    if lastHost == totesHost {
        return 0 
    } else {
        lastHost++
        log.Println("lasthost is " + strconv.Itoa(lastHost))
        return lastHost
    }
} 
