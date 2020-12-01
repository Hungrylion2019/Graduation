package main

import (
	"bytes"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
)

func main() {
var pubPEMData = []byte(`
-----BEGIN CERTIFICATE-----
MIICQzCCAemgAwIBAgIQCYgQgpuT9FK4QvzMOwQz4zAKBggqhkjOPQQDAjBzMQsw
CQYDVQQGEwJVUzETMBEGA1UECBMKQ2FsaWZvcm5pYTEWMBQGA1UEBxMNU2FuIEZy
YW5jaXNjbzEZMBcGA1UEChMQb3JnMS5leGFtcGxlLmNvbTEcMBoGA1UEAxMTY2Eu
b3JnMS5leGFtcGxlLmNvbTAeFw0yMDEwMjYwMTM2NTlaFw0zMDEwMjQwMTM2NTla
MHMxCzAJBgNVBAYTAlVTMRMwEQYDVQQIEwpDYWxpZm9ybmlhMRYwFAYDVQQHEw1T
YW4gRnJhbmNpc2NvMRkwFwYDVQQKExBvcmcxLmV4YW1wbGUuY29tMRwwGgYDVQQD
ExNjYS5vcmcxLmV4YW1wbGUuY29tMFkwEwYHKoZIzj0CAQYIKoZIzj0DAQcDQgAE
ODtE0xfVNL2tpRt6c+cOnbj2j03DZRLyz+QOATGg9p9sJk75M3+UJxig/5YfoRHu
b20HOlK7+mrWa9xqC2VlqqNfMF0wDgYDVR0PAQH/BAQDAgGmMA8GA1UdJQQIMAYG
BFUdJQAwDwYDVR0TAQH/BAUwAwEB/zApBgNVHQ4EIgQgOnypZftqSJ7LJwlNTHqf
ardJNw4Z4ePbYHbHnp/ZEQcwCgYIKoZIzj0EAwIDSAAwRQIhALDvjNkO5Km2I2jn
RKppK6d7SU4pO3Hqh2U+JyvjrZ7SAiADXKHgqIV0wAuoqFJ0+aXHzOGlUBZMwhmw
Ffj/NwIj6Q==
-----END CERTIFICATE-----
`)
b:=[]byte("-----BEGIN CERTIFICATE-----")
n:=bytes.Index(pubPEMData,b)

block, _ := pem.Decode(pubPEMData[n:])
if block == nil {
log.Fatal("failed to decode PEM block containing public key")
}

pub, err := x509.ParseCertificate(block.Bytes)
if err != nil {
log.Fatal(err)
}
fmt.Println(string(pub.Signature))
fmt.Println(pub.Subject.String())
//fmt.Println(pub.Subject.CommonName)
//name:=pub.Subject.CommonName
//
//memberName := name[(strings.Index(name, "@") + 1):strings.LastIndex(name, ".simple-network.com")]
//fmt.Println(memberName)
}