package middleware

//TODO: implement a gzip compression middleware
//handler that compresses the response stream
//using a gzip.Writer, if the client says it
//can handle that encoding. Check for a request
//header named Accept-Encoding, and if it's value
//contains the string "gzip", you can compress the
//response. Otherwise, don't compress the response
