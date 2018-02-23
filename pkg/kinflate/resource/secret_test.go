/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package resource

import (
	"encoding/base64"
	"reflect"
	"testing"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	manifest "k8s.io/kubectl/pkg/apis/manifest/v1alpha1"
)

func TestNewFromSecretGenerator(t *testing.T) {
	secret := manifest.SecretGenerator{
		Name: "secret",
		Commands: map[string]string{
			"DB_USERNAME": "printf admin",
			"DB_PASSWORD": "printf somepw",
		},
		Type: "Opaque",
	}
	re, err := NewFromSecretGenerator(secret, ".")

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	expected := &unstructured.Unstructured{
		Object: map[string]interface{}{
			"apiVersion": "v1",
			"kind":       "Secret",
			"metadata": map[string]interface{}{
				"name":              "secret",
				"creationTimestamp": nil,
			},
			"type": string(corev1.SecretTypeOpaque),
			"data": map[string]interface{}{
				"DB_USERNAME": base64.StdEncoding.EncodeToString([]byte("admin")),
				"DB_PASSWORD": base64.StdEncoding.EncodeToString([]byte("somepw")),
			},
		},
	}

	if !reflect.DeepEqual(*re.Data, *expected) {
		t.Fatalf("%#v\ndoesn't match expected:\n%#v", *re.Data, *expected)
	}
}
