package mergekeys

import (
	"log"
	"strings"
)

func (k *KEY) MergeWith(mk *KEY) error {
	lastBifId := uint32(len(k.bifs))
	for i, be := range mk.bifs {
		name, err := mk.GetBifPath(uint32(i))
		if err != nil {
			log.Printf("Err: %+v", err)
		}
		k.bifs = append(k.bifs, keyBifValue{Length: be.Length, Filename: strings.Replace(name, "data", "mod", 1)})
	}
	for _, mr := range mk.resources {
		kur := keyUniqueResource{Name: mr.CleanName(), Type: mr.Type}
		newBifId := mr.GetBifId() + lastBifId

		idx, ok := k.files[kur]
		if ok {
			or := k.resources[idx]
			or.SetBifId(newBifId)
			or.SetResourceId(mr.GetResourceId())
			k.resources[idx] = or
		} else {
			nr := mr
			nr.SetBifId(newBifId)
			k.resources = append(k.resources, nr)
			k.files[kur] = len(k.resources) - 1
		}
	}
	return nil
}
