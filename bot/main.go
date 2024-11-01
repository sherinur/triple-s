package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Response struct {
	StatusCode int
	Body       string
}

var waitingForFile map[int64]string // Хранит контекст команды `put_object`, ожидая файл

func main() {
	bot, err := tgbotapi.NewBotAPI("...")
	if err != nil {
		log.Panic(err)
	}

	waitingForFile = make(map[int64]string) // Инициализируем map для отслеживания команд `put_object`

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message != nil {
			if update.Message.IsCommand() {
				command := update.Message.Command()
				args := update.Message.CommandArguments()
				handleCommand(command, bot, update.Message.Chat.ID, args)
			} else if fileKey, waiting := waitingForFile[update.Message.Chat.ID]; waiting && update.Message.Document != nil {
				// Передаем objectKey в handleFileUpload
				objectKey := fileKey // Если у вас есть логика для objectKey
				handleFileUpload(bot, update.Message.Chat.ID, update.Message.Document.FileID, fileKey, objectKey)
				delete(waitingForFile, update.Message.Chat.ID) // Удаляем ожидание после загрузки файла
			} else {
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Введите команду для выполнения действия или используйте /start для просмотра доступных команд.")
				bot.Send(msg)
			}
		}
	}
}

// // Вызов функции в main
// if fileKey, waiting := waitingForFile[update.Message.Chat.ID]; waiting && update.Message.Document != nil {
// 	handleFileUpload(bot, update.Message.Chat.ID, update.Message.Document.FileID, fileKey, args)
// 	delete(waitingForFile, update.Message.Chat.ID) // Удаляем ожидание после загрузки файла
// }

func handleCommand(command string, bot *tgbotapi.BotAPI, chatID int64, args string) {
	var response Response // Используем новую структуру Response
	var err error

	switch command {
	case "start":
		response = Response{Body: `Добро пожаловать! Вот доступные команды:
			
	/health — проверка состояния сервера
	/list_buckets — получить список всех бакетов
	/create_bucket <имя> — создать новый bucket с указанным именем
	/delete_bucket <имя> — удалить bucket с указанным именем
	/put_object <bucket_name/object_name> — загрузить объект в указанный bucket (затем отправьте файл)
	/get_object <bucket_name/object_name> — получить объект из bucket
	/delete_object <bucket_name/object_name> — удалить объект из bucket

	Просто введите команду или нажмите на соответствующую кнопку.`}
	case "health":
		response, err = makeRequest("GET", "http://127.0.0.1:4400/health", "")
	case "list_buckets":
		response, err = makeRequest("GET", "http://127.0.0.1:4400/", "")
	case "create_bucket":
		if args == "" {
			response = Response{Body: "Пожалуйста, укажите имя bucket-а после команды."}
		} else {
			response, err = makeRequest("PUT", "http://127.0.0.1:4400/"+args, "")
		}
	case "delete_bucket":
		response, err = makeRequest("DELETE", "http://127.0.0.1:4400/"+args, "")
	case "put_object":
		if args == "" {
			response = Response{Body: "Пожалуйста, укажите путь bucket-а и имя объекта после команды."}
		} else {
			response = Response{Body: "Отправьте файл, который вы хотите загрузить в bucket."}
			waitingForFile[chatID] = args // Сохраняем контекст для ожидания файла
		}
	case "get_object":
		// response, err = makeRequest("GET", "http://127.0.0.1:4400/"+args, "")
		// if err != nil {
		// 	log.Printf("Ошибка запроса: %v", err)
		// 	return
		// }

		// // Проверка на успешный статус ответа
		// if response.StatusCode != http.StatusOK {
		// 	log.Printf("Ошибка: статус ответа %d", response.StatusCode)
		// 	return
		// }

		// // Чтение бинарного тела ответа
		// body := []byte(response.Body) // Преобразуем строку в срез байтов

		// // Получаем имя объекта из аргументов
		// objectName := args

		// // Создание временного файла для сохранения полученного содержимого
		// tmpFile, err := os.Create(objectName) // Используйте `objectName` для имени файла
		// if err != nil {
		// 	log.Printf("Ошибка создания файла: %v", err)
		// 	return
		// }
		// defer tmpFile.Close()

		// // Запись бинарного содержимого в файл
		// if _, err := tmpFile.Write(body); err != nil {
		// 	log.Printf("Ошибка записи в файл: %v", err)
		// 	return
		// }

		// // После создания временного файла
		// document := tgbotapi.NewDocument(chatID, tgbotapi.File{
		// 	File: tgbotapi.FileReader{
		// 		Name:   tmpFile.Name(), // Имя файла
		// 		Reader: tmpFile,        // Чтение из временного файла
		// 	},
		// })

		// // Отправка документа
		// if _, err := bot.Send(document); err != nil {
		// 	log.Printf("Ошибка отправки файла: %v", err)
		// }

	case "delete_object":
		response, err = makeRequest("DELETE", "http://127.0.0.1:4400/"+args, "")
	default:
		response = Response{Body: "Неизвестная команда. Введите /start для списка команд."}
	}

	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Ошибка: "+err.Error())
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatID, response.Body)
		bot.Send(msg)
	}
}

func handleFileUpload(bot *tgbotapi.BotAPI, chatID int64, fileID string, bucketName string, objectKey string) {
	file, err := bot.GetFileDirectURL(fileID)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Ошибка получения файла: "+err.Error())
		bot.Send(msg)
		return
	}

	// Загружаем содержимое файла в буфер
	fileResp, err := http.Get(file)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Ошибка скачивания файла: "+err.Error())
		bot.Send(msg)
		return
	}
	defer fileResp.Body.Close()

	// Чтение содержимого файла
	fileContent, err := io.ReadAll(fileResp.Body)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Ошибка чтения файла: "+err.Error())
		bot.Send(msg)
		return
	}

	// Проверка размера файла и первые несколько байт
	log.Printf("Загружен файл размером: %d байт", len(fileContent))
	log.Printf("Первые байты файла: %v", fileContent[:10]) // Логируем первые 10 байт

	// Отправка файла на сервер
	serverResponse, err := uploadFileToServer("http://127.0.0.1:4400/"+objectKey, bucketName, objectKey, fileContent)
	if err != nil {
		msg := tgbotapi.NewMessage(chatID, "Ошибка отправки файла на сервер: "+err.Error())
		bot.Send(msg)
	} else {
		msg := tgbotapi.NewMessage(chatID, "Файл успешно загружен: "+serverResponse)
		bot.Send(msg)
	}
}

// uploadFileToServer отправляет файл на сервер с необходимыми метаданными.
func uploadFileToServer(url string, bucketName string, objectKey string, fileContent []byte) (string, error) {
	// Получаем размер файла
	size := int64(len(fileContent))
	contentType := "application/octet-stream" // Устанавливаем значение по умолчанию для MIME-типа

	// Создаем HTTP запрос
	req, err := http.NewRequest("PUT", url, bytes.NewBuffer(fileContent))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", contentType)               // Устанавливаем тип содержимого
	req.Header.Set("Bucket-Name", bucketName)                 // Устанавливаем название бакета
	req.Header.Set("Content-Length", fmt.Sprintf("%d", size)) // Устанавливаем длину содержимого
	req.Header.Set("Object-Key", objectKey)                   // Устанавливаем ключ объекта

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("сервер вернул ошибку: %s, текст: %s", resp.Status, string(respBody))
	}

	return string(respBody), nil
}

// makeRequest выполняет HTTP-запрос с указанным методом, URL и данными в формате XML.
func makeRequest(method, urlString, data string) (Response, error) {
	client := &http.Client{}

	// Создаем новый запрос
	req, err := http.NewRequest(method, urlString, strings.NewReader(data))
	if err != nil {
		return Response{}, err
	}

	// Устанавливаем заголовок Content-Type для XML
	req.Header.Set("Content-Type", "application/xml")

	// Выполняем запрос
	resp, err := client.Do(req)
	if err != nil {
		return Response{}, err
	}
	defer resp.Body.Close()

	// Читаем тело ответа
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Response{}, err
	}

	// Возвращаем ответ с статусом и телом
	return Response{
		StatusCode: resp.StatusCode,
		Body:       string(body),
	}, nil
}
