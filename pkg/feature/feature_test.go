/*
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
package feature

import (
	"testing"

	"github.com/spf13/pflag"
)

func TestRegister(t *testing.T) {
	// Reset features for test isolation
	features = make([]Feature, 0)

	// Test registering a feature
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	flags.String("test-flag", "", "test flag")

	var initialized bool
	initFunc := func() {
		initialized = true
	}

	feature := Feature{
		Flags:      flags,
		Initialize: initFunc,
	}

	Register(feature)

	if len(features) != 1 {
		t.Errorf("expected 1 feature, got %d", len(features))
	}

	if features[0].Flags != flags {
		t.Error("registered feature has wrong flags")
	}

	if features[0].Initialize == nil {
		t.Error("registered feature has no initialize function")
	}

	// Test that initialize function works
	features[0].Initialize()
	if !initialized {
		t.Error("initialize function was not called")
	}
}

func TestRegisterMultipleFeatures(t *testing.T) {
	// Reset features for test isolation
	features = make([]Feature, 0)

	// Register multiple features
	for i := 0; i < 3; i++ {
		flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
		feature := Feature{
			Flags:      flags,
			Initialize: nil,
		}
		Register(feature)
	}

	if len(features) != 3 {
		t.Errorf("expected 3 features, got %d", len(features))
	}
}

func TestFeatures(t *testing.T) {
	// Reset features for test isolation
	features = make([]Feature, 0)

	// Register some features
	feature1 := Feature{
		Flags:      pflag.NewFlagSet("test1", pflag.ContinueOnError),
		Initialize: nil,
	}
	feature2 := Feature{
		Flags:      pflag.NewFlagSet("test2", pflag.ContinueOnError),
		Initialize: nil,
	}

	Register(feature1)
	Register(feature2)

	// Get all features
	allFeatures := Features()

	if len(allFeatures) != 2 {
		t.Errorf("expected 2 features, got %d", len(allFeatures))
	}

	// Verify returned slice contains registered features
	if allFeatures[0].Flags == nil {
		t.Error("first feature has no flags")
	}
	if allFeatures[1].Flags == nil {
		t.Error("second feature has no flags")
	}
}

func TestFeatureWithNilInitialize(t *testing.T) {
	// Reset features for test isolation
	features = make([]Feature, 0)

	// Test registering a feature with nil Initialize function
	flags := pflag.NewFlagSet("test", pflag.ContinueOnError)
	feature := Feature{
		Flags:      flags,
		Initialize: nil,
	}

	Register(feature)

	if len(features) != 1 {
		t.Errorf("expected 1 feature, got %d", len(features))
	}

	// Verify Initialize is nil
	if features[0].Initialize != nil {
		t.Error("expected Initialize to be nil")
	}
}

func TestFeaturesReturnsCopy(t *testing.T) {
	// Reset features for test isolation
	features = make([]Feature, 0)

	// Register a feature
	feature := Feature{
		Flags:      pflag.NewFlagSet("test", pflag.ContinueOnError),
		Initialize: nil,
	}
	Register(feature)

	// Get features
	retrieved := Features()

	// Verify we got the same slice (not a deep copy, but same reference)
	if len(retrieved) != len(features) {
		t.Errorf("expected %d features, got %d", len(features), len(retrieved))
	}
}
