package main

import (
	"bytes"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"strings"
)

type RESREF struct {
	Name [8]byte
}

func NewResref(name string) RESREF {
	r := RESREF{}
	copy(r.Name[:], []byte(name))
	return r
}

func (r *RESREF) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (r *RESREF) Valid() bool {
	return r.String() != ""
}

func (r *RESREF) String() string {
	str := strings.Split(string(r.Name[0:]), "\x00")[0]
	return str
}

type keyHeader struct {
	Signature, Version [4]byte
	BifCount           uint32
	ResourceCount      uint32
	BifOffset          uint32
	ResourceOffset     uint32
}

type keyBifEntry struct {
	Length         uint32
	OffsetFilename uint32
	LengthFilename uint16
	FileLocation   uint16
}

type keyBifValue struct {
	Length       uint32
	Filename     string
	FileLocation uint16
}

type keyResourceEntry struct {
	Name     RESREF
	Type     uint16
	Location uint32
}
type keyUniqueResource struct {
	Name string
	Type uint16
}

type KEY struct {
	header    keyHeader
	bifs      []keyBifValue
	resources []keyResourceEntry
	r         io.ReadSeeker
	root      string
	files     map[keyUniqueResource]int
}

var fileTypes = map[string]int{
	"bmp":  1,
	"mve":  2,
	"tga":  3,
	"wav":  4,
	"wfx":  5,
	"plt":  6,
	"bam":  1000,
	"wed":  1001,
	"chu":  1002,
	"tis":  1003,
	"mos":  1004,
	"itm":  1005,
	"spl":  1006,
	"bcs":  1007,
	"ids":  1008,
	"cre":  1009,
	"are":  1010,
	"dlg":  1011,
	"2da":  1012,
	"gam":  1013,
	"sto":  1014,
	"wmp":  1015,
	"eff":  1016,
	"bs":   1017,
	"chr":  1018,
	"vvc":  1019,
	"vef":  1020,
	"pro":  1021,
	"bio":  1022,
	"wbm":  1023,
	"fnt":  1024,
	"gui":  1026,
	"sql":  1027,
	"pvrz": 1028,
	"glsl": 1029,
	"tot":  1030,
	"toh":  1031,
	"menu": 1032,
	"lua":  1033,
	"ttf":  1034,
	"ini":  2050,
}

var fileTypesExt = map[int]string{}

func init() {
	for ext, num := range fileTypes {
		fileTypesExt[num] = ext
	}
}

func (res *keyResourceEntry) SetBifId(id uint32) {
	res.Location = (res.Location & 0xFFFFF) | id<<20
}

func (res *keyResourceEntry) SetResourceId(id uint32) {
	res.Location = res.Location&0xFFFFC000 | id&0x3fff
}

func (res *keyResourceEntry) GetBifId() uint32 {
	return res.Location >> 20
}
func (res *keyResourceEntry) GetResourceId() uint32 {
	return res.Location & 0x3fff
}
func (res *keyResourceEntry) GetTilesetId() uint32 {
	return (res.Location & 0x000FC000) >> 14
}
func (res *keyResourceEntry) CleanName() string {
	return strings.ToUpper(res.Name.String())
}

func (res *keyResourceEntry) String() string {
	return fmt.Sprintf("%s.%s: %d %d", res.CleanName(), fileTypesExt[int(res.Type)], res.GetBifId(), res.GetResourceId())
}

func OpenKEY(r io.ReadSeeker, root string) (*KEY, error) {
	key := &KEY{r: r, root: root}

	r.Seek(0, os.SEEK_SET)
	err := binary.Read(r, binary.LittleEndian, &key.header)
	if err != nil {
		return nil, err
	}

	r.Seek(int64(key.header.BifOffset), os.SEEK_SET)
	bifs := make([]keyBifEntry, key.header.BifCount)
	err = binary.Read(r, binary.LittleEndian, &bifs)
	if err != nil {
		return nil, err
	}
	for _, bifEntry := range bifs {
		_, err := key.r.Seek(int64(bifEntry.OffsetFilename), os.SEEK_SET)
		if err != nil {
			return nil, err
		}
		bufStr := make([]byte, bifEntry.LengthFilename)
		nBytes, err := io.ReadAtLeast(key.r, bufStr, int(bifEntry.LengthFilename))
		if err != nil {
			return nil, err
		}
		key.bifs = append(key.bifs, keyBifValue{Length: bifEntry.Length, Filename: path.Clean(strings.Replace(strings.Trim(string(bufStr[0:nBytes]), "\000"), "\\", "/", -1)), FileLocation: bifEntry.FileLocation})
	}
	r.Seek(int64(key.header.ResourceOffset), os.SEEK_SET)
	key.resources = make([]keyResourceEntry, key.header.ResourceCount)
	err = binary.Read(r, binary.LittleEndian, &key.resources)
	if err != nil {
		return nil, err
	}
	key.files = make(map[keyUniqueResource]int)
	for idx, res := range key.resources {
		kur := keyUniqueResource{Name: res.CleanName(), Type: res.Type}
		key.files[kur] = idx
	}
	return key, nil
}

func (key *KEY) GetBifPath(bifId uint32) (string, error) {
	if int(bifId) > len(key.bifs) {
		return "", errors.New("Invalid bifId")
	}
	return key.bifs[bifId].Filename, nil
}

func (key *KEY) TypeToExt(ext uint16) string {
	return fileTypesExt[int(ext)]
}
func (key *KEY) ExtToType(ext string) int {
	fileExt := strings.Trim(ext, ".")
	return fileTypes[fileExt]
}

func (key *KEY) GetResourceName(biffId uint32, resourceId uint32) (string, error) {
	nID := uint32((biffId << 20) | (resourceId & 0x3fff))
	for _, res := range key.resources {
		if res.Location == nID {
			name := string(res.CleanName()) + "." + key.TypeToExt(res.Type)
			return name, nil
		}
	}
	return "", errors.New("Resource not found")
}

func (key *KEY) Write(w io.Writer) error {
	var bifFilenames bytes.Buffer

	for _, bv := range key.bifs {
		bifFilenames.WriteString(bv.Filename)
		bifFilenames.WriteByte(0)
	}
	h := keyHeader{Signature: [4]byte{'K', 'E', 'Y', ' '}, Version: [4]byte{'V', '1', ' ', ' '}}
	h.BifCount = uint32(len(key.bifs))
	h.ResourceCount = uint32(len(key.resources))
	h.BifOffset = uint32(binary.Size(h))
	h.ResourceOffset = h.BifOffset + uint32(binary.Size(keyBifEntry{}))*h.BifCount + uint32(bifFilenames.Len())

	err := binary.Write(w, binary.LittleEndian, &h)
	if err != nil {
		return err
	}

	var bifs []keyBifEntry
	offset := h.ResourceOffset - uint32(bifFilenames.Len())
	for _, bv := range key.bifs {
		length := len(bv.Filename) + 1
		bifs = append(bifs, keyBifEntry{Length: bv.Length, OffsetFilename: offset, LengthFilename: uint16(length), FileLocation: 1})

		offset += uint32(length)
	}

	err = binary.Write(w, binary.LittleEndian, &bifs)
	if err != nil {
		return err
	}

	_, err = bifFilenames.WriteTo(w)
	if err != nil {
		return err
	}

	err = binary.Write(w, binary.LittleEndian, key.resources)
	if err != nil {
		return err
	}

	return nil
}

func (key *KEY) Validate() {
	log.Printf("Header: %+v", key.header)
	for idx, bif := range key.bifs {
		bifPath, _ := key.GetBifPath(uint32(idx))

		fmt.Printf("Idx: %d Path: %s Location: %d\n", idx, bifPath, bif.FileLocation)
	}
	for _, resource := range key.resources {
		bifPath, err := key.GetBifPath(resource.GetBifId())
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Res: %s in %s", resource.String(), bifPath)
		diskPath := path.Join(key.root, bifPath)
		_, err = os.Stat(diskPath)
		if err != nil {
			//fmt.Printf("Can't find bif: %s\n", diskPath)
		}

	}
}
