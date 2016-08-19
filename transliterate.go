package main

import (
	"fmt"
	"github.com/fiam/gounidecode/unidecode"
	"io/ioutil"
	"os"
)

func main() {
	dir := "D:/mp3 summer 2016/"
	files, _ := ioutil.ReadDir(dir)
	for _, f := range files {
		if f.Name() != unidecode.Unidecode(f.Name()) {
			fmt.Println(unidecode.Unidecode(f.Name()))

			err := os.Rename(dir+f.Name(), dir+unidecode.Unidecode(f.Name()))
			if err != nil {
				fmt.Println(err)
				return
			}
		}
	}
	fmt.Println("DONE!")
}
