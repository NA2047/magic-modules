package redis

import (
	"testing"
)

func TestSecondaryIpDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"empty strings": {
			Old:                "",
			New:                "",
			ExpectDiffSuppress: true,
		},
		"auto range": {
			Old:                "",
			New:                "auto",
			ExpectDiffSuppress: false,
		},
		"auto on already applied range": {
			Old:                "10.0.0.0/28",
			New:                "auto",
			ExpectDiffSuppress: true,
		},
		"same ranges": {
			Old:                "10.0.0.0/28",
			New:                "10.0.0.0/28",
			ExpectDiffSuppress: true,
		},
		"different ranges": {
			Old:                "10.0.0.0/28",
			New:                "10.1.2.3/28",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if secondaryIpDiffSuppress("whatever", tc.Old, tc.New, nil) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}

func TestUnitRedisInstance_redisVersionIsDecreasing(t *testing.T) {
	t.Parallel()
	type testcase struct {
		name       string
		old        interface{}
		new        interface{}
		decreasing bool
	}
	tcs := []testcase{
		{
			name:       "stays the same",
			old:        "REDIS_4_0",
			new:        "REDIS_4_0",
			decreasing: false,
		},
		{
			name:       "increases",
			old:        "REDIS_4_0",
			new:        "REDIS_5_0",
			decreasing: false,
		},
		{
			name:       "nil vals",
			old:        nil,
			new:        "REDIS_4_0",
			decreasing: false,
		},
		{
			name:       "corrupted",
			old:        "REDIS_4_0",
			new:        "REDIS_banana",
			decreasing: false,
		},
		{
			name:       "decreases",
			old:        "REDIS_6_0",
			new:        "REDIS_4_0",
			decreasing: true,
		},
	}

	for _, tc := range tcs {
		decreasing := isRedisVersionDecreasingFunc(tc.old, tc.new)
		if decreasing != tc.decreasing {
			t.Errorf("%s: expected decreasing to be %v, but was %v", tc.name, tc.decreasing, decreasing)
		}
	}
}

func TestMaintenanceVersionDiffSuppress(t *testing.T) {
	cases := map[string]struct {
		Old, New           string
		ExpectDiffSuppress bool
	}{
		"same versions": {
			Old:                "20250326_00_00",
			New:                "20250326_00_00",
			ExpectDiffSuppress: true,
		},
		"config is older than API (should suppress - instance was upgraded)": {
			Old:                "20250701_00_01",
			New:                "20250326_00_00",
			ExpectDiffSuppress: true,
		},
		"config is newer than API (should NOT suppress - user wants upgrade)": {
			Old:                "20250326_00_00",
			New:                "20250701_00_01",
			ExpectDiffSuppress: false,
		},
		"same date different patch version - config older": {
			Old:                "20250326_00_01",
			New:                "20250326_00_00",
			ExpectDiffSuppress: true,
		},
		"same date different patch version - config newer": {
			Old:                "20250326_00_00",
			New:                "20250326_00_01",
			ExpectDiffSuppress: false,
		},
		"different dates - config much older": {
			Old:                "20251007_00_00",
			New:                "20250326_00_00",
			ExpectDiffSuppress: true,
		},
		"different dates - config much newer": {
			Old:                "20250326_00_00",
			New:                "20251007_00_00",
			ExpectDiffSuppress: false,
		},
		"empty old value": {
			Old:                "",
			New:                "20250326_00_00",
			ExpectDiffSuppress: false,
		},
		"empty new value": {
			Old:                "20250326_00_00",
			New:                "",
			ExpectDiffSuppress: false,
		},
		"both empty": {
			Old:                "",
			New:                "",
			ExpectDiffSuppress: false,
		},
		"invalid format - should not suppress": {
			Old:                "invalid",
			New:                "20250326_00_00",
			ExpectDiffSuppress: false,
		},
		"minor version difference - config older": {
			Old:                "20250806_01_00",
			New:                "20250806_00_00",
			ExpectDiffSuppress: true,
		},
		"minor version difference - config newer": {
			Old:                "20250806_00_00",
			New:                "20250806_01_00",
			ExpectDiffSuppress: false,
		},
		"minor and patch version greater than 60": {
			Old:                "20250806_00_00",
			New:                "20250806_70_90",
			ExpectDiffSuppress: false,
		},
	}

	for tn, tc := range cases {
		if maintenanceVersionDiffSuppress(tc.Old, tc.New) != tc.ExpectDiffSuppress {
			t.Fatalf("bad: %s, '%s' => '%s' expect %t", tn, tc.Old, tc.New, tc.ExpectDiffSuppress)
		}
	}
}
