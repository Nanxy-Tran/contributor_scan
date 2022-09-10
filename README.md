# Git owner file generator 🔎
## Description
Auto generate output contains file paths and owners of that files (copy paste to your repo without reformat). 
Maximum 2 owners per file, specific number is under working :wink

### Installation
#### Make sure you have installed GO 1.18 on your machine
- Install locally: <br>
  1. Clone this repo <br>
  2. At cloned repo: `go install` to install at your GO PATH <br>
- Install globally:
  1. **Set your GOPATH** in `zshrc`, `bash_rc`, etc,..
  2. `go install github.com/nanxy-tran/contributor_scan` <br>
### Usage
At your desired directory, run: <br>
`contributor_scan git`

To retrieve number of files, run: <br>
`contributor_scan`

### Examples
Git generated output for Github owners: <br>

![img.png](img.png)

Total file count: <br>

![img_1.png](img_1.png)
