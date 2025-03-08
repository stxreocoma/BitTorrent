package torrentfile

const Port uint16 = 6881

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Length      int
	Name        string
}

func (bto *bencodeTorrent) toTorrentFile() (*TorrentFile, error) {

}
