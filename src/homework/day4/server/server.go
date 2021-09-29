package main

import (
	"fmt"
	"github.com/golang/glog"
	"net"
	"net/http"
	"os"
	"strconv"
)

func RemoteIp(req *http.Request) string {
	remoteAddr := req.RemoteAddr
	if ip := req.Header.Get("XRealIP"); ip != "" {
		remoteAddr = ip
	} else if ip = req.Header.Get("XForwardedFor"); ip != "" {
		remoteAddr = ip
	} else {
		remoteAddr, _, _ = net.SplitHostPort(remoteAddr)
	}

	if remoteAddr == "::1" {
		remoteAddr = "127.0.0.1"
	}

	return remoteAddr
}

func healthz(w http.ResponseWriter, r *http.Request) {
	_, err := fmt.Fprintf(w, strconv.Itoa(http.StatusOK))
	if err != nil {
		glog.Error("wrong")
	}

}

func headers(w http.ResponseWriter, r *http.Request) {
	os.Setenv("version", "1.17.1")
	for name, headers := range r.Header {
		for _, h := range headers {
			fmt.Fprintf(w, "%v: %v \n", name, h)
		}
	}
	fmt.Fprintf(w, "version is :%v \n", os.Getenv("version"))
	clientIp := RemoteIp(r)
	glog.Info(clientIp,strconv.Itoa(http.StatusOK))

}

func notfound(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/headers", http.StatusFound)
	}
}

func main() {
    http.HandleFunc("/",notfound)
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/headers", headers)
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		glog.Error("server is not started")
	}

}
