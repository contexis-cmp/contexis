package security

import "github.com/prometheus/client_golang/prometheus"

var (
    PromptInjectionDetections = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "cmp_security_prompt_injection_detections_total",
        Help: "Total number of detected potential prompt injection attempts.",
    })
    PolicyViolations = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "cmp_security_policy_violations_total",
        Help: "Total number of responses blocked due to policy violations.",
    })
    BlockedResponses = prometheus.NewCounter(prometheus.CounterOpts{
        Name: "cmp_security_blocked_responses_total",
        Help: "Total number of blocked responses.",
    })
)


