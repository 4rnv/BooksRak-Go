package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
)

type Book struct {
	Author_name   string  `json:"author_name"`
	Avg_rating    float32 `json:"avg_rating"`
	Blurb         string  `json:"blurb"`
	Book_id       int     `json:"book_id"`
	Book_name     string  `json:"book_name"`
	Img_url       string  `json:"img_url"`
	Rating_1      int     `json:"rating_1"`
	Rating_2      int     `json:"rating_2"`
	Rating_3      int     `json:"rating_3"`
	Rating_4      int     `json:"rating_4"`
	Rating_5      int     `json:"rating_5"`
	Ratings_count int     `json:"ratings_count"`
	Year_pub      string  `json:"year_pub"`
}

func readjson(filename string) ([]Book, error) {
	byteValue, err := os.ReadFile(filename)
	var books []Book
	if err != nil {
		fmt.Println(err)
	}
	if err := json.Unmarshal(byteValue, &books); err != nil {
		fmt.Println(err)
	}
	return books, nil
}

func ichi(w http.ResponseWriter, r *http.Request) {
	books, _ := readjson("books.json")
	//fmt.Println(books)
	plate := template.Must(template.ParseFiles("booksrak.html"))
	plate.Execute(w, books)
}

func ni(w http.ResponseWriter, r *http.Request) {
	fmt.Println("HTMX Post was triggered.", r.Header.Get("HX-Request"))
	query := r.PostFormValue("search-query")
	fmt.Println(query)

	books, _ := readjson("books.json")
	var results []Book
	for _, book := range books {
		if (strings.Contains(strings.ToLower(book.Author_name), strings.ToLower(query))) || (strings.Contains(strings.ToLower(book.Book_name), strings.ToLower(query))) {
			results = append(results, book)
		}
	}
	fmt.Println(results)
	if err := renderResults(w, results); err != nil {
		fmt.Println(err)
	}
}

func renderResults(w http.ResponseWriter, books []Book) error {
	tmpl := `    
	{{ range . }}
    <div class="book">
    <h2>{{.Book_name}}</h2>
    <p>Author: {{.Author_name}}</p>
    <p>Year: {{.Year_pub}}</p>
    <p>Rating: {{.Avg_rating}}</p>
    </div>
    {{ end }}`

	t, err := template.New("results").Parse(tmpl)
	if err != nil {
		return err
	}

	return t.Execute(w, books)
}

func main() {
	fmt.Println("BooksRak")
	http.HandleFunc("/", ichi)
	http.HandleFunc("/search/", ni)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	port := ":8800"
	err := http.ListenAndServe(port, nil)
	if err != nil {
		fmt.Println("Error running server", err)
		log.Fatal(err)
	} else {
		fmt.Println("Server running on port ", port)
	}
}
