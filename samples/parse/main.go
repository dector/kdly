package main

import (
	"fmt"

	"github.com/dector/kdly"
)

func main() {
	doc, err := kdly.Parse(`
      server www.example.com {
        root * path=#"/www/example.com"#
        file_server

        header Access-Control-Allow-Origin "*"
        header Cache-Control max-age=3600
        header Cache-Control max-age=604800 {
          // Apply it only to the specific path
          path #"/css/*"#
        }

        // Enable for debuging
        /-debug {
          log level=info
        }
      }
	`)
	if err != nil {
		fmt.Println("Failed to parse config:")
		panic(err)
	}

	fmt.Printf("Loaded config:\n\n%s\n", kdly.ToKDL(doc))

	fmt.Println()
}
