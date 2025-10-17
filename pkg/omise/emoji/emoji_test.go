package emoji

import "testing"

func TestGetDeterministic(t *testing.T) {
	emojis := []string{"🍣", "🍙", "🥢", "🍥"}

	// Test that same key returns same emoji
	key := "test-node"
	result1 := GetDeterministic(key, emojis)
	result2 := GetDeterministic(key, emojis)

	if result1 != result2 {
		t.Errorf("Expected same emoji for same key, got %s and %s", result1, result2)
	}

	// Test that result is from the list
	found := false
	for _, e := range emojis {
		if result1 == e {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected emoji from list, got %s", result1)
	}

	// Test empty list
	empty := GetDeterministic("key", []string{})
	if empty != "" {
		t.Errorf("Expected empty string for empty list, got %s", empty)
	}
}

func TestGetSushi(t *testing.T) {
	// Test that it returns a sushi emoji
	result := GetSushi("test")
	found := false
	for _, e := range Sushi {
		if result == e {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected sushi emoji, got %s", result)
	}

	// Test deterministic behavior
	if GetSushi("same-key") != GetSushi("same-key") {
		t.Error("GetSushi should be deterministic")
	}
}

func TestDifferentKeysProduceDifferentEmojis(t *testing.T) {
	// With enough different keys, we should see variety
	seen := make(map[string]bool)
	keys := []string{"node1", "node2", "node3", "node4", "node5", "node6", "node7", "node8"}

	for _, key := range keys {
		emoji := GetSushi(key)
		seen[emoji] = true
	}

	// With 8 keys and 4 emojis, we should see at least 2 different emojis
	if len(seen) < 2 {
		t.Error("Expected variety in emoji selection")
	}
}
