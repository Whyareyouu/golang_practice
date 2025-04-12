package main

import (
	"bufio"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

type Account struct {
	URL string
	Login string
	Password string
}

func main () {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables");
	}
	secret_key := os.Getenv("SECRET_KEY");
	accounts := []Account{};
	fmt.Println(secret_key)
	fmt.Println(accounts)

	getMenu(&accounts, secret_key)
}


func encryptPass(plainText string, key []byte) (string, error) {
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptPass(cipherText string, key []byte) (string, error) {
	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if(err != nil){
		return "", err
	}
	block, err := aes.NewCipher(key)

	if(err != nil){
		return "", err
	}

	gcm, err := cipher.NewGCM(block)

	if(err != nil){
		return "", err
	}

	nonceSize := gcm.NonceSize()
	
	if len(cipherTextBytes) < nonceSize {
		return "", fmt.Errorf("ciphertext too short")
	}
	nonce, ciphertext := cipherTextBytes[:nonceSize], cipherTextBytes[nonceSize:]
	plainText, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}
	
func readLine (inputText string) (string, error){
	reader := bufio.NewReader(os.Stdin);

	fmt.Println(inputText)
	fmt.Print("-> ")

	text, err := reader.ReadString('\n')

	if(err != nil){
		log.Fatal(err)
		return "", err
	}
	return text, nil
}

func getAcc(accounts []Account, field string, searchQuery string) Account {
	for _, account := range accounts {
		switch field {
		case "url":
			if account.URL == searchQuery {
				return account
			}
		case "login":
			if account.Login == searchQuery {
				return account
			}
		}
	}
	return Account{}
}


func getMenu(accounts* []Account, secret_key string) []Account {
	

	fmt.Println(`__Менеджер паролей__
	1. Добавить аккаунт
	2. Посмотреть все аккаунты
	3. Найти аккаунт по url
	4. Найти аккаунт по логину
	5. Удалить аккаунт
	6. Выход`)

	userChoise, _ := readLine("")
	fmt.Printf("userChoise: %v", userChoise)
	
	switch(strings.TrimSpace(userChoise)){
	case "1":
		url, _ := readLine("Введите url: ")
		login, _ := readLine("Введите ваш логин: ")
		password , _:= readLine("Введите ваш пароль: ")
		
		encryptPassword, _ := encryptPass(strings.TrimSpace(password), []byte(secret_key))

		newMap := map[string]string{
			"url": strings.TrimSpace(url),
			"login": strings.TrimSpace(login),
			"password": encryptPassword,
		}

		newAccount := Account{
			URL: newMap["url"],
			Login: newMap["login"],
			Password: newMap["password"],
		}
		
		*accounts = append(*accounts, newAccount)
		
		return *accounts
	case "2":
		for index, account := range *accounts {
			decryptedPass, _ := decryptPass(account.Password, []byte(secret_key));

			fmt.Printf("%v. URL: %s, Login: %s, Password: %s\n", index, account.URL, account.Login, decryptedPass)
		}
	case "3":
		searchQuery, _ := readLine("Введите url: ")
		account := getAcc(*accounts, "url", strings.TrimSpace(searchQuery))
		if(account.Password != ""){
			decryptedPass, _ := decryptPass(account.Password, []byte(secret_key))
			fmt.Printf("URL: %s, Login: %s, Password: %s\n", account.URL, account.Login, decryptedPass)
		}else {
			fmt.Println("Аккаунт не найден")
		}
		
	case "4":
		searchQuery, _ := readLine("Введите login: ")
		account := getAcc(*accounts, "login", strings.TrimSpace(searchQuery))
		if(account.Password != ""){
			decryptedPass, _ := decryptPass(account.Password, []byte(secret_key))
			fmt.Printf("\nURL: %s, Login: %s, Password: %s\n", account.URL, account.Login, decryptedPass)
		}else {
			fmt.Println("Аккаунт не найден")
		}
	case "5":
		for index, account := range *accounts {
			fmt.Printf("%v. URL: %s, Login: %s\n", index, account.URL, account.Login)
		}
		deleteIndex, _ := readLine("Введите номер аккаунта для удаления: ")
		
		transformedStrToInt, err := strconv.Atoi(strings.TrimSpace(deleteIndex))
		if err != nil || transformedStrToInt < 0 || transformedStrToInt >= len(*accounts) {
			fmt.Println("Неверный номер")
			break
		}

		accs := *accounts
		*accounts = append(accs[:transformedStrToInt], accs[transformedStrToInt+1:]...)
		fmt.Println("Аккаунт успешно удалён")
		return *accounts
	
	case "6":
	default:
		fmt.Println("Такого дейсвия нет")
	}
	return *accounts
}