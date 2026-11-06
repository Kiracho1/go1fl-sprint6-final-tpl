package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/Yandex-Practicum/go1fl-sprint6-final/internal/service"
)

// функция для корневого эндпоинта, которая возвращает HTML-страницу из файла index.html
func Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	http.ServeFile(w, r, "index.html")
}

// функция для эндпоинта /upload, которая конвертирует содержимое
// и сохраняет результат в локальный файл
func Upload(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// 1. Проверяем, что размер загружаемого файла не превышает 10 МБ
	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 2. Получаем файл из формы
	uploaded, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 3. Закрываем файл после завершения функции
	defer uploaded.Close()

	// 4. Читаем содержимое файла
	data, err := io.ReadAll(uploaded)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	// 5. Конвертируем содержимое файла
	converted, err := service.Convert(string(data))
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 6. Создаем локальное имя файла для сохранения результата
	ext := filepath.Ext(header.Filename)
	name := safeFileName(time.Now().UTC().String()) + ext
	out, err := os.Create(name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer out.Close()

	// 7. Записываем конвертированный контент в файл
	if _, err := out.WriteString(converted); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	// 8. Отправляем ответ клиенту с конвертированным содержимым
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, converted)
}

// safeFileName заменяет недопустимые символы в имени файла на допустимые.
// Без этого os.Create может вернуть ошибку, если имя файла содержит недопустимые символы.
func safeFileName(name string) string {
	replacer := strings.NewReplacer(
		":", "-",
		" ", "_",
		"+", "",
	)
	return replacer.Replace(name)
}