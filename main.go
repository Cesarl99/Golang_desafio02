package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

var wg sync.WaitGroup
var mx sync.Mutex

func main() {
	f := make(map[int]map[string]map[string]int)
	for {
		scanner := bufio.NewScanner(os.Stdin)
		fmt.Println("Qual posição que você deseja saber o n° da sequência de fibonacci:")
		scanner.Scan()
		numero := scanner.Text()
		num, err := strconv.Atoi(numero)
		if err != nil {
			log.Fatal(err)
		}
		if f[num] != nil {
			fmt.Println("Valor Já executado")
			fmt.Println(f[num])
		} else {
			wg.Add(4)
			go func() {
				fibtotal := make(map[string]map[string]int)
				canalfib := make(chan int, 1)
				canaltempo := make(chan int, 1)
				go func() {
					start := time.Now()
					canalfib <- fib(num)
					final1 := time.Since(start)
					timestring := final1.String()
					t, _ := time.ParseDuration(timestring)
					milisegundotime := t.Milliseconds()
					canaltempo <- int(milisegundotime)
					wg.Done()
				}()
				select {
				case <-time.After(time.Millisecond * 500):
					go func() {
						caminho := strconv.Itoa(num)
						caminhofinal := "/" + caminho
						app := fiber.New()
						app.Get(caminhofinal, func(c *fiber.Ctx) error {
							return c.JSON(fiber.Map{"Done": "False"})
						})
						teste := ":" + numero + "00"
						fmt.Printf("ouvido em : Localhost%v %v\n", teste, caminhofinal)
						app.Listen(teste)
					}()
					break
				case resultado := <-canalfib:
					go func() {
						tempoexecução := <-canaltempo
						caminho := strconv.Itoa(num)
						caminhofinal := "/" + caminho
						fibtotal = funcmapdone(resultado, num, tempoexecução)
						apc := fiber.New()
						go func() {
							apc.Get(caminhofinal, func(c *fiber.Ctx) error {
								return c.JSON(fibtotal)
							})
						}()
						f[num] = fibtotal
						wg.Done()
						teste := ":" + numero + "00"
						fmt.Printf("ouvido em : Localhost%v %v\n", teste, caminhofinal)
						apc.Listen(teste)
					}()
				}
				wg.Done()
				wg.Wait()
			}()
			time.Sleep(1000)
			scanner := bufio.NewScanner(os.Stdin)
			fmt.Println("Digite S para continuar e N para sair :")
			scanner.Scan()
			txt := scanner.Text()
			txtu := strings.ToUpper(txt)
			if txtu == "S" {
				continue
			} else if txtu == "N" {
				http.HandleFunc("/Final", func(w http.ResponseWriter, r *http.Request) {
					encoder := json.NewEncoder(w)
					encoder.Encode(f)
				})
				http.ListenAndServe(":8000", nil)
				fmt.Println("Todos os valores executados disponivel em:Localhost:8000/Final")
				scanner := bufio.NewScanner(os.Stdin)
				fmt.Println("Aperte o  Enter para sair :")
				scanner.Scan()
				txte := scanner.Text()
				if txte == "" {
					break
				} else {
					fmt.Println("opção invalida")
				}

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

func funcmapdone(resultado, num, tempo int) map[string]map[string]int {
	fibtotal := map[string]map[string]int{
		"done": {
			"input":    num,
			"output":   resultado,
			"duration": tempo,
		},
	}
	return fibtotal
}
