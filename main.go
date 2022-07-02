package main

import (
	"image"
	"log"
	"sync"

	"github.com/disintegration/imaging"
	flag "github.com/ogier/pflag"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg"
)

func main() {
	srcFile := flag.String("src", "", "the source file")
	outFile := flag.String("out", "./inverted.jpg", "the outputted file")
	profile := flag.String("profile", "", "the film profile")
	preset := flag.String("preset", "mid", "the brightness preset: light,dark,mid")
	histogram := flag.Bool("histogram", false, "whether to generate a histogram")
	flag.Parse()

	src, err := imaging.Open(*srcFile)
	if err != nil {
		log.Fatalf("failed to open image: %v", err)
	}

	if *histogram {
		var wg sync.WaitGroup
		wg.Add(1)

		go func() {
			defer wg.Done()
			BuildHistogram(src)
		}()

		wg.Wait()
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

func BuildHistogram(image image.Image) {
	histValues := imaging.Histogram(image)
	var values plotter.Values
	for _, i := range histValues {
		values = append(values, i)
	}
	histPlot(values)
}

func histPlot(values plotter.Values) {
	p := plot.New()

	p.Title.Text = "histogram plot"

	hist, err := plotter.NewHist(values, 20)
	if err != nil {
		panic(err)
	}

	p.Add(hist)

	if err := p.Save(3*vg.Inch, 3*vg.Inch, "hist.png"); err != nil {
		panic(err)
	}
}
