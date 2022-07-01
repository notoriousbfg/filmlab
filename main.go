package main

import (
	"log"

	"github.com/disintegration/imaging"
	flag "github.com/ogier/pflag"
)

func main() {
	srcFile := flag.String("src", "", "the source file")
	outFile := flag.String("out", "./inverted.jpg", "the outputted file")
	profile := flag.String("profile", "", "the film profile")
	preset := flag.String("preset", "mid", "the brightness preset: light,dark,mid")
	flag.Parse()

	src, err := imaging.Open(*srcFile)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	inverted := imaging.Invert(src)

	if len(*profile) > 0 {
		profiles := map[string]Profile{
			"Portra160": Portra160{},
			"HP5Plus":   HP5Plus{},
			"ColorPlus": ColorPlus{},
		}

		imageProfile, profileExists := profiles[*profile]
		if profileExists {
			inverted = imageProfile.Adjust(inverted, *preset)
		} else {
			log.Fatal("profile not found")
		}
	}

	err = imaging.Save(inverted, *outFile)
	if err != nil {
		log.Fatalf("failed to save image: %v", err)
	}

	log.Printf("image generated")
}

func adjustColourValue(colour int, amount int) int {
	colour = colour + amount

	if colour > 255 {
		return 255
	} else if colour < 0 {
		return 0
	}

	return colour
}
