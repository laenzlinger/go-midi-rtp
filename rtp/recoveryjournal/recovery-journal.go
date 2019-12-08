package recoveryjournal

import (
	"github.com/laenzlinger/go-midi-rtp/rtp"
)

// CheckpointHistory contains the history of sent packets of a stream
// since the start of the checkopoint.
type CheckpointHistory struct {
	SentMessages []rtp.MIDIMessage
}

// RecoveryJournal contains the internal structure of the complete
// sender recovery journal
type RecoveryJournal struct {
	// SequNum contains the extended sequence number, or 0.
	// SequNum = 0 codes empty journal
	CheckpointPackageSeqNum uint32

	// ChannelJournal contains the channel part of the history
	ChannelJournal ChannelJournal

	// TODO add system journal
}

/*

   0                   1                   2
   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3
  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
  |S|Y|A|H|TOTCHAN|   Checkpoint Packet Seqnum    |
  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

		Figure 8 -- Recovery Journal Header

*/
const (
	headerSFlag = 0x80 // Single Package Loss
	headerYFlag = 0x40 // System Journal present
	headerAFlag = 0x20 // Channel Journal present
	headerHFlag = 0x10 // Use enhanced Chapter C encoding
	totChanMask = 0xf  // Total Channels
)

// Encode will write the recovery journal to a package
func (j *RecoveryJournal) Encode() {

}
