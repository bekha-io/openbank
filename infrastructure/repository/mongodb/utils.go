package mongodb

import (
	"github.com/bekha-io/vaultonomy/domain/repository"
	"go.mongodb.org/mongo-driver/bson"
)

func ParseFilter(flt repository.Filter) bson.E {
	if flt.EqualTo != nil {
		return bson.E{Key: flt.Key, Value: flt.EqualTo}
	} else if flt.GreaterThan != nil {
		return bson.E{Key: "$gt", Value: flt.GreaterThan}
	} else if flt.LessThan != nil {
		return bson.E{Key: "$lt", Value: flt.LessThan}
	} else if flt.Like != nil {
		return bson.E{Key: "$like", Value: flt.Like}
	}

	return bson.E{}
}
