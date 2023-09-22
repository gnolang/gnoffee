// Package gnoffee provides a transpiler that extends the Go language
// with additional, custom keywords. These keywords offer enhanced
// functionality, aiming to make Go programming even more efficient
// and expressive.
//
// Current supported keywords and transformations:
//   - `export <structName> as <interfaceName>`:
//     This allows for the automatic generation of top-level functions
//     in the package that call methods on a specific instance of the struct.
//     It's a way to "expose" or "proxy" methods of a struct via free functions.
//
// How Gnoffee Works:
// Gnoffee operates in multiple stages. The first stage transforms
// gnoffee-specific keywords into their comment directive equivalents,
// paving the way for the second stage to handle the transpiling logic.
//
// Future Changes:
// As the Go and Gno ecosystems and requirements evolve, gnoffee might see the
// introduction of new keywords or alterations to its current functionality.
// Always refer to the package documentation for the most up-to-date details.
package gnoffee
