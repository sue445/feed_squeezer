package main

import (
	"flag"
	"fmt"
	"github.com/getsentry/sentry-go"
	sentryhttp "github.com/getsentry/sentry-go/http"
	"log"
	"net/http"
	"os"
	"time"
)

var (
	// Version represents app version (injected from ldflags)
	Version string

	// Revision represents app revision (injected from ldflags)
	Revision string

	isPrintVersion bool
)

func printVersion() {
	fmt.Println(getVersion())
}

func getVersion() string {
	return fmt.Sprintf("feed_squeezer %s (revision %s)", Version, Revision)
}

func main() {
	flag.BoolVar(&isPrintVersion, "version", false, "Whether showing version")

	flag.Parse()

	if isPrintVersion {
		printVersion()
		return
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	err := sentry.Init(sentry.ClientOptions{
		IgnoreErrors: []string{
			// NOTE: Suppress errors when https://www.youtube.com/feeds/videos.xml is called too often
			// e.g.
			// <html><head><meta http-equiv="content-type" content="text/html; charset=utf-8"/><title>Sorry...</title><style> body { font-family: verdana, arial, sans-serif; background-color: #fff; color: #000; }</style></head><body><div><table><tr><td><b><font face=sans-serif size=10><font color=#4285f4>G</font><font color=#ea4335>o</font><font color=#fbbc05>o</font><font color=#4285f4>g</font><font color=#34a853>l</font><font color=#ea4335>e</font></font></b></td><td style="text-align: left; vertical-align: bottom; padding-bottom: 15px; width: 50%"><div style="border-bottom: 1px solid #dfdfdf;">Sorry...</div></td></tr></table></div><div style="margin-left: 4em;"><h1>We're sorry...</h1><p>... but your computer or network may be sending automated queries. To protect our users, we can't process your request right now.</p></div><div style="margin-left: 4em;">See <a href="https://support.google.com/websearch/answer/86640">Google Help</a> for more information.<br/><br/></div><div style="text-align: center; border-top: 1px solid #dfdfdf;"><a href="https://www.google.com">Google Home</a></div></body></html>
			"429 Too Many Requests",
		},
	})

	if err != nil {
		log.Fatal(err)
	}

	// Flush buffered events before the program terminates.
	// Set the timeout to the maximum duration the program can afford to wait.
	defer sentry.Flush(2 * time.Second)

	sentryHandler := sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	})

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", sentryHandler.HandleFunc(IndexHandler))
	mux.HandleFunc("GET /api/feed", sentryHandler.HandleFunc(feedHandler))
	mux.HandleFunc("GET /api/version", sentryHandler.HandleFunc(VersionHandler))

	fmt.Printf("feed_squeezer started: port=%s\n", port)

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		panic(err)
	}
}
