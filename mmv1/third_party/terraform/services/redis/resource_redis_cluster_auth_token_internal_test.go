package redis

import (
	"testing"

	"github.com/hashicorp/go-cty/cty"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

type redisClusterAuthTokenDiffMock struct {
	values      map[string]interface{}
	unknownKeys map[string]bool
	computed    map[string]bool
	rawConfig   cty.Value
}

func (m *redisClusterAuthTokenDiffMock) Get(key string) interface{} {
	return m.values[key]
}

func (m *redisClusterAuthTokenDiffMock) GetOk(key string) (interface{}, bool) {
	v, ok := m.values[key]
	if !ok || v == nil || v == "" {
		return "", false
	}
	return v, true
}

func (m *redisClusterAuthTokenDiffMock) GetRawConfig() cty.Value {
	return m.rawConfig
}

func (m *redisClusterAuthTokenDiffMock) NewValueKnown(key string) bool {
	return !m.unknownKeys[key]
}

func (m *redisClusterAuthTokenDiffMock) SetNew(key string, value interface{}) error {
	if m.values == nil {
		m.values = map[string]interface{}{}
	}
	m.values[key] = value
	return nil
}

func (m *redisClusterAuthTokenDiffMock) SetNewComputed(key string) error {
	if m.computed == nil {
		m.computed = map[string]bool{}
	}
	m.computed[key] = true
	return nil
}

func TestUnitRedisClusterAuthToken_latestRedisClusterAuthToken(t *testing.T) {
	t.Parallel()

	cases := map[string]struct {
		input    interface{}
		wantName string
		wantErr  bool
	}{
		"nil slice": {
			input:   nil,
			wantErr: true,
		},
		"empty slice": {
			input:   []interface{}{},
			wantErr: true,
		},
		"missing name": {
			input: []interface{}{
				map[string]interface{}{"token": "abc"},
				map[string]interface{}{"name": ""},
			},
			wantErr: true,
		},
		"single version-1": {
			input: []interface{}{
				map[string]interface{}{"name": "projects/p/locations/r/clusters/c/tokenAuthUsers/u/authTokens/version-1"},
			},
			wantName: "projects/p/locations/r/clusters/c/tokenAuthUsers/u/authTokens/version-1",
		},
		"multiple versions ascending": {
			input: []interface{}{
				map[string]interface{}{"name": "projects/p/locations/r/clusters/c/tokenAuthUsers/u/authTokens/version-1"},
				map[string]interface{}{"name": "projects/p/locations/r/clusters/c/tokenAuthUsers/u/authTokens/version-2"},
			},
			wantName: "projects/p/locations/r/clusters/c/tokenAuthUsers/u/authTokens/version-2",
		},
		"numeric ordering past version-9": {
			input: []interface{}{
				map[string]interface{}{"name": "projects/p/locations/r/clusters/c/tokenAuthUsers/u/authTokens/version-10"},
				map[string]interface{}{"name": "projects/p/locations/r/clusters/c/tokenAuthUsers/u/authTokens/version-9"},
				map[string]interface{}{"name": "projects/p/locations/r/clusters/c/tokenAuthUsers/u/authTokens/version-2"},
			},
			wantName: "projects/p/locations/r/clusters/c/tokenAuthUsers/u/authTokens/version-10",
		},
		"prefers versioned over non-versioned": {
			input: []interface{}{
				map[string]interface{}{"name": "projects/p/locations/r/clusters/c/tokenAuthUsers/u/authTokens/custom-id"},
				map[string]interface{}{"name": "projects/p/locations/r/clusters/c/tokenAuthUsers/u/authTokens/version-1"},
			},
			wantName: "projects/p/locations/r/clusters/c/tokenAuthUsers/u/authTokens/version-1",
		},
		"falls back to non-versioned": {
			input: []interface{}{
				map[string]interface{}{"name": "projects/p/locations/r/clusters/c/tokenAuthUsers/u/authTokens/custom-id"},
			},
			wantName: "projects/p/locations/r/clusters/c/tokenAuthUsers/u/authTokens/custom-id",
		},
	}

	for tn, tc := range cases {
		t.Run(tn, func(t *testing.T) {
			got, err := latestRedisClusterAuthToken(tc.input)
			if tc.wantErr {
				if err == nil {
					t.Fatalf("expected error, got %#v", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got["name"] != tc.wantName {
				t.Fatalf("got name %v, want %q", got["name"], tc.wantName)
			}
		})
	}
}

func TestUnitRedisClusterAuthToken_dataSourceSchema(t *testing.T) {
	t.Parallel()

	ds := DataSourceRedisClusterAuthToken()
	if !ds.UseJSONNumber {
		t.Errorf("expected DataSourceRedisClusterAuthToken UseJSONNumber to be true")
	}
	tokenSchema, ok := ds.Schema["token"]
	if !ok {
		t.Fatalf("expected token field in DataSourceRedisClusterAuthToken schema")
	}
	if !tokenSchema.Sensitive {
		t.Errorf("expected DataSourceRedisClusterAuthToken token field to preserve Sensitive: true")
	}
}

func TestUnitRedisClusterAuthToken_projectRegionDiff(t *testing.T) {
	t.Parallel()

	t.Run("unknown cluster marks unset project and region computed", func(t *testing.T) {
		d := &redisClusterAuthTokenDiffMock{
			values:      map[string]interface{}{},
			unknownKeys: map[string]bool{"cluster": true},
		}
		unresolvedProj, unresolvedReg, err := redisClusterAuthTokenResolveClusterScopeDiff(d)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !unresolvedProj || !d.computed["project"] {
			t.Errorf("expected project to be marked SetNewComputed and unresolved, got unresolved=%v computed=%v", unresolvedProj, d.computed["project"])
		}
		if !unresolvedReg || !d.computed["region"] {
			t.Errorf("expected region to be marked SetNewComputed and unresolved, got unresolved=%v computed=%v", unresolvedReg, d.computed["region"])
		}
	})

	t.Run("unknown cluster preserves explicit region", func(t *testing.T) {
		d := &redisClusterAuthTokenDiffMock{
			values:      map[string]interface{}{"region": "europe-west4"},
			unknownKeys: map[string]bool{"cluster": true},
		}
		unresolvedProj, unresolvedReg, err := redisClusterAuthTokenResolveClusterScopeDiff(d)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !unresolvedProj || !d.computed["project"] {
			t.Errorf("expected project to be marked SetNewComputed")
		}
		if unresolvedReg || d.computed["region"] {
			t.Errorf("expected explicit region not to be marked SetNewComputed")
		}
	})

	t.Run("known full cluster extracts unset project and region", func(t *testing.T) {
		d := &redisClusterAuthTokenDiffMock{
			values: map[string]interface{}{
				"cluster": "projects/proj-a/locations/europe-west4/clusters/cluster-a",
			},
		}
		skipProj, skipReg, err := redisClusterAuthTokenResolveClusterScopeDiff(d)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !skipProj || !skipReg {
			t.Errorf("expected skipProj=true, skipReg=true, got %v, %v", skipProj, skipReg)
		}
		if d.values["project"] != "proj-a" {
			t.Errorf("got project=%v, want proj-a", d.values["project"])
		}
		if d.values["region"] != "europe-west4" {
			t.Errorf("got region=%v, want europe-west4", d.values["region"])
		}
	})

	t.Run("apply-time resolved cluster extracts null rawConfig project and region even when prior plan marked them computed", func(t *testing.T) {
		d := &redisClusterAuthTokenDiffMock{
			values: map[string]interface{}{
				"cluster": "projects/proj-a/locations/europe-west4/clusters/cluster-a",
			},
			unknownKeys: map[string]bool{
				"project": true,
				"region":  true,
			},
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"cluster": cty.StringVal("projects/proj-a/locations/europe-west4/clusters/cluster-a"),
				"project": cty.NullVal(cty.String),
				"region":  cty.NullVal(cty.String),
			}),
		}
		skipProj, skipReg, err := redisClusterAuthTokenResolveClusterScopeDiff(d)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if !skipProj || !skipReg {
			t.Errorf("expected skipProj=true, skipReg=true, got %v, %v", skipProj, skipReg)
		}
		if d.values["project"] != "proj-a" {
			t.Errorf("got project=%v, want proj-a", d.values["project"])
		}
		if d.values["region"] != "europe-west4" {
			t.Errorf("got region=%v, want europe-west4", d.values["region"])
		}
	})

	t.Run("known full cluster respects explicit rawConfig project and region", func(t *testing.T) {
		d := &redisClusterAuthTokenDiffMock{
			values: map[string]interface{}{
				"cluster": "projects/proj-a/locations/europe-west4/clusters/cluster-a",
				"project": "explicit-proj",
				"region":  "explicit-region",
			},
			rawConfig: cty.ObjectVal(map[string]cty.Value{
				"cluster": cty.StringVal("projects/proj-a/locations/europe-west4/clusters/cluster-a"),
				"project": cty.StringVal("explicit-proj"),
				"region":  cty.StringVal("explicit-region"),
			}),
		}
		if _, _, err := redisClusterAuthTokenResolveClusterScopeDiff(d); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if d.values["project"] != "explicit-proj" {
			t.Errorf("got project=%v, want explicit-proj", d.values["project"])
		}
		if d.values["region"] != "explicit-region" {
			t.Errorf("got region=%v, want explicit-region", d.values["region"])
		}
	})

	t.Run("known short cluster leaves project and region for provider defaults", func(t *testing.T) {
		d := &redisClusterAuthTokenDiffMock{
			values: map[string]interface{}{
				"cluster": "cluster-a",
			},
		}
		unresolvedProj, unresolvedReg, err := redisClusterAuthTokenResolveClusterScopeDiff(d)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if unresolvedProj || unresolvedReg {
			t.Errorf("expected false, false, got %v, %v", unresolvedProj, unresolvedReg)
		}
		if _, exists := d.values["project"]; exists {
			t.Errorf("expected project unset before DefaultProviderProject, got %v", d.values["project"])
		}
		if _, exists := d.values["region"]; exists {
			t.Errorf("expected region unset before DefaultProviderRegion, got %v", d.values["region"])
		}
	})
}

func TestUnitRedisClusterAuthToken_diffSuppressors(t *testing.T) {
	t.Parallel()

	newResourceData := func(t *testing.T, attrs map[string]interface{}) *schema.ResourceData {
		t.Helper()
		return schema.TestResourceDataRaw(t, ResourceRedisClusterAuthToken().Schema, attrs)
	}

	baseAttrs := map[string]interface{}{
		"project":  "proj-a",
		"region":   "europe-west4",
		"cluster":  "cluster-a",
		"user_id":  "default",
		"token_id": "version-1",
	}

	t.Run("cluster suppressor", func(t *testing.T) {
		d := newResourceData(t, baseAttrs)
		if !redisClusterAuthTokenClusterDiffSuppress("cluster", "cluster-a", "projects/proj-a/locations/europe-west4/clusters/cluster-a", d) {
			t.Errorf("expected same-scope cluster self-link diff to be suppressed")
		}
		if redisClusterAuthTokenClusterDiffSuppress("cluster", "cluster-a", "projects/proj-b/locations/europe-west4/clusters/cluster-a", d) {
			t.Errorf("expected cross-project cluster diff NOT to be suppressed")
		}
		if redisClusterAuthTokenClusterDiffSuppress("cluster", "cluster-a", "projects/proj-a/locations/us-central1/clusters/cluster-a", d) {
			t.Errorf("expected cross-region cluster diff NOT to be suppressed")
		}
		if redisClusterAuthTokenClusterDiffSuppress("cluster", "cluster-a", "projects/proj-a/locations/europe-west4/clusters/cluster-b", d) {
			t.Errorf("expected different cluster name NOT to be suppressed")
		}
		if redisClusterAuthTokenClusterDiffSuppress("cluster", "cluster-a", "projects/proj-a/locations/europe-west4/clusters/cluster-a", nil) {
			t.Errorf("expected nil ResourceData NOT to suppress short vs full path")
		}
	})

	t.Run("user_id suppressor", func(t *testing.T) {
		d := newResourceData(t, baseAttrs)
		if !redisClusterAuthTokenUserDiffSuppress("user_id", "default", "projects/proj-a/locations/europe-west4/clusters/cluster-a/tokenAuthUsers/default", d) {
			t.Errorf("expected same-scope user_id self-link diff to be suppressed")
		}
		if redisClusterAuthTokenUserDiffSuppress("user_id", "default", "projects/proj-a/locations/europe-west4/clusters/cluster-b/tokenAuthUsers/default", d) {
			t.Errorf("expected cross-cluster 'default' user_id diff NOT to be suppressed")
		}
		if redisClusterAuthTokenUserDiffSuppress("user_id", "default", "projects/proj-a/locations/us-central1/clusters/cluster-a/tokenAuthUsers/default", d) {
			t.Errorf("expected cross-region 'default' user_id diff NOT to be suppressed")
		}
		if redisClusterAuthTokenUserDiffSuppress("user_id", "default", "projects/proj-b/locations/europe-west4/clusters/cluster-a/tokenAuthUsers/default", d) {
			t.Errorf("expected cross-project 'default' user_id diff NOT to be suppressed")
		}
	})

	t.Run("token_id suppressor", func(t *testing.T) {
		d := newResourceData(t, baseAttrs)
		if !redisClusterAuthTokenTokenDiffSuppress("token_id", "version-1", "projects/proj-a/locations/europe-west4/clusters/cluster-a/tokenAuthUsers/default/authTokens/version-1", d) {
			t.Errorf("expected same-scope token_id self-link diff to be suppressed")
		}
		if redisClusterAuthTokenTokenDiffSuppress("token_id", "version-1", "projects/proj-a/locations/europe-west4/clusters/cluster-a/tokenAuthUsers/other-user/authTokens/version-1", d) {
			t.Errorf("expected cross-user 'version-1' token_id diff NOT to be suppressed")
		}
		if redisClusterAuthTokenTokenDiffSuppress("token_id", "version-1", "projects/proj-a/locations/europe-west4/clusters/cluster-b/tokenAuthUsers/default/authTokens/version-1", d) {
			t.Errorf("expected cross-cluster 'version-1' token_id diff NOT to be suppressed")
		}
		if redisClusterAuthTokenTokenDiffSuppress("token_id", "version-1", "projects/proj-b/locations/europe-west4/clusters/cluster-a/tokenAuthUsers/default/authTokens/version-1", d) {
			t.Errorf("expected cross-project 'version-1' token_id diff NOT to be suppressed")
		}
	})
}
