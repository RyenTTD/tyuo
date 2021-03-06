package context
import (
    "github.com/juju/loggo"
    "math/rand"
    "time"
)

var logger = loggo.GetLogger("context")

type void struct{}
var voidInstance = void{}

//used to denote the end of a sentence, so it will never be a valid ID
const undefinedDictionaryId = -2147483648 //int32 minimum; should constrain database sizing

//TODO: these should be flags
const rescaleThreshold = 1000.0
const rescaleDecimator = 3.0
//NOTE: when decimating, use math.RoundToEven to get to 0 faster to slimite obsolete paths

var rng = rand.New(rand.NewSource(time.Now().Unix()))
