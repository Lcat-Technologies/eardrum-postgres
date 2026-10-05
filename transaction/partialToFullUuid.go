package transaction

import (
    pgerror "errors"
    "os"
    "strconv"
    "github.com/Lcat-Technologies/eardrum-interfaces/errors"
    "github.com/google/uuid"
    "github.com/rs/zerolog/log"
    "gorm.io/gorm"
    
    "github.com/Lcat-Technologies/eardrum-postgres/user" 
)

// FacialMatchThreshold is resolved dynamically from the environment.
// Fallback default is set to 0.75 if FACIAL_MATCH_THRESHOLD is unset or invalid.
var FacialMatchThreshold float32 = getEnvAsFloat32("FACIAL_MATCH_THRESHOLD", 0.75)

func PartialToFullUuid(db *gorm.DB, partialScan string, facialEmbedding string) (uuid.UUID, *errors.PublicError) {
    if len(partialScan) < 10 {
        err := errors.New(errors.EARTxUserAccountNotFound, pgerror.New("Invalid QR code. Please scan again"))
        err.Log()
        return uuid.Nil, err
    }

    // 1. Fetch candidates matching the fuzzy QR code criteria
    type candidate struct {
        user.User
        Score float32 `gorm:"column:score"`
    }
    var results []candidate

    err := db.Transaction(func(tx *gorm.DB) error {
        tx.Exec("SET LOCAL pg_trgm.word_similarity_threshold = 0.6")
        return tx.Raw(`
            SELECT *, word_similarity(?, qr_code::text) as score 
            FROM users 
            WHERE ? <% qr_code::text 
            ORDER BY score DESC`, partialScan, partialScan).Scan(&results).Error
    })

    if err != nil {
		err1 := errors.New(errors.EARInternalError, err)
		err1.Log()
        return uuid.Nil, err1
    }

    // 2. Filter by strict subsequence AND collect face matches
    var matchingCandidates []candidate
    subsequencePassed := []map[string]interface{}{}

    for _, res := range results {
        if isStrictSubsequence(partialScan, res.QrCode.String()) {
            subsequencePassed = append(subsequencePassed, map[string]interface{}{
                "qr": res.QrCode.String(), "score": res.Score,
            })

            // Check if this user's face matches the scanned face using our env threshold
            if res.User.MatchFace(facialEmbedding, FacialMatchThreshold) {
                matchingCandidates = append(matchingCandidates, res)
            }
        }
    }

    // LOG: Filter details
    log.Info().
        Interface("subsequence_passed", subsequencePassed).
        Int("matching_faces_count", len(matchingCandidates)).
        Float32("configured_threshold", FacialMatchThreshold).
        Msg("Filter results executed")

    // 3. Collision and Disqualification Check
    if len(matchingCandidates) == 0 {
        log.Info().Str("input", partialScan).Msg("No accounts found matching both QR sequence and face profile")
        err1 := errors.New(errors.EARTxUserAccountNotFound, pgerror.New("Invalid QR code or biometric match failed. Please try again"))
		err1.Log()
		return uuid.Nil, err1
    }

    if len(matchingCandidates) > 1 {
        collisionLog := []string{}
        for _, mc := range matchingCandidates {
            collisionLog = append(collisionLog, mc.QrCode.String())
        }
        log.Warn().
            Interface("colliding_uuids", collisionLog).
            Msg("Security Disqualification: Facial embedding matched multiple candidates")
            
        return uuid.Nil, errors.New(errors.EARTxUserAccountNotFound, pgerror.New("Ambiguous match detected"))
    }

    // Exactly one clear match found
    matchedUser := matchingCandidates[0]
    log.Info().
        Str("qr", matchedUser.QrCode.String()).
        Float32("score", matchedUser.Score).
        Msg("Single secure face match identified successfully")

    return matchedUser.QrCode, nil
}

// isStrictSubsequence checks if partial is a subset of full in the correct order
func isStrictSubsequence(partial, full string) bool {
    pIdx, fIdx := 0, 0
    for pIdx < len(partial) && fIdx < len(full) {
        if partial[pIdx] == full[fIdx] {
            pIdx++
        }
        fIdx++
    }
    return pIdx == len(partial)
}

// Helper to look up an env variable and map it safely to a float32
func getEnvAsFloat32(key string, defaultVal float32) float32 {
    valueStr, exists := os.LookupEnv(key)
    if !exists {
        return defaultVal
    }
    
    value, err := strconv.ParseFloat(valueStr, 32)
    if err != nil {
        log.Warn().Err(err).Str("key", key).Str("value", valueStr).Msg("Invalid env float format, using default fallback")
        return defaultVal
    }
    
    return float32(value)
}