# Geckoboard Test

I've used an internal Weighted graph datastructure to hold the learned corpus.
The graph is implemented using an adjacency list, I've implemented this myself rather than use an external library for the datastructure. I would likely opt for a library in a real setting, unless there was a compelling reason to write my own datastructure. I felt this use case was simple and quick enough to write it myself.

The adjacency list uses weighted nodes to represent the frequency a particular trigram occours following the previous trigram.
This requires the generate function to first visit each adjacent node to the one being processesed, to obtain the frequencies, and then randomly select one with respect to the frequency. This adds to the time complexity, with the tradeoff being reduced space complexity.

This could have been implemented by creating additional nodes for each duplicated trigram, which would likely baloon the space required to store the structure, but would optimise for fast walking through the graph, as the a randomly selected node would inheritly be chosen with a frequency matching that of the occurrence.

## TODO

- Make concurrency safe, probably use mutexes to protect the low level structures
- Tests
- Push learning into a queue, and return from handler immediately with 202
- Make sure initial learning text is cleaned up correctly, no double spaces etc
- Godoc all exported functions
- Add comments to anywhere that needs it
- How can we test the frequency is accurate?
- Test it on the whole of project gutenberg, talk about tradeoffs of design
- Make graph more generic by having the 'value' be an interface{}?
- Finish readme
- Table based testing?
- As well as sending us your code, we would like to hear about any design considerations or technical details you gave particular consideration to. If you had to make any particular trade-offs, or if you would have done things differently given more time, then let us know.
- Golint, go vet, etc...
