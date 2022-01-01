# GoDupeDetector
GoDupeDetector (goduped) is a code lcone detection tool for go written in go.  It utilizes the standard library for most accurate parsing of go source code.  Detection is performed by applying the longest-common-subsequence algorithm on normalized and pretty-printed source code, and reporting those code fragments meeting a minimum threshold of identical source lines.

GoDupedDetector is in early development, and currently only detects go functions/methods with similar bodies.

# Version History
v0.0.1 - Base implementation 

# Features in Planning
- [ ] Identifier and literal value normalizations for improved near-miss detection.
- [ ] Function signature and type/interface-based detection for functional clone detection.
- [ ] Configurable parameters for detection (size, thresholds, normalizations).
- [ ] Clone metadata (similarity, type, size, etc) and prioritization.
- [ ] Inclusion and exclusion filters on go source files within input directory.
- [ ] Automatic detection of vendored and generated files to be ignored or specially labeled in detection results.
- [ ] Clustering of clone pairs into clone classes.
- [ ] Alternate input specifications, such as file lists.
- [ ] Alternate output formats.
- ...
