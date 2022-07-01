package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

var wg sync.WaitGroup
var mx sync.Mutex

func main() {
	f := make(map[int]map[string]int)
	app := fiber.New()
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Qual posição que você deseja saber o numero da sequencia de fibonacci (digite 1 para sair e para gerar um mostrar um total de consultas): ")
		scanner.Scan()
		numero := scanner.Text()
		num, err := strconv.Atoi(numero)
		if err != nil {
			log.Fatal(err)
		}
		if num != 1 {
			if f[num] != nil {
				fmt.Println("Valor Já executado")
				fmt.Println(f[num])
			} else {
				wg.Add(2)
				go func() {
					fibtotal := make(map[string]int)
					erro := make(map[string]string)
					canalfib := make(chan int, 1)
					go func() {
						mx.Lock()
						canalfib <- fib(num)
						mx.Unlock()
						wg.Done()
					}()
					select {
					case <-time.After(time.Millisecond * 500):
						erro["done"] = "false"
						caminho := strconv.Itoa(num)
						caminhofinal := "/" + caminho
						app.Get(caminhofinal, HHandler)
						fmt.Printf("ouvindo em : localhost:8000%v\n", caminhofinal)
						wg.Done()
						app.Listen(":5000")
					case resultado := <-canalfib:
						fibtotal["input"] = num
						fibtotal["result"] = resultado
						fibtotal["duration"] = 0
						caminho := strconv.Itoa(num)
						caminhofinal := "/" + caminho
						http.HandleFunc(caminhofinal, func(w http.ResponseWriter, r *http.Request) {
							encoder := json.NewEncoder(w)
							encoder.Encode(fibtotal)
						})
						f[num] = fibtotal
						fmt.Printf("ouvindo em : localhost:8000%v\n", caminhofinal)
						http.ListenAndServe(":8000", nil)
					}
					wg.Done()
					wg.Wait()
				}()
			}
		} else {
			http.HandleFunc("/Final", func(w http.ResponseWriter, r *http.Request) {
				encoder := json.NewEncoder(w)
				encoder.Encode(f)
			})
			fmt.Println("Todos os valores pesquisados estão em :LocalHost/8000/Final")
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Println("Digite ENTER para sair")
			scanner.Scan()
			txt := scanner.Text()
			if txt == "" {
				break
			}

		}
	}
}

func fib(n int) int {
	if n <= 1 {
		return n
	}
	return fib(n-1) + fib(n-2)
}

func HHandler(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{"Done": "False"})
}
