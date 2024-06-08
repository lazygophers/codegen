package utils

import "github.com/spf13/cobra"

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
