package docker

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/opencontainers/go-digest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/goodwithtech/dockle/pkg/deckoder/utils"
	"github.com/goodwithtech/dockle/pkg/types"
)

// ociDescriptor is a minimal subset of the OCI content descriptor.
type ociDescriptor struct {
	MediaType string        `json:"mediaType"`
	Digest    digest.Digest `json:"digest"`
	Size      int64         `json:"size"`
}

// buildOCILayout writes a minimal but valid OCI image layout into dir. The
// layout holds a single tar layer carrying the given files, gzipped when
// compress is true. It returns the digest of the layer blob, which is what
// LayerIDs reports and ExtractLayerFiles expects.
func buildOCILayout(t *testing.T, dir string, files map[string][]byte, compress bool) digest.Digest {
	t.Helper()

	blobsDir := filepath.Join(dir, "blobs", "sha256")
	require.NoError(t, os.MkdirAll(blobsDir, 0755))

	writeBlob := func(data []byte) (digest.Digest, int64) {
		dig := digest.FromBytes(data)
		require.NoError(t, os.WriteFile(filepath.Join(blobsDir, dig.Encoded()), data, 0644))
		return dig, int64(len(data))
	}

	// Build the tar layer. Its digest is the diff_id, i.e. the digest of the
	// uncompressed bytes.
	var tarBuf bytes.Buffer
	tw := tar.NewWriter(&tarBuf)
	for name, body := range files {
		require.NoError(t, tw.WriteHeader(&tar.Header{
			Name:     name,
			Mode:     0644,
			Size:     int64(len(body)),
			Typeflag: tar.TypeReg,
		}))
		_, err := tw.Write(body)
		require.NoError(t, err)
	}
	require.NoError(t, tw.Close())
	diffID := digest.FromBytes(tarBuf.Bytes())

	// The blob stored in the layout is the compressed layer when requested, so
	// its digest can differ from the diff_id.
	layerMediaType := "application/vnd.oci.image.layer.v1.tar"
	layerBlob := tarBuf.Bytes()
	if compress {
		layerMediaType = "application/vnd.oci.image.layer.v1.tar+gzip"
		var gzBuf bytes.Buffer
		gw := gzip.NewWriter(&gzBuf)
		_, err := gw.Write(tarBuf.Bytes())
		require.NoError(t, err)
		require.NoError(t, gw.Close())
		layerBlob = gzBuf.Bytes()
	}
	layerDigest, layerSize := writeBlob(layerBlob)

	config := map[string]interface{}{
		"architecture": "amd64",
		"os":           "linux",
		"rootfs": map[string]interface{}{
			"type":     "layers",
			"diff_ids": []string{diffID.String()},
		},
	}
	configBytes, err := json.Marshal(config)
	require.NoError(t, err)
	configDigest, configSize := writeBlob(configBytes)

	// Build the image manifest referencing the config and the layer.
	manifest := map[string]interface{}{
		"schemaVersion": 2,
		"mediaType":     "application/vnd.oci.image.manifest.v1+json",
		"config": ociDescriptor{
			MediaType: "application/vnd.oci.image.config.v1+json",
			Digest:    configDigest,
			Size:      configSize,
		},
		"layers": []ociDescriptor{
			{
				MediaType: layerMediaType,
				Digest:    layerDigest,
				Size:      layerSize,
			},
		},
	}
	manifestBytes, err := json.Marshal(manifest)
	require.NoError(t, err)
	manifestDigest, manifestSize := writeBlob(manifestBytes)

	// Build the top-level index.json pointing at the manifest.
	index := map[string]interface{}{
		"schemaVersion": 2,
		"mediaType":     "application/vnd.oci.image.index.v1+json",
		"manifests": []ociDescriptor{
			{
				MediaType: "application/vnd.oci.image.manifest.v1+json",
				Digest:    manifestDigest,
				Size:      manifestSize,
			},
		},
	}
	indexBytes, err := json.Marshal(index)
	require.NoError(t, err)
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.json"), indexBytes, 0644))

	// Mark the directory as an OCI image layout.
	require.NoError(t, os.WriteFile(
		filepath.Join(dir, "oci-layout"),
		[]byte(`{"imageLayoutVersion":"1.0.0"}`),
		0644,
	))

	return layerDigest
}

func TestNewOCIDirExtractor_ExtractLayerFiles(t *testing.T) {
	testcases := map[string]struct {
		compress bool
	}{
		// The OCI image spec allows uncompressed layers, and they do show up in
		// layouts built without a Docker daemon.
		"uncompressed layer": {compress: false},
		"gzipped layer":      {compress: true},
	}

	wantPath := "etc/hostname"

	for name, tc := range testcases {
		t.Run(name, func(t *testing.T) {
			dir := t.TempDir()
			layerDigest := buildOCILayout(t, dir, map[string][]byte{
				"etc/hostname": []byte("localhost\n"),
			}, tc.compress)

			d, cleanup, err := NewOCIDirExtractor(context.Background(), dir, types.DockerOption{})
			require.NoError(t, err)
			defer cleanup()

			// The layer digest reported by the layout must match the blob we wrote.
			assert.Equal(t, []string{layerDigest.String()}, d.LayerIDs())

			filter := utils.CreateFilterPathFunc([]string{wantPath})
			actualFileMap, opqDirs, whFiles, err := d.ExtractLayerFiles(context.Background(), layerDigest, filter)
			require.NoError(t, err)

			assert.Equal(t, types.FileMap{
				wantPath: types.FileData{
					Body:     []byte("localhost\n"),
					FileMode: NormalFileMode,
				},
			}, actualFileMap)
			assert.Empty(t, opqDirs)
			assert.Empty(t, whFiles)
		})
	}
}

// TestNewOCIDirExtractor_NotALayout checks that a directory which is not an OCI
// image layout is reported as an error rather than silently scanning nothing.
func TestNewOCIDirExtractor_NotALayout(t *testing.T) {
	_, _, err := NewOCIDirExtractor(context.Background(), t.TempDir(), types.DockerOption{})
	assert.Error(t, err)
}
