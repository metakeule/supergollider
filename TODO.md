# TODO

1. fix ending: always need and end event and do the freeing of all at the end event, raise an error if no end event is set, maybe register all voices inside the stage and mute them at the end, start and end should be methods of track
2.  generalize sampleplayer to be based on a string for matching the sample file name to freq etc. (or taking a function for this conversion)
3.  make a sc tool to manually find the real frequency and offset of a sample and save it in its meta file
   - eventuell:  aubiopitch -i datei und dann entweder einen durchschnitt nehmen, oder um den offset herum den durchschnitt nehmen
   - eventuell: aubioonset -i datei um den offset zu bestimmen
   
   ERGEBNIS: am besten funktioniert es manuell mit audacity und frequenzanalyse (erweiterte auto korrelation/hanning fenster und größe zwischen 1024 und 4096 variieren, dann die maus in die nähe der spitze setzen und spitze ablesen), als plausibilität vergleiche man den angegebenen ton für die spitze (z.b. C2) mit dem in der sampledatei definierten
   für den offset in audacity das entsprechende erste maximum suchen und 
   auf millisekunde genau ablesen (z.b. mit zoom), zur plausibilität mit
   rechnerisch ermitteltem offset vergleichen

   ergebnis: es gibt keinen einfachen weg, da die abspielgeschwindigkeit
   variiert wird, um die tonhöhe zu ändern. damit ändern sich jedoch auch
   attack und release usw. sowie die ganze länge des samples und der offset, so dass das alles nicht mehr hinhaut. vielleicht kann man das lösen, indem
   man immer durch den scale faktor teilt (muss man ausprobieren)

   vielleicht genügt es für den offset, statt an einem bestimmten stelle
   in sekunden zu starten, an einem bestimmten frame zu starten und zwar mit
   BufFrames, siehe http://new-supercollider-mailing-lists-forums-use-these.2681727.n2.nabble.com/Play-a-Sample-with-random-start-td7613852.html
   (ach ne kann nicht funktionieren, da ich den offset ja im verhältnis zur music brauche (d.h. er soll früher zu spielen anfangen))


0. wir brauchen eine möglichkeit, während des abspielens:
   1. audio aufnehmen zu können
   2. midi aufnehmen zu können
   3. text abspielen zu können
   4. den original quelltext abspielen zu können (den von einem gewählten track)
   5. aufgenommenes an der richtigen stelle (taktweise) in den quelltext einfügen zu können (in dem man den original quelltext ausgibt und die aufgenommenen midiereignisse entsprechend einfügt)
   6. synthesizer basierend auf midiereignissen abzuspielen
   7. aufgenommenes audio durch effekte zu verändern
   8. an bestimmten stellen zurückzuspringen für widerholungen etc. 
      vgl. you might copy all the buffers into a unique buffer and use the index to jump to a position by using BufRd and Phasor.
      (http://new-supercollider-mailing-lists-forums-use-these.2681727.n2.nabble.com/Playing-one-Buffer-immediately-after-another-td7611493.html)

   wahrscheinlich ist der einfachste weg, sclang dafür zu verwenden, d.h. den fertigen sample als buffer einzuladen und entsprechende aufnahme uns ausgabemöglichkeiten wegzuabstrahieren.

   da wir ja nicht wissen, wie groß die generierte datei wird, ist vielleicht folgendes interessant (aus dem Handbuch, http://danielnouri.org/docs/SuperColliderHelp/Tutorials/Getting-Started/Buffers.html):

   Streaming a File in From Disk

In some cases, for instance when working with very large files, you might not want to load a sound completely into memory. Instead, you can stream it in from disk a bit at a time, using the UGen DiskIn, and Buffer's 'cueSoundFile' method:

        (
        SynthDef("tutorial-Buffer-cue",{ arg out=0,bufnum;
            Out.ar(out,
                DiskIn.ar( 1, bufnum )
            )
        }).send(s);
        )
        
        b = Buffer.cueSoundFile(s,"sounds/a11wlk01-44_1.aiff", 0, 1);
        y = Synth.new("tutorial-Buffer-cue", [\bufnum,b.bufnum], s);

        b.free; y.free;

This is not as flexible as PlayBuf (no rate control), but can save memory.


1. außerdem brauchen wir folgende export möglichkeiten:
   - midifile
   - ardour file
   - lilypond file

https://github.com/rakyll/portmidi 




