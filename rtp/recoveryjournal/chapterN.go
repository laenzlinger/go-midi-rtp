package recoveryjournal

/*

    0                   1                   2                   3
    0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 8 0 1
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |B|     LEN     |  LOW  | HIGH  |S|   NOTENUM   |Y|  VELOCITY   |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |S|   NOTENUM   |Y|  VELOCITY   |             ....              |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
   |    OFFBITS    |    OFFBITS    |     ....      |    OFFBITS    |
   +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

                   Figure A.6.1 -- Chapter N Format

*/

// ChapterN is responsible for MIDI NoteOff (0x8), NoteOn (0x9) commands
type ChapterN struct {
	NoteSeqNumber uint32    // most recent note seqnum, or 0
	NoteTimestamp uint32    // NoteOn execution timestamp
	NoteOn        []NoteOn  // Max. 128 NoteOn messages
	NoteOff       []NoteOff // Max. 16 OffBit messages
}

/*

   0                   1
   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5
  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
  |S|   NOTENUM   |Y|  VELOCITY   |
  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

 Figure A.6.3 -- Chapter N Note Log
*/
type NoteOn struct {
	NoteNum            uint8
	Velocity           uint8 // never 0
	PlayRecommendation bool  // Y=1: play Y=0: skip
}

type NoteOff struct {
	NoteNum uint8
}
