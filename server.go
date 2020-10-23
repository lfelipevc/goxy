package main

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptrace"
	"net/http/httputil"
	"net/url"
	"strings"
)

func addOptionsAllowOriginHeaders(header http.Header, origin string, allow string) {
	header.Set("Access-Control-Allow-Origin", origin)
	header.Set("Access-Control-Allow-Methods", allow)
	header.Set("Access-Control-Allow-Headers", "*")
	// header.Set("Access-Control-Allow-Credentials", "true")
}

func addAllowOriginHeaders(header http.Header) {
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Methods", "*")
	header.Set("Access-Control-Allow-Headers", "*")

}
func removeOrigin(header *http.Header) {
	header.Del("Origin")
	header.Del("Referer")
}

func removeAcessControl(header *http.Header) {
	header.Del("Access-Control-Request-Headers")
	header.Del("Access-Control-Request-Method")
}

func debug(data []byte, err error) {
	if err == nil {
		fmt.Printf("%s\n", data)
	} else {
		log.Fatalf("%s\n", err)
	}
}

func handleUserRequest(clientRequest *http.Request, w http.ResponseWriter, r *http.Request) {
	client := &http.Client{
		// Timeout: 5 * time.Second,
	}
	// log.Println("Proxy URL Request:\t", clientRequest.URL)
	debug(httputil.DumpRequestOut(clientRequest, true))
	trace := &httptrace.ClientTrace{
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			fmt.Printf("DNS Info: %+v\n", dnsInfo)
		},
		GotConn: func(connInfo httptrace.GotConnInfo) {
			fmt.Printf("Got Conn: %+v\n", connInfo)
		},
		TLSHandshakeDone: func(connState tls.ConnectionState, err error) {
			if err == nil {
				fmt.Printf("TLSHand Done: %+v\n", connState)
			} else {
				fmt.Println("error", err)
			}
		},
		WroteHeaderField: func(key string, value []string) {
			fmt.Println("HeaderField: ", key, value)
		},
		WroteRequest: func(info httptrace.WroteRequestInfo) {
			fmt.Printf("Request info: %+v\n", info)
		},
	}
	clientRequest = clientRequest.WithContext(httptrace.WithClientTrace(clientRequest.Context(), trace))
	if _, err := http.DefaultTransport.RoundTrip(clientRequest); err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(clientRequest)
	if err != nil {
		log.Fatalln(err)
		w.Write([]byte("error fetch url"))
		return
	}
	defer resp.Body.Close()
	if resp != nil {
		debug(httputil.DumpResponse(resp, false))
		resp.Close = true
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			for key, values := range resp.Header {
				for _, value := range values {
					w.Header().Set(key, value)
				}
				if clientRequest.Method == "OPTIONS" {
					addOptionsAllowOriginHeaders(w.Header(), r.Header.Get("Origin"), resp.Header.Get("Allow"))
				} else {
					addAllowOriginHeaders(w.Header())
				}
			}
			w.Write(body)

			// w.Write(resp.Body)

		} else {
			w.Write([]byte("error fetch url"))
		}
	}

}

func handler(w http.ResponseWriter, r *http.Request) {

	userQueryURL := r.URL.Query().Get("url")
	// log.Println("User Request: ", r.URL.String(), "\t Method: ", r.Method, " \tHeaders: ", r.Header)

	if len(strings.TrimSpace(userQueryURL)) != 0 {
		userRequestURL, err := url.Parse(userQueryURL)
		if err == nil {

			userRequestQuery := userRequestURL.Query()
			for key, queryValues := range r.URL.Query() {
				for _, queryValue := range queryValues {
					if key != "url" {
						userRequestQuery.Set(key, queryValue)
					}
				}
			}

			userRequestURL.RawQuery = userRequestQuery.Encode()
			clientRequest, _ := http.NewRequest(r.Method, userRequestURL.String(), r.Body)

			for key, values := range r.Header {
				for _, value := range values {
					// log.Println("header user request :\t" + key)
					clientRequest.Header.Set(key, value)
				}
			}
			if r.Method == "OPTIONS" {
				removeAcessControl(&clientRequest.Header)
			}
			removeOrigin(&clientRequest.Header)
			defer handleUserRequest(clientRequest, w, r)

		} else {
			w.Write([]byte("url format not valid"))
			log.Fatalln(err)
		}

	} else {
		w.Write([]byte("no url parameter"))
	}

}

func main() {
	log.Println("starting proxy ... port: 8000")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))

}
