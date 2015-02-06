- Most frequent messages are types 1, 3, 4, 5, 18. We should decode all of them. Current we decode 1, 3, 4.
- DONE: Decode MMSI (includes country or other info as well)
- DONE: Implement a generic router where we feed it messages and it returns message type and payload or error.
- DONE: implement a parser for payloads spanning across 2 or more AIS messages. (tricky: variable length, out of order maybe, maybe we should expire if we don't receive all parts after some time)


Notes:
- Type 1/2/3/4/9/11/18/22 -        168bits -       1 sentence
- Type 5                  -        424bits -       2 sentences
- Type 6/8/12/14          - up to 1008bits - up to 5 sentences
- Type 7/13               -     72-168bits -       1 sentence
- Type 10                 -         72bits -       1 sentence
- Type 15                 -     88-160bits -       1 sentence
- Type 16                 -  96 or 144bits -       1 sentence
- Type 17                 -  80 to 816bits - up to 4 sentences
- Type 19                 -        312bits -       2 sentences
- Type 20                 -     72-160bits -       1 sentence
- Type 21                 -    272-360bits -       2 sentences
- Type 23                 -        160bits -       1 sentence
- Type 24                 -    328-336bits -       2 sentences
- Type 25                 -  up to 168bits -       1 sentence
- Type 26                 -    60-1064bits - up to 5 sentences
- Type 27                 -  96 or 168bits -       1 sentence
