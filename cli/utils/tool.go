package utils

import "github.com/spf13/cobra"

func Changed(key string, c *cobra.Command) bool {
	if flag := c.Flag(key); flag != nil {
		return flag.Changed
	}
	return false
}

func GetBool(key string, c *cobra.Command) (ok bool) {
	ok, _ = c.Flags().GetBool(key)
	return ok
}

func GetString(key string, c *cobra.Command) string {
	value, _ := c.Flags().GetString(key)
	return value
}

func GetStringSlice(key string, c *cobra.Command) []string {
	value, _ := c.Flags().GetStringSlice(key)
	return value
}
