# Modmerge
Merge BGEE/BG2EE v2.0 mod zip files back into the main game(making them weidu compatible).

# Downloads/Binaries
http://github.com/ScottBrooks/modmerge/releases

# Source build
 - Install golang from http://golang.org
 - go build github.com/ScottBrooks/modmerge/modmerge

# Details

modmerge is a tool to help people merge the mod back into the main game, making it weidu compatible.

modmerge does the following
 - Backs up your chitin.key to chitin.key.bak, just in case
 - Searches dlc/ or your root directory for sod-dlc.zip
 - Unzips the sod-dlc.zip, copying over top of the existing lang/, movies/, music/ folders.  The data folder inside of the zip will be renamed to no overwrite files in the existing data directory.
 - Updates your chitin.key with so any resources from inside of the zip file will point to their appropriate bif files.
 - Renames your sod-dlc.zip to sod-dlc.disabled
 
The end result, assuming everything is successful, is you will have an updated chitin.key that points to the resources from the base game, and the resources from the mod, that weidu, near infinity, DLTCEP, etc will all be able to load.

# Redistribution

You are free to package or redistribute modmerge.

# License

Copyright (c) 2016 IdeaSpark Labs.

Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the "Software"), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.



