package transport

import (
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

// CORSHeadersPolicies set headers for support CORS
func CORSHeadersPolicies(allMethods []string) mux.MiddlewareFunc {
	const (
		HeaderAllowOrigin  = "Access-Control-Allow-Origin"
		HeaderAllowMethods = "Access-Control-Allow-Methods"
		HeaderAllowHeaders = "Access-Control-Allow-Headers"
		AllowedHeaders     = "Origin, Referer, Accept, Accept-Encoding, Accept-Language, x-requested-with, Content-Type, Content-Length, Authorization"
	)

	//var AllowedOrigins = []string{"*", "localhost"}
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			//var allowedOrigin = ""
			var allowedOrigin = req.Header.Get("Origin")
			if allowedOrigin == "" {
				allowedOrigin = req.Header.Get("origin")
			}
			// for _, v := range AllowedOrigins {
			// 	if strings.Contains(ref, v) {
			// 		allowedOrigin = ref
			// 	}
			// }
			for _, v := range allMethods {
				if v == http.MethodOptions {
					w.Header().Add(HeaderAllowOrigin, allowedOrigin)
					if req.Method == http.MethodOptions {
						w.Header().Add(HeaderAllowHeaders, AllowedHeaders)
						w.Header().Add(HeaderAllowMethods, strings.Join(allMethods, ","))
					}
				}
			}
			next.ServeHTTP(w, req)
		})
	}
}
