package keys

import (
	"bufio"
	"errors"
	"os"

	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/flags"
	"github.com/cosmos/cosmos-sdk/client/input"
	"github.com/cosmos/cosmos-sdk/crypto/keys"

	"github.com/spf13/cobra"
)

const (
	flagYes   = "yes"
	flagForce = "force"
)

func deleteKeyCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete <name>",
		Short: "Delete the given key",
		Long: `Delete a key from the store.

Note that removing offline or ledger keys will remove
only the public key references stored locally, i.e.
private keys stored in a ledger device cannot be deleted with the CLI.
`,
		RunE: runDeleteCmd,
		Args: cobra.ExactArgs(1),
	}

	cmd.Flags().BoolP(flagYes, "y", false,
		"Skip confirmation prompt when deleting offline or ledger key references")
	cmd.Flags().BoolP(flagForce, "f", false,
		"Remove the key unconditionally without asking for the passphrase")
	cmd.Flags().Bool(flags.FlagLegacy, false, "Use legacy secret store")
	return cmd
}

func runDeleteCmd(cmd *cobra.Command, args []string) error {
	name := args[0]
	var kb keys.Keybase
	buf := bufio.NewReader(cmd.InOrStdin())

	if viper.GetBool(flags.FlagLegacy) {
		os.Stderr.WriteString("IMPORTANT: using deprecated secret store. This will be removed in a future release.")
		var err error
		kb, err = NewKeyBaseFromHomeFlag()
		if err != nil {
			return err
		}
	} else {
		kb = NewKeyringKeybase(buf)
	}

	info, err := kb.Get(name)
	if err != nil {
		return err
	}

	if info.GetType() == keys.TypeLedger || info.GetType() == keys.TypeOffline {
		if !viper.GetBool(flagYes) {
			if err := confirmDeletion(buf); err != nil {
				return err
			}
		}
		if err := kb.Delete(name, "", true); err != nil {
			return err
		}
		cmd.PrintErrln("Public key reference deleted")
		return nil
	}

	// skip passphrase check if run with --force
	skipPass := viper.GetBool(flagForce)
	var oldpass string
	//you don't need a passphrase in the new keystore, the keyring is already unlocked
	if !skipPass && viper.GetBool(flags.FlagLegacy) {
		if oldpass, err = input.GetPassword(
			"DANGER - enter password to permanently delete key:", buf); err != nil {
			return err
		}
	}

	err = kb.Delete(name, oldpass, skipPass)
	if err != nil {
		return err
	}
	cmd.PrintErrln("Key deleted forever (uh oh!)")
	return nil
}

func confirmDeletion(buf *bufio.Reader) error {
	answer, err := input.GetConfirmation("Key reference will be deleted. Continue?", buf)
	if err != nil {
		return err
	}
	if !answer {
		return errors.New("aborted")
	}
	return nil
}
