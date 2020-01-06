package main

import(
    "fmt"
    "strconv"
    "os"
    "time"
    "math/rand"
    "math"
)

type Token struct {
 	data string
 	recipient int
 	ttl int
} 

func send(prev <-chan Token, next chan<- Token, id int) {
				for {	
					token := <- prev				
					if token.ttl == 0 {
						next <- token
						close(next)
						return
					}else{
    				    fmt.Printf("got by %d, ttl = %d\n", id, token.ttl)
					    token.ttl--
                    }			

					if token.recipient == id {
						fmt.Printf("Received by %d!\n", id)
                        token.ttl = 0
						next <- token	
						close(next)
						return
					} else if  token.ttl == 0 {
						fmt.Printf("Didn't reach the end. Received by %d!\n", id)
						next <- token	
						close(next)
						return
					} else {
						next <- token			  				
					}
				}
}
			

func main() {
	if (len(os.Args) < 2) {
		fmt.Println("USAGE: go run 'TokenRing.go' <Number of goroutines>")
		os.Exit(0);
	}

	temp, err := strconv.Atoi(os.Args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}
    N := int(math.Abs(float64(temp)))

	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	msg := Token{"HW", r.Intn(N)+1, r.Intn(N*2)}
    fmt.Println("Message: '", msg.data, "' to ", msg.recipient, " with ttl=", msg.ttl)

	first := make(chan Token)
	last := first
	for i := 1; i <= N; i++ {
		current := make(chan Token)
		go send(last, current, i);
		last = current;
	}	 
	
	first <- msg
	go send(last, first, 0)	
    <- last
}
