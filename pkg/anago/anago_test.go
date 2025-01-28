/*
Copyright 2020 The Kubernetes Authors.

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

package anago_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"

	"sigs.k8s.io/release-sdk/git"

	"k8s.io/release/pkg/anago"
	"k8s.io/release/pkg/anago/anagofakes"
	"k8s.io/release/pkg/release"
)

var err = errors.New("error")

var testVersionTag = "v1.20.0"

type testStateParameters struct {
	versionsTag         *string
	createReleaseBranch *bool
}

func mockGenerateReleaseVersionStage(mock *anagofakes.FakeStageClient) {
	mock.GenerateReleaseVersionReturns(nil)
}

func mockGenerateReleaseVersionRelease(mock *anagofakes.FakeReleaseClient) {
	mock.GenerateReleaseVersionReturns(nil)
}

//nolint:dupl // duplications in those tests are intentional
func TestRunStage(t *testing.T) {
	for _, tc := range []struct {
		prepare     func(*anagofakes.FakeStageClient)
		shouldError bool
	}{
		{ // success
			prepare: func(mock *anagofakes.FakeStageClient) {
				mockGenerateReleaseVersionStage(mock)
			},
			shouldError: false,
		},
		{ // InitLogFile fails
			prepare: func(mock *anagofakes.FakeStageClient) {
				mock.InitLogFileReturns(err)
			},
			shouldError: true,
		},
		{ // ValidateOptions fails
			prepare: func(mock *anagofakes.FakeStageClient) {
				mock.ValidateOptionsReturns(err)
			},
			shouldError: true,
		},
		{ // CheckPrerequisites fails
			prepare: func(mock *anagofakes.FakeStageClient) {
				mock.CheckPrerequisitesReturns(err)
			},
			shouldError: true,
		},
		{ // CheckReleaseBranchState fails
			prepare: func(mock *anagofakes.FakeStageClient) {
				mock.CheckReleaseBranchStateReturns(err)
			},
			shouldError: true,
		},
		{ // GenerateReleaseVersion fails
			prepare: func(mock *anagofakes.FakeStageClient) {
				mock.GenerateReleaseVersionReturns(err)
			},
			shouldError: true,
		},
		{ // PrepareWorkspace fails
			prepare: func(mock *anagofakes.FakeStageClient) {
				mockGenerateReleaseVersionStage(mock)
				mock.PrepareWorkspaceReturns(err)
			},
			shouldError: true,
		},
		{ // Build fails
			prepare: func(mock *anagofakes.FakeStageClient) {
				mockGenerateReleaseVersionStage(mock)
				mock.BuildReturns(err)
			},
			shouldError: true,
		},
		{ // GenerateChangelog fails
			prepare: func(mock *anagofakes.FakeStageClient) {
				mockGenerateReleaseVersionStage(mock)
				mock.GenerateChangelogReturns(err)
			},
			shouldError: true,
		},
		{ // StageArtifacts fails
			prepare: func(mock *anagofakes.FakeStageClient) {
				mockGenerateReleaseVersionStage(mock)
				mock.StageArtifactsReturns(err)
			},
			shouldError: true,
		},
	} {
		opts := anago.DefaultStageOptions()
		sut := anago.NewStage(opts)
		mock := &anagofakes.FakeStageClient{}
		tc.prepare(mock)
		sut.SetClient(mock)

		err := sut.Run()
		if tc.shouldError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

//nolint:dupl // duplications in those tests are intentional
func TestRunRelease(t *testing.T) {
	for _, tc := range []struct {
		prepare     func(*anagofakes.FakeReleaseClient)
		shouldError bool
	}{
		{ // success
			prepare: func(mock *anagofakes.FakeReleaseClient) {
				mockGenerateReleaseVersionRelease(mock)
			},
			shouldError: false,
		},
		{ // InitLogFile fails
			prepare: func(mock *anagofakes.FakeReleaseClient) {
				mock.InitLogFileReturns(err)
			},
			shouldError: true,
		},
		{ // CheckPrerequisites fails
			prepare: func(mock *anagofakes.FakeReleaseClient) {
				mock.CheckPrerequisitesReturns(err)
			},
			shouldError: true,
		},
		{ // ValidateOptions fails
			prepare: func(mock *anagofakes.FakeReleaseClient) {
				mock.ValidateOptionsReturns(err)
			},
			shouldError: true,
		},
		{ // CheckReleaseBranchState fails
			prepare: func(mock *anagofakes.FakeReleaseClient) {
				mock.CheckReleaseBranchStateReturns(err)
			},
			shouldError: true,
		},
		{ // GenerateReleaseVersion fails
			prepare: func(mock *anagofakes.FakeReleaseClient) {
				mock.GenerateReleaseVersionReturns(err)
			},
			shouldError: true,
		},
		{ // PrepareWorkspace fails
			prepare: func(mock *anagofakes.FakeReleaseClient) {
				mockGenerateReleaseVersionRelease(mock)
				mock.PrepareWorkspaceReturns(err)
			},
			shouldError: true,
		},
		{ // PushArtifacts fails
			prepare: func(mock *anagofakes.FakeReleaseClient) {
				mockGenerateReleaseVersionRelease(mock)
				mock.PushArtifactsReturns(err)
			},
			shouldError: true,
		},
		{ // PushGitObjects fails
			prepare: func(mock *anagofakes.FakeReleaseClient) {
				mockGenerateReleaseVersionRelease(mock)
				mock.PushGitObjectsReturns(err)
			},
			shouldError: true,
		},
		{ // CreateAnnouncement fails
			prepare: func(mock *anagofakes.FakeReleaseClient) {
				mockGenerateReleaseVersionRelease(mock)
				mock.CreateAnnouncementReturns(err)
			},
			shouldError: true,
		},
	} {
		opts := anago.DefaultReleaseOptions()
		sut := anago.NewRelease(opts)
		mock := &anagofakes.FakeReleaseClient{}
		tc.prepare(mock)
		sut.SetClient(mock)

		err := sut.Run()
		if tc.shouldError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestValidateOptions(t *testing.T) {
	for _, tc := range []struct {
		provided    *anago.Options
		shouldError bool
	}{
		{ // success
			provided: &anago.Options{
				ReleaseType:   release.ReleaseTypeAlpha,
				ReleaseBranch: git.DefaultBranch,
			},
			shouldError: false,
		},
		{ // invalid release type
			provided: &anago.Options{
				ReleaseType: "invalid",
			},
			shouldError: true,
		},
		{ // invalid release branch
			provided: &anago.Options{
				ReleaseType:   release.ReleaseTypeAlpha,
				ReleaseBranch: "invalid",
			},
			shouldError: true,
		},
	} {
		err := tc.provided.Validate()
		if tc.shouldError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestValidateBuildVersion(t *testing.T) {
	for _, tc := range []struct {
		provided    *anago.Options
		shouldError bool
	}{
		{ // success
			provided: &anago.Options{
				ReleaseType:   release.ReleaseTypeAlpha,
				ReleaseBranch: git.DefaultBranch,
				BuildVersion:  "v1.20.0-beta.1.203+8f6ffb24df9896",
			},
			shouldError: false,
		},
		{ // invalid build version
			provided: &anago.Options{
				ReleaseType:   release.ReleaseTypeAlpha,
				ReleaseBranch: git.DefaultBranch,
				BuildVersion:  "invalid",
			},
			shouldError: true,
		},
	} {
		state := anago.DefaultState()

		err := tc.provided.ValidateBuildVersion(state)
		if tc.shouldError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestStagingOptionsValidate(t *testing.T) {
	for _, tc := range []struct {
		provided    *anago.StageOptions
		shouldError bool
	}{
		{ // valid build version should validate
			provided: &anago.StageOptions{
				&anago.Options{
					ReleaseType:   release.ReleaseTypeAlpha,
					ReleaseBranch: git.DefaultBranch,
					BuildVersion:  "v1.20.0-beta.1.203+8f6ffb24df9896",
				},
			},
			shouldError: false,
		},
		{ // empty build version should validate
			provided: &anago.StageOptions{
				&anago.Options{
					ReleaseType:   release.ReleaseTypeAlpha,
					ReleaseBranch: git.DefaultBranch,
				},
			},
			shouldError: false,
		},
		{ // invalid build version should not validate
			provided: &anago.StageOptions{
				&anago.Options{
					ReleaseType:   release.ReleaseTypeAlpha,
					ReleaseBranch: git.DefaultBranch,
					BuildVersion:  "decaf-bad",
				},
			},
			shouldError: true,
		},
	} {
		state := anago.DefaultState()

		err := tc.provided.Validate(state)
		if tc.shouldError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestSubmitStage(t *testing.T) {
	for _, tc := range []struct {
		prepare     func(*anagofakes.FakeStageClient)
		shouldError bool
	}{
		{ // success
			prepare: func(mock *anagofakes.FakeStageClient) {
				mockGenerateReleaseVersionStage(mock)
			},
			shouldError: false,
		},
		{ // Submit fails
			prepare: func(mock *anagofakes.FakeStageClient) {
				mock.SubmitReturns(err)
			},
			shouldError: true,
		},
	} {
		opts := anago.DefaultStageOptions()
		sut := anago.NewStage(opts)
		mock := &anagofakes.FakeStageClient{}
		tc.prepare(mock)
		sut.SetClient(mock)

		err := sut.Submit(false)
		if tc.shouldError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}

func TestSubmitRelease(t *testing.T) {
	for _, tc := range []struct {
		prepare     func(*anagofakes.FakeReleaseClient)
		shouldError bool
	}{
		{ // success
			prepare: func(mock *anagofakes.FakeReleaseClient) {
				mockGenerateReleaseVersionRelease(mock)
			},
			shouldError: false,
		},
		{ // Submit fails
			prepare: func(mock *anagofakes.FakeReleaseClient) {
				mock.SubmitReturns(err)
			},
			shouldError: true,
		},
	} {
		opts := anago.DefaultReleaseOptions()
		sut := anago.NewRelease(opts)
		mock := &anagofakes.FakeReleaseClient{}
		tc.prepare(mock)
		sut.SetClient(mock)

		err := sut.Submit(false)
		if tc.shouldError {
			require.Error(t, err)
		} else {
			require.NoError(t, err)
		}
	}
}
