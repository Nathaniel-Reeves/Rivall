package password_recovery

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
)

type RecoveryOTP struct {
	Email     string    `json:"email"`
	Code      string    `json:"code"`
	ExpiresAt time.Time `json:"expires_at"`
}

type Verifier interface {
	VerifyOTP(otp string) bool
}

type RecoveryRetentionMap map[string]RecoveryOTP

const CODE_TIMEOUT = time.Second * 60

func NewRecoveryRetentionMap(ctx context.Context) *RecoveryRetentionMap {
	rm := make(RecoveryRetentionMap)

	go rm.Retention(ctx)

	return &rm
}

func (rm RecoveryRetentionMap) NewRecoveryOTP(email string) RecoveryOTP {
	o := RecoveryOTP{
		Email:     email,
		Code:      strings.ToUpper(uuid.New().String()[0:6]),
		ExpiresAt: time.Now().Add(CODE_TIMEOUT),
	}

	rm[email] = o
	return o
}

func (rm RecoveryRetentionMap) VerifyRecoveryOTP(code string, email string) bool {

	// Verify OTP is existing
	if _, ok := rm[email]; !ok {
		// otp does not exist
		log.Debug().Msg("OTP does not exist")
		return false
	}

	// Verify OTP is correct
	if strings.ToUpper(rm[email].Code) != strings.ToUpper(code) {
		// otp is incorrect
		log.Debug().Msg("OTP is incorrect")
		return false
	}

	delete(rm, email)
	return true
}

func (rm RecoveryRetentionMap) Retention(ctx context.Context) {
	ticker := time.NewTicker(time.Millisecond * 500)
	for {
		select {
		case <-ticker.C:
			for _, otp := range rm {
				if otp.ExpiresAt.Before(time.Now()) {
					delete(rm, otp.Email)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
