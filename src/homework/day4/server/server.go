package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/cncamp/golang/httpserver/metrics"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"io"
	"math/rand"
	"net"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

func randInt(min int, max int) int {
	rand.Seed(time.Now().UTC().UnixNano())
	return min + rand.Intn(max-min)
}

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

func rootHandler(w http.ResponseWriter, r *http.Request) {
	glog.V(4).Info("entering root handler")
	timer := metrics.NewTimer()
	defer timer.ObserveTotal()
	user := r.URL.Query().Get("user")
	delay := randInt(10, 2000)
	time.Sleep(time.Millisecond * time.Duration(delay))
	if user != "" {
		io.WriteString(w, fmt.Sprintf("hello [%s]\n", user))
	} else {
		io.WriteString(w, "hello [stranger]\n")
	}
	io.WriteString(w, "===================Details of the http request header:============\n")
	for k, v := range r.Header {
		io.WriteString(w, fmt.Sprintf("%s=%s\n", k, v))
	}
	glog.V(4).Infof("Respond in %d ms", delay)
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
	metrics.Register()
	http.HandleFunc("/", notfound)
	http.HandleFunc("/root", rootHandler)
	http.HandleFunc("/healthz", healthz)
	http.HandleFunc("/headers", headers)
	server := &http.Server{
		Addr: ":8080",
	}
	go server.ListenAndServe()
	listenSignal(context.Background(), server)
	defer glog.Flush()
}
