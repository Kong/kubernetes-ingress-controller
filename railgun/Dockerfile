# NOTE: at the time of writing Kubebuilder does not tolerate moving or removing the `Dockerfile`
#       and if present will validate that the line "COPY api/ api/" exists. If you're looking
#       for the actual image build, look at ../Dockerfile.railgun.
#       Otherwise this file simply exists as a shim to work around the upstream functionality.
#       See https://github.com/kubernetes-sigs/kubebuilder/issues/2165 for follow-up.
COPY apis/ apis/
