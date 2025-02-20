# MiniRL

- Where should MiniRL run?
	- Host machine (alongside the actual backend server)
	- Separate machine

- Which DataStores does it support?
	- Redis
	- SQLite

- How can the user configure MiniRL as per their usecase?
	- YAML config file

- What can the user configure
	- Backend Endpoint
	- DataStore
	- Retention period of the data
	- Rate Limiting Thresholds
	- ...

- How can the user get started?
	- Create a config file.
	- Spin up the Docker Container

- Other Key Features
	- Request Count: Track number of requests within a time window
	- Burst Handling: Determines how many requests are allowed in short bursts. This is important for handling sudden spikes in traffic.
	- Client Identification: Tracks requests per client using an IP address, user ID, or API key.
	- Customizable Limits: Allows setting different rate limits for different APIs or endpoints.

- Key Rate-Limiting Factors
	- IP Address
	- User Identity: API Keys/Token, User ID/Account ID, Session ID
	- Request Metadata: User-Agent, Referer Header, Geolocation, Device Fingerprinting
	- Request Context: Endpoint/Route, HTTP Method, Payload Size, Action Type
	- Behavioral Patterns: Failed Requests, Request Velocity, Access Patterns
	- Temporal Patterns: Time of Day, Sliding Window