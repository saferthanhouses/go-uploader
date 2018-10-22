#### File Uploader Module

Golang service for handling file uploads

- Authenticates upload requests using JSON Web Tokens
- Upload files to any aws compatible object storage (e.g. Digital Ocean)

####Directory Structure

├── README.md
<br>├── config
<br>│   └── config.go   # load configuration from environmental variables
<br>├── main.go         # start server
<br>├── middleware        
<br>│   ├── auth.go     # authenticate requests
<br>│   ├── cors.go     # make requests work with Cross Origin Resource Sharing
<br>│   └── logging.go  # TODO
<br>├── routes
<br>│   └── upload.go   # Basic upload route
<br>├── upload-client
<br>│   └── client.go   # Wrapper around AWS-style object storage client 
<br>└── utils
<br>    ├── file.go     # bork
<br>    └── url.go      # bonk