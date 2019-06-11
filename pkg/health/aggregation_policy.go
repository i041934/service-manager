/*
 * Copyright 2018 The Service Manager Authors
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package health

import "github.com/InVisionApp/go-health"

// DefaultAggregationPolicy aggregates the healths by constructing a new Health based on the given
// where the overall health status is negative if one of the healths is negative and positive if all are positive
type DefaultAggregationPolicy struct {
}

// Apply aggregates the given healths
func (*DefaultAggregationPolicy) Apply(healths map[string]health.State) *Health {
	if len(healths) == 0 {
		return New().WithDetail("error", "no health indicators registered").Unknown()
	}
	overallStatus := StatusUp
	for _, health := range healths {
		if health.Status == "failed" && health.ContiguousFailures > 3 {
			overallStatus = StatusDown
			break
		}
	}
	details := make(map[string]interface{})
	for k, v := range healths {
		details[k] = convertStatus(v.Status)
	}
	return New().WithStatus(overallStatus).WithDetails(details)
}

func convertStatus(status string) Status {
	switch status {
	case "ok":
		return StatusUp
	case "failed":
		return StatusDown
	default:
		return StatusUnknown
	}
}
