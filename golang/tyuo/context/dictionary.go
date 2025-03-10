package context
import (
    "bufio"
    "os"
    "strings"
)

type ParsedToken struct {
    Base string
    Variant string
}


func processBoringTokens(listPath string) (map[string]void, error) {
    file, err := os.Open(listPath)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    output := make(map[string]void)
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        token := strings.ToLower(strings.TrimSpace(scanner.Text()))
        if len(token) > 0 {
            output[token] = voidInstance
        }
    }
    if err := scanner.Err(); err != nil {
        return nil, err
    }
    logger.Debugf("loaded %d language-level boring tokens", len(output))
    return output, nil
}


type DictionaryToken struct {
    id int
    baseOccurrences int
    baseRepresentation string
    variantForms map[string]int
}
func (dt *DictionaryToken) GetId() (int) {
    return dt.id
}
func (dt *DictionaryToken) Represent(baseRepresentationThreshold float32) (string, bool) {
    sum := float32(dt.baseOccurrences)
    for _, count := range dt.variantForms {
        sum += float32(count)
    }
    
    if float32(dt.baseOccurrences) / sum > baseRepresentationThreshold {
        return dt.baseRepresentation, true
    } else {
        var mostRepresented string
        var mostRepresentedCount int = 0
        for representation, count := range dt.variantForms {
            if count > mostRepresentedCount {
                mostRepresentedCount = count
                mostRepresented = representation
            }
        }
        return mostRepresented, false
    }
}
func (dt *DictionaryToken) rescale(rescaleThreshold int,  rescaleDecimator int) {
    rescaleNeeded := false
    for _, count := range dt.variantForms{
        if count > rescaleThreshold {
            rescaleNeeded = true
            break
        }
    }
    if rescaleNeeded {
        for variant, count := range dt.variantForms {
            count /= rescaleDecimator
            if count > 0 {
                dt.variantForms[variant] = count
            } else {
                delete(dt.variantForms, variant)
            }
        }
    }
}

type dictionary struct {
    database *database
    
    nextIdentifier int
}
func prepareDictionary(database *database) (*dictionary, error) {
    nextIdentifier, err := database.dictionaryGetNextIdentifier()
    if err != nil {
        return nil, err
    }
    
    return &dictionary{
        database: database,
        
        nextIdentifier: nextIdentifier,
    }, nil
}
func (d *dictionary) getSliceByToken(tokens stringset) (map[string]DictionaryToken, error) {
    dictionaryTokens, err := d.database.dictionaryGetTokensByToken(tokens)
    if err != nil {
        return nil, err
    }
    dictionarySlice := make(map[string]DictionaryToken, len(dictionaryTokens))
    for _, dt := range dictionaryTokens {
        dictionarySlice[dt.baseRepresentation] = dt
    }
    return dictionarySlice, nil
}
func (d *dictionary) getSliceById(ids intset) (map[int]DictionaryToken, error) {
    dictionaryTokens, err := d.database.dictionaryGetTokensById(ids)
    if err != nil {
        return nil, err
    }
    dictionarySlice := make(map[int]DictionaryToken, len(dictionaryTokens))
    for _, dt := range dictionaryTokens {
        dictionarySlice[dt.id] = dt
    }
    return dictionarySlice, nil
}

func (d *dictionary) getIdsByToken(tokens stringset) ([]int, error) {
    return d.database.dictionaryEnumerateIdsByToken(tokens)
}

func (d *dictionary) learnTokens(tokens []ParsedToken, rescaleThreshold int,  rescaleDecimator int) ([]DictionaryToken, error) {
    //get any existing entries from the database
    tokenSet := make(stringset, len(tokens))
    for _, token := range tokens {
        tokenSet[token.Base] = false
    }
    dictionarySlice, err := d.getSliceByToken(tokenSet)
    if err != nil {
        return make([]DictionaryToken, 0), err
    }
    
    //update the slice with changes
    for _, token := range tokens {
        dt, defined := dictionarySlice[token.Base]
        if !defined {
            d.nextIdentifier++
            dt = DictionaryToken{
                id: d.nextIdentifier,
                baseOccurrences: 0,
                baseRepresentation: token.Base,
                variantForms: make(map[string]int),
            }
        }
        
        if token.Base == token.Variant {
            dt.baseOccurrences += 1
        } else {
            count, _ := dt.variantForms[token.Variant] //default is 0, so it doesn't matter if it's undefined
            dt.variantForms[token.Variant] = count + 1
        }
        
        dictionarySlice[token.Base] = dt
    }
    
    //update the database
    newTokens := make([]DictionaryToken, 0, len(tokenSet))
    for _, dt := range dictionarySlice {
        newTokens = append(newTokens, dt)
    }
    return newTokens, d.database.dictionarySetTokens(newTokens, rescaleThreshold, rescaleDecimator)
}
