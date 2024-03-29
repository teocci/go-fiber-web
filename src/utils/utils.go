// Package utils
// Created by RTT.
// Author: teocci@yandex.com on 2022-12월-25
package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

func Contains(list []string, k string) bool {
	for _, v := range list {
		if strings.ToLower(v) == strings.ToLower(k) {
			return true
		}
	}

	return false
}

func ContainsString(list []string, k string) bool {
	sk := strings.ToLower(k)
	for _, v := range list {
		sv := strings.ToLower(v)
		if strings.Contains(sk, sv) {
			return true
		}
	}

	return false
}

func MaxInt(a, b int) int {
	if a > b {
		return a
	}

	return b
}

func MinInt(a, b int) int {
	if a < b {
		return a
	}

	return b
}

// Ternary evaluates a condition and returns one of two values based on the condition.
// If the condition is true, it returns the first value (va); otherwise, it returns the second value (vb).
// Useful to avoid an if statement when initializing variables
// The function is generic and works with any type T.
//
// Example usage:
// result := Ternary(true, "Yes", "No") // result will be "Yes"
// result := Ternary(false, 42, 0)     // result will be 0
func Ternary[T any](cond bool, va, vb T) T {
	if cond {
		return va
	}
	return vb
}

// Ptr returns a pointer to the passed value.
//
// Useful when you have a value and need a pointer, e.g.:
//
//	func f() string { return "foo" }
//
//	foo := struct{
//	    Bar *string
//	}{
//	    Bar: Ptr(f()),
//	}
func Ptr[T any](v T) *T {
	return &v
}

// Must takes 2 arguments, the second being an error.
// If err is not nil, Must panics. Else the first argument is returned.
//
// Useful when inputs to some function are provided in the source code,
// and you are sure they are valid (if not, it's OK to panic).
// For example:
//
//	t := Must(time.Parse("2006-01-02", "2022-04-20"))
func Must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

// Value returns the first argument.
// Useful when you want to use the first result of a function call that has more than one return values
// (e.g. in a composite literal or in a condition).
//
// For example:
//
//	func f() (i, j, k int, s string, f float64) { return }
//
//	p := image.Point{
//	    X: Value(f()),
//	}
func Value[T any](v T, _ ...any) T {
	return v
}

// IsOK returns the second argument.
// Useful when you want to use the second result of a function call that has more than one return values
// (e.g. in a composite literal or in a condition).
//
// For example:
//
//	func f() (i, j, k int, s string, f float64) { return }
//
//	p := image.Point{
//	    X: IsOK(f()),
//	}
func IsOK[T any](_ any, ok T, _ ...any) T {
	return ok
}

// Third returns the third argument.
// Useful when you want to use the third result of a function call that has more than one return values
// (e.g. in a composite literal or in a condition).
//
// For example:
//
//	func f() (i, j, k int, s string, f float64) { return }
//
//	p := image.Point{
//	    X: Third(f()),
//	}
func Third[T any](_, _ any, third T, _ ...any) T {
	return third
}

// Coalesce returns the first non-zero value from listed arguments.
// Returns the zero value of the type parameter if no arguments are given or all are the zero value.
// Useful when you want to initialize a variable to the first non-zero value from a list of fallback values.
//
// For example:
//
//	hostVal := Coalesce(hostName, os.Getenv("HOST"), "localhost")
func Coalesce[T comparable](values ...T) (v T) {
	var zero T
	for _, v = range values {
		if v != zero {
			return
		}
	}
	return
}

func FloatToUnix(t float64) int64 {
	sec, dec := math.Modf(t)
	return time.Unix(int64(sec), int64(dec*(1e9))).Unix()
}

func StringToUInt64(s string) (n uint64) {
	n, err := strconv.ParseUint(s, 10, 64)
	if err != nil {
		return 0
	}

	return
}

func EmptyUserPass(username, password string) bool {
	return strings.Trim(username, " ") == "" || strings.Trim(password, " ") == ""
}

func AnyToString(v any) string {
	return fmt.Sprintf("%v", v)
}

func StringToTime(s string) (t time.Time) {
	t, err := time.Parse("2006-01-02", s)
	if err != nil {
		return time.Now()
	}

	return
}

func StringToBool(s string) (b bool) {
	b, err := strconv.ParseBool(s)
	if err != nil {
		return false
	}

	return
}

func StringToFloat(s string) (n float64) {
	n, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0
	}

	return
}

func StringToInt(s string) (n int) {
	n, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return
}
func GetFormattedDate(date time.Time) string {
	return date.Format("2006-01-02")
}

func GetTodayDate() string {
	today := time.Now()
	return GetFormattedDate(today)
}

func GetLastWeekDate() string {
	sevenDaysAgo := time.Now().AddDate(0, 0, -7)
	return GetFormattedDate(sevenDaysAgo)
}

func GetYesterdayDate() string {
	sevenDaysAgo := time.Now().AddDate(0, 0, -1)
	return GetFormattedDate(sevenDaysAgo)
}

// GetWithCookiesAndHeaders sends a GET request with custom cookies added to the request.
func GetWithCookiesAndHeaders(urlStr string, cookies []*http.Cookie, headers map[string]string) (resp *http.Response, err error) {
	// Create a new request with GET method
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return
	}

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Add custom cookies to the request
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// Send the request
	client := &http.Client{}
	resp, err = client.Do(req)

	if err != nil {
		return
	}

	code := resp.StatusCode
	if code != 200 {
		fmt.Printf("Error status code: %d\n", code)
	}

	return
}

// GetWithCookies sends a GET request with custom cookies added to the request.
func GetWithCookies(urlStr string, cookies []*http.Cookie) (resp *http.Response, err error) {
	// Create a new request with GET method
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return
	}

	// Add custom cookies to the request
	for _, cookie := range cookies {
		req.AddCookie(cookie)
	}

	// Send the request
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	return
}

// GetWithHeaders performs an HTTP GET request with custom headers.
func GetWithHeaders(urlStr string, headers map[string]string) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", urlStr, nil)
	if err != nil {
		return
	}

	// Set custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	// Send request
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	return
}

// PostWithCookies sends a POST request with custom cookies added to the request.
func PostWithCookies(urlStr string, cookies map[string]string, formData url.Values) (resp *http.Response, err error) {
	// Encode the form data
	data := formData.Encode()

	// Create a new request with POST method
	req, err := http.NewRequest("POST", urlStr, strings.NewReader(data))
	if err != nil {
		return
	}

	// Set the Content-Type header for form data
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Add custom cookies to the request
	for key, value := range cookies {
		req.AddCookie(&http.Cookie{Name: key, Value: value})
	}

	// Send the request
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		return
	}

	return
}

func PostRequester(urlStr string, jsonBody interface{}) (*http.Response, error) {
	jsonBytes, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	client := http.Client{}
	return client.Do(req)
}

func PostWithHeaders(urlStr string, jsonBody interface{}, headers map[string]string) (*http.Response, error) {
	jsonBytes, err := json.Marshal(jsonBody)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	client := http.Client{}
	return client.Do(req)
}
