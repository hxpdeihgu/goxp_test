package main
import (
	"fmt"
	"net/http"
	"github.com/hxpdeihgu/goxp"
)

func main() {
	d := goxp.Incubate()

	d.Get("/", func() string {
		return "Hello world!"
	})

	d.Get("/hxp", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello hxp!")
	})
	d.Run()
}
