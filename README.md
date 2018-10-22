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
│   ├── auth.go
│   ├── cors.go
│   └── logging.go
├── prisma
│   ├── datamodel.prisma
│   └── prisma.yml
├── prisma-client
│   ├── lib.go
│   └── prisma.go
├── routes
│   └── upload.go
├── upload-client
│   └── client.go
└── utils
    ├── file.go
    └── url.go