# Geckoboard Test

## Solution

The repo contains a HTTP API that will 'learn' from an uploaded body of text, and generate new sentences on demand.
It contains two endpoints, `POST /learn`, which accepts a plain text body (`Content-Type: text/plain`),
and `GET /generate`, which returns a newly generated sentence.

It generates sentences based on trigrams from the learned text, for example the sentence "To be or not to be, that is the question" contains the following trigrams
```
[to, be, or]
[be, or, not]
[or, not, to]
[not, to, be]
[to, be, that]
[be, that, is]
[that, is, the]
[is, the, question]
```

If we start with the words "to be", then we have a choice of the next word being "or", or "that", with an equal probability. We then take the next two words, either "be or", or "be that" and repeat the process. The next word is chosen with the same frequency that occurs in the learned text. For example, if the word "stormy" follows the words "dark and" 9 times out of 10, then the generated text will maintain that frequency.

I've used weighted graph datastructure to hold the corpus, this allows for producing newly generated sentences by randomly walking through adjacent nodes.
The weights of each edge represent the frequency of occurrence in the original body of text, with the sum of these weights being stored in the relevant graph node. 
This allows us to chose a random number 1 <= n <= totalWeight, iterate over the array of edges, summing the weights of the passed edges. We stop at the point this sum exceeds the randomly chosen value. An alternative method of preserving the frequencies would be to create duplicate nodes for each trigram that appears in the text. This would allow us to randomly select an edge, without using weighted values, reducing the time complexity, but potentially massively increasing the space complexity if there were lots of reoccurring trigrams in the text.

The method I chose trades a slightly worse time complexity, for an improved space complexity. In the worst case, for every node in the path, we're to visiting every adjacent node to chose one with the correct frequency. I believe this is an acceptable trade-off to keep the memory usage of the graph to a minimum.

The graph is implemented using an adjacency list, I've implemented this myself rather than use an external library for the datastructure. I would likely opt for a library in a real setting, unless there was a compelling reason to write my own datastructure. I felt this use case was simple and quick enough to write it myself.

## Assumptions

- I'm assuming the uploaded body of text needs cleaning and normalising first. To do this I've evened out the whitespace by ensuring each word is separated by a single space. I then removed any characters outside of a-z A-Z, spaces, commas, and fullstops.
- When generating the new sentences, I've made the assumption that sentences must start with capital letters, and end with full stops. Commas are also treated as part of the word, e.g. trigrams `this, then that`, and `this then that` are treated differently.

## Improvements

Given more time, further improvements I would implement include:

- Expanding the test suites. As it stands, only the graph structure includes some basic unit tests; I would have liked to expand these using table based testing to cover more cases.
- The HTTP handlers need tests, implementing parallel testing would be nice in order to test concurrent access to the datastructure.
- Implementing a queue for the `/learn` endpoint. Learning the new text is potentially a time consuming task, scaling with the size of the text uploaded. As long as the caller isn't relying on a confirmation of the successful processing of the text, we could push the uploaded text onto a queue return a `202` to the caller immediately.
