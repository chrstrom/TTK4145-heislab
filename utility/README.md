Utility
================
This package contains utility that extends golangs stock capabilities, that don't really fit in a specific
module in the elevator system.

Also contains a shell script to clear the logs (that can grow to be decently large ): )

### Interface

StringArray2BoolArray(s []string) []bool   
FindMostCommonElement(s []string) (element string, count int)  
IsStringInSlice(s string, slice []string) bool  
