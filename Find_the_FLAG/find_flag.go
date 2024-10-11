package main

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"sync"
)


func decodeImage(filename string) (image.Image, error){
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}

	return img, err
}


func encodeImage(filename string, img image.Image) (error) {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	return png.Encode(f, img)
}

func inner_loop(i int, height int, img1 img, img2 img, img_out *img, wg *sync.WaitGroup){
	defer wg.Done()
	for j:= range height{
		//var r1, r2, g1, g2, b1, b2 uint32
		r1, g1, b1, _ := randomImage.At(i, j).RGBA()
		r2, g2, b2, _ := encryptedImage.At(i, j).RGBA()
		
		
		// if r1 == r2 && g1 == g2 && b1 == b2{
		// 	fmt.Println(i, j)
		// 	fmt.Println(randomImage.At(i,j).RGBA())
		// 	fmt.Println(encryptedImage.At(i,j).RGBA())
		// }
		col := color.RGBA{uint8(r2 + r1), uint8(g2 + g1), uint8(b2 + b1), 255}
		decryptedImage.Set(i, j, col)
	}
}

func main(){

	randomImage, err := decodeImage("random-image.png")
	if err != nil{
		fmt.Println("Failed to decode image:", err)
	}
	encryptedImage, err := decodeImage("encrypted1.png")
	if err != nil{
		fmt.Println("Failed to decode image:", err)
	}

	bounds := randomImage.Bounds()
	width, height := bounds.Max.X, bounds.Max.Y

	var wg sync.WaitGroup
	wg.Add(width)
	
	decryptedImage := image.NewRGBA(bounds)

	for i := range width{
		// for j:= range height{
		// 	//var r1, r2, g1, g2, b1, b2 uint32
		// 	r1, g1, b1, _ := randomImage.At(i, j).RGBA()
		// 	r2, g2, b2, _ := encryptedImage.At(i, j).RGBA()


		// 	// if r1 == r2 && g1 == g2 && b1 == b2{
		// 	// 	fmt.Println(i, j)
		// 	// 	fmt.Println(randomImage.At(i,j).RGBA())
		// 	// 	fmt.Println(encryptedImage.At(i,j).RGBA())
		// 	// }
		// 	col := color.RGBA{uint8(r2 + r1), uint8(g2 + g1), uint8(b2 + b1), 255}
		// 	decryptedImage.Set(i, j, col)
		// }
	}

	
	err = encodeImage("decoded1.png", decryptedImage)
	if err != nil{
		fmt.Println("Failed to save image:", err)
	}

	
	return
}

