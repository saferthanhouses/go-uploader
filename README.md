#### File Uploader Module

Golang service for handling file uploads

- Authenticates upload requests using JSON Web Tokens
- Upload files to any aws compatible object storage (e.g. Digital Ocean)

####Directory Structure

├── README.md
├── config
│   └── config.go   # load configuration from environmental variables
├── main.go         # start server
├── middleware        
│   ├── auth.go     # authenticate requests
│   ├── cors.go     # make requests work with Cross Origin Resource Sharing
│   └── logging.go  # TODO
├── routes
│   └── upload.go   # Basic upload route
├── upload-client
│   └── client.go   # Wrapper around AWS-style object storage client 
└── utils
    ├── file.go     # bork
    └── url.go      # bonk