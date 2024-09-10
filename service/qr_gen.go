package service

import (
	"fmt"
	"image/png"
	"os"

	// barcode scale
	"github.com/boombuler/barcode"
	// Build the QRcode for the text
	"github.com/boombuler/barcode/qr"
)

func QrCodeGen(t string, filename string) (*os.File, barcode.Barcode, error) {
	// Create the barcode
	qrCode, err := qr.Encode(t, qr.M, qr.Auto)
	if err != nil {
		fmt.Println("could not generate the qr code : ", err)
		return nil, nil, err
	}

	// Scale the barcode to 200x200 pixels
	qrCode, err = barcode.Scale(qrCode, 2000, 2000)
	if err != nil {
		fmt.Println("could not scale the qr code : ", err)
		return nil, nil, err
	}

	// create the output file
	file, err := os.Create(filename + ".png")
	if err != nil {
		fmt.Println("could not create the file : ", err)
		return nil, nil, err
	}

	defer file.Close()

	// encode the barcode as png
	png.Encode(file, qrCode)

	return file, qrCode, err
}
