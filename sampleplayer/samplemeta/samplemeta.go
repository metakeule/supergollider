package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/metakeule/sclang/scserver/lib/scserver"
)

func walkDir(dir string) {
	files, err3 := ioutil.ReadDir(dir)
	if err3 != nil {
		fmt.Printf("Error: %s\n", err3)
		return
	}

	fmt.Printf("\nentering %s\n", dir)

	for _, f := range files {
		if !f.IsDir() {
			ext := strings.ToLower(filepath.Ext(f.Name()))
			// if ext
			if ext == ".wav" || ext == ".aiff" {
				time.Sleep(time.Millisecond * 100)
				// fmt.Println(f.Name())
				fmt.Print(".")
				//fmt.Println(filepath.Join(wd, f.Name()))
				code := scCodeForOffset(filepath.Join(dir, f.Name()))
				res, err := http.Post("http://"+scserver.Address+scserver.RunURL, "application/octet-stream", strings.NewReader(code))
				if err != nil {
					fmt.Printf("Error: %s\n", err.Error())
					continue
				}
				d, err2 := ioutil.ReadAll(res.Body)

				if err2 != nil {
					fmt.Printf("Error: %s\n", err2.Error())
					continue
				}

				if string(d) != "ok" {
					fmt.Printf("Error: %s\n", string(d))
					continue
				}
			}
		} else {
			walkDir(filepath.Join(dir, f.Name()))
		}
	}
}

func main() {

	_, err := http.Get("http://" + scserver.Address + scserver.OutURL)
	if err != nil {
		fmt.Println("Server not running")
		go func() {
			fmt.Println("starting server")
			scserver.Run()
		}()
		time.Sleep(time.Second * 3)
	}

	wd, err2 := os.Getwd()

	if err2 != nil {
		panic(err2)
	}

	walkDir(wd)

}

func scCodeForOffset(file string) string {
	return fmt.Sprintf(
		//		`f = SoundFile.new; f.openRead(%#v); a = FloatArray.newClear(f.numFrames * f.numChannels); f.readData(a); ~time = a.abs.maxIndex / a.size * f.duration; ~time.postln;`,
		//`p = %#v.pathMatch;  p.do({ |y|  try {  f = SoundFile.new; f.openRead(y); a = FloatArray.newClear(f.numFrames * f.numChannels); f.readData(a); ~time = a.abs.maxIndex / a.size * f.duration; f.close(); o = ();  o.offset = ~time; g=File(y+".offset", "w+"); g.write(JSON.stringify(o)); g.close; } { |error| 	}; });`,

		// f.openRead("/media/usb0/philharmonic-orchestra/cello/cello_Fs4_1_pianissimo_arco-normal.wav"); a = FloatArray.newClear(f.numFrames * f.numChannels); f.readData(a); ~time = a.maxIndex() / a.size * f.duration; ~index =  a.maxIndex(); f.close(); o = ();  o.offset = ~time; o.maxAmp = a.at(~index); o.postln;

		// `try {  f = SoundFile.new; f.openRead(%#v); a = FloatArray.newClear(f.numFrames * f.numChannels); f.readData(a); ~time = a.abs.maxIndex / a.size * f.duration; f.close(); o = ();  o.offset = ~time; g=File(%#v, "w+"); g.write(JSON.stringify(o)); g.close; } { |error| 	};`,
		// `try {  f = SoundFile.new; f.openRead(%#v); a = FloatArray.newClear(f.numFrames * f.numChannels); f.readData(a); ~time = a.maxIndex() / a.size * f.duration; ~index =  a.maxIndex(); f.close(); o = ();  o.offset = ~time; o.maxAmp = a.at(~index); g=File(%#v, "w+"); g.write(JSON.stringify(o)); g.close; } { |error| 	};`,
		`try {  f = SoundFile.new; f.openRead(%#v); a = FloatArray.newClear(f.numFrames * f.numChannels); f.readData(a); ~time = a.maxIndex() / a.size * f.duration; ~index =  a.maxIndex(); f.close(); o = ();  o.offset = ~time; o.maxAmp = a.at(~index); o.channels = f.numChannels; o.sampleRate = f.sampleRate; o.headerFormat = f.headerFormat; o.duration = f.duration; o.sampleFormat = f.sampleFormat; o.numFrames = f.numFrames; g=File(%#v, "w+"); g.write(JSON.stringify(o)); g.close; } { |error| 	};`,
		file,
		file+".meta",
	)
}
