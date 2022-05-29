# Contribution scanner 🔎
## Counting member contribution and who contributed most will buy team beers 🍻😘 . 

### Description
*This side project is an example of my learning method, reading should come along with practice, because it's fun and help us strongly grasp the lessons*.<br> 
Scan through repository and using git blame to identify author contributions. <br>
With the concurrency pattern we only need ~50s to complete a big project (I believe so), but without concurrency, we may need about ~2m30s.<br>
After writing this READ, I got some new ideas about flag args for CLI to count how many lines of code in the repo. Let's do it

### Usage
**CLI**
    Install globally: <br>
        1. Make sure you have installed GO 1.18 on your machine <br>
        2. Clone this repo <br>
        3. At root directory: run `go install` <br>
        4. At your desired directory: run `scanner` and see the result.
