# AIS repository #

This repo contains 3 sub-repos.

### 1. aislib

A library for reading, decoding and processing AIS messages. WIP but it is in a good and usable state.

### 2. ais-receptor

This will be the program that will read AIS data from various sources. Its job is simple; to read from
all sources, log all AIS sentences, do some preliminary processing like weed out duplicates, sentences
with wrong checksum or incomplete sequences and then forward the good ones to another service for
processing.

This program should be always up so that we will not miss any data.


### 3. ais-processor

This program will get verified AIS payloads from the ais-receptor, decode them and process them as we
see fit, like submitting them to elasticsearch or to mongodb.

Since ais-receptor will have to be 24/7 available, ais-processor will be our playground, where we can
try things without losing historical data.
