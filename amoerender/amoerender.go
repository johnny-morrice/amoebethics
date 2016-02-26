package main

import (
    "bufio"
    "flag"
    "fmt"
    "image"
    "image/color"
    "image/draw"
    "image/png"
    "math"
    "os"
    "runtime"
    "strings"
    "sync"
    d2d "github.com/llgcode/draw2d"
    dimg "github.com/llgcode/draw2d/draw2dimg"
    core "github.com/johnny-morrice/amoebethics/libamoebethics"
    ani "github.com/johnny-morrice/amoebethics/libamoenator"
)

type args struct {
    w uint
    h uint
    jobs uint
}

func main() {
    params := readArgs()

    pool := newGoPool(params.jobs)

    sc := bufio.NewScanner(os.Stdin)
    i := 0
    pool.run()
    for sc.Scan() {
        frdat := sc.Text()
        pool.add(func() {
            compile(params, frdat, i)
        })
        i++
    }

    pool.wait()
    e := sc.Err()
    if e != nil {
        fmt.Fprintf(os.Stderr, "Read warning: %v\n", e)
    }
}

type goPool struct {
    hold sync.WaitGroup
    work chan func()
    ready chan bool
}

func newGoPool(jobs uint) *goPool {
    pool := &goPool{}
    pool.ready = make(chan bool, jobs)
    pool.work = make(chan func())
    for i := uint(0); i < jobs; i++ {
        pool.ready<- true
    }
    return pool
}

func (pool *goPool) run() {
    go func() {
        for range pool.ready {
            f := <-pool.work
            go func() {
                f()
                pool.ready<- true
            }()
        }
    }()
}

func (pool *goPool) add(f func()) {
    pool.hold.Add(1)
    pool.work<- func() {
        f()
        pool.hold.Done()
    }
}

func (pool *goPool) wait() {
    pool.hold.Wait()
    close(pool.work)
    close(pool.ready)
}

func readArgs() args {
    defjobs := uint(runtime.NumCPU())
    params := args{}
    flag.UintVar(&params.w, "width", 1920, "width")
    flag.UintVar(&params.h, "height", 1080, "height")
    flag.UintVar(&params.jobs, "jobs", defjobs, "jobs")
    flag.Parse()
    return params
}

func compile(params args, frdat string, fnum int) {
    r := strings.NewReader(frdat)

    fr, rerr := ani.ReadFrame(r)
    if rerr != nil {
        fmt.Fprintf(os.Stderr,
            "Error reading frame %v: %v\nFrame follows\n%v\n",
            fnum, rerr, frdat)
        os.Exit(1)
    }

    img := render(fr, fnum, params)

    fname := fmt.Sprintf("%06d.png", fnum)
    file, ferr := os.Create(fname)

    if ferr != nil {
        fmt.Fprintf(os.Stderr,
            "Error creating output png (%v): %v\n", fname, ferr)
        os.Exit(1)
    }

    defer file.Close()

    encerr := png.Encode(file, img)

    if encerr != nil {
        fmt.Fprintf(os.Stderr,
            "Error encoding PNG: %v", encerr)
    }
}

func render(fr ani.Frame, fnum int, params args) image.Image {
    w, h := int(params.w), int(params.h)
    tr := trans(w, h, fr.Torus)

    // Blank white image
    bounds := image.Rect(0,0, w, h)
    dest := image.NewRGBA(bounds)
    draw.Draw(dest, bounds, &image.Uniform{white}, image.ZP, draw.Src)

    shapes := []string {
        "circle",
        "square",
    }
    plts := []plot {}

    for _, sh := range shapes {
        plts = append(plts, plotShape(sh, fr))
    }

    plts = append(plts, plotExplode(fr))

    gc := dimg.NewGraphicContext(dest)

    gc.SetStrokeColor(black)
    gc.SetFillColor(black)

    gc.SetMatrixTransform(tr)

    for _, plot := range plts {
        plot.draw(gc)
    }

    return dest
}

type plot interface {
    draw(gc d2d.GraphicContext)
}

type exploder struct {
    bombs []ani.Explosion
    palette []ani.Color
}

var _ plot = exploder{}

func plotExplode(fr ani.Frame) exploder {
    ex := exploder{}
    ex.bombs = fr.Explosions
    ex.palette = fr.Palette
    return ex
}

func (ex exploder) draw(gc d2d.GraphicContext) {
    for _, b := range ex.bombs {
        c := ex.palette[b.Color]
        alpha := b.Intensity / 2
        trans := color.RGBA{c.R, c.G, c.B, alpha}
        gc.SetStrokeColor(trans)
        gc.SetFillColor(trans)

        lineCircle(gc, b.P.X, b.P.Y, b.Radius)

        fmt.Printf("Color was %v: %v\n", b.Color, trans)
        gc.Close()
        gc.Fill()
    }
}

type boxplot struct {
    boxes []ani.ColorBox
    palette []ani.Color
}

type circleplot struct {
    boxplot
}
var _ plot = circleplot{}

var black = color.RGBA{0, 0, 0, 0xff}
var white = color.RGBA{0xff, 0xff, 0xff, 0xff}
var clear = color.RGBA{0, 0, 0, 0}

func lineCircle(gc d2d.GraphicContext, cx, cy, r float64) {
    const blowup = 100.0
    step := 1.0 / (r * blowup)

    gc.MoveTo(cx + r, cy)

    max := math.Pi * 2
    for t := 0.0; t < max; t += step {
        x := cx + (math.Cos(t) * r)
        y := cy + (math.Sin(t) * r)
        gc.LineTo(x, y)
    }
}

func (c circleplot) draw(gc d2d.GraphicContext) {
    for _, cb := range c.boxes {
        gc.SetStrokeColor(black)
        gc.SetFillColor(clear)
        gc.SetLineWidth(1.0)

        lineCircle(gc, cb.P.X, cb.P.Y, cb.Radius)

        gc.Close()
        gc.FillStroke()

    }
}

type squareplot struct {
    boxplot
}
var _ plot = squareplot{}

func (c squareplot) draw(gc d2d.GraphicContext) {
    for _, cb := range c.boxes {
        gc.SetStrokeColor(black)
        gc.SetFillColor(clear)
        gc.SetLineWidth(1.0)

        side := cb.Radius / 2.0
        xmin, ymin := cb.P.X - side, cb.P.Y - side
        xmax, ymax := cb.P.X + side, cb.P.Y + side

        gc.MoveTo(xmin, ymin)
        gc.LineTo(xmin, ymax)
        gc.LineTo(xmax, ymax)
        gc.LineTo(xmax, ymin)
        gc.LineTo(xmin, ymin)
        gc.Close()
        gc.FillStroke()
    }
}

func plotShape(sh string, fr ani.Frame) plot {
    boxp := boxplot{}
    boxp.boxes = fr.Shapes[sh]
    boxp.palette = fr.Palette

    switch sh {
    case "circle":
        return circleplot{boxplot: boxp}
    case "square":
        return squareplot{boxplot: boxp}
    default:
        panic(fmt.Sprintf("Unknown shape type: %v", sh))
        return nil
    }
}

func trans(w, h int, t core.Torus) d2d.Matrix {
    fw, fh := float64(w), float64(h)
    left := -t.W / 2
    right := t.W / 2
    bottom := - t.H / 2
    top := t.H / 2
    from := [4]float64{left, bottom, right, top}
    to := [4]float64{0, fh, fw, 0}
    return d2d.NewMatrixFromRects(from, to)
}