# Das Sample-SynthDef Problem

     SynthDef("sample2", 
        { |
            attac=0.0000005, 
            release=0.0000005, 
            releasenode=1, 
            length=1, 
            gate=1,
            bufnum = 0,
            amp=1, 
            out=0, 
            pan=0, 
            rate=1,
            skip=0
            | 
        var z;
        z =  EnvGen.kr(
                Env([0, 1, 0], [attac,length,release],\linear, releaseNode: releasenode,doneAction: 2),
                gate
            ) * PlayBuf.ar(2, bufnum, BufRateScale.kr(bufnum) * rate, loop: 0, startPos: skip);
        FreeSelfWhenDone.kr(z);
        Out.ar(out, Pan2.ar(z, pos: pan, level: amp));
        } 
    ).writeDefFile;

Ich will einen Synthdef mit den folgenden Eigenschaften:

- nachdem der Sample abgespielt ist, wird der Synth gelöscht
- panorama-panning soll parameter sein (pan)
- gesamtlautstärke soll parameter sein (amp)
- abspielgeschwindigkeit soll einstellbar sein (rate)
- keine loop
- es soll angegeben werden können, wieviel übersprungen wird (skip), allerdings soll die abspielrate dabei berücksichtigt werden
- der ausschwingvorgang soll über gate ausgelöst werden
- attac und release sollen angegeben werden können, release soll erst über gate ausgelöst werden

also was funktioniert:

- skip ist tatsächlich unabhängig von der abspielgeschwindigkeit

was nicht funktioniert:
- dauer des asr ist

s.boot;
b = Buffer.read(s, "/home/benny/Musik/supercollider/testsample/testsample.wav"); 
c = Buffer.read(s, "/home/benny/Musik/supercollider/testsample/testsample2.wav"); 
d = Buffer.read(s, "/home/benny/Musik/supercollider/testsample/testsample3.wav");

SynthDef("sample2", 
        { |
            attac=0.0000005, 
            release=0.2, 
            gate=1,
            bufnum=0,
            amp=1, 
            out=0, 
            pan=0, 
            rate=1,
            skip=0,
            curve='linear'
            | 
        var z;
        z =  EnvGen.kr(
        //  Env.adsr, Env.asr
            Env.asr(attackTime: attac, sustainLevel: 1, releaseTime: release, curve: curve),
            //    Env([0, 1, 0], [attac,length,release],\linear, releaseNode: releasenode),
                gate
            ) * PlayBuf.ar(2, bufnum, BufRateScale.kr(bufnum) * rate, loop: 0, startPos: skip, doneAction: 2);
        FreeSelfWhenDone.kr(z);
        Out.ar(out, Pan2.ar(z, pos: pan, level: amp));
        } 
).add();


z = Synth("sample2", [out: 0, bufnum: d, attac: 0.002, release: 0.01, skip: 106500, rate: 1, amp: 1]);
z.set(\gate, -0.0000001);


was noch nicht funktioniert:

- hüllkurve soll maximal die länge vom sample - release haben
- skip frame soll in abhängigkeit von der abspielrate sein








