package main

import (
	"fmt"
	"io"
	"net/http"
	"strconv"
)

const chunkCount = 6
const maxFileSize = 11 << 30 // 11 GiB

// Обработчик загрузки файла
func uploadHandler(w http.ResponseWriter, r *http.Request) {
	// Получаем информацию о загружаемом файле
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error reading file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Получаем адреса серверов хранения из Consul
	storageServers, err := getStorageServers()
	if err != nil || len(storageServers) < chunkCount {
		http.Error(w, "Not enough storage servers available", http.StatusInternalServerError)
		return
	}

	// Ограничиваем количество серверов до нужного числа (6)
	storageServers = storageServers[:chunkCount]

	// Вычисляем примерный размер каждого куска
	fileSize := header.Size
	chunkSize := fileSize / chunkCount

	// Чтение и отправка частей файла во время загрузки
	for i, server := range storageServers {
		// Рассчитываем размер текущего куска
		var currentChunkSize int64
		if i == chunkCount-1 { // Последний кусок может быть больше/меньше из-за деления
			currentChunkSize = fileSize - (chunkSize * int64(i))
		} else {
			currentChunkSize = chunkSize
		}

		// Ограничиваем чтение файла для отправки одной части
		limitedReader := io.LimitReader(file, currentChunkSize)

		// Отправляем текущий кусок на сервер хранения
		err := sendChunkStream(server, header.Filename, i, limitedReader)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error sending chunk to storage server %s", server), http.StatusInternalServerError)
			return
		}
	}

	fmt.Fprintf(w, "File uploaded successfully")
}

// Функция для отправки куска файла на сервер хранения
func sendChunkStream(serverURL, filename string, chunkIndex int, chunkData io.Reader) error {
	// Создаем POST-запрос
	url := fmt.Sprintf("%s/upload", serverURL)
	req, err := http.NewRequest("POST", url, chunkData)
	if err != nil {
		return err
	}

	// Устанавливаем заголовки
	req.Header.Set("Content-Type", "application/octet-stream")
	req.Header.Set("File-Name", filename)
	req.Header.Set("Chunk-Index", strconv.Itoa(chunkIndex))

	// Отправляем запрос
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("received non-OK response code: %d", resp.StatusCode)
	}

	return nil
}

func main() {
	http.HandleFunc("/upload", uploadHandler)

	fmt.Println("Server is running on port 8080")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
	}
}
