# nzbindexer
 Proof of concept for a nzb indexer to index news groups from scratch (back from first message to last message)

 Still needs much more sophisticated parsing of the subject to better account for all the very different subject formats used for file posts.

 This is only the "indexer" part. There is currently no frontend to search for headers and to generate a NZB file from the information stored in the DB.

 Needs golang v1.18 to compile (I am lazy and used the new slices package available with v1.18 instead of programming a "contains" routine myself...)
