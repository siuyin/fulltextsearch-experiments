// Code generated by "stringer -type=docIndex"; DO NOT EDIT.

package doc

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[ShowID-0]
	_ = x[Type-1]
	_ = x[Title-2]
	_ = x[Director-3]
	_ = x[Cast-4]
	_ = x[Country-5]
	_ = x[DateAdded-6]
	_ = x[ReleaseYear-7]
	_ = x[Rating-8]
	_ = x[Duration-9]
	_ = x[ListedIn-10]
	_ = x[Description-11]
}

const _docIndex_name = "ShowIDTypeTitleDirectorCastCountryDateAddedReleaseYearRatingDurationListedInDescription"

var _docIndex_index = [...]uint8{0, 6, 10, 15, 23, 27, 34, 43, 54, 60, 68, 76, 87}

func (i docIndex) String() string {
	if i < 0 || i >= docIndex(len(_docIndex_index)-1) {
		return "docIndex(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _docIndex_name[_docIndex_index[i]:_docIndex_index[i+1]]
}
