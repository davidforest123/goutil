package anycloud

type (
	Gossip struct {
	}
)

func (g *Gossip) Broadcast(b []byte) error {
	return nil
}

func (g *Gossip) ListenRadio() ([]byte, error) {
	return nil, nil
}
