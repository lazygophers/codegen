package cli

import "github.com/spf13/cobra"

func getBool(key string, c *cobra.Command) (ok bool) {
	ok, _ = c.Flags().GetBool(key)
	return ok
}

func getString(key string, c *cobra.Command) string {
	value, _ := c.Flags().GetString(key)
	return value
}
