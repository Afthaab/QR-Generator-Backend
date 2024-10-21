package service

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

func QrCodeGen(t string, filename string) (*os.File, barcode.Barcode, error) {
	// Create the barcode
	qrCode, err := qr.Encode(t, qr.M, qr.Auto)
	if err != nil {
		return nil, nil, fmt.Errorf("could not generate the qr code: %w", err)
	}

	// Scale the barcode to 2000x2000 pixels
	qrCode, err = barcode.Scale(qrCode, 2000, 2000)
	if err != nil {
		return nil, nil, fmt.Errorf("could not scale the qr code: %w", err)
	}

	// Define the output directory
	outputDir := filepath.Join("QR-Generator-UI", "QR_Codes")

	// Ensure the directory exists
	err = os.MkdirAll(outputDir, os.ModePerm)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create directory: %w", err)
	}

	// Create the output file in the specified directory
	filePath := filepath.Join(outputDir, filename+".png")
	file, err := os.Create(filePath)
	if err != nil {
		return nil, nil, fmt.Errorf("could not create the file: %w", err)
	}

	// Ensure the file is closed after writing
	defer file.Close()

	// Encode the barcode as PNG
	if err = png.Encode(file, qrCode); err != nil {
		return nil, nil, fmt.Errorf("could not encode the png: %w", err)
	}

	return file, qrCode, nil
}
