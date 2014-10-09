package music

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

var audioFile = flag.String("out", "", "file to put the audio to (must be .aiff)")
var scoreFile = flag.String("score", "", "file to put the score to (must be .scd)")
var writeSynthDefs = flag.Bool("write-synthdefs", false, "write synthdefs in scserver mode")
var loadSamples = flag.Bool("load-samples", false, "load samples in scserver mode")

var _ generator = &Stage{}

func New() *Stage {
	s := &Stage{
		instrNumber:      2000,
		sampleInstNumber: 0,
		synthdefs:        map[string][]byte{},
		AudioFile:        "",
		busNumber:        16,
		groupNumber:      1012,
		groups:           map[int]int{},
	}

	flag.Parse()
	if audioFile != nil {
		s.AudioFile = *audioFile
	}
	if scoreFile != nil {
		s.ScoreFile = *scoreFile
	}
	if writeSynthDefs != nil {
		s.WriteSynthDefs = *writeSynthDefs
	}
	if loadSamples != nil {
		s.LoadSamples = *loadSamples
	}
	return s
}

func (s *Stage) newBusId() int {
	s.busNumber++
	return s.busNumber
}

func (s *Stage) newGroupId() int {
	s.groupNumber++
	return s.groupNumber
}

func (s *Stage) newNodeId() int {
	s.instrNumber++
	return s.instrNumber
}

func (s *Stage) newSampleBuffer() int {
	s.sampleInstNumber++
	return s.sampleInstNumber
}

type Stage struct {
	buffer           *bytes.Buffer
	eventbuffer      *bytes.Buffer
	synthdefs        map[string][]byte
	instrNumber      int
	sampleInstNumber int
	AudioFile        string
	ScoreFile        string
	busNumber        int
	WriteSynthDefs   bool
	LoadSamples      bool
	groupNumber      int
	scServerOnline   bool
	tracks           []*Track
	synthDefDirs     []string
	instruments      []instrument
	groups           map[int]int
}

func (s *Stage) Track(bar string, tempo Tempo) *Track {
	tr := newTrack(tempo, M(bar))
	s.tracks = append(s.tracks, tr)
	return tr
}

func (s *Stage) Instrument(name string, dir string, numVoices int) []*Voice {
	if numVoices < 1 {
		panic("minimum for numVoices: 1")
	}
	path := filepath.Join(dir, name+".scd")
	// fmt.Println("instrument", path)
	vs := newSCInstrument(s, name, path, numVoices)
	s.instruments = append(s.instruments, vs[0].instrument)
	return vs
}

func (s *Stage) SampleInstrument(instrument string, sampleLib SampleLibrary, numVoices int) []*Voice {
	vs := newSCSampleInstrument(s, instrument, sampleLib, numVoices)
	s.instruments = append(s.instruments, vs[0].instrument)
	return vs
}

func (s *Stage) Route(name, dir string, numVoices int) []*Voice {
	if numVoices < 1 {
		panic("minimum for numVoices: 1")
	}
	path := filepath.Join(dir, name+".scd")
	// fmt.Println("route", path)
	vs := newRoute(s, name, path, numVoices)
	s.instruments = append(s.instruments, vs[0].instrument)
	return vs
}

func (s *Stage) Sample(path string, numVoices int) []*Voice {
	if numVoices < 1 {
		panic("minimum for numVoices: 1")
	}
	vs := newSCSample(s, path, numVoices)
	s.instruments = append(s.instruments, vs[0].instrument)
	return vs
}

func (s *Stage) SampleFreq(path string, freq float64, numVoices int) []*Voice {
	if numVoices < 1 {
		panic("minimum for numVoices: 1")
	}
	vs := newSCSampleFreq(s, path, freq, numVoices)
	s.instruments = append(s.instruments, vs[0].instrument)
	return vs
}

func (s *Stage) Group(parent int) *Voice {
	v := newGroup(s)
	s.groups[v.Group] = parent
	return v
}

func (s *Stage) Bus(name string, numchannels int) *Voice {
	v := newBus(s, name)
	s.busNumber += numchannels
	return v
}

var home = os.Getenv("HOME")
var SynthDefPool = filepath.Join(home, ".local/share/SuperCollider/quarks/SynthDefPool/pool")

func (s *Stage) writeSynthDefs(w io.Writer) {
	for _, instr := range s.instruments {
		switch t := instr.(type) {
		case *sCInstrument:
			if t.IsUsed() {
				sdef := t.LoadCode()
				fmt.Fprintf(w, strings.TrimSpace(string(sdef))+".writeDefFile;"+"\n")
			}
		}
	}
}

var sampleSynthDef = `SynthDef("sample%d", { |gate=1,bufnum = 0,amp=1, out=0, pan=0, rate=1| var z;
	z =  EnvGen.kr(Env.perc,gate) * PlayBuf.ar(%d, bufnum, BufRateScale.kr(bufnum) * rate, loop: 0, doneAction: 2);
	FreeSelfWhenDone.kr(z);
	Out.ar(out, Pan2.ar(z, pos: pan, level: amp));
} ).writeDefFile;`

func (s *Stage) writeLoadSamples(w io.Writer) {
	channelPlayers := map[uint]struct{}{}

	for _, instr := range s.instruments {
		switch t := instr.(type) {
		case *sCSample:
			if t.IsUsed() {
				if _, has := channelPlayers[t.Sample.Channels]; !has {
					channelPlayers[t.Sample.Channels] = struct{}{}
					fmt.Fprintf(w, strings.TrimSpace(sampleSynthDef)+"\n", t.Sample.Channels, t.Sample.Channels)
				}
			}
		case *sCSampleInstrument:
			if len(t.Samples) > 0 {
				for _, ch := range t.Channels() {
					if _, has := channelPlayers[uint(ch)]; !has {
						channelPlayers[uint(ch)] = struct{}{}
						fmt.Fprintf(w, strings.TrimSpace(sampleSynthDef)+"\n", ch, ch)
					}
				}
			}
		}
	}

}

type eventWriterOptions struct {
	startOffset  uint
	startTick    uint
	tickNegative int
	ticksSorted  []int
	finTick      uint
	tickMapped   map[int][]*Event
}

func (s *Stage) writeEvents(w io.Writer, opts eventWriterOptions) (skipSecs float32) {
	t := 0
	// withStartTick := 1.0
	skipSecs = float32(0.0)

	beginOffset := float32(opts.startOffset) / float32(1000)

	if opts.startTick != 0 {
		skipSecs = getSeconds(int(opts.startTick), opts.tickNegative, beginOffset)
		// skipSecs = tickToSeconds(int(startTick)+(tickNegative*(-1))) + 0.000001 + beginOffset
	}

	_ = skipSecs

	for _, ti := range opts.ticksSorted {
		if opts.finTick != 0 && int(opts.finTick) <= ti {
			t = int(opts.finTick)
			break
		}
		inSecs := getSeconds(ti, opts.tickNegative, beginOffset)

		// inSecs := tickToSeconds(ti+(tickNegative*(-1))) + 0.000001 + beginOffset

		var tickBf bytes.Buffer

		// if len(opts.tickMapped[ti]) > 0 {
		for _, ev := range opts.tickMapped[ti] {
			code := ev.Voice.getCode(ev)
			if code != "" {
				tickBf.WriteString(code)
				// fmt.Fprintf(w, code)
			}
			// ev.Runner(ev)
			//fmt.Fprintf(w, ev.sccode.String())
		}
		if tbf := tickBf.String(); tbf != "" {
			fmt.Fprintf(w, `  [%0.6f%s`, inSecs, tbf)
			fmt.Fprintf(w, "],\n")
		}
		// }
		t = ti
	}

	fmt.Fprintf(w, "  [%0.6f, [\\g_deepFree, 1], [\\c_set, 0, 0]]];\n", float32(t)/float32(1000000000))
	return
}

func (s *Stage) writeAtPosZero(w io.Writer) {
	// fmt.Printf("used samples: %#v\n", s.usedSamples)
	fmt.Fprintf(w, `  [%0.6f, `, 0.0)

	// create the bus routing group
	fmt.Fprintf(w, fmt.Sprintf("\n"+`[\g_new, %d, 0, 0],`, 1200))
	// create the instruments group
	fmt.Fprintf(w, fmt.Sprintf("\n"+`[\g_new, %d, 0, 0],`, 1010))

	for groupId, groupParent := range s.groups {
		fmt.Fprintf(w, "\n"+`[\g_new, %d, 1, %d], `, groupId, groupParent)
	}

	first_sample := true

	for _, instr := range s.instruments {
		switch t := instr.(type) {
		case *sCSample:
			if t.IsUsed() {
				if !first_sample {
					fmt.Fprintf(w, ", ")
				}
				first_sample = false
				fmt.Fprintf(w, fmt.Sprintf("\n"+`[\b_allocRead, %d, "%s"]`, t.Sample.sCBuffer, t.Sample.Path))
			}
		case *sCSampleInstrument:
			for _, sample := range t.Samples {
				if !first_sample {
					fmt.Fprintf(w, ", ")
				}
				first_sample = false
				fmt.Fprintf(w, fmt.Sprintf("\n"+`[\b_allocRead, %d, "%s"]`, sample.sCBuffer, sample.Path))
			}
		}

	}
	fmt.Fprintf(w, "],\n")

}

// startOffset is in milliseconds and must be positive
func (s *Stage) Play(startOffset uint) {

	evts := []*Event{}

	for _, tr := range s.tracks {
		tr.compile()
		evts = append(evts, tr.Events...)
	}

	sortedEvents := eventsSorted(evts)
	sort.Sort(sortedEvents)

	dir, err := ioutil.TempDir("/tmp", "go-sc-music-generator")
	if err != nil {
		panic(err.Error())
	}

	defer os.RemoveAll(dir)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for sig := range c {
			os.RemoveAll(dir)
			os.Exit(2)
			_ = sig
		}
	}()

	sclangCodeFile := s.ScoreFile
	if sclangCodeFile == "" {
		sclangCodeFile = filepath.Join(dir, "sclang-code.scd")
	}
	oscCodeFile := filepath.Join(dir, "sclang-compiled.osc")
	audioFile := s.AudioFile
	if audioFile == "" {
		audioFile = filepath.Join(dir, "out.aiff")
	}
	libraryPath := "/usr/local/share/SuperCollider/SCClassLibrary"

	tickMapped := map[int][]*Event{}
	ticksSorted := []int{}

	finTick := uint(0)
	startTick := uint(0)

	for _, ev := range sortedEvents {
		currTick := int(ev.tick)
		// custom events are only run after the big sorting in the voice
		//if ev.Type != "CUSTOM" {

		// do it before ev.runner(ev) since that sets ev.Voice.lastEvent
		if ev.type_ == "ON" && ev.Voice.lastEvent != nil {
			stopper := OffEvent(ev.Voice)
			stopper.reference = ev.Voice.lastEvent
			stopper.offset = ev.Voice.offset
			tickMapped[int(currTick)] = append(tickMapped[int(currTick)], stopper)
			freeer := newFreeEvent(ev.Voice, ev.Voice.lastEvent)
			freeer.offset = ev.Voice.offset
			freeTick := MillisecsToTick(20) + int(currTick)
			tickMapped[freeTick] = append(tickMapped[freeTick], freeer)
		}

		if ev.runner != nil {
			// fmt.Printf("nil runner: %#v\n", ev.type_)
			ev.runner(ev)
		}
		//}
		if ev.type_ == "ON" {
			currTick = MillisecsToTick(ev.offset) + currTick
		}
		if ev.type_ == "CHANGE" {
			currTick = MillisecsToTick(ev.offset) + currTick
		}

		if ev.type_ == "ON" || ev.type_ == "CHANGE" || ev.type_ == "OFF" || ev.type_ == "MUTE" || ev.type_ == "UNMUTE" || ev.type_ == "fin" || ev.type_ == "CUSTOM" {
			tickMapped[int(currTick)] = append(tickMapped[int(currTick)], ev)
		}

		if ev.type_ == "fin" {
			if finTick == 0 || finTick > ev.tick {
				finTick = ev.tick
			}
		}
		if ev.type_ == "start" {
			if startTick == 0 || startTick < ev.tick {
				startTick = ev.tick
			}
		}
	}

	var tickNegative int = 0

	for ti := range tickMapped {
		if ti < tickNegative {
			tickNegative = ti
		}
		ticksSorted = append(ticksSorted, ti)
	}

	sort.Ints(ticksSorted)

	opts := eventWriterOptions{}
	opts.startOffset = startOffset
	opts.startTick = startTick
	opts.tickNegative = tickNegative
	opts.ticksSorted = ticksSorted
	opts.finTick = finTick
	opts.tickMapped = tickMapped

	s.eventbuffer = &bytes.Buffer{}

	// we have to do the write events before the synthdefs and loadsamples because
	// we need to run everything to know which samples and synthdefs are needed
	skipSecs := s.writeEvents(s.eventbuffer, opts)

	s.buffer = &bytes.Buffer{}
	fmt.Fprintf(s.buffer, "(\n")

	s.checkForScServer()

	if s.WriteSynthDefs {
		s.writeSynthDefs(s.buffer)
	}

	if !s.scServerOnline || s.LoadSamples {
		s.writeLoadSamples(s.buffer)
	}

	fmt.Fprintf(s.buffer, "TempoClock.default.tempo = 1; \n")
	fmt.Fprintf(s.buffer, "x = [\n")
	s.writeAtPosZero(s.buffer)

	io.Copy(s.buffer, s.eventbuffer)

	// TODO change the generating code, so that the online server is reused
	if s.scServerOnline {
		println("server is online")
		fmt.Fprintf(s.buffer, "\n\nScore.play(x); )")
		err := s.runBulkScServerCode(s.buffer.String())
		if err != nil {
			panic(err)
		}
		//time.Sleep(time.Millisecond * 500)
		// time.Sleep(time.Second * 2)
		return
	}

	println("server is NOT online")
	fmt.Fprintf(s.buffer, `Score.write(x, "`+oscCodeFile+`");`+"\n")
	fmt.Fprintf(s.buffer, "\n\n"+` "quitting".postln; 0.exit; )`)
	err = ioutil.WriteFile(sclangCodeFile, s.buffer.Bytes(), 0644)
	if err != nil {
		fmt.Printf("Error: %s\n", err.Error())
		return
	}

	fmt.Printf("tempfile %s written\n", sclangCodeFile)
	now := time.Now()
	if !s.mkOSCFile(libraryPath, sclangCodeFile) {
		println("could not write osc file")
		return
	}

	// fileWriteTime := time.Since(now)

	SclangTime := time.Since(now)
	now = time.Now()
	var exportFloat bool

	if s.AudioFile == "" {
		exportFloat = true
	}

	if !s.mkAudiofile(oscCodeFile, audioFile, exportFloat) {
		println("could not write audio file")
		return
	}

	println("audio file written")
	ScsynthTime := time.Since(now)

	// fmt.Printf("Time:\nwrite file: %s\nsclang: %s\nScsynth: %s\n", fileWriteTime, SclangTime, ScsynthTime)
	fmt.Printf("Time:\nsclang: %s\nScsynth: %s\n", SclangTime, ScsynthTime)

	if s.AudioFile == "" {
		playFile(audioFile, skipSecs)
	}
}

func getSeconds(tick int, negativeOffset int, offset float32) float32 {
	return TickToSeconds(tick+(negativeOffset*(-1))) + 0.000001 + offset
}

func (s *Stage) runBulkScServerCode(code string) error {
	return s.runScServerCode(strings.Replace(code, "\n", "", -1))
}

func (s *Stage) runScServerCode(code string) error {
	res, err := http.Post("http://localhost:9999/run", "application/octet-stream", strings.NewReader(code))
	if err == nil {
		defer res.Body.Close()
		b, err2 := ioutil.ReadAll(res.Body)
		if err2 == nil {
			if string(b) == "ok" {
				return nil
			} else {
				return fmt.Errorf(string(b))
			}
		} else {
			return err2
		}
	} else {
		return err
	}
}

func (s *Stage) checkForScServer() {
	if s.runScServerCode(`"Go music script".postln;`) == nil {
		s.scServerOnline = true
	}
}

func (s *Stage) mkOSCFile(libraryPath, sclangCodeFile string) (ok bool) {
	cmd := exec.Command(
		"sclang",
		"-r",
		"-s",
		"-l",
		libraryPath,
		sclangCodeFile,
	)
	out, err := cmd.CombinedOutput()

	if err != nil {
		fmt.Println("ERROR running sclang")
		fmt.Printf("%s\n", out)
		fmt.Println(err)
		return false
	}
	return true
}

func (s *Stage) mkAudiofile(oscCodeFile, audioFile string, exportFloat bool) (ok bool) {
	// sample rate
	// channels
	// file format
	// bit depth

	//	cmd = exec.Command("scsynth", "-N", oscCodeFile, "_", audioFile, "44100", "AIFF", "int16", "-o", "2")
	// sampleFormat "int8", "int16", "int24", "int32", "mulaw", "alaw","float"
	// from http://doc.sccode.org/Classes/SoundFile.html#-sampleFormat
	// headerFormat
	// from http://doc.sccode.org/Classes/SoundFile.html#-headerFormat
	/*
	   "AIFF"	Apple/SGI AIFF format
	   "WAV","WAVE", "RIFF"	Microsoft WAV format
	   "Sun", "NeXT"	Sun/NeXT AU format
	   "SD2"	Sound Designer 2
	   "IRCAM"	Berkeley/IRCAM/CARL
	   "raw"	no header = raw data
	   "MAT4"	Matlab (tm) V4.2 / GNU Octave 2.0
	   "MAT5"	Matlab (tm) V5.0 / GNU Octave 2.1
	   "PAF"	Ensoniq PARIS file format
	   "SVX"	Amiga IFF / SVX8 / SV16 format
	   "NIST"	Sphere NIST format
	   "VOC"	VOC files
	   "W64"	Sonic Foundry's 64 bit RIFF/WAV
	   "PVF"	Portable Voice Format
	   "XI"	Fasttracker 2 Extended Instrument
	   "HTK"	HMM Tool Kit format
	   "SDS"	Midi Sample Dump Standard
	   "AVR"	Audio Visual Research
	   "FLAC"	FLAC lossless file format
	   "CAF"	Core Audio File format
	*/
	//cmd = exec.Command("scsynth", "-N", oscCodeFile, "_", audioFile, "48000", "AIFF", "int16", "-o", "2")
	//cmd = exec.Command("scsynth", "-N", oscCodeFile, "_", audioFile, "48000", "AIFF", "float", "-o", "2")

	format := "int32"
	if exportFloat {
		format = "float"
	}

	cmd := exec.Command(
		"scsynth",
		"-N",
		oscCodeFile,
		"_",
		audioFile,
		"96000",
		"AIFF",
		format,
		"-o",
		"2",
	)

	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Println("ERROR running scsynth")
		fmt.Println(err)
		fmt.Printf("%s\n", out)
		return false
	}
	return true
}

func playFile(audioFile string, skipSecs float32) (ok bool) {
	// S16_BE
	// --channels=2 --file-type raw|au|voc|wav --rate=48000 --format=S16_BE
	//cmd = exec.Command("aplay", "--rate=48000", "-f", "cdr", audioFile)
	//cmd = exec.Command("aplay", "--rate=48000", "-f", "U24_BE", audioFile)
	//cmd = exec.Command("aplay", "-f", "S16_BE", "-c2", "--rate=48000", audioFile)
	// "--start-delay=1000"

	// cmd = exec.Command("aplay", "-f", "FLOAT_BE", "-c2", "--rate=96000", audioFile)
	//cmd = exec.Command("aplay", "-f", "S32_BE", "-c2", "--rate=48000", audioFile)
	// -f S16_BE -c2 -f44100

	cmd := exec.Command(
		"play",
		"-q",
		audioFile,
		"trim",
		fmt.Sprintf(`%0.6f`, skipSecs),
	)

	if err := cmd.Run(); err != nil {
		fmt.Println("ERROR running play")
		fmt.Println(err)
		return false
	}
	return true
}
