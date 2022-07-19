package main

/*
Esta API usa a web framework Gin.
Ela será uma API de BookStore. Será capaz de armazenar livros, fazer checkIn e checkOut,
adicionar e exccluir livros, ver todos os livros e buscar um livro por sua ID.
*/

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

//Para armazenar nossos livros, usaremos uma struct com a representação json.
type book struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}


//Aqui será o slice de livros. Ela será uma estrutura de dados que irá representar todos os nossos livros.
var books = []book{
	{ID: "1", Title: "Capitães da Areia", Author: "Jorge Amado", Quantity: 2},
	{ID: "2", Title: "Dom Casmurro", Author: "Machado de Assis", Quantity: 5},
	{ID: "3", Title: "A droga da obediência", Author: "Pedro Bandeira", Quantity: 6},
}


func getBooks(c *gin.Context) { 
	c.IndentedJSON(http.StatusOK, books)
}

func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found."})
		return
	}
	c.IndentedJSON(http.StatusOK, book)
}

func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok { //check se temos uma id ou nao
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}
	book, err := getBookById(id) //tentar pegar o book por sua id

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not found."})
		return
	}

	if book.Quantity <= 0 { //checkar a quantidade
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available."})
		return
	}
	//se nada acima acontecer, ira ser feito o checkout do livro, diminuindo 1 de seu total
	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok { //check se temos uma id ou nao
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}
	book, err := getBookById(id) //tentar pegar o book por sua id

	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not found."})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

//Get book by id do tipo string, retornando um ponteiro para book e retornando error
func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id { //checar se o bookId é igual ao id
			return &books[i], nil
		}
	}
	return nil, errors.New("Book not found.")
}

//Função para adicionar livros à slice.
func createBook(c *gin.Context) {
	var newBook book //criamos uma variavel newBook que é do tipo book.

	if err := c.BindJSON(&newBook); err != nil { //tentar ligar o JSON ao newBook struct. Passamos um ponteiro
		return //em new books pra modificar os valores e depois vamos checar
	} //se temos um erro ou nao. Se o erro nao for igual a nill, vamos dar
	//um return, e o BindJSON vai lidar com enviar o erro.
	books = append(books, newBook)              //Caso nao tenha erro, vamos criar um livro com append no books slice.
	c.IndentedJSON(http.StatusCreated, newBook) //Vamos retornar o livro criado com o Status code Created.

}

func deleteBookById(c *gin.Context) {
	id := c.Param("id")

	for i, a := range books {
		if a.ID == id {
			aux := books[i+1:]
			books = append(books[:i], aux...)

			c.IndentedJSON(http.StatusOK, books)
			return
		}
	}
}

/*
O router será responsável por lidar com diferentes routs e diferentes end points da api.
A variável router vai rotear uma rout específica a uma função.
*/
func main() {
	router := gin.Default()                 //aqui a criação do router. O router vem do Gin.
	router.GET("/books", getBooks)          //a rout que iremos lidar aqui é /books.
	router.GET("books/:id", bookById)       //GET: getting information
	router.POST("/books", createBook)       //POST: adding information
	router.PATCH("/checkout", checkoutBook) //PATCH: updating something
	router.PATCH("/return", returnBook)
	router.DELETE("/books/:id", deleteBookById)
	router.Run("localhost:8080") //quando entrarmos na porta 8080/books, irá chamar a função getBooks.
}
