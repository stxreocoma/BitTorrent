package torrent

type TorrentFile struct {
	Announce    string
	InfoHash    [20]byte
	PieceHashes [][20]byte
	PieceLength int
	Name        string
}

func (bto bencodeTorrent) toTorrentFile() (TorrentFile, error) {

}
