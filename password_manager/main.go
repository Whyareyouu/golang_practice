package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	encryption "password-manager-app/pkgs"

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

	isExit := false;

	_, err := os.Stat("accounts.json")
	if err == nil {
		accounts, err = loadAccountsFromFile("accounts.json")
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Данные восстановлены: %v\n", len(accounts))
	}

	http.HandleFunc("/account", func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "application/json")
	
		switch req.Method {
		case "GET":
			jsonData, err := json.Marshal(accounts)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			w.Write(jsonData)
	
		case "POST":
			var newAccount Account
			err := json.NewDecoder(req.Body).Decode(&newAccount)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			encryptedPass, err := encryption.EncryptPass(newAccount.Password, []byte(secret_key))
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
	
			newAccount.Password = encryptedPass
			accounts = append(accounts, newAccount)
	
			w.WriteHeader(http.StatusCreated)
			json.NewEncoder(w).Encode(newAccount)
			
			saveAccountsToFile(accounts, "accounts.json")
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `{"error": "Method not allowed"}`)
		}
	})

	fmt.Println(`Выберите вариант
	1. Web
	2. Terminal`)

	userChoise, _ := readLine("Выберите вариант взаимодействия")
	if(strings.TrimSpace(userChoise) == "1"){
		log.Println(http.ListenAndServe(":9090", nil))
	}
	if(strings.TrimSpace(userChoise) == "2"){
		for {
			if(isExit) {
				saveAccountsToFile(accounts, "accounts.json")
				break
			}
			getMenu(&accounts, secret_key, &isExit)
		}
	}
}


func saveAccountsToFile(accounts []Account, filename string) error {
	data, err := json.MarshalIndent(accounts, "", "  ")
	if err != nil {
		return err
	}

	err = os.WriteFile(filename, data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func loadAccountsFromFile(filename string) ([]Account, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var accounts []Account
	err = json.Unmarshal(data, &accounts)
	if err != nil {
		return nil, err
	}
	return accounts, nil
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
			if strings.Contains(account.URL,searchQuery) {
				return account
			}
		case "login":
			if strings.Contains(account.Login,searchQuery)  {
				return account
			}
		}
	}
	return Account{}
}


func getMenu(accounts* []Account, secret_key string, isExit *bool) []Account {
	

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
		
		EncryptPassword, _ := encryption.EncryptPass(strings.TrimSpace(password), []byte(secret_key))

		newMap := map[string]string{
			"url": strings.TrimSpace(url),
			"login": strings.TrimSpace(login),
			"password": EncryptPassword,
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
			decryptedPass, _ := encryption.DecryptPass(account.Password, []byte(secret_key));

			fmt.Printf("%v. URL: %s, Login: %s, Password: %s\n", index, account.URL, account.Login, decryptedPass)
		}
	case "3":
		searchQuery, _ := readLine("Введите url: ")
		account := getAcc(*accounts, "url", strings.TrimSpace(searchQuery))
		if(account.Password != ""){
			decryptedPass, _ := encryption.DecryptPass(account.Password, []byte(secret_key))
			fmt.Printf("URL: %s, Login: %s, Password: %s\n", account.URL, account.Login, decryptedPass)
		}else {
			fmt.Println("Аккаунт не найден")
		}
		
	case "4":
		searchQuery, _ := readLine("Введите login: ")
		account := getAcc(*accounts, "login", strings.TrimSpace(searchQuery))
		if(account.Password != ""){
			decryptedPass, _ := encryption.DecryptPass(account.Password, []byte(secret_key))
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
		*isExit = true;
		return *accounts
	default:
		fmt.Println("Такого дейсвия нет")
	}
	return *accounts
}