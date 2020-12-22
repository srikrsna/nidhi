package nidhigen

import "github.com/srikrsna/nidhi"

type StringField nidhi.OrderByString

type IntField = nidhi.OrderByInt

type FloatField = nidhi.OrderByFloat

type BoolField = UnorderedField

type TimeField = nidhi.OrderByTime

type UnorderedField string
