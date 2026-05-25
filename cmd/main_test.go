package main

import "testing"

func TestMainCommandRegistersSubcommands(t *testing.T) {
	for _, name := range []string{"status", "updatetap", "whatwhere"} {
		cmd, _, err := MAIN.Find([]string{name})
		if err != nil {
			t.Fatalf("find command %q: %v", name, err)
		}
		if cmd == nil || cmd.Name() != name {
			t.Fatalf("command %q not registered; got %#v", name, cmd)
		}
	}
}
