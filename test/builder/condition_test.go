/*
Copyright 2019 The Tekton Authors

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

package builder_test

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/tektoncd/pipeline/pkg/apis/pipeline/v1alpha1"
	tb "github.com/tektoncd/pipeline/test/builder"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestCondition(t *testing.T) {
	condition := tb.Condition("cond-name",
		tb.ConditionLabels(
			map[string]string{
				"label-1": "label-value-1",
				"label-2": "label-value-2",
			}),
		tb.ConditionSpec(tb.ConditionSpecCheck("", "ubuntu", tb.Command("exit 0")),
			tb.ConditionParamSpec("param-1", v1alpha1.ParamTypeString,
				tb.ParamSpecDefault("default"),
				tb.ParamSpecDescription("desc")),
			tb.ConditionResource("git-resource", v1alpha1.PipelineResourceTypeGit),
			tb.ConditionResource("pr", v1alpha1.PipelineResourceTypePullRequest),
		),
	)

	expected := &v1alpha1.Condition{
		ObjectMeta: metav1.ObjectMeta{
			Name: "cond-name",
			Labels: map[string]string{
				"label-1": "label-value-1",
				"label-2": "label-value-2",
			},
		},
		Spec: v1alpha1.ConditionSpec{
			Check: corev1.Container{
				Image:   "ubuntu",
				Command: []string{"exit 0"},
			},
			Params: []v1alpha1.ParamSpec{{
				Name:        "param-1",
				Type:        v1alpha1.ParamTypeString,
				Description: "desc",
				Default: &v1alpha1.ArrayOrString{
					Type:      v1alpha1.ParamTypeString,
					StringVal: "default",
				}}},
			Resources: []v1alpha1.ResourceDeclaration{{
				Name: "git-resource",
				Type: "git",
			}, {
				Name: "pr",
				Type: "pullRequest",
			}},
		},
	}

	if d := cmp.Diff(expected, condition); d != "" {
		t.Fatalf("Condition diff -want, +got: %v", d)
	}
}
