# Modmerge
Merge BGEE/BG2EE v2.0 mod zip files back into the main game(making them weidu
compatible).

# Pre-requisites

You must have **bought** and **downloaded**:

- Baldur's Gate 1: Enhanced Edition
- Baldur's Gate 1: Siege of Dragonspear
- Baldur's Gate 2: Enhanced Edition

# Running

1. Download binary from http://github.com/ScottBrooks/modmerge/releases
2. Unzip.
   1. MacOSX only: Rename `/path/to/modmerge-osx` to `/path/to/modmerge`.
3. Navigate to `/path/to/BG1EE`.
   - There should be a `chitin.key` file in the directory.
   - If downloaded with Steam on MacOSX, this should be at `~/Library/Application Support/Steam/steamapps/common/Baldur's Gate Enhanced Edition`.
   - If downloaded with Steam on Windows, this should be at `C:\Program Files (x86)\Steam\steamapps\common\Baldur's Gate Enhanced Edition`.
4. Run `/path/to/modmerge`. If you see `Conversion complete.`, it worked!

# Details

modmerge is a tool to help people merge the mod back into the main game, making it weidu compatible.

modmerge does the following
 - Backs up your chitin.key to chitin.key.bak, just in case
 - Searches dlc/ or your root directory for sod-dlc.zip
 - Unzips the sod-dlc.zip, copying over top of the existing lang/, movies/, music/ folders.  The data folder inside of the zip will be renamed to no overwrite files in the existing data directory.
 - Updates your chitin.key with so any resources from inside of the zip file will point to their appropriate bif files.
 - Renames your sod-dlc.zip to sod-dlc.disabled
 
The end result, assuming everything is successful, is you will have an updated chitin.key that points to the resources from the base game, and the resources from the mod, that weidu, near infinity, DLTCEP, etc will all be able to load.

# Building from source

```
cd /path/to/modmerge
go build ./...
```

# Redistribution

You are free to package or redistribute modmerge.
