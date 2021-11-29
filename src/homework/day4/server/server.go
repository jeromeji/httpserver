package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/golang/glog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
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
		glog.V(4).Info("wrong")
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
	glog.V(2).Info(clientIp, strconv.Itoa(http.StatusOK))

}

func notfound(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/" {
		http.Redirect(w, r, "/headers", http.StatusFound)
	}
}

func listenSignal(ctx context.Context, httpSrv *http.Server) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	select {
	case <-sigs:
		timeoutCtx, _ := context.WithTimeout(ctx, 3*time.Second)
		glog.V(2).Info("notify sigs")
		httpSrv.Shutdown(timeoutCtx)
		glog.V(2).Info("http shutdown")
	}
}

func main() {
	glog.V(2).Info("....." + os.Getenv("version") + ".....")
	flag.Parse()
	http.HandleFunc("/", notfound)
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/headers", headers)
	server := &http.Server{
		Addr: ":8080",
	}
	go server.ListenAndServe()
	listenSignal(context.Background(), server)
	defer glog.Flush()
}
