package pod

import (
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-containerregistry/pkg/name"
	corev1 "k8s.io/api/core/v1"
)

const validDigest = "sha256:aec27421b7a64a63b5dbf3c62b4a1d44f0bda5632cc74b256651df920d61e09b"

func TestResolveEntrypoints(t *testing.T) {
	cache := fakeCache{
		"image-by-digest@" + validDigest: data{
			ep:     []string{"my", "entrypoint"},
			digest: validDigest,
		},
		"my-image": data{
			ep:     []string{"cool", "entrypoint", "bro"},
			digest: validDigest,
		},
	}

	got, err := resolveEntrypoints(cache, "namespace", "serviceAccountName", []corev1.Container{{
		Image:   "fully-specified",
		Command: []string{"specified", "command"}, // nothing to resolve
	}, {
		Image: "my-image",
	}, {
		Image: "image-by-digest@" + validDigest, // specified by digest.
	}})
	if err != nil {
		t.Fatalf("Error resolving entrypoints: %v", err)
	}

	want := []corev1.Container{{
		Image:   "fully-specified",
		Command: []string{"specified", "command"},
	}, {
		Image:   "index.docker.io/library/my-image@" + validDigest, // digest was resolved when looking up entrypoint.
		Command: []string{"cool", "entrypoint", "bro"},
	}, {
		Image:   "index.docker.io/library/image-by-digest@" + validDigest,
		Command: []string{"my", "entrypoint"},
	}}
	if d := cmp.Diff(want, got); d != "" {
		t.Fatalf("Diff (-want, +got): %s", d)
	}
}

type fakeCache map[string]data
type data struct {
	ep     []string
	digest string
}

func (f fakeCache) Get(imageName, _, _ string) ([]string, name.Digest, error) {
	ref, err := name.ParseReference(imageName, name.WeakValidation)
	if err != nil {
		return nil, name.Digest{}, err
	}
	if d, ok := ref.(name.Digest); ok {
		if data, found := f[d.String()]; found {
			return data.ep, d, nil
		}
	}

	d, found := f[imageName]
	if !found {
		return nil, name.Digest{}, errors.New("not found")
	}
	dig, err := name.NewDigest(ref.Context().RepositoryStr()+"@"+d.digest, name.WeakValidation)
	if err != nil {
		return nil, name.Digest{}, err
	}
	return d.ep, dig, nil
}
