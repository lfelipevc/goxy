package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

func addAllowOriginHeaders(header http.Header, origin string) {
	header.Set("Access-Control-Allow-Origin", "*")
	header.Set("Access-Control-Allow-Methods", "*")

}
func removeOrigin(header *http.Header) {
	header.Del("Origin")
	header.Del("Referer")
}

func handleUserRequest(clientRequest *http.Request, w http.ResponseWriter) {
	client := &http.Client{}
	log.Println("Proxy URL Request:\t", clientRequest.URL)
	resp, err := client.Do(clientRequest)
	if err != nil {
		log.Fatalln(err)
	} else {
		// defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			log.Println("Proxy Reponse Headers:\t", resp.Header)
			for key, values := range resp.Header {
				for _, value := range values {
					w.Header().Set(key, value)
				}
				addAllowOriginHeaders(w.Header(), clientRequest.Host)
			}
			// w.Write(resp.Body)
			w.Write(body)
			defer resp.Body.Close()
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {

	userQueryURL := r.URL.Query().Get("url")
	log.Println("User Request: ", r.URL.String(), "\t Method: ", r.Method, " \tHeaders: ", r.Header)

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
					clientRequest.Header.Set(key, value)
				}
			}
			removeOrigin(&clientRequest.Header)
			handleUserRequest(clientRequest, w)

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
