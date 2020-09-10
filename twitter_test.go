package bot

import (
	"testing"
)

func TestGreek(t *testing.T) {

	t.Run("returns false for no Greek in text", func(t *testing.T) {
		englishTxt := "It was 65 years ago, at night, when the pogrom against the Greek community of Istanbul started. We remember those whose very lives were lost, and their way of life was destroyed by a preplanned act of terror built on a lie."
		got := greek(englishTxt)
		want := false

		assertEqualBooleans(t, got, want)
	})

	t.Run("returns true for Greek in text", func(t *testing.T) {
		greekTxt := "Πριν από 65 χρόνια, τη νύχτα της 6ης Σεπτεμβρίου, ξεκινούσε το πογκρόμ της Ελληνικής κοινότητας της Κωνσταντινούπολης. Σήμερα θυμόμαστε αυτούς που χάθηκαν, εκείνες και εκείνους που η ζωή τους καταστράφηκε από μια προμελετημένη τρομοκρατία, που στήθηκε πάνω σε ένα ψέμα."
		got := greek(greekTxt)
		want := true

		assertEqualBooleans(t, got, want)
	})
}
