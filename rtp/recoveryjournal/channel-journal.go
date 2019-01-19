package recoveryjournal

// ChannelJournal maintains the state of a channel.
type ChannelJournal struct {
}

/*

   0                   1                   2                   3
   0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
  |S| CHAN  |H|      LENGTH       |P|C|M|W|N|E|T|A|  Chapters ... |
  +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+

                 Figure 9 -- Channel Journal Format

*/

// first 2 header octetts
const (
	channelSFlag      = 0x8000 // Single Package Loss
	channelMask       = 0x7800 // Channel Mask
	channelHFlag      = 0x0400 // Use enhanced Chapter C encoding
	channelLengthMask = 0x03f  //length mask
)

// chapter Table of Content (TOC) (3rd octett)
const (
	chapterP = 0x80 // Chapter P present
	chapterC = 0x40 // Chapter C present
	chapterM = 0x20 // Chapter M present
	chapterW = 0x10 // Chapter W present
	chapterN = 0x08 // Chapter N present
	chapterE = 0x04 // Chapter E present
	chapterT = 0x02 // Chapter T present
	chapterA = 0x01 // Chapter A present
)

// Encode will write the channel journal to a package
func (j *ChannelJournal) Encode() {

}
