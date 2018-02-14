<p align="center">
    <img align="center" src="https://user-images.githubusercontent.com/1072598/31840406-1fe37936-b59a-11e7-939a-71d36e584fc9.png">
    <h3 align="center">CertStream-Go</h3>
    <p align="center">See SSL certs as they're issued live.</p>
</p>

**Certstream-go** is a library for interacting with the [certstream network](https://certstream.calidog.io/) to monitor an aggregated feed from a collection of [Certificate Transparency Logs](https://www.certificate-transparency.org/known-logs).

It leverages the excellent [websocket](https://github.com/gorilla/websocket) and [jsonq](https://github.com/jmoiron/jsonq) libraries and supports reconnecting automatically. 

# Installing

```
go get github.com/CaliDog/certstream-go
```

# Usage

`CertStreamEventStream` returns two channels, a certificate stream with a [JsonQuery](https://godoc.org/github.com/jmoiron/jsonq#JsonQuery) structure (which makes working with the JSON returned from CertStream much easier to work with), and an error stream. 

```go
package main

import (
	"github.com/CaliDog/certstream-go"
	logging "github.com/op/go-logging"
)

var log = logging.MustGetLogger("example")

func main() {
	// The false flag specifies that we want heartbeat messages.
	stream, errStream := certstream.CertStreamEventStream(false)
	for {
		select {
			case jq := <-stream:
				messageType, err := jq.String("message_type")

				if err != nil{
					log.Fatal("Error decoding jq string")
				}

				log.Info("Message type -> ", messageType)
				log.Info("recv: ", jq)
      
			case err := <-errStream:
				log.Error(err)
		}
	}
}
```

# Example data structure

The data structure coming from CertStream looks like this:

```
{
    "message_type": "certificate_update",
    "data": {
        "update_type": "X509LogEntry",
        "leaf_cert": {
            "subject": {
                "aggregated": "/CN=rawlivingvibrantenergy.com",
                "C": null,
                "ST": null,
                "L": null,
                "O": null,
                "OU": null,
                "CN": "rawlivingvibrantenergy.com"
            },
            "extensions": {
                "keyUsage": "Digital Signature, Key Encipherment",
                "extendedKeyUsage": "TLS Web Server Authentication, TLS Web Client Authentication",
                "basicConstraints": "CA:FALSE",
                "subjectKeyIdentifier": "A0:F9:ED:C3:6F:60:50:80:A4:20:5E:9F:75:1C:72:5C:B2:65:3E:50",
                "authorityKeyIdentifier": "keyid:A8:4A:6A:63:04:7D:DD:BA:E6:D1:39:B7:A6:45:65:EF:F3:A8:EC:A1\n",
                "authorityInfoAccess": "OCSP - URI:http://ocsp.int-x3.letsencrypt.org\nCA Issuers - URI:http://cert.int-x3.letsencrypt.org/\n",
                "subjectAltName": "DNS:rawlivingvibrantenergy.com",
                "certificatePolicies": "Policy: 2.23.140.1.2.1\nPolicy: 1.3.6.1.4.1.44947.1.1.1\n  CPS: http://cps.letsencrypt.org\n  User Notice:\n    Explicit Text: This Certificate may only be relied upon by Relying Parties and only in accordance with the Certificate Policy found at https://letsencrypt.org/repository/\n"
            },
            "not_before": 1509717735.0,
            "not_after": 1517493735.0,
            "serial_number": "3782aa9e6eb6abf940a126f2369bc3d2b16",
            "fingerprint": "33:0B:EC:7E:F4:27:1C:98:AB:F3:62:60:44:59:31:B0:2A:5D:6F:59",
            "as_der": "MIIFFzCCA/+gAwIBAgISA3gqqebrar+UChJvI2m8PSsWMA0GCSqGSIb3DQEBCwUAMEoxCzAJBgNVBAYTAlVTMRYwFAYDVQQKEw1MZXQncyBFbmNyeXB0MSMwIQYDVQQDExpMZXQncyBFbmNyeXB0IEF1dGhvcml0eSBYMzAeFw0xNzExMDMxNDAyMTVaFw0xODAyMDExNDAyMTVaMCUxIzAhBgNVBAMTGnJhd2xpdmluZ3ZpYnJhbnRlbmVyZ3kuY29tMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAvUKvYVE/02jC2InwNdFoN25Fmup9+irAcK1+pLdhsidBk1Gc+YMTL7iY9K3J9Ca9IH5N/9a3aNCs1lP29qeaE9VPi1ardYgNsx58f6WJzSu6rk0xaeW+M6MIEta/PzL5An8RlAiJ4UMq1vliqdaRmXGr7I3sfV7/i0wjKFGupTNtYhex3WWTYD/VQhpz2K+eMcND2phn/N0OVEFvCbKWB1g6JFlNkcCA//WQxZ7zZ6Qhsg5BsFwYnmR7cpLLLhEQo12Kb4YeblpYLTdCiOfGRxnHviPidxlfPU0yt9okmtFYYF4Cn9OyvUq2dBvyJmqT82wWtUCj7BPEN3XpCpm41QIDAQABo4ICGjCCAhYwDgYDVR0PAQH/BAQDAgWgMB0GA1UdJQQWMBQGCCsGAQUFBwMBBggrBgEFBQcDAjAMBgNVHRMBAf8EAjAAMB0GA1UdDgQWBBSg+e3Db2BQgKQgXp91HHJcsmU+UDAfBgNVHSMEGDAWgBSoSmpjBH3duubRObemRWXv86jsoTBvBggrBgEFBQcBAQRjMGEwLgYIKwYBBQUHMAGGImh0dHA6Ly9vY3NwLmludC14My5sZXRzZW5jcnlwdC5vcmcwLwYIKwYBBQUHMAKGI2h0dHA6Ly9jZXJ0LmludC14My5sZXRzZW5jcnlwdC5vcmcvMCUGA1UdEQQeMByCGnJhd2xpdmluZ3ZpYnJhbnRlbmVyZ3kuY29tMIH+BgNVHSAEgfYwgfMwCAYGZ4EMAQIBMIHmBgsrBgEEAYLfEwEBATCB1jAmBggrBgEFBQcCARYaaHR0cDovL2Nwcy5sZXRzZW5jcnlwdC5vcmcwgasGCCsGAQUFBwICMIGeDIGbVGhpcyBDZXJ0aWZpY2F0ZSBtYXkgb25seSBiZSByZWxpZWQgdXBvbiBieSBSZWx5aW5nIFBhcnRpZXMgYW5kIG9ubHkgaW4gYWNjb3JkYW5jZSB3aXRoIHRoZSBDZXJ0aWZpY2F0ZSBQb2xpY3kgZm91bmQgYXQgaHR0cHM6Ly9sZXRzZW5jcnlwdC5vcmcvcmVwb3NpdG9yeS8wDQYJKoZIhvcNAQELBQADggEBAD6q1wDWkma5OT5k1AQNEeNypIGh+28ACL0TgH9NeBeJgXlsceSKUroiqoZrs7LS7p88R3iuFUuLjDoOhMQ4B+lV7eRyRJSseJf3gLWD04QuNS3D6/C/eqtsQHXxrrMEPiNuuFEkePiHI41SA/BQnn/Ba97OMXRK0vxc28e6Gy0BI2cqMVoRFNEXm33oJQIabgnAp14GvvRuVxRMAhqoVjHectVyrwgw9rpHCafz0eXlvrPhjoGRZ0xlnqlGI3w7hTXraHgnnDpLFNnsPN+N1qvOt0Wybyv+krldGoEJfTxhKURx3sGQAvzPJGELLBsIDKnk0Whzkt9XGRU2wEg+y3c=",
            "all_domains": [
                "rawlivingvibrantenergy.com"
            ]
        },
        "chain": [
            {
                "subject": {
                    "aggregated": "/C=US/O=Let's Encrypt/CN=Let's Encrypt Authority X3",
                    "C": "US",
                    "ST": null,
                    "L": null,
                    "O": "Let's Encrypt",
                    "OU": null,
                    "CN": "Let's Encrypt Authority X3"
                },
                "extensions": {
                    "basicConstraints": "CA:TRUE, pathlen:0",
                    "keyUsage": "Digital Signature, Certificate Sign, CRL Sign",
                    "authorityInfoAccess": "OCSP - URI:http://isrg.trustid.ocsp.identrust.com\nCA Issuers - URI:http://apps.identrust.com/roots/dstrootcax3.p7c\n",
                    "authorityKeyIdentifier": "keyid:C4:A7:B1:A4:7B:2C:71:FA:DB:E1:4B:90:75:FF:C4:15:60:85:89:10\n",
                    "certificatePolicies": "Policy: 2.23.140.1.2.1\nPolicy: 1.3.6.1.4.1.44947.1.1.1\n  CPS: http://cps.root-x1.letsencrypt.org\n",
                    "crlDistributionPoints": "\nFull Name:\n  URI:http://crl.identrust.com/DSTROOTCAX3CRL.crl\n",
                    "subjectKeyIdentifier": "A8:4A:6A:63:04:7D:DD:BA:E6:D1:39:B7:A6:45:65:EF:F3:A8:EC:A1"
                },
                "not_before": 1458232846.0,
                "not_after": 1615999246.0,
                "serial_number": "a0141420000015385736a0b85eca708",
                "fingerprint": "E6:A3:B4:5B:06:2D:50:9B:33:82:28:2D:19:6E:FE:97:D5:95:6C:CB",
                "as_der": "MIIEkjCCA3qgAwIBAgIQCgFBQgAAAVOFc2oLheynCDANBgkqhkiG9w0BAQsFADA/MSQwIgYDVQQKExtEaWdpdGFsIFNpZ25hdHVyZSBUcnVzdCBDby4xFzAVBgNVBAMTDkRTVCBSb290IENBIFgzMB4XDTE2MDMxNzE2NDA0NloXDTIxMDMxNzE2NDA0NlowSjELMAkGA1UEBhMCVVMxFjAUBgNVBAoTDUxldCdzIEVuY3J5cHQxIzAhBgNVBAMTGkxldCdzIEVuY3J5cHQgQXV0aG9yaXR5IFgzMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnNMM8FrlLke3cl03g7NoYzDq1zUmGSXhvb418XCSL7e4S0EFq6meNQhY7LEqxGiHC6PjdeTm86dicbp5gWAf15Gan/PQeGdxyGkOlZHP/uaZ6WA8SMx+yk13EiSdRxta67nsHjcAHJyse6cF6s5K671B5TaYucv9bTyWaN8jKkKQDIZ0Z8h/pZq4UmEUEz9l6YKHy9v6Dlb2honzhT+Xhq+w3Brvaw2VFn3EK6BlspkENnWAa6xK8xuQSXgvopZPKiAlKQTGdMDQMc2PMTiVFrqoM7hD8bEfwzB/onkxEz0tNvjj/PIzark5McWvxI0NHWQWM6r6hCm21AvA2H3DkwIDAQABo4IBfTCCAXkwEgYDVR0TAQH/BAgwBgEB/wIBADAOBgNVHQ8BAf8EBAMCAYYwfwYIKwYBBQUHAQEEczBxMDIGCCsGAQUFBzABhiZodHRwOi8vaXNyZy50cnVzdGlkLm9jc3AuaWRlbnRydXN0LmNvbTA7BggrBgEFBQcwAoYvaHR0cDovL2FwcHMuaWRlbnRydXN0LmNvbS9yb290cy9kc3Ryb290Y2F4My5wN2MwHwYDVR0jBBgwFoAUxKexpHsscfrb4UuQdf/EFWCFiRAwVAYDVR0gBE0wSzAIBgZngQwBAgEwPwYLKwYBBAGC3xMBAQEwMDAuBggrBgEFBQcCARYiaHR0cDovL2Nwcy5yb290LXgxLmxldHNlbmNyeXB0Lm9yZzA8BgNVHR8ENTAzMDGgL6AthitodHRwOi8vY3JsLmlkZW50cnVzdC5jb20vRFNUUk9PVENBWDNDUkwuY3JsMB0GA1UdDgQWBBSoSmpjBH3duubRObemRWXv86jsoTANBgkqhkiG9w0BAQsFAAOCAQEA3TPXEfNjWDjdGBX7CVW+dla5cEilaUcne8IkCJLxWh9KEik3JHRRHGJouM2VcGfl96S8TihRzZvoroed6ti6WqEBmtzw3Wodatg+VyOeph4EYpr/1wXKtx8/wApIvJSwtmVi4MFU5aMqrSDE6ea73Mj2tcMyo5jMd6jmeWUHK8so/joWUoHOUgwuX4Po1QYz+3dszkDqMp4fklxBwXRsW10KXzPMTZ+sOPAveyxindmjkW8lGy+QsRlGPfZ+G6Z6h7mjem0Y+iWlkYcV4PIWL1iwBi8saCbGS5jN2p8M+X+Q7UNKEkROb3N6KOqkqm57TH2H3eDJAkSnh6/DNFu0Qg=="
            },
            {
                "subject": {
                    "aggregated": "/O=Digital Signature Trust Co./CN=DST Root CA X3",
                    "C": null,
                    "ST": null,
                    "L": null,
                    "O": "Digital Signature Trust Co.",
                    "OU": null,
                    "CN": "DST Root CA X3"
                },
                "extensions": {
                    "basicConstraints": "CA:TRUE",
                    "keyUsage": "Certificate Sign, CRL Sign",
                    "subjectKeyIdentifier": "C4:A7:B1:A4:7B:2C:71:FA:DB:E1:4B:90:75:FF:C4:15:60:85:89:10"
                },
                "not_before": 970348339.0,
                "not_after": 1633010475.0,
                "serial_number": "44afb080d6a327ba893039862ef8406b",
                "fingerprint": "DA:C9:02:4F:54:D8:F6:DF:94:93:5F:B1:73:26:38:CA:6A:D7:7C:13",
                "as_der": "MIIDSjCCAjKgAwIBAgIQRK+wgNajJ7qJMDmGLvhAazANBgkqhkiG9w0BAQUFADA/MSQwIgYDVQQKExtEaWdpdGFsIFNpZ25hdHVyZSBUcnVzdCBDby4xFzAVBgNVBAMTDkRTVCBSb290IENBIFgzMB4XDTAwMDkzMDIxMTIxOVoXDTIxMDkzMDE0MDExNVowPzEkMCIGA1UEChMbRGlnaXRhbCBTaWduYXR1cmUgVHJ1c3QgQ28uMRcwFQYDVQQDEw5EU1QgUm9vdCBDQSBYMzCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoCggEBAN+v6ZdQCINXtMxiZfaQguzH0yxrMMpb7NnDfcdAwRgUi+DoM3ZJKuM/IUmTrE4Orz5Iy2Xu/NMhD2XSKtkyj4zl93ewEnu1lcCJo6m67XMuegwGMoOifooUMM0RoOEqOLl5CjH9UL2AZd+3UWODyOKIYepLYYHsUmu5ouJLGiifSKOeDNoJjj4XLh7dIN9bxiqKqy69cK3FCxolkHRyxXtqqzTWMIn/5WgTe1QLyNau7Fqckh49ZLOMxt+/yUFw7BZy1SbsOFU5Q9D8/RhcQPGX69Wam40dutolucbY38EVAjqr2m7xPi71XAicPNaDaeQQmxkqtilX4+U9m5/wAl0CAwEAAaNCMEAwDwYDVR0TAQH/BAUwAwEB/zAOBgNVHQ8BAf8EBAMCAQYwHQYDVR0OBBYEFMSnsaR7LHH62+FLkHX/xBVghYkQMA0GCSqGSIb3DQEBBQUAA4IBAQCjGiybFwBcqR7uKGY3Or+Dxz9LwwmglSBd49lZRNI+DT69ikugdB/OEIKcdBodfpga3csTS7MgROSR6cz8faXbauX+5v3gTt23ADq1cEmv8uXrAvHRAosZy5Q6XkjEGB5YGV8eAlrwDPGxrancWYaLbumR9YbK+rlmM6pZW87ipxZzR8srzJmwN0jP41ZL9c8PDHIyh8bwRLtTcm1D9SZImlJnt1ir/md2cXjbDaJWFBM5JDGFoqgCWjBH4d1QB7wCCZAA62RjYJsWvIjJEubSfZGL+T0yjWW06XyxV3bqxbYoOb8VZRzI9neWagqNdwvYkQsEjgfbKbYK7p2CNTUQ"
            }
        ],
        "cert_index": 163947593,
        "seen": 1510095720.5917807,
        "source": {
            "url": "ct.googleapis.com/rocketeer",
            "name": "Google 'Rocketeer' log"
        }
    }
}
```
