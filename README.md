# RTP-MIDI implementation in go

The final goal is to provide a [RTP-MIDI](https://en.wikipedia.org/wiki/RTP-MIDI) implemation in go.

The implementation is currently only tested with the Apple MIDI Network Driver and is restricted to 
Apple's specific session initiation protocol.

This work is inspired and based on the the following open source code:

* [raveloxmidi](https://github.com/ravelox/pimidi/tree/master/raveloxmidi)
* [node-rtpmidi](https://github.com/jdachtera/node-rtpmidi)

The project depends on [zeroconf](https://github.com/grandcat/zeroconf) to support service discovery with mDSN

## Supported features
* Act as session listener
* Single and mulitple MIDI commands per message with delta time


## TODO

WARNING: THIS IMPLEMENTATION IS INCOMPLETE AND WORK IN PROGRESS

The API is not yet stable and will change in future.

The implementation is planned to continue with the following tasks

## Act as session listener
* Send recovery journal
  * Support closed-loop sending policy
  * Support channel-journal
    * Chapter-N
    * Other Chapters
  * Support system-journal
* Support receiving midi payload
  * Receive recovery journal
* Keep-alive message (empty data)
* Improve error handling
* Merge multiple streams
* Hide implementation details (Slimmer API)
* Support phantom bit
* Support enhanced Chapter C encoding

## Act as session initiator
* initiate a new connection to a remote session

