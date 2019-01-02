# RTP-MIDI implementation in go

The final goal is to provide a [RTP-MIDI](https://en.wikipedia.org/wiki/RTP-MIDI) (aka Apple Midi) implemation in go.

This work is inspired and based on the the following open source code:

* [reaveloxmidi](https://github.com/ravelox/pimidi/tree/master/raveloxmidi)
* [node-rtpmidi](https://github.com/jdachtera/node-rtpmidi)

The project depends on [zeroconf](https://github.com/grandcat/zeroconf) to support service discovery with mDSN
