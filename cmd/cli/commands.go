package cli

import (
	"fmt"
	"slices"

	"github.com/alowayed/go-univers/pkg/univers"
)

func compare[V univers.Version[V], VR univers.VersionRange[V]](
	e univers.Ecosystem[V, VR],
	args []string,
) (int, error) {
	if len(args) != 2 {
		err := fmt.Errorf("compare requires exactly 2 version arguments")
		return 0, err
	}

	vl := args[0]
	vr := args[1]

	verl, err := e.NewVersion(vl)
	if err != nil {
		err = fmt.Errorf("invalid version '%s': %w", vl, err)
		return 0, err
	}
	verr, err := e.NewVersion(vr)
	if err != nil {
		err = fmt.Errorf("invalid version '%s': %w", vr, err)
		return 0, err
	}

	return verl.Compare(verr), nil
}

func sort[V univers.Version[V], VR univers.VersionRange[V]](
	e univers.Ecosystem[V, VR],
	args []string,
) ([]string, error) {
	if len(args) == 0 {
		return nil, fmt.Errorf("sort requires at least 1 version argument")
	}

	versions := make([]V, 0, len(args))
	for _, vStr := range args {
		v, err := e.NewVersion(vStr)
		if err != nil {
			return nil, fmt.Errorf("invalid version '%s': %w", vStr, err)
		}
		versions = append(versions, v)
	}

	slices.SortFunc(versions, V.Compare)

	sortedversions := make([]string, 0, len(versions))
	for _, ver := range versions {
		sortedversions = append(sortedversions, ver.String())
	}

	return sortedversions, nil
}

func contains[V univers.Version[V], VR univers.VersionRange[V]](
	e univers.Ecosystem[V, VR],
	args []string,
) (bool, error) {
	if len(args) != 2 {
		return false, fmt.Errorf("contains requires exactly 2 arguments: <version> <range>")
	}

	rangeStr := args[0]
	versionStr := args[1]

	r, err := e.NewVersionRange(rangeStr)
	if err != nil {
		return false, fmt.Errorf("invalid range '%s': %w", rangeStr, err)
	}

	v, err := e.NewVersion(versionStr)
	if err != nil {
		return false, fmt.Errorf("invalid version '%s': %w", versionStr, err)
	}

	return r.Contains(v), nil
}
