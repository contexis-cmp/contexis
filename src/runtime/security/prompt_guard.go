package security

import (
    "regexp"
    "strings"
)

// PromptRisk represents a simple risk classification for inputs
type PromptRisk string

const (
    RiskLow    PromptRisk = "low"
    RiskMedium PromptRisk = "medium"
    RiskHigh   PromptRisk = "high"
)

// Heuristic patterns for prompt injection and authority spoofing
var (
    reIgnorePrev    = regexp.MustCompile(`(?i)ignore\s+(previous|prior)\s+(instructions|context|rules)`)
    reRevealSystem  = regexp.MustCompile(`(?i)(reveal|show)\s+(system\s+prompt|hidden\s+prompt)`)
    reAsAdmin       = regexp.MustCompile(`(?i)(as\s+an?\s+admin|as\s+your\s+manager|override\s+policy)`)
    reChangeRules   = regexp.MustCompile(`(?i)(disregard|bypass|break)\s+(rules|policy|guardrails)`)
    reDataExfil     = regexp.MustCompile(`(?i)(leak|dump|exfiltrate)\s+(data|secrets|keys|passwords)`) 
    // PII patterns (heuristics)
    reEmail         = regexp.MustCompile(`(?i)[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}`)
    rePhone         = regexp.MustCompile(`(?i)(\+?\d[\d\s\-]{7,}\d)`)
    reSSN           = regexp.MustCompile(`\b\d{3}-\d{2}-\d{4}\b`)
)

// ClassifyPromptRisk applies simple heuristics to detect likely injection attempts
func ClassifyPromptRisk(text string, extras ...string) PromptRisk {
    combined := strings.ToLower(strings.TrimSpace(text + "\n" + strings.Join(extras, "\n")))
    matches := 0
    for _, re := range []*regexp.Regexp{reIgnorePrev, reRevealSystem, reAsAdmin, reChangeRules, reDataExfil} {
        if re.MatchString(combined) {
            matches++
        }
    }
    switch {
    case matches >= 2:
        return RiskHigh
    case matches == 1:
        return RiskMedium
    default:
        return RiskLow
    }
}

// SanitizeUserInput strips or neutralizes model-directed instructions in user-provided text
func SanitizeUserInput(text string) string {
    s := text
    // Neutralize common directive verbs
    replacements := []struct{ from, to string }{
        {"Ignore previous", "[user]"},
        {"ignore previous", "[user]"},
        {"Disregard", "[user]"},
        {"disregard", "[user]"},
        {"bypass", "[user]"},
        {"Break", "[user]"},
        {"break", "[user]"},
    }
    for _, r := range replacements {
        s = strings.ReplaceAll(s, r.from, r.to)
    }
    return s
}

// DetectPII returns true if content likely contains PII
func DetectPII(s string) bool {
    if reEmail.FindStringIndex(s) != nil { return true }
    if rePhone.FindStringIndex(s) != nil { return true }
    if reSSN.FindStringIndex(s) != nil { return true }
    return false
}

// RedactPII replaces likely PII with placeholders
func RedactPII(s string) string {
    s = reEmail.ReplaceAllString(s, "[REDACTED_EMAIL]")
    s = rePhone.ReplaceAllString(s, "[REDACTED_PHONE]")
    s = reSSN.ReplaceAllString(s, "[REDACTED_SSN]")
    return s
}


