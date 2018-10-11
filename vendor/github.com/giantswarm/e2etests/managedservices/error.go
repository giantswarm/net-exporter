package managedservices

import "github.com/giantswarm/microerror"

var invalidConfigError = &microerror.Error{
	Kind: "invalidConfigError",
}

// IsInvalidConfig asserts invalidConfigError.
func IsInvalidConfig(err error) bool {
	return microerror.Cause(err) == invalidConfigError
}

var invalidLabelsError = &microerror.Error{
	Kind: "invalidLabelsError",
}

// IsInvalidLabels asserts invalidLabelsError.
func IsInvalidLabels(err error) bool {
	return microerror.Cause(err) == invalidLabelsError
}

var invalidReplicasError = &microerror.Error{
	Kind: "invalidReplicasError",
}

// IsInvalidReplicas asserts invalidReplicasError.
func IsInvalidReplicas(err error) bool {
	return microerror.Cause(err) == invalidReplicasError
}

var notFoundError = &microerror.Error{
	Kind: "notFoundError",
}

// IsNotFound asserts NotFoundError.
func IsNotFound(err error) bool {
	return microerror.Cause(err) == notFoundError
}
