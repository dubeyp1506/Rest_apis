package middlewares

import "net/http"

func SecurityHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-DNS-Prefetch-Control", "off")    //disable dns prefetch
		w.Header().Set("X-Frame-Options", "DENY")          //block the webpage to display on the iframe of other websites
		w.Header().Set("X-XSS-Protection", "1:mode=block") //instrct browser to block the Xss attack
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("Stirct Transport Security", "max-age=6307200:includeSubeDomains:preload") // only connect via https
		w.Header().Set("Content-Security-Policy", "default-src 'self'")                           // only load resorces from same origin
		w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")                      // only send referrer info to same origin
		w.Header().Set("Permissions-Policy", "microphone=(self);camera=(self)")                   // only allow microphone and camera access to same origin
		w.Header().Set("Feature-Policy", "microphone=(self);camera=(self)")
		w.Header().Set("X-Permitted-Cross-Domain-Policies", "none")
		w.Header().Set("X-Content-Security-Policy", "default-src 'self'")
		next.ServeHTTP(w, r)
	})
}
