/*
Copyright 2025.

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

package v1alpha1

import (
	"bytes"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"kdex.dev/crds/base"
)

const (
	importStatementTemplate = `<script type="module">
  %s
</script>`
)

// KDexScriptLibrarySpec defines the desired state of KDexScriptLibrary
// +kubebuilder:validation:XValidation:rule="[has(self.scripts), has(self.packageReference)].filter(x, x).size() > 0",message="at least one of scripts or packageReference must be specified"
type KDexScriptLibrarySpec struct {
	// packageReference specifies the name and version of an NPM package that contains the script. The package.json must describe an ES module.
	// +optional
	PackageReference *PackageReference `json:"packageReference,omitempty"`

	// scripts is a set of script references. They may contain URLs that point to resources hosted at some public address, npm module references or they may contain tag contents.
	// +kubebuilder:validation:MaxItems=32
	// +optional
	Scripts []Script `json:"scripts,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:scope=Namespaced,shortName=kdex-sl

// KDexScriptLibrary is the Schema for the kdexscriptlibraries API
//
// A KDexScriptLibrary is a reusable collection of JavaScript for powering the imperative aspects of KDexPageBindings.
// Most other components of the model are able to reference KDexScriptLibrary as well in order to encapsulate component
// specific logic.
//
// +kubebuilder:printcolumn:name="Ready",type=string,JSONPath=`.status.conditions[?(@.type=="Ready")].status`,description="The state of the Ready condition"
type KDexScriptLibrary struct {
	base.KDexObject `json:",inline"`

	// spec defines the desired state of KDexScriptLibrary
	// +required
	Spec KDexScriptLibrarySpec `json:"spec"`
}

// +kubebuilder:object:root=true

// KDexScriptLibraryList contains a list of KDexScriptLibrary
type KDexScriptLibraryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []KDexScriptLibrary `json:"items"`
}

// PackageReference specifies the name and version of an NPM package that contains the micro-frontend application.
type PackageReference struct {
	// exportMapping is a mapping of the module's exports that will be used when the module import is written. e.g. `import [exportMapping] from [module_name];`. If exportMapping is not provided the module will be written as `import [module_name];`
	// +optional
	ExportMapping string `json:"exportMapping,omitempty"`

	// name contains a scoped npm package name.
	// +kubebuilder:validation:Required
	Name string `json:"name"`

	// secretRef is a reference to a secret containing authentication credentials for the NPM registry that holds the package.
	// +optional
	SecretRef *corev1.LocalObjectReference `json:"secretRef,omitempty"`

	// version contains a specific npm package version.
	// +kubebuilder:validation:Required
	Version string `json:"version"`
}

func (p *PackageReference) ToImportStatement() string {
	var buffer bytes.Buffer

	buffer.WriteString("import ")
	if p.ExportMapping != "" {
		buffer.WriteString(p.ExportMapping)
		buffer.WriteString(" from ")
	}
	buffer.WriteString("\"")
	buffer.WriteString(p.Name)
	buffer.WriteString("\";")

	return buffer.String()
}

func (p *PackageReference) ToScriptTag() string {
	return fmt.Sprintf(importStatementTemplate, p.ToImportStatement())
}

// +kubebuilder:validation:XValidation:rule="[has(self.script), has(self.scriptSrc)].filter(x, x).size() == 1",message="script and scriptSrc are mutually exclusive"
type Script struct {
	// attributes are key/value pairs that will be added to the element when rendered.
	// +optional
	Attributes map[string]string `json:"attributes,omitempty"`

	// footScript is a flag for script or scriptSrc that indicates if the tag should be added in the head of the page or at the foot. The default is false (add to head). To add the script to the foot of the page set footScript to true.
	// +optional
	// +kubebuilder:default:=false
	FootScript bool `json:"footScript,omitempty"`

	// script is the text content to be added into a <script> element when rendered.
	// +optional
	Script string `json:"script,omitempty"`

	// scriptSrc must be an absolute URL with a protocol and host which can be used in a src attribute.
	// +optional
	ScriptSrc string `json:"scriptSrc,omitempty"`
}

func (s *Script) ToScriptTag(footScript bool) string {
	if !s.FootScript && footScript {
		return ""
	}

	var buffer bytes.Buffer

	if s.ScriptSrc != "" {
		buffer.WriteString("<script")
		for key, value := range s.Attributes {
			if key == src {
				continue
			}
			buffer.WriteRune(' ')
			buffer.WriteString(key)
			buffer.WriteString("=\"")
			buffer.WriteString(value)
			buffer.WriteRune('"')
		}
		buffer.WriteString(" src=\"")
		buffer.WriteString(s.ScriptSrc)
		buffer.WriteString("\"></script>")
	} else if s.Script != "" {
		buffer.WriteString("<script")
		for key, value := range s.Attributes {
			if key == src {
				continue
			}
			buffer.WriteRune(' ')
			buffer.WriteString(key)
			buffer.WriteString("=\"")
			buffer.WriteString(value)
			buffer.WriteRune('"')
		}
		buffer.WriteString(">\n")
		buffer.WriteString(s.Script)
		buffer.WriteString("\n</script>")
	}

	return buffer.String()
}

func init() {
	SchemeBuilder.Register(&KDexScriptLibrary{}, &KDexScriptLibraryList{})
}
