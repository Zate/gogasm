//  Copyright (c) 2014 Marty Schoch
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package gouchstore

type DocumentWalkFun func(db *Gouchstore, di *DocumentInfo, doc *Document) error

// Walk the DB from a specific location including the complete docs.
func (db *Gouchstore) WalkDocs(startkey, endkey string, callback DocumentWalkFun) error {

	return db.AllDocuments(startkey, endkey, func(fdb *Gouchstore, di *DocumentInfo, context interface{}) error {
		doc, err := fdb.DocumentByDocumentInfo(di)
		if err != nil {
			return err
		}
		return callback(fdb, di, doc)
	}, nil)

}
