# Sonic

[Lookout](https://github.com/src-d/lookout) analyzer that uses [Sonic Pi](https://sonic-pi.net/) to generate sound from PRs.

The analyzer part uses [bblfsh](https://bblf.sh) to extract the UAST nodes from old and new code separating the nodes in deleted and added. It gets the type of node, token (name), size (characters) and hash (used for comparing).

Sound generation uses a preconfigured Sonic Pi setup with a synth and listens to [OSC](https://en.wikipedia.org/wiki/Open_Sound_Control) for note and synth configuration data. Currently the data sent is:

* Note
* Duration (sustain time in seconds)
* Cutoff (filter parameter for the synth)
* Attack (how much it takes the note to reach volume peak)
* Release (how much time it sounds after "releasing" the key)

Currently only Note and Duration are changed and the rest of parameters are always the same. Two note sequences are generated per file, one for deleted and another one for added:

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

Here's the information on how to configure Sonic Pi and the development environment. The project as is cannot be used to send data to github as the sound is played where Sonic Pi is running. Currently the OSC endpoint is hardcoded to `localhost`.

### Prerequisites

* [bblfshd](https://bblf.sh)

```
docker run -d --name bblfshd --privileged -p 9432:9432 -v /var/lib/bblfshd:/var/lib/bblfshd bblfsh/bblfshd:latest
```

* [lookout sdk binary](https://github.com/src-d/lookout/releases)

* [Sonic Pi](https://sonic-pi.net/)



### Configure Sonic Pi

After installing Sonic Pi check with some examples that works and it produces sound. If you are using Linux and it does not work try using jack 2 instead of jack 1 and changing its configuration.

We need to activate OSC input messages. To do so open preferences, go to IO tab and make sure that "Receive remote OSC messages" is checked.

The last step is setting the configuration script in the editor. Delete whatever is in it and paste this script:

```ruby
live_loop :foo do
  use_real_time
  note, duration, cutoff, attack, release = sync "/osc/trigger/prophet"
  synth :prophet,
    note: note,
    cutoff: cutoff,
    sustain: duration,
    release: release,
    attack: attack,
    amp: 8
end
```

Press `Run` button.

## Running it

### Sound tests

There's a sound test with hardcoded data that can be run with:

```
go run cmd/osc/main.go
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

