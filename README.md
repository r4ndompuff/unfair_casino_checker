# Unfair Casino Checker
Example of how with basic Markov Chains with fixed sliding window can help you to check the fairness of casino.

## Files
- "GOLANG/" - different fast scrappers written on Go (by my brother). Excluded one lib, so you will need to rewrite part with sockets.
-  "data.xlsx" - main data file with different info about every bet.
-  "onlineCasino.scv" - small data file (you can check EDA here: https://www.kaggle.com/code/andreylovyagin/online-unfair-casino-example)
-  "dataGO.csv" - small data file needed for Markov Chains.
-  "parser.py" - parser for data from go scrappers (if you need).
-  "gamer.py" - only kernel for markov chains.
 
## Additional info
I've deleted some self-written libs from golang scrappers (they are not published yet and they are not mine).
What is more, I've deleted some site-control functions from "gamer.py" (selenium-based), because this repository only an example that you can check simple gaming-casinos for fairness.
