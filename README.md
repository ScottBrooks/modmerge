# Modmerge
Merge BGEE/BG2EE v2.0 mod zip files back into the main game(making them weidu compatible).

# Downloads/Binaries
http://www.scottbrooks.ca/modmerge/

# Source build
 - Install golang from http://golang.org
 - go build github.com/ScottBrooks/modmerge/modmerge

# Details

modmerge is a tool to help people merge the mod back into the main game, making it weidu compatible.

modmerge does the following
 - Backs up your chitin.key to chitin.key.bak, just in case
 - Unzips your dlc/sod-dlc.zip, copying over top of the existing lang/, movies/, music/ folders.  The data folder inside of the zip will be renamed to no overwrite files in the existing data directory.
 - Updates your chitin.key with so any resources from inside of the zip file will point to their appropriate bif files.
 - Renames your dlc/sod-dlc.zip to dlc/sod-dlc.disabled
 
The end result, assuming everything is successful, is you will have an updated chitin.key that points to the resources from the base game, and the resources from the mod, that weidu, near infinity, DLTCEP, etc will all be able to load.

