package main

import (
	"fmt"
	"log"
	"os"
	"sinta/internal/config"
	"sinta/internal/pdf"
	"sinta/internal/sinta"
	"strings"
)

func main() {
	err := config.LoadEnv()
	if err != nil {
		log.Fatalf("error loading .env file: %v", err)
	}

	unipdfKey := config.GetEnvVar("UNIPDF_KEY")
	if unipdfKey == "" {
		log.Fatal("UNIPDF_KEY belum diisi di .env file")
	}

	err = config.SetUnipdfKey(unipdfKey)
	if err != nil {
		log.Fatalf("gagal setting UNIPDF_KEY: %v", err)
	}

	files, err := os.ReadDir("./pdf")
	if err != nil {
		log.Fatalf("error baca directory: %v\n", err)
	}

	if len(files) == 0 {
		fmt.Println("tidak terdapat file dalam folder pdf, program dihentikan")
		return
	}

	for i, file := range files {
		if strings.HasSuffix(file.Name(), ".pdf") {
			filePath := "./pdf/" + file.Name()
			fmt.Printf("%s%d. proses file: %s\n", func() string {
				if i == 0 {
					return ""
				} else {
					return "\n"
				}
			}(), i+1, filePath)

			issnNumbers, err := pdf.ExtractISSNNumbers(filePath)
			if err != nil {
				fmt.Printf("error extract issn dari %s: %v\n", filePath, err)
				continue
			}

			if len(issnNumbers) == 0 {
				fmt.Printf("tidek ditemukan ISSN valid di file %s\n", filePath)
				continue
			}

			hasSinta := false
			var accreditation string
			for _, issn := range issnNumbers {
				fmt.Printf("cek status untuk ISSN: %s\n", issn)
				accreditation = sinta.CheckSintaStatus(issn)
				if accreditation != "" {
					hasSinta = true
					break
				}
			}

			if hasSinta {
				newFileName := fmt.Sprintf("./pdf/%s - %s.pdf", strings.TrimLeft(accreditation, " "), strings.TrimSuffix(file.Name(), ".pdf"))
				err := os.Rename(filePath, newFileName)
				if err != nil {
					log.Printf("error ganti nama file %s jadi %s: %v\n", filePath, newFileName, err)
				} else {
					fmt.Printf("file sudah diganti nama jadi: %s\n", newFileName)
				}
			} else {
				fmt.Printf("file tidak terdaftar di SINTA: %s. File tetap disimpan tanpa perubahan nama\n", filePath)
			}
		}
	}
}
