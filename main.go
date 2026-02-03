package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os/exec"
	"strconv"
	"strings"
)

type VersionType string

const (
	Major VersionType = "major"
	Minor VersionType = "minor"
	Patch VersionType = "patch"
)

func (v *VersionType) String() string {
	return string(*v)
}

func (v *VersionType) Set(s string) error {
	switch s {
	case string(Major), string(Minor), string(Patch):
		*v = VersionType(s)
		return nil
	default:
		return fmt.Errorf("invalid version type: %s. Available: major, minor or patch", s)
	}
}

type TagVersion struct {
	Major int
	Minor int
	Patch int
}

func (t *TagVersion) String() string {
	return fmt.Sprintf(
		"v%d.%d.%d",
		t.Major,
		t.Minor,
		t.Patch,
	)
}

func (t *TagVersion) Bump(bumpType VersionType) {
	switch bumpType {
	case Major:
		t.Major++
		t.Minor = 0
		t.Patch = 0
	case Minor:
		t.Minor++
		t.Patch = 0
	case Patch:
		t.Patch++
	}
}

func mustAtoi(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		panic(err)
	}
	return i
}

func NewTag(s string) *TagVersion {
	s = strings.TrimSpace(s)
	_p := strings.Split(s, "v")
	versions := strings.Split(_p[1], ".")
	return &TagVersion{
		Major: mustAtoi(versions[0]),
		Minor: mustAtoi(versions[1]),
		Patch: mustAtoi(versions[2]),
	}
}

func gitTag(tag TagVersion) {
	cmd := exec.Command(
		"git", "tag", tag.String(),
	)
	_, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error(
			"failed to git tag",
			"err", err,
		)
		panic(err)
	}
}

func gitTagsPush() {
	cmd := exec.Command(
		"git", "push", "origin", "--tags",
	)
	_, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error(
			"failed to git push tags",
			"err", err,
		)
		panic(err)
	}
}

func getTags() []TagVersion {
	var tags []TagVersion
	cmd := exec.Command(
		"git", "tag",
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		slog.Error(
			"getting tags failed",
			"err", err,
		)
		panic(err)
	}
	tagsRaw := string(out)
	for _, line := range strings.Split(tagsRaw, "\n") {
		if strings.TrimSpace(line) == "" {
			continue
		}
		tags = append(
			tags,
			*NewTag(line),
		)
	}
	return tags
}

func getLatestTag(tags []TagVersion) TagVersion {
	var latest TagVersion
	highestMajor := -1
	for _, t := range tags {
		if t.Major > highestMajor {
			highestMajor = t.Major
			latest = t
		}
	}
	highestMinor := -1
	for _, t := range tags {
		if highestMajor > t.Major {
			continue
		}
		if t.Minor > highestMinor {
			highestMinor = t.Minor
			latest = t
		}
	}
	highestPatch := -1
	for _, t := range tags {
		if highestMajor > t.Major || highestMinor > t.Minor {
			continue
		}
		if t.Patch > highestPatch {
			highestPatch = t.Patch
			latest = t
		}
	}
	return latest
}

func main() {
	var version VersionType = Patch
	flag.Var(&version, "b", "What to bump: major, minor or patch (default: patch)")
	isDryRun := flag.Bool("dry", false, "If true, only logs will be shown, an actual git tag and push won't be executed")
	flag.Parse()
	if *isDryRun {
		slog.Info("dry run found, won't git tag nor git push")
	}
	tags := getTags()
	latest := getLatestTag(tags)
	slog.Debug(
		"git information",
		"found tags", len(tags),
		"latest tag", latest.String(),
	)
	latest.Bump(version)
	slog.Debug(
		"git tagging",
		"version", latest.String(),
	)
	if !*isDryRun {
		gitTag(latest)
	}
	slog.Debug(
		"git pushing",
	)
	if !*isDryRun {
		gitTagsPush()
	}
}
