/*
 * Licensed under the Apache License, Version 2.0 (the “License”);
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *         http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an “AS IS” BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package middlewares

import (
	"errors"
	"net/http"

	"crapi.proj/goservice/api/auth"
	"crapi.proj/goservice/api/responses"
	"github.com/jinzhu/gorm"

	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

//SetMiddlewareJSON set content type and options
func SetMiddlewareJSON(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
		if r.Method == "OPTIONS" {
			return
		}
		next(w, r)
	}
}

//AccessControlMiddleware set content type of header
func AccessControlMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS,PUT")
		w.Header().Set("Access-Control-Allow-Headers", "authorization,content-type")

		if r.Method == "OPTIONS" {
			return
		}

		next.ServeHTTP(w, r)
	})
}

//SetMiddlewareAuthentication checks Authentication token for each request.
func SetMiddlewareAuthentication(next http.HandlerFunc, db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, err := auth.ExtractTokenID(r, db)
		if err != nil {
			responses.ERROR(w, http.StatusUnauthorized, errors.New("Unauthorized"))
			return
		}
		next(w, r)
	}
}


// Log HTTP Requests
var (
	logFile *os.File
	err     error
)

func init() {
	logFilePath := "/home/request_logs.txt"
	logFile, err = os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
}

// LogRequest logs the HTTP request data to a file
func LogRequest(r *http.Request) {
	data := fmt.Sprintf(
		"Time: %s, Method: %s, URL: %s, UserAgent: %s, Cookie: %s, Payload: %s, ContentType: %s, ContentLanguage: %s, Origin: %s, Authorization: %s\n",
		time.Now().Format(time.RFC3339),
		r.Method,
		r.URL.String(),
		r.UserAgent(),
		r.Header.Get("Cookie"),
		getPayload(r),
		r.Header.Get("Content-Type"),
		r.Header.Get("Content-Language"),
		r.Header.Get("Origin"),
		r.Header.Get("Authorization"),
	)

	if _, err := logFile.WriteString(data); err != nil {
		log.Printf("Failed to write log: %v", err)
	}
}

// getPayload reads the request payload
func getPayload(r *http.Request) string {
	if r.Body == nil {
		return ""
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("Failed to read request body: %v", err)
		return ""
	}
	r.Body = ioutil.NopCloser(bytes.NewBuffer(body)) // Reset the body for further reading
	return string(body)
}

// LoggingMiddleware logs all HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		LogRequest(r)
		next.ServeHTTP(w, r)
	})
}
