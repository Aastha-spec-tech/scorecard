// Copyright 2024 OpenSSF Scorecard Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package csproj

import (
	"errors"
	"testing"
)

func TestIsRestoreLockedModeEnabled(t *testing.T) {
	t.Parallel()

	//nolint:govet
	tests := []struct {
		name    string
		content []byte
		wantErr error
		want    bool
	}{
		{
			name: "RestoreLockedMode enabled",
			content: []byte(`<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <RestoreLockedMode>true</RestoreLockedMode>
  </PropertyGroup>
</Project>`),
			want:    true,
			wantErr: nil,
		},
		{
			name: "RestoreLockedMode disabled explicitly",
			content: []byte(`<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <RestoreLockedMode>false</RestoreLockedMode>
  </PropertyGroup>
</Project>`),
			want:    false,
			wantErr: nil,
		},
		{
			name: "RestoreLockedMode omitted implicitly",
			content: []byte(`<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <TargetFramework>netstandard2.0</TargetFramework>
  </PropertyGroup>
</Project>`),
			want:    false,
			wantErr: nil,
		},
		{
			name: "multiple PropertyGroups with RestoreLockedMode enabled in second group",
			content: []byte(`<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
  </PropertyGroup>
  <PropertyGroup>
    <RestoreLockedMode>true</RestoreLockedMode>
  </PropertyGroup>
</Project>`),
			want:    true,
			wantErr: nil,
		},
		{
			name: "multiple PropertyGroups without RestoreLockedMode",
			content: []byte(`<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
  </PropertyGroup>
  <PropertyGroup>
    <Nullable>enable</Nullable>
  </PropertyGroup>
</Project>`),
			want:    false,
			wantErr: nil,
		},
		{
			name: "empty PropertyGroup",
			content: []byte(`<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup></PropertyGroup>
</Project>`),
			want:    false,
			wantErr: nil,
		},
		{
			name:    "empty Project tag",
			content: []byte(`<Project></Project>`),
			want:    false,
			wantErr: nil,
		},
		{
			name: "extra tags and attributes with RestoreLockedMode enabled",
			content: []byte(`<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <OutputType>Exe</OutputType>
    <TargetFramework>net8.0</TargetFramework>
    <RestoreLockedMode>true</RestoreLockedMode>
  </PropertyGroup>
  <ItemGroup>
    <PackageReference Include="Newtonsoft.Json" Version="13.0.3" />
  </ItemGroup>
</Project>`),
			want:    true,
			wantErr: nil,
		},
		{
			name:    "malformed XML unclosed tag",
			content: []byte(`<Project <PropertyGroup>`),
			want:    false,
			wantErr: errInvalidCsProjFile,
		},
		{
			name:    "invalid XML content",
			content: []byte(`not-xml-content`),
			want:    false,
			wantErr: errInvalidCsProjFile,
		},
		{
			name:    "empty byte content",
			content: []byte(``),
			want:    false,
			wantErr: errInvalidCsProjFile,
		},
		{
			name: "invalid boolean value in RestoreLockedMode",
			content: []byte(`<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <RestoreLockedMode>invalid-bool</RestoreLockedMode>
  </PropertyGroup>
</Project>`),
			want:    false,
			wantErr: errInvalidCsProjFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := IsRestoreLockedModeEnabled(tt.content)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("IsRestoreLockedModeEnabled() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("IsRestoreLockedModeEnabled() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestIsAllowUnsafeBlocksEnabled(t *testing.T) {
	t.Parallel()

	//nolint:govet
	tests := []struct {
		name    string
		content []byte
		wantErr error
		want    bool
	}{
		{
			name: "AllowUnsafeBlocks enabled",
			content: []byte(`<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <AllowUnsafeBlocks>true</AllowUnsafeBlocks>
  </PropertyGroup>
</Project>`),
			want:    true,
			wantErr: nil,
		},
		{
			name: "AllowUnsafeBlocks disabled explicitly",
			content: []byte(`<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <AllowUnsafeBlocks>false</AllowUnsafeBlocks>
  </PropertyGroup>
</Project>`),
			want:    false,
			wantErr: nil,
		},
		{
			name: "AllowUnsafeBlocks omitted implicitly",
			content: []byte(`<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
  </PropertyGroup>
</Project>`),
			want:    false,
			wantErr: nil,
		},
		{
			name: "multiple PropertyGroups with AllowUnsafeBlocks in second group",
			content: []byte(`<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
  </PropertyGroup>
  <PropertyGroup>
    <AllowUnsafeBlocks>true</AllowUnsafeBlocks>
  </PropertyGroup>
</Project>`),
			want:    true,
			wantErr: nil,
		},
		{
			name: "multiple PropertyGroups without AllowUnsafeBlocks",
			content: []byte(`<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <TargetFramework>net8.0</TargetFramework>
  </PropertyGroup>
  <PropertyGroup>
    <Nullable>enable</Nullable>
  </PropertyGroup>
</Project>`),
			want:    false,
			wantErr: nil,
		},
		{
			name: "empty PropertyGroup",
			content: []byte(`<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup></PropertyGroup>
</Project>`),
			want:    false,
			wantErr: nil,
		},
		{
			name:    "empty Project tag",
			content: []byte(`<Project></Project>`),
			want:    false,
			wantErr: nil,
		},
		{
			name: "extra tags and attributes with AllowUnsafeBlocks enabled",
			content: []byte(`<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <OutputType>Exe</OutputType>
    <TargetFramework>net8.0</TargetFramework>
    <AllowUnsafeBlocks>true</AllowUnsafeBlocks>
  </PropertyGroup>
  <ItemGroup>
    <PackageReference Include="Newtonsoft.Json" Version="13.0.3" />
  </ItemGroup>
</Project>`),
			want:    true,
			wantErr: nil,
		},
		{
			name:    "malformed XML unclosed tag",
			content: []byte(`<Project <PropertyGroup>`),
			want:    false,
			wantErr: errInvalidCsProjFile,
		},
		{
			name:    "invalid XML content",
			content: []byte(`not-xml-content`),
			want:    false,
			wantErr: errInvalidCsProjFile,
		},
		{
			name:    "empty byte content",
			content: []byte(``),
			want:    false,
			wantErr: errInvalidCsProjFile,
		},
		{
			name: "invalid boolean value in AllowUnsafeBlocks",
			content: []byte(`<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <AllowUnsafeBlocks>invalid-bool</AllowUnsafeBlocks>
  </PropertyGroup>
</Project>`),
			want:    false,
			wantErr: errInvalidCsProjFile,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := IsAllowUnsafeBlocksEnabled(tt.content)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("IsAllowUnsafeBlocksEnabled() error = %v, wantErr = %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("IsAllowUnsafeBlocksEnabled() got = %v, want = %v", got, tt.want)
			}
		})
	}
}

func TestBothPropertiesInDifferentPropertyGroups(t *testing.T) {
	t.Parallel()

	content := []byte(`<Project Sdk="Microsoft.NET.Sdk">
  <PropertyGroup>
    <RestoreLockedMode>true</RestoreLockedMode>
  </PropertyGroup>
  <PropertyGroup>
    <AllowUnsafeBlocks>true</AllowUnsafeBlocks>
  </PropertyGroup>
</Project>`)

	locked, err := IsRestoreLockedModeEnabled(content)
	if err != nil {
		t.Fatalf("unexpected error for IsRestoreLockedModeEnabled: %v", err)
	}
	if !locked {
		t.Errorf("expected RestoreLockedMode to be true, got false")
	}

	unsafe, err := IsAllowUnsafeBlocksEnabled(content)
	if err != nil {
		t.Fatalf("unexpected error for IsAllowUnsafeBlocksEnabled: %v", err)
	}
	if !unsafe {
		t.Errorf("expected AllowUnsafeBlocks to be true, got false")
	}
}
