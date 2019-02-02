//  Copyright (c) 2014 Marty Schoch
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

/*
Package mergesort is a library for performing an on-disk merge sort.

You provide a method to read, write, and compare records.  The library will take a file
containing unsorted records, and produce a file containing sorted records.
*/
package mergesort

import (
	"io"
	"io/ioutil"
	"os"
	"sort"
)

// ReadRecord is responsible for reading a single record from the file.  If an error occurs,
// return nil and the error.  If the end of the file is reached, return io.EOF.  The record
// object can be anything, but WriteRecord and CompareRecord calls in the same MergeSort
// operation must be able to work with it.  The optional context object passed to the
// original MergeSort operation is provided.
type ReadRecord func(file *os.File, context interface{}) (interface{}, error)

// WriteRecord is responsible for writing a single record to the file.  If an error occurs,
// return the error, otherwise nil.  The record object could be anything produced by another
// ReadRecord call in the same MergeSort operation.  The optional context object passed to
// the original MergeSort operation is provided.
type WriteRecord func(file *os.File, record interface{}, context interface{}) error

// CompareRecords is responsible for ordering records.  Two records are provided, if rec1 is
// less than rec2, return -1.  If rec1 is greater than rec2, return 1.  If they are the same,
// return 0.  The record objects could be anything produced by another ReadRecord call in the
// same MergeSort operation.  The optional context object passed ot the original MergeSort
// operation is provided.
type CompareRecords func(rec1, rec2 interface{}, context interface{}) int

// MergeSort takes records from the unsortedFile, sorts them, and produces a sortedFile
// containing the same records in sorted order.  To do this, the call provides three
// methods ReadRecord, WriteRecord, and CompareRecords.  You can optionally provide a
// context object, this object will be included in all calls to ReadRecord, WriteRecord,
// CompareRecords.  The blockSize parameter will limit the total number of records that
// are sorted in memory at a single time.
//
// NOTE: sortedFile is allowed to be the same as unsortedFile.
//
// NOTE: Be sure to rewind unsortedFile first if just wrote it out.  See
// https://github.com/mschoch/mergesort/issues/1
func MergeSort(unsortedFile, sortedFile *os.File, read ReadRecord, write WriteRecord, compare CompareRecords, context interface{}, blockSize int) error {
	var err error
	sourceTape := make([]tape, 2)
	record := make([]interface{}, 2)

	// create temporary files sourceTape[0] and sourceTape[1]
	sourceTape[0].fp, err = ioutil.TempFile("", "goms")
	if err != nil {
		return err
	}
	defer os.Remove(sourceTape[0].fp.Name())
	sourceTape[1].fp, err = ioutil.TempFile("", "goms")
	if err != nil {
		return err
	}
	defer os.Remove(sourceTape[1].fp.Name())

	// read blocks, sort them in memory, and write the alternately to tapes 0 and 1
	blockCount := 0
	destination := 0
	list := newRecordsList(blockSize, compare, context)
	for {
		record[0], err = read(unsortedFile, context)
		if err != nil && err != io.EOF {
			// error reading, return
			return err
		}
		if err == nil {
			// not EOF, add record to in memory list
			list.add(record[0])
			blockCount++
		}
		if blockCount == blockSize || err == io.EOF && blockCount != 0 {
			// sort the in memory list
			sort.Sort(list)
			// now write them out
			for _, rec := range list.records {
				err := write(sourceTape[destination].fp, rec, context)
				if err != nil {
					return err
				}
				sourceTape[destination].count++
			}
			list = newRecordsList(blockSize, compare, context)
			destination ^= 1 // toggle tape
			blockCount = 0
		}
		if err == io.EOF {
			break // all done
		}
	}
	if sortedFile == unsortedFile {
		unsortedFile.Seek(0, os.SEEK_SET)
	}
	sourceTape[0].fp.Seek(0, os.SEEK_SET)
	sourceTape[1].fp.Seek(0, os.SEEK_SET)

	// FIXME (what?) delete the unsorted file here, if required (see instructions)

	if sourceTape[1].count == 0 {
		// handle case where memory sort is all that is required
		err = sourceTape[1].fp.Close()
		if err != nil {
			return err
		}
		sourceTape[1] = sourceTape[0]
		sourceTape[0].fp = sortedFile
		for sourceTape[1].count != 0 {
			record[0], err = read(sourceTape[1].fp, context)
			if err != nil {
				return err
			}
			err := write(sourceTape[0].fp, record[0], context)
			if err != nil {
				return err
			}
			sourceTape[1].count--
		}
	} else {
		// merge tapes, two by two, until every record is in source_tape[0]
		for sourceTape[1].count != 0 {
			destination := 0
			destinationTape := make([]tape, 2)
			if sourceTape[0].count <= blockSize {
				destinationTape[0].fp = sortedFile
			} else {
				destinationTape[0].fp, err = ioutil.TempFile("", "goms")
				if err != nil {
					return err
				}
				defer os.Remove(destinationTape[0].fp.Name())
			}
			destinationTape[1].fp, err = ioutil.TempFile("", "goms")
			if err != nil {
				return err
			}
			defer os.Remove(destinationTape[1].fp.Name())
			record[0], err = read(sourceTape[0].fp, context)
			if err != nil {
				return err
			}
			record[1], err = read(sourceTape[1].fp, context)
			if err != nil {
				return err
			}
			for sourceTape[0].count != 0 {
				count := make([]int, 2)
				count[0] = sourceTape[0].count
				if count[0] > blockSize {
					count[0] = blockSize
				}
				count[1] = sourceTape[1].count
				if count[1] > blockSize {
					count[1] = blockSize
				}
				for count[0]+count[1] != 0 {
					sel := 0
					if count[0] == 0 {
						sel = 1
					} else if count[1] == 0 {
						sel = 0
					} else if compare(record[0], record[1], context) < 0 {
						sel = 0
					} else {
						sel = 1
					}
					err = write(destinationTape[destination].fp, record[sel], context)
					if err != nil {
						return err
					}
					if sourceTape[sel].count > 1 {
						record[sel], err = read(sourceTape[sel].fp, context)
						if err != nil {
							return err
						}
					}
					sourceTape[sel].count--
					count[sel]--
					destinationTape[destination].count++
				}
				destination ^= 1
			}
			sourceTape[0].fp.Close()
			sourceTape[1].fp.Close()
			destinationTape[0].fp.Seek(0, os.SEEK_SET)
			destinationTape[1].fp.Seek(0, os.SEEK_SET)
			sourceTape[0] = destinationTape[0]
			sourceTape[1] = destinationTape[1]
			blockSize <<= 1
		}
	}
	sourceTape[1].fp.Close()
	return nil
}

type tape struct {
	fp    *os.File
	count int
}
