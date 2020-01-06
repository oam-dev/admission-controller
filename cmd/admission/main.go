package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/oam-dev/oam-go-sdk/apis/core.oam.dev/v1alpha1"
	"k8s.io/client-go/kubernetes/scheme"

	"github.com/oam-dev/admission-controller/common"
	"github.com/oam-dev/admission-controller/pkg/admit"

	"k8s.io/api/admission/v1beta1"
	"k8s.io/klog"
)

// admitFunc is the type we use for all of our validators and mutators
type admitFunc func(v1beta1.AdmissionReview) *v1beta1.AdmissionResponse

// Serve handles the http portion of a request prior to handing to an admit
// function
func Serve(admit admitFunc) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var body []byte
		if r.Body != nil {
			if data, err := ioutil.ReadAll(r.Body); err == nil {
				body = data
			}
		}

		// verify the content type is accurate
		contentType := r.Header.Get("Content-Type")
		if contentType != "application/json" {
			klog.Errorf("contentType=%s, expect application/json", contentType)
			return
		}

		klog.Info(fmt.Sprintf("handling request: %s", body))

		// The AdmissionReview that was sent to the webhook
		requestedAdmissionReview := v1beta1.AdmissionReview{}

		// The AdmissionReview that will be returned
		responseAdmissionReview := v1beta1.AdmissionReview{}

		deserializer := common.Codecs.UniversalDeserializer()
		if _, _, err := deserializer.Decode(body, nil, &requestedAdmissionReview); err != nil {
			klog.Error(err)
			responseAdmissionReview.Response = common.ToErrorResponse(err)
		} else {
			// pass to admitFunc
			responseAdmissionReview.Response = admit(requestedAdmissionReview)
		}

		// Return the same UID
		responseAdmissionReview.Response.UID = requestedAdmissionReview.Request.UID

		klog.V(2).Info(fmt.Sprintf("sending response: %v", responseAdmissionReview.Response))

		respBytes, err := json.Marshal(responseAdmissionReview)
		if err != nil {
			klog.Error(err)
		}
		if _, err := w.Write(respBytes); err != nil {
			klog.Error(err)
		}
	}
}

func main() {
	var healthAddr, addr string
	var config common.Config
	config.AddFlags()
	flag.StringVar(&healthAddr, "health-addr", ":9001", "health probe address")
	flag.StringVar(&addr, "listen-addr", ":443", "listen address")
	klog.InitFlags(nil)
	flag.Parse()

	v1alpha1.AddToScheme(scheme.Scheme)

	adm, err := admit.New()
	if err != nil {
		klog.Errorf("new admit client error %v", err)
		return
	}
	stopCh := make(chan struct{})

	mux := http.NewServeMux()
	mux.HandleFunc("/mutating-appconfig", Serve(adm.MutateAppConfigSpec))

	mux.HandleFunc("/appconfig", Serve(adm.AppConfigSpec))
	mux.HandleFunc("/component", Serve(adm.ComponentSpec))
	mux.HandleFunc("/scope", Serve(adm.ScopeSpec))
	mux.HandleFunc("/traits", Serve(adm.TraitSpec))
	server := http.Server{
		Addr:      addr,
		Handler:   mux,
		TLSConfig: common.ConfigTLS(config),
	}
	go server.ListenAndServeTLS("", "")
	klog.Infof("server starting, listening on %s", addr)
	adm.Start(stopCh)
	go http.ListenAndServe(healthAddr, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("OK"))
	}))
	klog.Info("informer cache completed, serving...")
	WaitForInterrupt(func() {
		server.Shutdown(context.Background())
		close(stopCh)
	})
}

// WaitForInterrupt serves function until a system signal was received
func WaitForInterrupt(interrupt func()) {

	// Set up channel on which to send signal notifications.
	// We must use a buffered channel or risk missing the signal
	// if we're not ready to receive when the signal is sent.
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGTERM, os.Interrupt, syscall.SIGQUIT)

	// Block until a signal is received.
	s := <-c

	klog.Infof("Receiving signal: %v, stopping...", s)

	interrupt()
}
