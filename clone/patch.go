package clone

import (
	"crypto/dsa"
	"math/big"
	"reflect"

	"github.com/bouk/monkey"
)

// PatchDsa2048 will bluntly patch some methods in the crypto/dsa and math/big
// packages to disable the unsupported DSA with key size 2048 checks in the
// golang.org/x/crypto/ssh package.
func PatchDsa2048() {
	// This is the DSA Verify method, and will always return true on verification
	// of DSA keys
	monkey.Patch(dsa.Verify, func(pub *dsa.PublicKey, hash []byte, r, s *big.Int) bool {
		Info("Ignore signature check of dsa public key (Verify)")
		return true
	})

	// In golang.org/x/crypto/ssh/keys there is a private method "checkDSAParams"
	// that checks if the size of the key is 2048 bits. This check needs to be
	// disabled, but since the method is private, the BitLen method is patched
	// instead.
	var bi *big.Int
	var guard *monkey.PatchGuard
	guard = monkey.PatchInstanceMethod(reflect.TypeOf(bi), "BitLen", func(x *big.Int) int {
		guard.Unpatch()
		defer guard.Restore()
		l := x.BitLen()
		if l == 2048 {
			Info("Ignore dsa param check (BitLen)")
			return 1024
		}
		return l
	})
}
