# ci-bench-noise

How often does a benchmark gate on a GitHub-hosted runner mistake noise for a regression, and how large a slowdown does it catch?

The workload (`main.go`) sorts a fixed set of numbers. `slowdown.txt` is embedded at build time and scales the work: `10` makes every run do 10% more of the same work. The workflow `.github/workflows/noise.yml` runs many jobs, each on a fresh runner, and in each job compares the commit with the working tree using [himorime](https://github.com/nao1215/himorime):

- interleaved: `himorime compare --against HEAD`, the gate the nao1215 CLIs run on pull requests;
- sequential: all runs of the base, then all runs of the head, judged offline.

A slowdown of `0` is an A/A test: the two sides are the same code, so every regression it reports is a false alarm. The CPU model and runner image of each job are recorded with its results.

```sh
gh workflow run noise.yml -f slowdowns=0,5,10,20 -f repeats=10 -f label=batch1
```

Results are uploaded as workflow artifacts. This repository is an experiment for a note on debimate.jp and is not maintained as a tool.
