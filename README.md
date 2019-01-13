# RTP-MIDI implementation in go

The final goal is to provide a [RTP-MIDI](https://en.wikipedia.org/wiki/RTP-MIDI) implemation in go.

The implementation is currently only tested with the Apple MIDI Network Driver and is restricted to 
Apple's specific session initiation protocol.

This work is inspired and based on the the following open source code:

* [reaveloxmidi](https://github.com/ravelox/pimidi/tree/master/raveloxmidi)
* [node-rtpmidi](https://github.com/jdachtera/node-rtpmidi)

The project depends on [zeroconf](https://github.com/grandcat/zeroconf) to support service discovery with mDSN

## TODO

WARNING: THIS IMPLEMENTATION IS INCOMPLETE AND WORK IN PROGRESS

The implementation is planned to continue with the following tasks

## Act as session listener
* Decode receiver feedbak
* Support journal
* Support receiving midi payload
* Support mulitple commands in one message

## Act as session initiator
* initiate a new connection to a remote session

