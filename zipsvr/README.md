# Zip Server

We will learn Go by creating a simple web server that supports the following APIs:

- `/hello`: responds with a simple "Hello, World!" plaintext message.
- `/hello/[name]`: responds with a plaintext greeting to `[name]`, which can be any string.
- `/zips/[city]`: responds with all zip codes and states associated with the given `[city]`, encoded as a JSON array of objects.

This last API will use the data in `zips.csv` to find zips for a given city name.