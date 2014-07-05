rekt-post
=========

A post recommendation engine written in Go

formula
=======

Uses [term frequency-inverse document frequency](http://en.wikipedia.org/wiki/Tf-idf), or td-idf.

Term frequency-inverse document frequency is broken into two parts, the term frequency, which finds how often a term appears in a specific document, and the inverse document frequency, which finds the importance of a term globally, across all documents. The term's global importance is multiplied by the term's frequency is a specific document to find how significant that term is to the document.

term frequency
==============

There are a couple ways to determine the term frequency. We use the augmented frequency, which normalize's a terms frequency with regards to the document's length, in order to keep longer documents from being favored by our formula.

We use the following formula to calculate the augmented frequency, which I have swiped from wikipedia. *t* is the term, *d* is the document, *f* is a function that takes a term and a document and returns the raw frequency (how many times term *t* appeared in *d*):

tf(t, d) = 0.5 + (0.5 * f(t, d)) / max{ f(w, d) : w in d }

inverse document frequency
==========================

The inverse document frequency gives you the importance of each word by telling you how rare or common it is across all of your documents. I have no idea why the formula is like this, but wikipedia says it works well as a heuristic, thought it has shakey theoretical foundations.

"It is the logarithmically scaled fraction of the documents that contain the word, obtained by dividing the total number of documents by the number of documents containing the term, and then taking the logarithm of that quotient." - wikipedia

idf(t, D) = log( |D| / (1 + |{ d in D : t in d}|))

tf-idf
=====

To obtain our tf-idf for each word, we use the following formula:

tfidf(t, d, D) = tf(t, d) * idf(t, D)
