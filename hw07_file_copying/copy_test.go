package main

import (
	"crypto/md5"
	"errors"
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	const (
		input  = "./testdata/input.txt"
		output = "./out.txt"
	)

	tests := []struct {
		validOut string
		offset   int64
		limit    int64
	}{
		{"./testdata/out_offset0_limit0.txt", 0, 0},
		{"./testdata/out_offset0_limit10.txt", 0, 10},
		{"./testdata/out_offset0_limit1000.txt", 0, 1000},
		{"./testdata/out_offset0_limit10000.txt", 0, 10000},
		{"./testdata/out_offset100_limit1000.txt", 100, 1000},
		{"./testdata/out_offset6000_limit1000.txt", 6000, 1000},
	}

	t.Run("Compare input and output files", func(t *testing.T) {
		for _, tc := range tests {
			err := Copy(input, output, tc.offset, tc.limit)
			require.NoError(t, err)
			defer os.Remove(output)

			hOutput := md5.New()
			hValidOut := md5.New()

			fileOutput, err := os.Open(output)
			require.NoError(t, err)
			defer fileOutput.Close()

			fileValid, err := os.Open(tc.validOut)
			require.NoError(t, err)
			defer fileValid.Close()

			// get hashes
			_, err = io.Copy(hOutput, fileOutput)
			require.NoError(t, err)

			_, err = io.Copy(hValidOut, fileValid)
			require.NoError(t, err)

			require.Equal(t, hOutput, hValidOut)
		}
	})

	t.Run("Negative case with big offset value", func(t *testing.T) {
		var offset int64 = 7000
		var limit int64

		err := Copy(input, output, offset, limit)
		defer os.Remove(output)

		require.Truef(t, errors.Is(err, ErrOffsetExceedsFileSize), "actual err - %v", err)
	})

	t.Run("Case with big limit value", func(t *testing.T) {
		var offset int64
		var limit int64 = 10000000

		err := Copy(input, output, offset, limit)
		require.NoError(t, err)
		defer os.Remove(output)

		fileIn, err := os.Stat(input)
		require.NoError(t, err)

		fileOut, err := os.Stat(output)
		require.NoError(t, err)

		require.Equal(t, fileIn.Size(), fileOut.Size())
	})

	t.Run("Negative case with unsupported source file type", func(t *testing.T) {
		const (
			random = "/dev/urandom"
			null   = "/dev/null"
		)

		var offset int64
		var limit int64

		err := Copy(random, null, offset, limit)

		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual err - %v", err)
	})

	t.Run("Negative case with unsupported destination file", func(t *testing.T) {
		out := ""
		var offset int64
		var limit int64

		err := Copy(input, out, offset, limit)

		require.Truef(t, errors.Is(err, ErrUnsupportedFile), "actual err - %v", err)
	})
}
