# go-cswizard

CSWizard is a CSV writer that doesn't stand in your way as your system evolves. Using `encoding/csv` directly is fine when you want to just do the thing and forget it. Long-living projects, however, are rarely done this way: every day business demands to add new columns, remove and reorder them, and under this conditions `encoding/csv` becomes too fragile due to nature of its API: you just can't change one thing and be sure that everything else would keep working. With CSWizard, you can.

So, in a nutshell, that's a small wrapper around `encoding/csv` for reports that change often in various ways.

See [example.go](example/example.go) for typical usage.
