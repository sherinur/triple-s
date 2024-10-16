package main

import (
	"flag"
	"log"
	"net/http"

	"triple-s/internal/logger"
	"triple-s/internal/server"
)

// TODO: Создать HTTP-сервер, который будет слушать определённый порт.

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

func main() {
	var err error

	port := flag.String("port", "4400", "Port number")
	dir := flag.String("dir", ".", "Path to the directory")

	flag.Parse()

	logger.PrintfInfoMsg("Starting server on port :" + *port)
	logger.PrintfInfoMsg("Path to the directory set: " + *dir)

	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", server.HandleHealth)
	mux.HandleFunc("GET /", server.HandleListBuckets)
	mux.HandleFunc("PUT /{BucketName}", server.HandleCreateBucket)
	mux.HandleFunc("DELETE /{BucketName}", server.HandleDeleteBucket)

	loggedMux := logger.LogRequestMiddleware(mux)

	err = http.ListenAndServe(":"+*port, loggedMux)
	if err != nil {
		log.Fatal(err)
	}
}
