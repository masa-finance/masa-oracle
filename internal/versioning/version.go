package versioning

var (
	ApplicationVersion string

	// XXX: Bump this value only when there are protocol changes that makes the oracle
	// incompatible between version!
	ProtocolVersion = `v0.7.0`
)
