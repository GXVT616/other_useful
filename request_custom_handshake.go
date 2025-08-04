package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"regexp"
	"strings"

	utls "github.com/refraction-networking/utls"
)

type responce struct {
	Responce_data []byte
	Body          []byte
	StatusCode    string
}

func spec_handshake(url string) responce {

	var result responce

	host := strings.Split(url, "/")[2]
	path := strings.Join(strings.Split(url, "/")[3:], "/")

	addr := fmt.Sprintf("%v:443", host)

	// TCP-соединение
	rawConn, err := net.Dial("tcp", addr)
	if err != nil {
		log.Printf("TCP dial error: %v", err)
		return result
	}

	defer rawConn.Close()

	// Конфиг TLS
	config := &utls.Config{
		ServerName: host,
	}

	// Создаем клиентский TLS-объект с имитацией Chrome 120
	client := utls.UClient(rawConn, config, utls.HelloChrome_120)
	defer client.Close()

	// Выполняем TLS handshake
	err = client.Handshake()
	if err != nil {
		log.Printf("TLS handshake error: %v", err)
		return result
	}

	// Отправляем HTTP GET-запрос
	_, err = io.WriteString(client, fmt.Sprintf("GET /%v HTTP/1.1\r\nHost: %v\r\nConnection: close\r\n\r\n", path, host))
	if err != nil {
		log.Printf("Ошибка при отправке запроса: %v", err)
		return result
	}

	// Читаем ответ
	buf := make([]byte, 4096)
	for {
		n, err := client.Read(buf)
		if n > 0 {
			result.Responce_data = append(result.Responce_data, buf[:n]...)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("Ошибка чтения: %v", err)
			break
		}
	}
	(&result).get_statusCode()
	(&result).get_Parse()
	return result

}

func (r *responce) get_Parse() {

	rege := regexp.MustCompile(`\r\n\r\n`)
	reg_find := rege.Split(string((*r).Responce_data), -1)

	if len(reg_find) != 0 {
		// fmt.Println(reg_find)
		(*r).Body = []byte(reg_find[len(reg_find)-1])
	}

}

func (r *responce) get_statusCode() {

	rege := regexp.MustCompile(`HTTP/\d\.\d{1,}\s\d{3}`)
	(*r).StatusCode = "0"
	reg_find := string(rege.Find((*r).Responce_data))

	if len(reg_find) != 0 {
		// fmt.Println(reg_find)
		(*r).StatusCode = string([]rune(reg_find)[len(reg_find)-3:])
	}

}
