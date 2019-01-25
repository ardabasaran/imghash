package main

import (
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io/ioutil"
	"os"
	"path/filepath"
	"similar-image-finder/imghash"
	"time"
)


func main() {
	files, err := ioutil.ReadDir("./resources")
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return
	}

	groups := make([][]imghash.ImageHashPair, 0)
	start := time.Now()
	for filenum, file := range files {
		fmt.Printf("(%d/%d) Working on file %v\n",filenum+1,len(files),file.Name())
		infile, err := os.Open("resources/" + file.Name())
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}
		src, _, err := image.Decode(infile)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			return
		}
		defer infile.Close()

		fileHash := imghash.Ahash(src)
		placed := false

		for i, group := range groups {
			fmt.Printf("\t(%d/%d) Comparing against group %d\n", i+1, len(groups), i+1)
			for _, imageHashPair := range group {
				distance := imghash.HammingDistance(fileHash, imageHashPair.Hash)
				if distance < 10 {
					fmt.Printf("\tPlacing to group %d\n", i+1)
					group = append(group, imghash.ImageHashPair{
						Filename:file.Name(),
						Hash:fileHash,
						Image:src,
					})
					placed = true
					groups[i] = group
					break
				}
			}
		}

		if !placed {
			fmt.Printf("\tNew group. Placing to the group %d\n", len(groups)+1)
			groups = append(groups, make([]imghash.ImageHashPair, 0))
			groups[len(groups)-1] = append(groups[len(groups)-1], imghash.ImageHashPair{
				Filename:file.Name(),
				Hash:fileHash,
				Image:src,
			})
		}
	}
	elapsed := time.Since(start).Seconds()
	fmt.Printf("Grouping took %v\n", elapsed)

	for i, group := range groups {
		for _, imageHashPair := range group {
			path := filepath.Join(".","groups",fmt.Sprintf("group%d",i+1))
			os.MkdirAll(path, os.ModePerm)
			saveImage(imageHashPair.Image, path + "/" + imageHashPair.Filename)
		}
	}
}

func saveImage(img image.Image, filename string) {
	out, err := os.Create("./" + filename)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	err = jpeg.Encode(out, img, &jpeg.Options{jpeg.DefaultQuality})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}