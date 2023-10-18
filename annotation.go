package exclusion

import (
	"encoding/json"

	"entgo.io/ent/schema"
)

type Annotation struct {
	SkipCreate bool
	SkipDelete bool
	SkipUpdate bool
	SkipList   bool
	SkipRead   bool
}

// Decode from ent.
func (a *Annotation) Decode(o interface{}) error {
	buf, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, a)
}

func (wa Annotation) Name() string {
	return "OasOperation"
}

// Merge implements ent.Merger interface.
func (a Annotation) Merge(o schema.Annotation) schema.Annotation {
	var ant Annotation
	switch o := o.(type) {
	case Annotation:
		ant = o
	case *Annotation:
		if o != nil {
			ant = *o
		}
	default:
		return a
	}

	if ant.SkipCreate {
		a.SkipCreate = ant.SkipCreate
	}

	if ant.SkipUpdate {
		a.SkipUpdate = ant.SkipUpdate
	}

	if ant.SkipDelete {
		a.SkipDelete = ant.SkipDelete
	}

	if ant.SkipList {
		a.SkipList = ant.SkipList
	}

	if ant.SkipRead {
		a.SkipRead = ant.SkipRead
	}

	return a
}

func SkipCreate() Annotation {
	return Annotation{SkipCreate: true}
}

func SkipDelete() Annotation {
	return Annotation{SkipDelete: true}
}

func SkipUpdate() Annotation {
	return Annotation{SkipUpdate: true}
}

func SkipList() Annotation {
	return Annotation{SkipList: true}
}

func SkipRead() Annotation {
	return Annotation{SkipRead: true}
}
