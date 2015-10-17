package jwt

import (
	"bytes"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"log"
)

var (
	ErrKeyMustBePEMEncoded = errors.New("Invalid Key: Key must be PEM encoded PKCS1 or PKCS8 private key")
	ErrNotRSAPrivateKey    = errors.New("Key is not a valid RSA private key")
)

// Parse PEM encoded PKCS1 or PKCS8 private key
func ParseRSAPrivateKeyFromPEM(key []byte) (*rsa.PrivateKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block
	if block, _ = pem.Decode(key); block == nil {
		return nil, ErrKeyMustBePEMEncoded
	}

	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKCS1PrivateKey(block.Bytes); err != nil {
		if parsedKey, err = x509.ParsePKCS8PrivateKey(block.Bytes); err != nil {
			return nil, err
		}
	}

	var pkey *rsa.PrivateKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PrivateKey); !ok {
		return nil, ErrNotRSAPrivateKey
	}

	return pkey, nil
}

// https://developers.google.com/identity/sign-in/web/backend-auth
//  https://www.googleapis.com/oauth2/v1/certs
var MappingToPEM = map[string]string{
	"070e3470adad2a897c648ee3d320f693cde6a1b0": "-----BEGIN CERTIFICATE-----\nMIIDJjCCAg6gAwIBAgIIMY08VrIzxEwwDQYJKoZIhvcNAQEFBQAwNjE0MDIGA1UE\nAxMrZmVkZXJhdGVkLXNpZ25vbi5zeXN0ZW0uZ3NlcnZpY2VhY2NvdW50LmNvbTAe\nFw0xNTEwMDkwNDEzMzRaFw0xNTEwMTAxNzEzMzRaMDYxNDAyBgNVBAMTK2ZlZGVy\nYXRlZC1zaWdub24uc3lzdGVtLmdzZXJ2aWNlYWNjb3VudC5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDVt1f6T91xxQrX5DOcFzPJqjClUbAE2ntO\nLWv2zVJU75ZDEFm/HR1AoK54+Z4xsDJr7EtEwcVZcIBZ8TwemvljuF6e0bGSjw+6\nzLqt3RtjB9WY19jcvsdCActYwBHe2Po0AHgeBjlQA3w5begTsG4TMQFPKXd4hWwE\nfIF0HfDsjoSiEnpADM/rtap/HjRUpRs4g5VMUl19yuIRil11PkeKrW7speK6d4wz\neTvuscNt1Rz75Mk19YSG4ZxbbaizWr0JtwpLgj9j2XdoD/VF+3s69oZaaV90IeQP\nIKJ7iPlJaaAYValxCrwLirjVtMvWEuRIekeMbZDNFjRSQ2uFS42vAgMBAAGjODA2\nMAwGA1UdEwEB/wQCMAAwDgYDVR0PAQH/BAQDAgeAMBYGA1UdJQEB/wQMMAoGCCsG\nAQUFBwMCMA0GCSqGSIb3DQEBBQUAA4IBAQC+VrTMe6i8rkY96dNU9rZOuqZ2Yx17\nGqWFeZFMzMzIlFXL36ZBm4HYhjLthloQ9D2E5X29iu/n6e86PEm8mccaqpxX1C0J\nruZlJhFXEZskv4xvJ8z4padjS/NaJURHWztiqYOVwXW78TtsFHRM6cUmvRXXxQMi\nChWxLY2aVqpo7IdyftsWDrojTw34yg8W8uHRSXS75ciPbacQOw8uWzk1v4toNezN\n8qGYl5TR8Yz/fYMg6oSS83OAHi9wjOuOUaDS2vr7HJGbGFdmr7iqq9qsGQ5pMJS1\nOY4FATPPraWwettGxb9jo2vg7ZxYwvEwZnocgZPfVfawWvMrWo83iWV5\n-----END CERTIFICATE-----\n",
	"3b386b69a9815a1651155a3427c6ce7c84fbfbb8": "-----BEGIN CERTIFICATE-----\nMIIDJjCCAg6gAwIBAgIIH0lHQ9hi80AwDQYJKoZIhvcNAQEFBQAwNjE0MDIGA1UE\nAxMrZmVkZXJhdGVkLXNpZ25vbi5zeXN0ZW0uZ3NlcnZpY2VhY2NvdW50LmNvbTAe\nFw0xNTEwMDcwNDQzMzRaFw0xNTEwMDgxNzQzMzRaMDYxNDAyBgNVBAMTK2ZlZGVy\nYXRlZC1zaWdub24uc3lzdGVtLmdzZXJ2aWNlYWNjb3VudC5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDBjExXqZM9+E7Qm3gn5PY+eYQql1oBv/+x\n4QsPTIaH+KQK7gyj7FQmdb82+tiCmj68yO9RQEs+gspLJSJvUainKToQJFRDvh7q\nXn3K4YOdLq1L9JwDyrjEIurPP0K4gX27gNF5FgefeBB4zBBUS0V1aPUOwePsbV6f\n7A4LM7MMssuN5/xg+S5gM5VLK/3G8yzS2MPkhR/jzTCuHCuqWM4BjNvG1OVjvqlF\ntzKwmu+q3KQWpmddSqAua48sMycNA5pNjpFt3tlR5dlgQUodlQDBmFOqOgHs5Ra7\nWUQ4/v9sjfufWxj/B87EeeSImd78DnsQvRXtMi8EJ23UywADbI0NAgMBAAGjODA2\nMAwGA1UdEwEB/wQCMAAwDgYDVR0PAQH/BAQDAgeAMBYGA1UdJQEB/wQMMAoGCCsG\nAQUFBwMCMA0GCSqGSIb3DQEBBQUAA4IBAQBE2vRKCM1PfUqVEId4U7sell6gUTO8\nYPqOOZSb7IDki/rtPFDBi83fpO7nYT1KF++fz2d92dRtQotJaJs7WMGwq9AvgbCK\nOwG5kK3avNXKDR2VjZ5CYqJdT2N/mvYEh/1rRC9n35Xhg2c6lyr9gLCBplvOvKYi\nnl6GpM9UhE6qdSNmtj8WMSJdCMCoarIGUkuD0Q5Hw4mLgmnm0B55KCCcH9WfpbM8\ngXlCG0uQY4RcjtGXVT0jbR1++kqpSuzLTVl3Ydm9hURtkE7GsVcZxffBA9QN41D+\nb2p8whvcgrLQEt5gj6DdQx1zW2jUtimkr3xv4y9ndPmeaLFIgfQYO6md\n-----END CERTIFICATE-----\n",
	"9015759ea37707cb6d325cca00e6299231b7f72f": "-----BEGIN CERTIFICATE-----\nMIIDJjCCAg6gAwIBAgIIF6DTc3dfzIQwDQYJKoZIhvcNAQEFBQAwNjE0MDIGA1UE\nAxMrZmVkZXJhdGVkLXNpZ25vbi5zeXN0ZW0uZ3NlcnZpY2VhY2NvdW50LmNvbTAe\nFw0xNTEwMDgwNDI4MzRaFw0xNTEwMDkxNzI4MzRaMDYxNDAyBgNVBAMTK2ZlZGVy\nYXRlZC1zaWdub24uc3lzdGVtLmdzZXJ2aWNlYWNjb3VudC5jb20wggEiMA0GCSqG\nSIb3DQEBAQUAA4IBDwAwggEKAoIBAQDerjHE9ShSDpOToUhwphHqgCCCUaFmqv4H\npanLmHicb/GRvYbzgnTx3fhVERTziiajR05kh0/6cZuzmKaqiDKJV4EMUI3LzbNc\nR3b38q2ZNhUjUR2xBoLnp1qN2+HI5fDrh3DpWjv5h8MAz+w8w94ZzfQlsONd3iCg\n3GUD4XXX0A88aqEj9ioW5CXXiMz644iopu3uWecvYKf2oyf0knR/S+EHz449k9Gd\n2Ao9Bz3kwWhGi1/Dj5Zbn1sbxY9pCVGLfZHyFHER6yZSe1XzFxCnyI3UZIl/kmM8\nIeAIf5fcrsHKzVHnVD2gUIRbtfmIwOmwfkWJ8j8VS+ctg3TWIIzhAgMBAAGjODA2\nMAwGA1UdEwEB/wQCMAAwDgYDVR0PAQH/BAQDAgeAMBYGA1UdJQEB/wQMMAoGCCsG\nAQUFBwMCMA0GCSqGSIb3DQEBBQUAA4IBAQCewKRGFViLC94Dnhc9YUZWE1oatp2F\nw1El7EHMbSR1fTWo7hcZIdiPc1bewHrkW+mMHnKx8zYKiFagRk6sZdoGvBHeu2Oz\nLWtcnzEdj51z8/piLDfkZD3FZaCZlF66NnekGs1vq+2zmJRBGRSCuM5X18/OQKkz\nkyKIM36OQXzlpsMCoep3BRyKUgBUV6zCMOIpsVNzOj2sPIxZGguNYUMk899eDrFe\nEdvL1K92XOJPtteQW2a7yD2tA5ln3wdBbLkiZqpumGa55SobCB0tXCMUlhlcgAXv\n5ChR4JyWZ88gHRXiRpxR+9rYVnSChMw8I8suonfTus/CqLc152FGMory\n-----END CERTIFICATE-----\n",
}

var ErrPEMMappingObsolete = fmt.Errorf("mapping to PEM key is obsolete")

// Parse PEM encoded PKCS1 or PKCS8 public key
func ParseRSAPublicKeyFromPEM(key []byte) (*rsa.PublicKey, error) {
	var err error

	// Parse PEM block
	var block *pem.Block

	if !bytes.Contains(key, []byte("BEGIN CERTIFICATE")) {
		skey := string(key)
		if sval, ok := MappingToPEM[skey]; ok {
			key = []byte(sval)
			log.Printf("\n=====================\nreplaced %v with %v", len(skey), len(MappingToPEM[skey]))
			log.Printf("%v", skey)
		} else {
			return nil, ErrPEMMappingObsolete
		}
	}

	if block, _ = pem.Decode(key); block == nil {
		log.Println("PEM_decoding failed")
		return nil, ErrKeyMustBePEMEncoded
	}

	// Parse the key
	var parsedKey interface{}
	if parsedKey, err = x509.ParsePKIXPublicKey(block.Bytes); err != nil {
		if cert, err := x509.ParseCertificate(block.Bytes); err == nil {
			parsedKey = cert.PublicKey
		} else {
			return nil, err
		}
	}

	var pkey *rsa.PublicKey
	var ok bool
	if pkey, ok = parsedKey.(*rsa.PublicKey); !ok {
		return nil, ErrNotRSAPrivateKey
	}

	return pkey, nil
}
