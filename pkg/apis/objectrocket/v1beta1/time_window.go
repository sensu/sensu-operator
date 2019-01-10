// Copyright 2017 The etcd-operator Authors
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

package v1beta1

// TimeWindowWhen defines the "when" attributes for time windows
type TimeWindowWhen struct {
	// Days is a hash of days
	Days TimeWindowDays `json:"days"`
}

// TimeWindowDays defines the days of a time window
type TimeWindowDays struct {
	All       []*TimeWindowTimeRange `json:"all,omitempty"`
	Sunday    []*TimeWindowTimeRange `json:"sunday,omitempty"`
	Monday    []*TimeWindowTimeRange `json:"monday,omitempty"`
	Tuesday   []*TimeWindowTimeRange `json:"tuesday,omitempty"`
	Wednesday []*TimeWindowTimeRange `json:"wednesday,omitempty"`
	Thursday  []*TimeWindowTimeRange `json:"thursday,omitempty"`
	Friday    []*TimeWindowTimeRange `json:"friday,omitempty"`
	Saturday  []*TimeWindowTimeRange `json:"saturday,omitempty"`
}

// TimeWindowTimeRange defines the time ranges of a time
type TimeWindowTimeRange struct {
	// Begin is the time which the time window should begin, in the format
	// '3:00PM', which satisfies the time.Kitchen format
	Begin string `json:"begin"`
	// End is the time which the filter should end, in the format '3:00PM', which
	// satisfies the time.Kitchen format
	End string `json:"end"`
}
