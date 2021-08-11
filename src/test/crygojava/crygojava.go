package main

import (
	"esnd/src/cry"
	"fmt"
)

var privateKey = "-----BEGIN PRIVATE KEY-----" +
	"MIIEowIBAAKCAQEA2Wfj5T9lwALueG15t62lO6517QxUyeU/dSdSIvucKSPh9wYD" +
	"iZfrez0T3HpjXSRBrIi5cnQH3f5BPBnDVy7WhK2i8QU2ebMHGMbRxbNuZAGhvZ/D" +
	"UfjBLQRMSzYaZvqXU3CCrnvpagvs9T87k1LJDwi9woTe4GlZ+pqkfODclvpo40fN" +
	"czFsh6IX8qrysBvIFOeUWExFUD+Dl/oOkGFc2M0MLQ6O1RhupHz/T1jAGlN3sNol" +
	"91mK1VAbjhMGY/uJvfbKJhdSm4U1gjfStL2P6Zzwjybf24A5IOY2oCMpiLTGqfuX" +
	"bveSVtkoAztiKCF6sqGSOiOIfXFIGlI4s3hYnQIDAQABAoIBAQC7NjZOLCi/jwax" +
	"l3wwCoz19sa/6VV+QjZB+SlGzKpt1uN356rWKodyKWdX/eBgzZ7sJxSilX5M0Ox2" +
	"B61p/wBlYmyk5htB80OSN2tetqPB5JHWC6STiwU2cbQNNDrKINJ83K779+JJGpnj" +
	"mp7/v1M56goWXnrafn4oSlCI5M2wB72Qzt+JiFFBSzhi9eLd3Qr8FwdS1wfVD9I9" +
	"gmH5woZBJ7G/gwL/b0dkYFJ9a9msbV1Cqpbb/eb1WvRvhajEHGKwrdVM6UmnZStE" +
	"PwCY3Sj+IdiJ9AAC92VXGkht+Xlgtsvbi6EvPmG0PeX8/HJKHBizkiHxVNd90DP0" +
	"zC0ZsZ4BAoGBAN4++4nw2ApUO+fRkL9RwMz7qqBG1Woo2xWVB9h851AYKdHXgBtu" +
	"lrRPcsRCf6Z2EffNdZjBCZ8J8LzV3HCyAUvcldPk+VbGm0FphOCYBwiBWgWOLt5h" +
	"yddvN/dMyKvTQFG9gd1183FaNdiDu3aAsbRjdwv5eUfZYv5TgZ3oSd0dAoGBAPps" +
	"t+7l9cY5v2hDjIgYwctZm0prdfIexURZgL/yzgmPeqgXSgruDSvTIkAay5hx8LPL" +
	"tuwEFFj063+a1PXF1z7zXMVE6UeYu6d0pFE2yuyfO1Ea9W3q9OS7uuMnTrfRCxG/" +
	"VpJ7MExJWe1xudCPtxIiGYwztMSowba6FMeRBRGBAoGALAHKqxC+pqTxS8DqaYfV" +
	"poE60wvTnHbEkux0pkBtSSXPuhZy7nuiacfFkOkd/6cnfar4Wyv2LMC6I5oxUTte" +
	"GFhwbonLeYxQF86+Gf7gfaWnXqw9yZkRb5A9Q8G3hpaJCOZ+fYyqjMpxGRNUnm1z" +
	"QqXjX8KhakG4YWXFp6/kWF0CgYA/feEEiPlPUMTewoGW3/AChq2AqM42nOaW9bpW" +
	"8FCcy+vlQkJbkw9z1QwSBLkp5XmJnFS8cixWgYJT0AW+anKwWzNiMJ6UsHyjcEdY" +
	"7/NzGswHPDaNr8x3UcGIZibnI/EShtiEOwd7z/0k3nimEEnyJwMjMNjcI405ruQl" +
	"1PbcAQKBgBKgiSV6PCKi52rUG/BeA4xSuXit2qdTDouevbhGskaNoZWeKTH59+DR" +
	"+SPRPmAO/jacuQ+N7iZtHRltU1wZm0An/PhC0ZFHGO6MGrpp9tFWvGzJ164MhqyJ" +
	"OBA5LHiQqUfM61wVsmYTRpi7+JMWQNN6tPTjZhVFovkJHjgNzGwu" +
	"-----END PRIVATE KEY-----"
var publicKey = "-----BEGIN PUBLIC KEY-----\n" +
	"MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2Wfj5T9lwALueG15t62l\n" +
	"O6517QxUyeU/dSdSIvucKSPh9wYDiZfrez0T3HpjXSRBrIi5cnQH3f5BPBnDVy7W\n" +
	"hK2i8QU2ebMHGMbRxbNuZAGhvZ/DUfjBLQRMSzYaZvqXU3CCrnvpagvs9T87k1LJ\n" +
	"Dwi9woTe4GlZ+pqkfODclvpo40fNczFsh6IX8qrysBvIFOeUWExFUD+Dl/oOkGFc\n" +
	"2M0MLQ6O1RhupHz/T1jAGlN3sNol91mK1VAbjhMGY/uJvfbKJhdSm4U1gjfStL2P\n" +
	"6Zzwjybf24A5IOY2oCMpiLTGqfuXbveSVtkoAztiKCF6sqGSOiOIfXFIGlI4s3hY\n" +
	"nQIDAQAB\n" +
	"-----END PUBLIC KEY-----\n"

func main() {
	fmt.Println(publicKey)
	e0, err := cry.RSA_encrypter(publicKey, []byte("{\"User\":\"root\",\"Pass\":\"changeMe\"}"))
	if err != nil {
		panic(err)
	}
	fmt.Println(string(e0))
}
