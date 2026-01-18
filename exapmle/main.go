package main

import (
	"fmt"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Welcome to Product API")
	})

	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Service running")
	})

	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "Here will be products")
	})

	mux.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
		id := r.URL.Path[len("/products/"):]
		fmt.Fprintf(w, "Product details with ID: %s", id)
	})

	fmt.Println("Start server on :8081")

	http.ListenAndServe(":8081", mux)
}

// package main

// import (
// 	"fmt"
// 	"net/http"
// )

// func main() {
// 	// 1. Создаем свой личный мультиплексор (роутер)
// 	mux := http.NewServeMux()

// 	// 2. Регистрируем обработчики ВНУТРИ mux, а не глобально
// 	mux.HandleFunc("/products", func(w http.ResponseWriter, r *http.Request) {
// 		fmt.Fprint(w, "Список продуктов")
// 	})

// 	mux.HandleFunc("/products/", func(w http.ResponseWriter, r *http.Request) {
// 		// r.URL.Path содержит полный путь, например "/products/123"
// 		id := r.URL.Path[len("/products/"):]
// 		fmt.Fprintf(w, "Детали продукта с ID: %s", id)
// 	})

// 	fmt.Println("Start server on :8081")
// 	// 3. Передаем наш mux в ListenAndServe ВМЕСТО nil
// 	http.ListenAndServe(":8081", mux)
// }

// package main

// import (
// 	"fmt"
// 	"log"
// 	"net/http"
// )

// func homeHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Printf("Запрос пришел на путь %s\n", r.URL.Path)

// 	w.WriteHeader(http.StatusOK)

// 	fmt.Fprintf(w, "Привет! Это твой первый Go сервер.")
// }

// func healthHandler(w http.ResponseWriter, r *http.Request) {
// 	fmt.Printf("Запрос пришел на путь %s\n", r.URL.Path)

// 	w.WriteHeader(http.StatusOK)

// 	fmt.Fprint(w, "Service is running")
// }

// func main() {
// 	http.HandleFunc("/", homeHandler)
// 	http.HandleFunc("/health", healthHandler)

// 	fmt.Println("Сервер запущен на порту :8080...")

// 	err := http.ListenAndServe(":8080", nil)
// 	if err != nil {
// 		log.Fatal("Ошибка запуска сервера: ", err)
// 	}
// }
