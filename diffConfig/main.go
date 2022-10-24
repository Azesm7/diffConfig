package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli"
)

var keypress map[string]string

var keypressfolse map[string]string

func MatchCheck(textFail1 map[string]string, textFail2 map[string]string, format string) map[string]string { // проверка файлов на соавпадение

	keypressfolse = make(map[string]string) // инцелизируем мапу

	for Key := range textFail1 { // перебераем каждый ключ из 1 фала и сравнивем их с ключами 2 файла

		if textFail1[Key] == textFail2[Key] {

			continue
		} else {
			if textFail2[Key] == "" { // если нет ключа
				if format == "json" {
					Textfolse := Key + "=" + textFail1[Key]
					keypressfolse[Textfolse] = "file 2 does not have this key"
				} else if format == "" {
					fmt.Println(Key, "=", textFail1[Key], "|", "file 2 does not have this key")
				}
			} else { // если значения разные
				if format == "json" {
					Textfolse := Key
					Textfolse2 := textFail1[Key] + "!=" + textFail2[Key]
					keypressfolse[Textfolse] = Textfolse2
				} else if format == "" {
					fmt.Println(Key, textFail1[Key], "!=", textFail2[Key])
				}
			}
		}

	}

	for Key2 := range textFail2 { // перебераем каждый ключ из 2 фала и сравнивем их с ключами 1 файла

		if textFail1[Key2] == textFail2[Key2] {

			continue
		} else {
			if textFail1[Key2] == "" { // если нет ключа
				if format == "json" {
					Textfolse := Key2 + "=" + textFail2[Key2]
					keypressfolse[Textfolse] = "file 1 does not have this key"
				} else if format == "" {
					fmt.Println(Key2, "=", textFail2[Key2], "|", "file 1 does not have this key")
				}
			} else { // если значения разные
				if format == "json" {
					Textfolse := Key2
					Textfolse2 := textFail1[Key2] + "!=" + textFail2[Key2]
					keypressfolse[Textfolse] = Textfolse2
				} else if format == "" {
					fmt.Println(Key2, textFail1[Key2], "!=", textFail2[Key2])
				}
			}
		}
	}
	return keypressfolse

}
func ReadFile(file string) map[string]string { //функция для считывание текста из файла
	keypress = make(map[string]string)
	f, err := os.Open(file) // открытие файла
	if err != nil {         //если ошибка
		log.Fatal(err)
	}
	defer f.Close() // закрытие файла
	scanner := bufio.NewScanner(f)
	for scanner.Scan() { // сканирование файла по строкам
		if strings.HasPrefix(scanner.Text(), "#") { // если это каментарий
			continue
		} else { // если нет то дабавляем их в мапу
			if strings.Count(scanner.Text(), "=") == 2 {
				s := scanner.Text()
				s1 := strings.Split(s, "=")
				str := fmt.Sprintf(`%s=%s`, s1[1], s1[2])
				keypress[string(s1[0])] = str
			} else {
				if strings.Count(scanner.Text(), "=") == 0 {
					continue
				} else {
					s := scanner.Text()
					s1 := strings.Split(s, "=")
					keypress[string(s1[0])] = string(s1[1])
				}
			}
		}
	}
	if err := scanner.Err(); err != nil { // если конец файла
		log.Fatal(err)
	}
	return keypress // возврат канала
}

func main() {
	var format string = ""
	var AddressFile2 string
	var AddressFile1 string
	var file1 string
	var file2 string
	var formatString string
	app := cli.NewApp() // инцелизируем флаг

	app.Flags = []cli.Flag{ // создание флага
		cli.StringFlag{
			Name:        "File1",         // имя флага
			Usage:       "Address File1", // вывод в help
			Destination: &file1,          // запись аргумента в переменную
		},
		cli.StringFlag{
			Name:        "File2",
			Usage:       "Address File2",
			Destination: &file2,
		},
		cli.StringFlag{
			Name:        "Format",
			Usage:       "json",
			Destination: &formatString,
		},
	}

	app.Action = func(c *cli.Context) error {
		AddressFile1 = file1
		AddressFile2 = file2
		format = formatString
		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
	textFail1 := ReadFile(AddressFile1) // передаём адрес где находиться файл и принимаем текст с файла
	textFail2 := ReadFile(AddressFile2) // передаём адрес где находиться файл и и принимаем текст с файла

	//-------------------------------------------------------------------------------------------------------------------------------
	//проверка на совпадение файлов
	FolseMap := MatchCheck(textFail1, textFail2, format)
	time.Sleep(2 * time.Second)
	if format == "json" {
		bytes, err := json.Marshal(FolseMap) // преобразование в json
		if err != nil {
			fmt.Println(err)
			return
		}
		text := string(bytes)
		fmt.Println(text)
	}

}
