package util

import "github.com/google/uuid"

func ContainsString(stack *[]string, needle *string) bool {
	for _, s := range *stack {
		if s == *needle {
			return true
		}
	}
	return false
}

func ContainsOneOfStrings(stack *[]string, needles *[]string) bool {
	for _, n := range *needles {
		for _, s := range *stack {
			if s == n {
				return true
			}
		}
	}
	return false
}

func ContainsUUID(stack *[]uuid.UUID, needle *uuid.UUID) bool {
	for _, s := range *stack {
		if s == *needle {
			return true
		}
	}
	return false
}

func ContainsOneOfUUIDs(stack *[]uuid.UUID, needles *[]uuid.UUID) bool {
	for _, n := range *needles {
		for _, s := range *stack {
			if s == n {
				return true
			}
		}
	}
	return false
}
