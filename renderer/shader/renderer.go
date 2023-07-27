package shader

import (
	"bytes"
	"fmt"
	"image"
	"image/png"
	"math"
	"os"
	"strings"

	lua "github.com/yuin/gopher-lua"
	"github.com/yuin/gopher-lua/parse"
)

var (
	defaultModuleSize = 5
	//defaultBackgroundColor = color.White
)

type Options struct {
	ModuleSize int
	//BackgroundColor color.Color
}

func Save(data []byte, shader string, filename string, opts *Options) error {
	img, err := Render(data, shader, opts)
	if err != nil {
		return err
	}

	w := bytes.NewBuffer(nil)

	errEncode := png.Encode(w, img)
	if errEncode != nil {
		return errEncode
	}

	return os.WriteFile(filename, w.Bytes(), 0644)
}

func Render(data []byte, shader string, opts *Options) (image.Image, error) {
	o := &Options{
		ModuleSize: defaultModuleSize,
		//BackgroundColor: defaultBackgroundColor,
	}
	if opts != nil {
		if o.ModuleSize > 0 {
			o.ModuleSize = opts.ModuleSize
		}
		//if opts.BackgroundColor != nil {
		//	o.BackgroundColor = opts.BackgroundColor
		//}
	}
	size := int(math.Sqrt(float64(len(data))))

	img := image.NewRGBA(image.Rect(0, 0, (size+8)*o.ModuleSize, (size+8)*o.ModuleSize))

	//draw.Draw(img, img.Bounds(), &image.Uniform{C: o.BackgroundColor}, image.Point{}, draw.Src)

	chunk, errParse := parse.Parse(strings.NewReader(shader), "shader.lua")
	if errParse != nil {
		return nil, fmt.Errorf("error parse shader, %w", errParse)
	}

	proto, errCompile := lua.Compile(chunk, "shader.lua")
	if errCompile != nil {
		return nil, fmt.Errorf("error compile shader, %w", errCompile)
	}

	state := lua.NewState()
	defer state.Close()

	errFunc := state.NewFunction(func(_ *lua.LState) int { return 1 })

	lFunc := state.NewFunctionFromProto(proto)
	state.Push(lFunc)
	errDo := state.PCall(0, lua.MultRet, errFunc)
	if errDo != nil {
		return nil, fmt.Errorf("error execute shader, %w", errDo)
	}

	f := state.Get(1)
	if f.Type() != lua.LTFunction {
		return nil, fmt.Errorf("shader must return a function, got %s", f.Type().String())
	}

	ctx := state.NewTable()

	state.Push(f)
	state.Push(ctx)
	errDoShader := state.PCall(1, lua.MultRet, errFunc)
	if errDoShader != nil {
		panic(errDoShader)
	}

	//var x, y int
	//for _, dotValue := range data {
	//	state := lua.NewState()
	//
	//	state.SetGlobal("get", state.NewFunction(func(state *lua.LState) int {
	//		x := state.CheckInt(1)
	//		y := state.CheckInt(2)
	//
	//		if x < 0 || y < 0 || x >= size || y >= size {
	//			state.Push(lua.LNumber(0))
	//			return 1
	//		}
	//
	//		if data[y*size+x] == 1 {
	//			state.Push(lua.LNumber(1))
	//			return 1
	//		}
	//
	//		state.Push(lua.LNumber(0))
	//
	//		return 1
	//	}))
	//
	//	currentModule := state.NewTable()
	//	currentModule.RawSetString("dotSize", lua.LNumber(r.dotSize))
	//	currentModule.RawSetString("x", lua.LNumber(x))
	//	currentModule.RawSetString("size", lua.LNumber(size))
	//	currentModule.RawSetString("y", lua.LNumber(y))
	//	currentModule.RawSetString("isSet", lua.LBool(dotValue == 1))
	//
	//	state.SetGlobal("currentModule", currentModule)
	//
	//	canvas := state.NewTable()
	//	canvas.RawSetString("set", state.NewFunction(func(state *lua.LState) int {
	//		x := state.CheckInt(1)
	//		y := state.CheckInt(2)
	//		c := state.CheckString(3)
	//
	//		x += 4 * r.dotSize
	//		y += 4 * r.dotSize
	//
	//		cc, err := strconv.ParseInt(c, 16, 64)
	//		if err != nil {
	//			state.RaiseError("error parse color, %v", err)
	//			return 0
	//		}
	//
	//		ccc := color.RGBA{R: uint8(cc >> 24), G: uint8(cc >> 16), B: uint8(cc >> 8), A: uint8(cc)}
	//
	//		img.Set(x, y, ccc)
	//
	//		return 0
	//	}))
	//
	//	canvas.RawSetString("rect", state.NewFunction(func(state *lua.LState) int {
	//		startX := state.CheckInt(1)
	//		startY := state.CheckInt(2)
	//		ww := state.CheckInt(3)
	//		h := state.CheckInt(4)
	//		c := state.CheckString(5)
	//
	//		startX += 4 * r.dotSize
	//		startY += 4 * r.dotSize
	//
	//		cc, err := strconv.ParseInt(c, 16, 64)
	//		if err != nil {
	//			state.RaiseError("error parse color, %v", err)
	//			return 0
	//		}
	//
	//		ccc := color.RGBA{R: uint8(cc >> 24), G: uint8(cc >> 16), B: uint8(cc >> 8), A: uint8(cc)}
	//
	//		draw.Draw(img, image.Rect(startX, startY, startX+ww, startY+h), &image.Uniform{C: ccc}, image.Point{}, draw.Src)
	//
	//		return 0
	//	}))
	//	state.SetGlobal("canvas", canvas)
	//
	//	lFunc := state.NewFunctionFromProto(proto)
	//	state.Push(lFunc)
	//	errDo := state.PCall(0, lua.MultRet, nil)
	//	if errDo != nil {
	//		fmt.Printf("error execute lua plugin, %v\n", errDo)
	//		state.Close()
	//		return nil, errDo
	//	}
	//	state.Close()
	//	x++
	//	if x >= size {
	//		x = 0
	//		y++
	//	}
	//}

	return img, nil
}
