package cry

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"os"
)

func Getkeys(name string) error {
	//create path
	ex, _ := PathExists(".esnd/crypto/private")
	if !ex {
		os.MkdirAll(".esnd/crypto/private", os.ModePerm)
	}
	ex, _ = PathExists(".esnd/crypto/public")
	if !ex {
		os.MkdirAll(".esnd/crypto/public", os.ModePerm)
	}
	//得到私钥
	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	x509_Privatekey := x509.MarshalPKCS1PrivateKey(privateKey)
	//创建一个用来保存私钥的以.pem结尾的文件
	fp, err := os.Create(".esnd/crypto/private/" + name + ".pem")
	if err != nil {
		return err
	}
	defer fp.Close()
	//将私钥字符串设置到pem格式块中
	pem_block := pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: x509_Privatekey,
	}
	//转码为pem并输出到文件中
	pem.Encode(fp, &pem_block)

	//处理公钥,公钥包含在私钥中
	publickKey := privateKey.PublicKey
	//接下来的处理方法同私钥
	//通过x509标准将得到的ras私钥序列化为ASN.1 的 DER编码字符串
	x509_PublicKey, _ := x509.MarshalPKIXPublicKey(&publickKey)
	pem_PublickKey := pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: x509_PublicKey,
	}
	file, err := os.Create(".esnd/crypto/public/" + name + ".pem")
	if err != nil {
		return nil
	}
	defer file.Close()
	//转码为pem并输出到文件中
	pem.Encode(file, &pem_PublickKey)
	return nil
}
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func RSA_encrypter(key string, msg []byte) ([]byte, error) {
	//下面的操作是与创建秘钥保存时相反的
	//pem解码
	block, _ := pem.Decode([]byte(key))
	//x509解码,得到一个interface类型的pub
	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//加密操作,需要将接口类型的pub进行类型断言得到公钥类型
	cipherText, err := rsa.EncryptPKCS1v15(rand.Reader, pub.(*rsa.PublicKey), msg)
	if err != nil {
		return nil, err
	}
	return cipherText, nil
}
func RSA_decrypter(key string, cipherText []byte) ([]byte, error) {
	block, _ := pem.Decode([]byte(key))
	PrivateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	//二次解码完毕，调用解密函数
	return rsa.DecryptPKCS1v15(rand.Reader, PrivateKey, cipherText)

}
