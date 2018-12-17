package faker

// Pair is to define a simple Key Value
type Pair struct {
	key   string
	value string
}

// Key returns key of the Pair
func (p *Pair) Key() string { return p.key }

// Value returns the value of the Pair
func (p *Pair) Value() string { return p.value }

// KeyValue returns both key & value
func (p *Pair) KeyValue() (string, string) { return p.key, p.value }

// NewPair returns a new key value Pair
func NewPair(key, value string) Pair { return Pair{key, value} }

// NewEmptyPair returns the empty pair where value is empty
func NewEmptyPair(key string) Pair { return Pair{key, ""} }
