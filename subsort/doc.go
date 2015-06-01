// Use case: We have a slice of any struct,
// and want to sort that slice by some field *in* the struct.
//
// We can standardize / generify this sorting,
// provided we stick to a few narrowing assumptions
//
// 1.) We want to sort slices of any struct
// 2.) We want to sort by a *string* or *int* field
// 3.) Float values could be string formatted %04.4d, with a huge base summand, avoiding negative values
// 4.) Because we need our slices in *several* sort orders at the same time,
//		we always want to *copy* the slice
// 5.) To minimize memory footprint, we only copy the sort field
//      and an int reference to the original array.
//      Thus we can retrieve other values from the original slice
// 6.) It's expensive to convert slices of various types to
//		slices of interfaces to slices of SortedByVal, that is:
//			[]TypeX to []interface{} to []SortedByVal
//		Therefore we resort to creating desired []SortedByVal outside.
// 7.) Similar mechanisms for float, time could be added
// 8.) Maintenance by file diffing

// We don't need to return pointers to the sorted slices,
// since slices are already references to the underlying array.

// SortedBy[String/Int]Val sadly needs to be exported,
// because the sorted slices will be given to the outside callee
// with type []SortedBy...Val

// I tried to cascade the sort interface implementations,
// using a "base" type for .Len() and have the others wrap it,
// but strangely, compiler demands implementation at the top level.
// Thus we end up with *four* identical implementations of Len(),Swap()
// and four *almost identical* implementations of Less()

package subsort
