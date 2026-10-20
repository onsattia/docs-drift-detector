## 1. Design review

- [ ] 1.1 Review the port and adapter design with the project owner
- [ ] 1.2 Confirm the endpoint spike in #8 is on the path before any adapter is built

## 2. Implementation (starts with #8)

- [ ] 2.1 Write the judge specs for the four verdicts and the no-evidence rule
- [ ] 2.2 Define `Judge`, `Claim`, `Evidence`, and `Verdict` in `internal/judge`
- [ ] 2.3 Write the verdict rules as tests against a fake `Judge`
- [ ] 2.4 Implement the verdict rules to pass 2.3
- [ ] 2.5 Add the approved-endpoint allow-list check with a test
- [ ] 2.6 Add one provider adapter after the endpoint spike is approved
