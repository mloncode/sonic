# Sonic

[Lookout](https://github.com/src-d/lookout) analyzer that generates sound from PRs.

The analyzer part uses [bblfsh](https://bblf.sh) to extract the UAST nodes from old and new code separating the nodes in deleted and added. It gets the type of node, token (name), size (characters) and hash (used for comparing).

Sound generation is compatible with any midi synthesizer which exposes midi device. Currently the data sent is:

* Note
* Duration (sustain time in seconds)

Two note sequences are generated per file, one for deleted and another one for added:

* file 1
  * deleted sequence
  * added sequence
* file 2
  * deleted sequence
  * added sequence
* ...

Note generation is done using Markov chains. The idea was taken from a [hackernoon post](https://hackernoon.com/generating-music-using-markov-chains-40c3f3f46405) and even two of its midi files borrowed to generate the two chains used in the project:

* `song1.midi`: [Yiruma - River Flows in You](https://github.com/omgimanerd/markov-music/blob/master/midi/river_flows.mid)
* `song2.midi`: [Chopin’s Concerto №1 in E Minor (Op. 11)](https://github.com/omgimanerd/markov-music/blob/master/midi/chopin_concerto_no1_e_minor_op11.mid)

The midi files are not used as is but converted first to JSON so they are easier to read in the project. [This web application](https://tonejs.github.io/MidiConvert/) was used to convert them.

Even if it's not really noticeable delete sequences use `song1.midi` Markov chain and added use `song2.midi`. To pick one of the weighted notes we use the [FNV](https://golang.org/pkg/hash/fnv/) hash of the token (name).

Duration is a direct conversion from the node size with a maximum of 250 milliseconds.

## Getting Started

Here's the information on how to configure midi synthesizer and the development environment. The project as is cannot be used to send data to github as the sound is played where synthesizer is running.

### Prerequisites

* [bblfshd](https://bblf.sh)

```
docker run -d --name bblfshd --privileged -p 9432:9432 -v /var/lib/bblfshd:/var/lib/bblfshd bblfsh/bblfshd:latest
```

* [lookout sdk binary](https://github.com/src-d/lookout/releases)

* [portmidi](http://portmedia.sourceforge.net/portmidi/)

```
apt-get install libportmidi-dev
# or
sudo pacman -S portmidi
# or
brew install portmidi
```

* Midi synthesizer

  - for macOS we recommend [SimpleSynth](http://notahat.com/simplesynth/)
  - for Linux, [FluidSynth](http://www.fluidsynth.org/). Documentation abouut how to install it can be found [here](https://wiki.archlinux.org/index.php/FluidSynth). Once installed and configured, the output of `aplaymidi -l` should look alike:

  ```
   Port    Client name                      Port name
   128:0    FLUID Synth (32076)              Synth input port (32076:0)
  ```




## Running it

### Sound tests

There's a sound test with hardcoded data that can be run with:

```
go run cmd/test/main.go
```

### Run analyzer

Start analyzer:

```
go run cmd/sonic/main.go
```

Analyze a change (last commit of a repo)

* Go to a directory with the commit we want to analyze
* `git checkout` to the specific dir in case we are not in it
* Execute lookout-sdk executable:

```
$PATH_TO_LOOKOUT_SDK_BINARY/lookout-sdk review ipv4://localhost:2001
```

## Authors

* Maxim Sukharev
* Javi Fontán

## License

This project is licensed under the MIT License - see the [LICENSE.md](LICENSE.md) file for details

## Acknowledgments

* Kuba Podgórski - project idea
* Alvin Lin - [post about Markov chains and music](https://hackernoon.com/generating-music-using-markov-chains-40c3f3f46405)
