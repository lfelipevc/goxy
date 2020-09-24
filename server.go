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

func handleUserRequest(client_request *http.Request, w http.ResponseWriter) {
	client := &http.Client{}
	log.Println("Proxy URL Request:\t", client_request.URL)
	resp, err := client.Do(client_request)
	if err != nil {
		log.Fatalln(err)
	} else {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			log.Println("Proxy Reponse Headers:\t", resp.Header)
			for key, values := range resp.Header {
				for _, value := range values {
					w.Header().Set(key, value)
				}
				addAllowOriginHeaders(w.Header(), client_request.Host)
			}
			w.Write(body)
		}
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	user_query_url := r.URL.Query().Get("url")
	log.Println("User Request: ", r.URL.RequestURI(), "\t Method: ", r.Method)
	if len(strings.TrimSpace(user_query_url)) != 0 {
		user_request_url, err := url.Parse(user_query_url)
		if err == nil {

			user_request_query := user_request_url.Query()
			for key, queryValues := range r.URL.Query() {
				for _, queryValue := range queryValues {
					if key != "url" {
						user_request_query.Set(key, queryValue)
					}
				}
			}
			user_request_url.RawQuery = user_request_query.Encode()
			client_request, _ := http.NewRequest(r.Method, user_request_url.String(), r.Body)
			log.Println("Request Headers:\t", r.Header)
			for key, values := range r.Header {
				for _, value := range values {
					client_request.Header.Set(key, value)
				}
			}
			removeOrigin(&client_request.Header)
			handleUserRequest(client_request, w)

		} else {
			w.Write([]byte("url format not valid"))
			log.Fatalln(err)
		}

	} else {
		w.Write([]byte("no url parameter"))
	}

}

func main() {
	log.Println("starting proxy ...")
	http.HandleFunc("/", handler)
	log.Fatal(http.ListenAndServe(":8000", nil))

}
