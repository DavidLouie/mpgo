package subsonic

import (
    "bytes"
    "crypto/md5"
    "encoding/hex"
    "encoding/xml"
    "errors"
    "fmt"
    "log"
    "net/http"
    "strconv"
)

const apiVersion = "1.16.1"
const username = "david"
const password = "sesame"
var errorMap = map[int]string{
    0:  "A generic error",
    10: "Required parameter is missing",
    20: "Incompatible Subsonic REST protocol version. Client must upgrade",
    30: "Incompatible Subsonic REST protocol version. Server must upgrade",
    40: "Wrong username or password",
    41: "Token authentication not supported for LDAP users",
    50: "User is not authorized for the given operation",
    70: "The requested data was not found",
}

type subParams struct {
    username string
    token    string
    salt     string
    version  string
    client   string
}

// Initialize server and REST endpoints
func Init() {
    http.HandleFunc("/rest/getMusicFolders", GetMusicFolders)
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// Parse and return URL parameters from the HTTP request
func parseParams(w http.ResponseWriter, r *http.Request) (*subParams, error) {
    q := r.URL.Query()
    if q.Get("u") == "" ||
            q.Get("t") == "" ||
            q.Get("s") == "" ||
            q.Get("v") == "" ||
            q.Get("c") == "" {
        const ec = 10
        sendError(w, ec)
        return nil, errors.New(errorMap[ec])
    }

    params := &subParams{
        username: q.Get("u"),
        token:    q.Get("t"),
        salt:     q.Get("s"),
        version:  q.Get("v"),
        client:   q.Get("c")}
    return params, nil
}

// Returns an error response to the endpoint
func sendError(w http.ResponseWriter, code int) {
    if _, ok := errorMap[code]; !ok {
        log.Fatalf("Illegal error code: %d", code)
    }

    response := &subResp{Status: "failed", Version: apiVersion}
    ec := &errorCode{Code: strconv.Itoa(code), Message: errorMap[code]}
    response.ErrorCode = ec

    encoded, err := xml.MarshalIndent(response, "  ", "    ")
    if err != nil {
        fmt.Println(err)
    }

    w.WriteHeader(http.StatusOK)
    w.Header().Set("Content-Type", "application/xml")
    w.Write(encoded)
}

// Authenticates user based on username, token, and salt
// token must match the MD5 hash of the user's password + salt
func authenticate(w http.ResponseWriter, params *subParams) error {
    computedToken := md5.Sum([]byte(password + params.salt))
    tokenBytes, err := hex.DecodeString(params.token)
    if err != nil {
        panic(err)
    }

    if username != params.username || !bytes.Equal(computedToken[:], tokenBytes) {
        const ec = 40
        sendError(w, ec)
        return errors.New(errorMap[ec])
    }
    return nil
}

