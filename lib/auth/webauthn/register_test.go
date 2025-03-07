// Copyright 2021 Gravitational, Inc
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package webauthn_test

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/gravitational/teleport/api/types"
	"github.com/gravitational/teleport/lib/auth/mocku2f"
	"github.com/gravitational/trace"
	"github.com/stretchr/testify/require"

	wanlib "github.com/gravitational/teleport/lib/auth/webauthn"
)

func TestRegistrationFlow_BeginFinish(t *testing.T) {
	const user = "llama"
	const rpID = "localhost"
	const origin = "https://localhost"
	identity := newFakeIdentity(user)
	webRegistration := wanlib.RegistrationFlow{
		Webauthn: &types.Webauthn{
			RPID: rpID,
		},
		Identity: identity,
	}

	ctx := context.Background()

	// Begin is the first step in registration.
	credentialCreation, err := webRegistration.Begin(ctx, user)
	require.NoError(t, err)
	// Assert some parts of the credentialCreation (it's framework-created, no
	// need to go too deep).
	require.NotNil(t, credentialCreation)
	require.NotEmpty(t, credentialCreation.Response.Challenge)
	require.Equal(t, webRegistration.Webauthn.RPID, credentialCreation.Response.RelyingParty.ID)
	// Did we record the SessionData in storage?
	require.NotEmpty(t, identity.SessionData)

	// Sign CredentialCreation, typically requires user interaction.
	dev, err := mocku2f.Create()
	require.NoError(t, err)
	ccr, err := dev.SignCredentialCreation(origin, credentialCreation)
	require.NoError(t, err)

	// Finish is the final step in registration.
	newDevice, err := webRegistration.Finish(ctx, user, "webauthn1", ccr)
	require.NoError(t, err)
	// Did we get a proper WebauthnDevice?
	gotDevice := newDevice.GetWebauthn()
	require.NotNil(t, gotDevice)
	require.NotEmpty(t, gotDevice.PublicKeyCbor) // validated indirectly via authentication
	wantDevice := &types.WebauthnDevice{
		CredentialId:     dev.KeyHandle,
		PublicKeyCbor:    gotDevice.PublicKeyCbor,
		AttestationType:  gotDevice.AttestationType,
		Aaguid:           make([]byte, 16), // 16 zeroes
		SignatureCounter: 0,
	}
	if diff := cmp.Diff(wantDevice, gotDevice); diff != "" {
		t.Errorf("Finish() mismatch (-want +got):\n%s", diff)
	}
	// SessionData was cleared?
	require.Empty(t, identity.SessionData)
	// Device created in storage?
	require.Len(t, identity.UpdatedDevices, 1)
	require.Equal(t, newDevice, identity.UpdatedDevices[0])
}

func TestRegistrationFlow_Begin_errors(t *testing.T) {
	const user = "llama"
	webRegistration := &wanlib.RegistrationFlow{
		Webauthn: &types.Webauthn{
			RPID: "localhost",
		},
		Identity: newFakeIdentity(user),
	}

	ctx := context.Background()
	_, err := webRegistration.Begin(ctx, "" /* user */)
	require.True(t, trace.IsBadParameter(err)) // user required
}

func TestRegistrationFlow_Finish_errors(t *testing.T) {
	const user = "llama"
	const webOrigin = "https://localhost"
	webRegistration := &wanlib.RegistrationFlow{
		Webauthn: &types.Webauthn{
			RPID: "localhost",
		},
		Identity: newFakeIdentity(user),
	}

	ctx := context.Background()

	cc, err := webRegistration.Begin(ctx, user)
	require.NoError(t, err)
	key, err := mocku2f.Create()
	require.NoError(t, err)
	okCCR, err := key.SignCredentialCreation(webOrigin, cc)
	require.NoError(t, err)

	tests := []struct {
		name             string
		user, deviceName string
		createResp       func() *wanlib.CredentialCreationResponse
		wantErr          string
	}{
		{
			name:       "NOK user empty",
			user:       "",
			deviceName: "webauthn2",
			createResp: func() *wanlib.CredentialCreationResponse { return okCCR },
			wantErr:    "user required",
		},
		{
			name:       "NOK device name empty",
			user:       user,
			deviceName: "",
			createResp: func() *wanlib.CredentialCreationResponse { return okCCR },
			wantErr:    "device name required",
		},
		{
			name:       "NOK credential response nil",
			user:       user,
			deviceName: "webauthn2",
			createResp: func() *wanlib.CredentialCreationResponse { return nil },
			wantErr:    "response required",
		},
		{
			name:       "NOK credential with bad origin",
			user:       user,
			deviceName: "webauthn2",
			createResp: func() *wanlib.CredentialCreationResponse {
				resp, err := key.SignCredentialCreation("https://alpacasarerad.com" /* origin */, cc)
				require.NoError(t, err)
				return resp
			},
			wantErr: "origin",
		},
		{
			name:       "NOK credential with bad RPID",
			user:       user,
			deviceName: "webauthn2",
			createResp: func() *wanlib.CredentialCreationResponse {
				cc, err := webRegistration.Begin(ctx, user)
				require.NoError(t, err)
				cc.Response.RelyingParty.ID = "badrpid.com"

				resp, err := key.SignCredentialCreation(webOrigin, cc)
				require.NoError(t, err)
				return resp
			},
			wantErr: "authenticator response",
		},
		{
			name:       "NOK credential with invalid signature",
			user:       user,
			deviceName: "webauthn2",
			createResp: func() *wanlib.CredentialCreationResponse {
				cc, err := webRegistration.Begin(ctx, user)
				require.NoError(t, err)
				// Flip a challenge bit, this should be enough to consistently fail
				// signature checking.
				cc.Response.Challenge[0] = 1 ^ cc.Response.Challenge[0]

				resp, err := key.SignCredentialCreation(webOrigin, cc)
				require.NoError(t, err)
				return resp
			},
			wantErr: "validating challenge",
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := webRegistration.Finish(ctx, test.user, test.deviceName, test.createResp())
			require.Error(t, err)
			require.Contains(t, err.Error(), test.wantErr)
		})
	}
}
