package images

import (
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"

	"github.com/nfnt/resize"
)

// ResizeImage is used to resize the image.
func ResizeImage(originalImage, resizedImage string, width, height uint) error {
	content, err := os.Open(originalImage)
	if err != nil {
		log.Println("Could not read image from amazon S3. Error :", err.Error())
		return err
	}
	defer content.Close()

	original, imgFmt, err := image.Decode(content)
	if err != nil {
		log.Println("Could not decode image from content Error :", err.Error())

		original, err = png.Decode(content)
		if err != nil {
			log.Println("Could not decode image from content Error :", err.Error())
			return err
		}

		imgFmt = "png"
	}

	log.Println("Format of the image is : ", imgFmt)
	thumb := resize.Resize(width, height, original, resize.Bilinear)

	// Destination
	dst, err := os.Create(resizedImage)
	if err != nil {
		return err
	}
	defer dst.Close()

	err = jpeg.Encode(dst, thumb, nil)
	if err != nil {
		log.Println("Could not encode thumb image . Error :", err.Error())
		return err
	}
	return nil
}
