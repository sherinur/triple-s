package main

import (
	"flag"
	"fmt"
	"os"

	"triple-s/internal/server"
)

// TODO: Реализовать маршрутизацию для обработки PUT-запросов на путь /{BucketName}.

// TODO: Извлечь имя бакета из пути запроса.

// TODO: Реализовать функцию для проверки валидности имени бакета:
// Имя должно быть длиной от 3 до 63 символов.
// Должно содержать только маленькие буквы, цифры, дефисы и точки.

// TODO: Реализовать проверку уникальности имени бакета:
// Сохранять созданные имена бакетов в память (например, используя map).
// Вернуть ошибку 409 Conflict, если бакет с таким именем уже существует.

// TODO: Создать директорию для нового бакета на файловой системе.
// Использовать os.Mkdir для создания директории.

// TODO: Реализовать функцию для записи метаданных бакета в CSV файл:
// Открыть или создать файл buckets.csv.
// Записать имя созданного бакета и его статус (например, "created").

// TODO: Если запись в CSV файл успешна, вернуть ответ 200 OK с информацией о созданном бакете.

// TODO: Если имя некорректное, вернуть 400 Bad Request.

// TODO: Если создание бакета не удалось (например, ошибка при создании директории), вернуть 500 Internal Server Error.

// TODO: Добавить обработку всех других HTTP-методов, чтобы они возвращали 405 Method Not Allowed, если запрос не является PUT.

var (
	configPath string
	port       string
	dir        string
)

func init() {
	flag.StringVar(&port, "port", "4400", "Port number")
	flag.StringVar(&dir, "dir", ".", "Path to the directory")
	flag.StringVar(&configPath, "cfg", "configs/server.yaml", "Path to config file")
}

func main() {
	flag.Parse()

	port = ":" + port

	apiServer := server.New(server.NewConfig(configPath, port, dir))
	err := apiServer.Start()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
