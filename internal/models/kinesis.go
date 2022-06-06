package models

import "reflect"

type Kinesis struct {
	Model
	ModuleReleaseID uint   `json:"-"`
	Name            string `json:"name"`
	Region          string `json:"region"`
	Shards          int32  `json:"shards"`
	Tags            string `json:"tags"`
	Revision        int    `json:"revision"`
}

func (k Kinesis) IsEmpty() bool {
	return reflect.DeepEqual(k, Kinesis{})
}
