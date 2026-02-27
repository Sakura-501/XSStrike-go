package main

import (
	"flag"
	"fmt"
	"html"
	"log"
	"net/http"
	"strings"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:18080", "listen address")
	flag.Parse()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/reflect/html", reflectHTMLHandler)
	mux.HandleFunc("/reflect/attr", reflectAttrHandler)
	mux.HandleFunc("/reflect/script", reflectScriptHandler)
	mux.HandleFunc("/reflect/sanitized", reflectSanitizedHandler)
	mux.HandleFunc("/dom", domHandler)
	mux.HandleFunc("/waf", wafHandler)

	log.Printf("xsstrike-lab listening on http://%s\n", *addr)
	if err := http.ListenAndServe(*addr, mux); err != nil {
		log.Fatal(err)
	}
}

func indexHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprint(w, `<!doctype html>
<html>
<head><meta charset="utf-8"><title>XSStrike-go Lab</title></head>
<body>
<h1>XSStrike-go Local Benchmark Lab</h1>
<ul>
  <li><a href="/reflect/html?q=test">/reflect/html</a></li>
  <li><a href="/reflect/attr?q=test">/reflect/attr</a></li>
  <li><a href="/reflect/script?q=test">/reflect/script</a></li>
  <li><a href="/reflect/sanitized?q=test">/reflect/sanitized</a></li>
  <li><a href="/dom?q=test">/dom</a></li>
  <li><a href="/waf?q=test">/waf</a></li>
</ul>
</body>
</html>`)
}

func reflectHTMLHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, "<html><body><h2>HTML reflection</h2><div id='echo'>%s</div></body></html>", q)
}

func reflectAttrHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, "<html><body><h2>Attribute reflection</h2><input value=\"%s\"></body></html>", q)
}

func reflectScriptHandler(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, "<html><body><h2>Script reflection</h2><script>var input = \"%s\";</script></body></html>", q)
}

func reflectSanitizedHandler(w http.ResponseWriter, r *http.Request) {
	q := html.EscapeString(r.URL.Query().Get("q"))
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, "<html><body><h2>Sanitized reflection</h2><div>%s</div></body></html>", q)
}

func domHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprint(w, `<!doctype html>
<html>
<body>
<div id="target"></div>
<script>
  const p = new URLSearchParams(window.location.search);
  const q = p.get('q') || '';
  document.querySelector('#target').innerHTML = q;
</script>
</body>
</html>`)
}

func wafHandler(w http.ResponseWriter, r *http.Request) {
	q := strings.ToLower(r.URL.Query().Get("q"))
	if strings.Contains(q, "<script") || strings.Contains(q, "alert(") || strings.Contains(q, "onerror=") {
		w.Header().Set("cf-ray", "7fa4f6a3dcf9a123-SJC")
		w.WriteHeader(http.StatusForbidden)
		_, _ = fmt.Fprint(w, "Attention Required! | Cloudflare")
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	_, _ = fmt.Fprintf(w, "<html><body>WAF pass-through: %s</body></html>", q)
}
