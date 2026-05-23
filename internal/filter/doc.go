// Package filter implements the core log-line filtering engine for logslice.
//
// A Filter is constructed from an Options struct that specifies:
//
//   - Include patterns — only lines matching at least one include pattern are kept.
//     When no include patterns are provided every line passes this stage.
//   - Exclude patterns — lines matching any exclude pattern are dropped before
//     the include check is applied.
//   - CaseSensitive — by default matching is case-insensitive; set this flag to
//     enable exact-case comparisons.
//   - UseRegex — by default patterns are treated as plain sub-strings; set this
//     flag to interpret patterns as Go regular expressions.
//
// Typical usage:
//
//	f, err := filter.New(filter.Options{
//		Include:  []string{"ERROR", "WARN"},
//		Exclude:  []string{"healthcheck"},
//		UseRegex: false,
//	})
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	scanner := bufio.NewScanner(os.Stdin)
//	for scanner.Scan() {
//		if f.Match(scanner.Text()) {
//			fmt.Println(scanner.Text())
//		}
//	}
package filter
