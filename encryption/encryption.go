package encryption

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"goweb/paillier"
	"math/big"
)

func Test_homomorphicCrypto() {

	// 生成一个128位的私钥
	privKey, _ := paillier.GenerateKey(rand.Reader, 128)

	// 公钥加密明文：数字15
	m15 := new(big.Int).SetInt64(15)
	c15, _ := paillier.Encrypt(&privKey.PublicKey, m15.Bytes())
	fmt.Println(c15)
	db, err := sql.Open("mysql", "root:1592933843zzz@tcp(127.0.0.1:3306)/sql_find")

	if err != nil {
		fmt.Printf("connect to db 127.0.0.1:3306 error: %v\n", err)
	}
	sqlStr := "insert into encryption(tid,x,y) values(?,?,?)"
	_, errr := db.Exec(sqlStr, 2, c15, c15)
	if errr != nil {
		fmt.Println("插入数据错误", err)
	}

	// 私钥解密密文：15
	d, _ := paillier.Decrypt(privKey, c15)
	plainText := new(big.Int).SetBytes(d)
	fmt.Println("Decryption Result of 15: ", plainText.String()) // 15

	// 公钥加密数字20
	m20 := new(big.Int).SetInt64(20)
	c20, _ := paillier.Encrypt(&privKey.PublicKey, m20.Bytes())

	// 将15的密文与20的密文相加（注意都是同一公钥）
	plusM16M20 := paillier.AddCipher(&privKey.PublicKey, c15, c20)
	// 使用私钥解密和的明文结果
	decryptedAddition, _ := paillier.Decrypt(privKey, plusM16M20)
	fmt.Println("Result of 15+20 after decryption: ",
		new(big.Int).SetBytes(decryptedAddition).String()) // 35!

	// 将15密文与10明文常数相加
	plusE15and10 := paillier.Add(&privKey.PublicKey, c15, new(big.Int).SetInt64(10).Bytes())
	decryptedAddition, _ = paillier.Decrypt(privKey, plusE15and10)
	fmt.Println("Result of 15+10 after decryption: ",
		new(big.Int).SetBytes(decryptedAddition).String()) // 25!

	// 将15密文与10明文常数相乘
	mulE15and10 := paillier.Mul(&privKey.PublicKey, c15, new(big.Int).SetInt64(10).Bytes())
	decryptedMul, _ := paillier.Decrypt(privKey, mulE15and10)
	fmt.Println("Result of 15*10 after decryption: ",
		new(big.Int).SetBytes(decryptedMul).String()) // 150!
}
