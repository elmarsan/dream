package main

import "log"

func main() {
	bg, err := newBackground("cat.gif")

	if err != nil {
		log.Fatal(err)
	}

	err = run(bg.animate)

	if err != nil {
		log.Fatal(err)
	}
}
