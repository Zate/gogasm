//  Copyright (c) 2014 Marty Schoch
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package mergesort

type recordsList struct {
	records []interface{}
	compare CompareRecords
	context interface{}
}

func newRecordsList(expectedSize int, compare CompareRecords, context interface{}) *recordsList {
	return &recordsList{
		records: make([]interface{}, 0, expectedSize),
		compare: compare,
		context: context,
	}
}

func (rl *recordsList) add(r interface{}) {
	rl.records = append(rl.records, r)
}

func (rl *recordsList) Len() int      { return len(rl.records) }
func (rl *recordsList) Swap(i, j int) { rl.records[i], rl.records[j] = rl.records[j], rl.records[i] }
func (rl *recordsList) Less(i, j int) bool {
	return rl.compare(rl.records[i], rl.records[j], rl.context) < 0
}
