package dimension

// Dimension defines the foundational interface for spatial environments,
// encapsulating the canonical nomenclature and unique protocol identifier.
type Dimension interface {
	// Name returns the formal string identifier of the environmental realm.
	Name() string
	// ID returns the unique 32-bit integer identifier utilized within the network protocol.
	ID() uint32
}
