package v1

import (
	"os"
	"path/filepath"
	"testing"

	"gopkg.in/yaml.v3"
)

func TestIssuerCRDsPreserveLegacyCABundleSchema(t *testing.T) {
	for _, filename := range []string{
		"adcs.certmanager.csf.nokia.com_adcsissuers.yaml",
		"adcs.certmanager.csf.nokia.com_clusteradcsissuers.yaml",
	} {
		t.Run(filename, func(t *testing.T) {
			data, err := os.ReadFile(filepath.Join("..", "..", "config", "crd", "bases", filename))
			if err != nil {
				t.Fatal(err)
			}

			var crd map[string]any
			if err := yaml.Unmarshal(data, &crd); err != nil {
				t.Fatal(err)
			}

			version := nestedMap(crd, "spec", "versions").([]any)[0].(map[string]any)
			schema := nestedMap(version, "schema", "openAPIV3Schema").(map[string]any)
			specSchema := nestedMap(schema, "properties", "spec").(map[string]any)
			properties := nestedMap(specSchema, "properties").(map[string]any)

			legacyCABundle := properties["caBundle"].(map[string]any)
			if legacyCABundle["type"] != "string" || legacyCABundle["format"] != "byte" {
				t.Fatalf("legacy caBundle schema changed: %#v", legacyCABundle)
			}
			if _, ok := properties["credentialsRef"]; !ok {
				t.Fatal("legacy credentialsRef field is missing")
			}
			if _, ok := properties["url"]; !ok {
				t.Fatal("legacy url field is missing")
			}

			ref := properties["caBundleRef"].(map[string]any)
			if ref["type"] != "object" {
				t.Fatalf("caBundleRef should be an object: %#v", ref)
			}
			refProperties := ref["properties"].(map[string]any)
			for _, field := range []string{"name", "kind", "key"} {
				if _, ok := refProperties[field]; !ok {
					t.Fatalf("caBundleRef.%s is missing", field)
				}
			}
			kind := refProperties["kind"].(map[string]any)
			if kind["type"] != "string" {
				t.Fatalf("caBundleRef.kind should be a string: %#v", kind)
			}
			keys := refProperties["keys"].(map[string]any)
			if keys["type"] != "array" || keys["minItems"] != 1 {
				t.Fatalf("caBundleRef.keys should be a non-empty array: %#v", keys)
			}
			if _, ok := ref["required"]; !ok {
				t.Fatal("caBundleRef.name must remain required when the reference is supplied")
			}

			required := nestedMap(specSchema, "required")
			for _, field := range []string{"credentialsRef", "url"} {
				if !containsString(required.([]any), field) {
					t.Fatalf("legacy required field %q is missing", field)
				}
			}
			if containsString(required.([]any), "caBundleRef") {
				t.Fatal("caBundleRef must remain optional for backward compatibility")
			}
		})
	}
}

func nestedMap(value map[string]any, keys ...string) any {
	var current any = value
	for _, key := range keys {
		current = current.(map[string]any)[key]
	}
	return current
}

func containsString(values []any, want string) bool {
	for _, value := range values {
		if value == want {
			return true
		}
	}
	return false
}
