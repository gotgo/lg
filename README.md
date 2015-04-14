# lg
Structured Logging Api


## Thoughts on Logging for Servers

The ideal case is that a well structured (i.e. JSON) messages gets recorded to a central log repository on a different
machine that provides search (i.e. Logstash and ElasticSearch).  As long as all log messages are guaranteed to be
well formatted this works great. However, an erroneous unformatted message can cause problems for readers expecting JSON. 
Therefore, a strategy to handle this situation is to only have unexpected messages (i.e. Panics) go to stderr and capture these 
unstructured errors in a separate file in the operational environment such as using runit's default logging for this.  
The rest of the structured logging can go to a structured only destination, such as a file or external system.

There are 2 primary contexts to consume logs.

1. A human logged into a machine an running commands locally on that machine. 
	* tail a human readable format
	* have a separate file for alerts such that they are not lost when there are many log messages
2. Centralized logging and search server.
	* Requires a consistent format. (i.e. all entries are JSON)
		- Structured Stream
		- Unstructured Stream
	* Can handle a large volume of messages, due to search.
