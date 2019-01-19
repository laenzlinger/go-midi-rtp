package recoveryjournal

// RecoveryJournal maintains the state of the journal.
type RecoveryJournal struct {
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
