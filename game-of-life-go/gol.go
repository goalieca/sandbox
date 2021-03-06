package main

// the package import list
import (
    "fmt"
    "os"
    "image"
    "image/color"
    "image/png"
)

//creates a board type
type Board struct {
    Image *image.RGBA
    evolution int
}

func NewBoard(r image.Rectangle) *Board {
    //initialize and return reference
    board := &Board{ image.NewRGBA(r), 0 }

    //iterate over image and initialize data
    //note: no range over 2D matrices,Image type
    for y := 0; y <r.Dy(); y++ {
        for x := 0; x <r.Dx(); x++ {
            //explicit casting
            var gray_x = uint8(x * 255 / r.Dx())
            var gray_y = uint8(y * 255 / r.Dy())
            var gray = uint8((x*255/r.Dx() + y*255/r.Dy())/2)
            //create color struct
            var col = color.RGBA{gray_x,gray_y,gray,255}
            //assignment, note no -> operator in go. 
            // . and -> knows which to do
            board.Image.Set(x,y,col)
        }
    }

    return board
}

func min(a,b int) int {
    if a < b {
        return a
    }
    return b
}

func max(a,b int) int {
    if a > b {
        return a
    }
    return b
}

//Evolve "member" for Board
func (b *Board) Evolve() {
    outImage := image.NewRGBA(b.Image.Rect)

    for y := 0; y < b.Image.Rect.Dy(); y++ {
        for x := 0; x < b.Image.Rect.Dx(); x++ {
            col := b.Image.At(x,y).(color.RGBA)

            var sumRed, sumGreen, sumBlue = 0, 0, 0

            //check neighbours
            for dy := -1; dy <= 1; dy++ {
                for dx := -1; dx <= 1; dx++ {
                    //returns black if not in bounds
                    col_n := b.Image.At(x+dx,y+dy).(color.RGBA)
                    sumRed += int(col_n.R)
                    sumGreen += int(col_n.G)
                    sumBlue += int(col_n.B)
                }
            }

            //too much casting, ewww
            if sumGreen > 255*3 {
                col.R = byte(max(255,int(col.G) + sumBlue/8))
            } else {
                col.R = byte(min(255,int(col.G) + sumBlue/8))
            }

            if sumBlue > 255*3 {
                col.G = byte(max(255,int(col.B) + sumRed/8))
            } else {
                col.G = byte(min(255,int(col.B) + sumRed/8))
            }

            if sumRed > 255*3 {
                col.B = byte(max(255,int(col.R) + sumGreen/8))
            } else {
                col.B = byte(min(255,int(col.R) + sumGreen/8))
            }

            //setter (assume checks bounds)
            outImage.Set(x,y,col)
        }
    }

    //replace old image with new generation
    b.Image = outImage

    b.evolution++
}

func (b* Board) Write() {
    var fileName = fmt.Sprintf("./%000d.png",b.evolution)

    //open writer and save image
    outFile, err := os.OpenFile(fileName, os.O_CREATE | os.O_WRONLY, 0666)
    if err != nil {
        panic(fmt.Sprintf("%v\n",err))
    }
    defer outFile.Close()

    err2 := png.Encode(outFile,b.Image)
    if err2 != nil {
        panic(fmt.Sprintf("%v\n",err2))
    }
}

func main() {

    const (
        height = 512
        width = 512
        generations = 100
    )

    b := NewBoard(image.Rect(0,0,width,height))
    for i := 0 ; i < generations ; i++ {
        b.Evolve()
    }

    b.Write()

}
