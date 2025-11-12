package cryptz

import (
	"testing"

	"github.com/welllog/golib/testz"
)

const prvKey = `-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQC3idKBBGck67xk
+ZGQdJYvGwiO7SYr+jxEp7NzeoY5LHU8kh6x63CfAVvuq1D+gxy1UInbKPxmk9mo
mh6m3lvxnOMJwzA7qDRjGRj9SlLOyFXHwuSBYtfGStHwtdv/bmD45imXO280b7X/
1yZERBHVvr//WpSEoFfAYs9UVLUc4Qcp2Jgb5h1fKVrJR09+khWresmjAyFH/Oid
s2DtsVGmGWBvzcvhwlUhxmJmctINqdYKj5/EBpwErpUURqVO3VqS13CoFeKjmRdi
ThiM1jWtXwg4IeAHpBe+sKr2q4Pp+Ug/UAnbxSwaD2+8SBCSeO5fI58WXpUChMrT
ENjFERQjAgMBAAECggEAFVu2AeqLwDmolk2OmXNfyMKKF+vT/dkka07irKUM+plI
WVCdvtGprO9BDgCkr8F9NUJWkIjv+lXpRdwNhhipNXAu4YNz2PVlh0Sz6kHcahDZ
HqJ46e+hMYOic7MOE2b6ZOyP2XgCpGT9lnSokNglBS0p9aLKVra8D3jQLL2gx+fa
zHZgDUZhPOnyKtpcZ8Y6+B01wwVtBlp9x8/dzfz1ycIJUfq5QAUr1eJgYbps1s7F
BSZXhNhvZPTMtLNDMReOPx7fQyPKLEM64+Mj+mkYcv6PJ2yxtZRC+TZIdF/S0OwO
JJDvXgWjqV00zcs93in5U6Jes1LSmmZlTUqyXgP3mQKBgQD+ts6tMiw1MVGlMeDP
8N+vqo9eU+2SQc1x+vcPrdPl9ts+0npCPKpMuAlpR2Zl2ZElWtu+vxu6D/kg5J7/
KOWIY/pOxGkKaIh3Wg8d7Wz9/eiXe29wvUDFAYQolH0fFi51OqrmIoBKGKBTXPBg
+89z0M9I//MjtB7mlhtIOV5QGQKBgQC4dwcDeqecxwDgcldeGWpnzf4AtWhM2pkP
rM0DXJs3h9c6gQtk5b8zw19HOoNQuJDMne7345HwOb/9Dmhcm4dcEjGGMAnqeQQz
PLGU8i9nHN/i+8ZaM4NPme+/ROJxFn6bnIitREDzBieXIKeA0NskRJexwMW1CY1C
UhKyEnndmwKBgQCJkapmmKaPxCdYlWvaYzos4m20gJfbWnbfjLBLY5MCrSUU9RDb
HXDNJsjOd6WydKOUDVKJ3yXWhDIFtfS50xjFZVoXmLUyzeqGq7lmbIllVPF+f0hd
F5oXzQ3X7Pr3Az/sSNdsnE21tz9ARv39I4OUBb8uqi5jpjDaUVBC3dk2YQKBgEES
hB/fEefFb/K9g0KHtriduz/mvr910c7sx3mrHnpNakiSI0HZpkSNZDwNUSuVoEb+
Y8GAvwe+Z5LOlVQt7Wc2Z9ANfEIBpNCqVX7UnJJEZlp9NPC3AViAVknj8/mu1MTy
SGjPDhZtEmRGubBAfT3jEJw9A8Gkd/dwWnYE/IaHAoGBANo7EMoqmXVSTwzBPVaC
ZDY2Q9sM7D/pNEA36IxhhTXbs5ajMCoJZXTHXJfRYtF/9PIGBhoIqGBRbQIPAFxG
DQ+7b8JKaC6HBoouXEEHw5ZA/RlNtiBL5sBVqvJi20qnJ1+Z2d++tYEQIJmPUxR7
XnLCLbi/Oj6ue3pFp4EiHKBp
-----END PRIVATE KEY-----`

const pubKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAt4nSgQRnJOu8ZPmRkHSW
LxsIju0mK/o8RKezc3qGOSx1PJIesetwnwFb7qtQ/oMctVCJ2yj8ZpPZqJoept5b
8ZzjCcMwO6g0YxkY/UpSzshVx8LkgWLXxkrR8LXb/25g+OYplztvNG+1/9cmREQR
1b6//1qUhKBXwGLPVFS1HOEHKdiYG+YdXylayUdPfpIVq3rJowMhR/zonbNg7bFR
phlgb83L4cJVIcZiZnLSDanWCo+fxAacBK6VFEalTt1aktdwqBXio5kXYk4YjNY1
rV8IOCHgB6QXvrCq9quD6flIP1AJ28UsGg9vvEgQknjuXyOfFl6VAoTK0xDYxREU
IwIDAQAB
-----END PUBLIC KEY-----`

func TestRsaOAEPEncryptDecrypt(t *testing.T) {
	var str = "dadsads"

	pub, err := ParseRsaPublicKey(pubKey)
	testz.Nil(t, err)

	prv, err := ParseRsaPrivateKey(prvKey)
	testz.Nil(t, err)

	enc, err := RsaOAEPEncrypt(str, []byte(nil), pub)
	testz.Nil(t, err)

	dec, err := RsaOAEPDecrypt(enc, "", prv)
	testz.Nil(t, err)

	testz.Equal(t, str, string(dec))
}

func TestRsaPKCS1v15EncryptDecrypt(t *testing.T) {
	var str = "dadsads"

	pub, err := ParseRsaPublicKey(pubKey)
	testz.Nil(t, err)

	prv, err := ParseRsaPrivateKey(prvKey)
	testz.Nil(t, err)

	enc, err := RsaPKCS1v15Encrypt(str, pub)
	testz.Nil(t, err)

	dec, err := RsaPKCS1v15Decrypt(enc, prv)
	testz.Nil(t, err)

	testz.Equal(t, str, string(dec))
}
