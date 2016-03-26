package main

import (
    "bufio"
    "flag"
    "fmt"
    "image"
    "image/color"
    "image/draw"
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
        buildTask := func(frnum int) func () {
            return func() {
                err := compile(frdat, frnum)
                if err != nil {
                    fmt.Fprintf(os.Stderr, "Error: %v\n", err)
                }
            }
        }
        pool.add(buildTask(i))
        i++
    }

    pool.wait()
    scanerr := sc.Err()
    if scanerr != nil {
        fmt.Fprintf(os.Stderr, "Read warning: %v\n", scanerr)
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
            f, ok := <-pool.work
            if ok {
                go func() {
                    f()
                    pool.ready<- true
                    pool.hold.Done()
                }()
            }
        }
    }()
}

func (pool *goPool) add(f func()) {
    pool.hold.Add(1)
    pool.work<- f
}

func (pool *goPool) wait() {
    pool.hold.Wait()
    close(pool.work)
    close(pool.ready)
}

func readArgs() args {
    defjobs := uint(runtime.NumCPU())
    params := args{}
    flag.UintVar(&params.jobs, "jobs", defjobs, "jobs")
    flag.Parse()
    return params
}

func compile(frdat string, fnum int) error {
    r := strings.NewReader(frdat)

    fr, rerr := ani.ReadFrame(r)
    if rerr != nil {
        return fmt.Errorf("Error reading frame %v: %v\nFrame follows\n%v\n",
            fnum, rerr, frdat)
    }

    img := render(fr, fnum)

    fname := fmt.Sprintf("%06d.png", fnum)
    err := dimg.SaveToPngFile(fname, img)

    return err
}

func render(fr ani.Frame, fnum int) image.Image {
    w, h := fr.SurfaceDims()

    // Blank white image
    bounds := image.Rect(0,0, int(w), int(h))
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
    palette ani.Palette
}

var _ plot = exploder{}

func plotExplode(fr ani.Frame) exploder {
    ex := exploder{}
    ex.bombs = fr.Explosions
    ex.palette = fr.Colors
    return ex
}

func (ex exploder) draw(gc d2d.GraphicContext) {
    for _, b := range ex.bombs {
        c := ex.palette.Get(b.Color)
//        alpha := b.Intensity / 2
        trans := color.RGBA{c.R, c.G, c.B, 255}
        gc.SetStrokeColor(trans)
        gc.SetFillColor(trans)
        gc.SetLineWidth(1.0)

        lineCircle(gc, b.P.X, b.P.Y, b.Radius)

        gc.Close()
        gc.Stroke()
    }
}

type boxplot struct {
    boxes []ani.ColorBox
    palette ani.Palette
}

type circleplot struct {
    boxplot
}
var _ plot = circleplot{}

var black = color.RGBA{0, 0, 0, 0xff}
var white = color.RGBA{0xff, 0xff, 0xff, 0xff}
var clear = color.RGBA{0, 0, 0, 0}

func lineCircle(gc d2d.GraphicContext, cx, cy, r float64) {
    count := int(1000.0 * r)

    trace := traceCircle(count, cx, cy, r)

    gc.MoveTo(cx + r, cy)

    for _, pos := range trace {
        gc.LineTo(pos.X, pos.Y)
    }
}

func (c circleplot) draw(gc d2d.GraphicContext) {
    for _, cb := range c.boxes {
        x := cb.P.X
        y := cb.P.Y
        gc.SetStrokeColor(black)
        gc.SetFillColor(clear)
        gc.SetLineWidth(1.0)

        lineCircle(gc, x, y, cb.Radius)

        gc.Close()
        gc.FillStroke()

        // Inner shades
        round := traceCircle(len(cb.Colors), x, y, cb.Radius / 2.0)

        for i, pos := range round {
            shade := cb.Colors[i]
            prim := c.palette.Get(shade)
            col := color.RGBA{prim.R, prim.G, prim.B, 255}
            gc.SetStrokeColor(col)
            gc.SetFillColor(clear)
            gc.SetLineWidth(1.0)

            lineCircle(gc, pos.X, pos.Y, cb.Radius / 10)
            gc.Close()
            gc.FillStroke()
        }
    }
}

func traceCircle(count int, cx, cy, r float64) []core.UserVec {
    trace := make([]core.UserVec, count)
    step := (math.Pi * 2) / float64(count)
    t := 0.0
    for i := 0; i < count; i ++ {
        x := cx + (math.Cos(t) * r)
        y := cy + (math.Sin(t) * r)
        trace[i] = core.UserVec{x, y}
        t += step
    }
    return trace
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
        lineSquare(gc, cb.P.X, cb.P.Y, side)
        gc.Close()
        gc.FillStroke()
    }
}

func lineSquare(gc d2d.GraphicContext, x, y, side float64) {
    xmin, ymin := x - side, y - side
    xmax, ymax := x + side, y + side

    gc.MoveTo(xmin, ymin)
    gc.LineTo(xmin, ymax)
    gc.LineTo(xmax, ymax)
    gc.LineTo(xmax, ymin)
    gc.LineTo(xmin, ymin)
}

func plotShape(sh string, fr ani.Frame) plot {
    boxp := boxplot{}
    boxp.boxes = fr.Shapes[sh]
    boxp.palette = fr.Colors

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
