package redbase_plugin

import (
	"context"
	"fmt"
	"time"
	"strings"
	"bufio"
	"net"
	"net/http"
)

type Config struct {
        RedbaseURL  string 	`json:"redbaseurl,omitempty"`
        DefaultURL  string 	`json:"defaulturl,omitempty"`
}

func CreateConfig() *Config {
	return &Config{}
}

type Redbase struct {
	next        http.Handler
	name        string
        redbaseurl  string
        defaulturl  string
}

func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {

	fmt.Println("Redbase plugin -- github.com/vtacquet/redbase-plugin")
	fmt.Println("Redbase daemon -- github.com/vtacquet/redbase")

	if len(config.RedbaseURL) == 0 {
		return nil, fmt.Errorf("Redbase 'redbaseurl' cannot be empty")
	}
	if len(config.DefaultURL) == 0 {
		return nil, fmt.Errorf("Redbase 'defaulturl' cannot be empty")
	}

	fmt.Println("Redbase redbase url [" + strings.ToLower(config.RedbaseURL) + "]")
	fmt.Println("Redbase default url [" + strings.ToLower(config.DefaultURL) + "]")

	return &Redbase{
		next:       next,
		name:       name,
		redbaseurl: strings.ToLower(config.RedbaseURL),
		defaulturl: strings.ToLower(config.DefaultURL),
	}, nil
}

func getFullURL(r *http.Request) string {

	var proto = "https://"
	if r.TLS == nil {
		proto = "http://"
 	}
	var host = r.URL.Host
	if len(r.URL.Host) == 0 {
		host = r.Host
	}
	var answer = proto + host + r.URL.Path
	return strings.ToLower(answer)
}

func (a *Redbase) ServeHTTP(rw http.ResponseWriter, req *http.Request) {

	var request = getFullURL(req)
	var answerurl = a.defaulturl
        dial, dialerr := net.DialTimeout("tcp", a.redbaseurl, 2*time.Second)
	defer dial.Close()
        if dialerr != nil {
		fmt.Println("Redbase daemon not reachable on " + a.redbaseurl)
		fmt.Println("Redbase default redirect: " + request + " -> " + answerurl)
	} else {
		fmt.Fprintf(dial, request+"\n")
                answer, _ := bufio.NewReader(dial).ReadString('\n')
                answertrim := strings.ToLower(strings.TrimSuffix(answer, "\n"))
		if answertrim != "@default" {
			answerurl = answertrim
			fmt.Println("Redbase database redirect: " + request + " -> " + answerurl)
		} else {
			fmt.Println("Redbase default redirect: " + request + " -> " + answerurl)
		}
	}
        http.Redirect(rw, req, answerurl, 302)
}
