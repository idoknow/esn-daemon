package main

import (
	"esnd/src/cry"
	"fmt"
	"io/ioutil"
	"strconv"
)

func main() {
	msg := "daiuohdaisubhf9qgb9gw09hn-9ncbb9asgv972gubhf9daiuohdaisubhf9qsubhf9da-9ncbohd9daiuohdaisubhf9qgb9gw09hn-9ncbohdaisubhf9qsubhf9daiuohdaisubhf9qgb9gw09hn-9ncbb9asgv972gf9qgb9gw09hn-9ncbb9asgv972gy890bneionb doacdvdvfscxcvxvcxvvcvxvvxvbfhgfbfdgdgfsdftgfrb cdrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrttttttgdfdaiuohdaisubhf9qgb9gw09hn-9ncbb9asgv972gy890bneionb doacdvdvfscxcvxvcxvvcvxvvxvbfhgfbfdgdgfsdftgfrb cdrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrrttttttgdf"

	fmt.Println("msg length:" + strconv.Itoa(len(msg)))

	cry.Getkeys("test")

	publicKey, err := ioutil.ReadFile(".esnd/crypto/public/test.pem")
	if err != nil {
		panic(err)
	}
	e0, err := cry.RSA_encrypter(string(publicKey), []byte(msg))
	if err != nil {
		panic(err)
	}
	fmt.Println(e0)

	privateKey, err := ioutil.ReadFile(".esnd/crypto/private/test.pem")
	if err != nil {
		panic(err)
	}
	d0, err := cry.RSA_decrypter(string(privateKey), e0)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(d0))
}
