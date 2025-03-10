package context
import (
    "github.com/juju/loggo"
    "math/rand"
    "time"
)

var logger = loggo.GetLogger("context")

type void struct{}
var voidInstance = void{}

//set-types where the value doesn't actually matter
type intset map[int]bool
type stringset map[string]bool

//used to denote the end of a sentence
const BoundaryId = -2147483648 //int32 minimum; should constrain database byte-sizing
const undefinedDictionaryId = BoundaryId + 2048 //int32 minimum, plus space for reserved tokens
//reservedIdsPunctuation = 32, from -2147483647 to -2147483615

var rng = rand.New(rand.NewSource(time.Now().Unix()))
