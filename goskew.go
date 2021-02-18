package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/docopt/docopt-go"
)

func sq(val float64) float64 {
	return val * val
}

func triangleError(base float64, left float64, right float64, ref float64) float64 {
	// calclate the height of the reference equilateral
	height := math.Sin(math.Pi/3) * ref

	// calculate the left angle of the printed triangle
	alpha := math.Acos((sq(left) + sq(base) - sq(right)) / (2 * left * base))

	// calculate the X distance to the point at <height> along the projected left side.
	delta := height / math.Tan(alpha)

	// calculate the error in x.
	xerr := delta - ref/2

	// return the xytan error.
	return xerr / height
}

func getCoord(in *float64, reg *regexp.Regexp, line string) {
	search := reg.FindStringSubmatch(line)

	if len(search) > 0 {
		*in, _ = strconv.ParseFloat(search[1], 64)
	}
}

func skew(input []byte, xytan float64, xztan float64, yztan float64) string {
	lines := strings.Split(string(input), "\n")

	// init the coords
	xin, yin, zin := 0.0, 0.0, 0.0

	// init the regular expressions
	xreg, _ := regexp.Compile(`[xX](-?\d*\.\d*)`)
	yreg, _ := regexp.Compile(`[yY](-?\d*\.\d*)`)
	zreg, _ := regexp.Compile(`[zZ](-?\d*\.\d*)`)

	for i, line := range lines {
		gmatch, _ := regexp.MatchString(`G[0-1]`, line)

		if gmatch {
			fmt.Printf(".")

			// find X, Y, and Y coords in line
			getCoord(&xin, xreg, line)
			getCoord(&yin, yreg, line)
			getCoord(&zin, zreg, line)

			// skew the X and Y. Z is unchanged to avoid squishing layers.
			xout := fmt.Sprintf("%.3f", (xin-yin*xytan)-zin*xztan)
			yout := fmt.Sprintf("%.3f", yin-zin*yztan)

			// replace the X and Y coords
			line = xreg.ReplaceAllString(line, "X"+xout)
			line = yreg.ReplaceAllString(line, "Y"+yout)

			lines[i] = line
		}
	}
	fmt.Printf("\n")
	return strings.Join(lines, "\n")
}

func main() {
	usage := `Go Skew.

Usage:
  %basename% err <xy> <xz> <yz> [--output=FILE] <file>
  %basename% tri <base> <left> <right> [--xz=ERROR --yz=ERROR --ref=LENGTH --output=FILE] [<file>]
  %basename% -h | --help

Options:	
  -o FILE, --output=FILE    The file name to write out to, by default Go Skew overwrites the original file. 
  --xz=ERROR                The error tangent in the XZ axis.
  --yz=ERROR                The error tangent in the YZ axis.       
  -r LENGTH, --ref=LENGTH   The side length of the reference equilateral triangle in millimeters. 100 by default.
  -h, --help
`
	basename := filepath.Base(os.Args[0])

	usage = strings.ReplaceAll(usage, "%basename%", basename)

	opts, _ := docopt.ParseDoc(usage)
	//fmt.Println(opts)

	// The tan error in all planes
	xy, xz, yz := 0.0, 0.0, 0.0

	// are we in triangle mode?
	triangle, _ := opts.Bool("tri")

	if triangle {
		fmt.Println("calculating xytan from given triangle")
		// extract options
		ref, _ := opts.Float64("--ref")
		base, _ := opts.Float64("<base>")
		left, _ := opts.Float64("<left>")
		right, _ := opts.Float64("<right>")

		// ref is 100mm by default
		if ref == 0 {
			ref = 100
		}

		// calculate the xytan
		xy = triangleError(base, left, right, ref)

		// get the xz and yz errors if given
		xz, _ = opts.Float64("--xz")
		yz, _ = opts.Float64("--yz")
	} else {
		// get the tan errors.
		xy, _ = opts.Float64("<xy>")
		xz, _ = opts.Float64("<xz>")
		yz, _ = opts.Float64("<yz>")
	}

	fmt.Printf("Error tangents:\nxytan: %0.7f, xztan: %0.7f, yztan: %0.7f\n", xy, xz, yz)

	iFile, _ := opts.String("<file>")
	oFile, _ := opts.String("--output")

	if iFile == "" {
		return
	}

	if oFile == "" {
		oFile = iFile
	}

	input, err := ioutil.ReadFile(iFile)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Printf("skewing")

	output := skew(input, xy, xz, yz)

	err = ioutil.WriteFile(oFile, []byte(output), 0644)

	if err != nil {
		log.Fatalln(err)
	}

	fmt.Println("Done!")
}
