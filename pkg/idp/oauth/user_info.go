// Copyright 2022 Paul Greenberg greenpau@outlook.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package oauth

import (
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
)

func (b *IdentityProvider) fetchUserInfo(tokenData, userData map[string]interface{}) error {
	// The fetching of user info happens only if the below conditions are
	// are met.
	if b.config.Driver != "generic" || !b.ScopeExists("openid") || b.userInfoURL == "" {
		return nil
	}
	if len(b.userInfoFields) == 0 {
		return nil
	}
	if tokenData == nil || userData == nil {
		return nil
	}
	if _, exists := tokenData["access_token"]; !exists {
		return fmt.Errorf("access_token not found")
	}

	// Initialize HTTP client.
	cli, err := b.newBrowser()
	if err != nil {
		return err
	}

	req, err := http.NewRequest("GET", b.userInfoURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Accept", "application/json")
	req.Header.Add("Authorization", "Bearer "+tokenData["access_token"].(string))

	// Fetch data from the URL.
	resp, err := cli.Do(req)
	if err != nil {
		return err
	}

	respBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}

	b.logger.Debug(
		"User info received",
		zap.Any("body", respBody),
		zap.String("url", b.userInfoURL),
	)

	userinfo := make(map[string]interface{})
	if err := json.Unmarshal(respBody, &userinfo); err != nil {
		return err
	}

	b.logger.Debug(
		"User info received",
		zap.Any("userinfo", userinfo),
		zap.String("url", b.userInfoURL),
	)

	if _, exists := b.userInfoFields["all"]; exists {
		userData["userinfo"] = userinfo
	} else {
		for k := range userinfo {
			if _, exists := b.userInfoFields[k]; !exists {
				delete(userinfo, k)
			}
		}
		if len(userinfo) > 0 {
			userData["userinfo"] = userinfo
		}
	}
	return nil
}
