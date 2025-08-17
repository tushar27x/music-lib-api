# Music Library API

A RESTful API for managing music library with albums, songs, and playlists built with Go, Gin, and GORM.

## 🚀 Features

- **User Authentication**: JWT-based authentication system
- **Album Management**: Full CRUD operations for music albums (artists only)
- **Song Management**: Complete song lifecycle management
- **Playlist Management**: Create and manage custom playlists
- **Database**: PostgreSQL with GORM ORM
- **API Documentation**: Auto-generated Swagger/OpenAPI 3.0 documentation
- **Docker Support**: Containerized deployment
- **Production Ready**: Koyeb deployment configuration included

## 📚 API Documentation

The API is fully documented using Swagger/OpenAPI 3.0. Once the server is running, you can access the interactive API documentation at:

**Swagger UI**: `http://localhost:8082/swagger/index.html`

**Production API**: `https://independent-carlene-tushar27x-a3461680.koyeb.app/swagger/index.html`

## 🛠️ Getting Started

### Prerequisites

- Go 1.23.5 or higher
- PostgreSQL database
- Git

### Installation

1. **Clone the repository:**
```bash
git clone <repository-url>
cd music-lib-api
```

2. **Install dependencies:**
```bash
go mod tidy
```

3. **Set up environment variables:**
   
   Copy the example environment file:
   ```bash
   cp env.example .env
   ```
   
   Update the `.env` file with your configuration:
   ```env
   # Environment Configuration
   ENVIRONMENT=development
   
   # Server Configuration
   PORT=8082
   DEV_HOST=localhost:8082
   PROD_HOST=api.yourdomain.com
   
   # Database Configuration
   DB_URI=postgresql://username:password@localhost:5432/music_library
   
   # JWT Configuration
   JWT_SECRET=your-super-secret-jwt-key-change-this-in-production
   
   # CORS Configuration (for production)
   CORS_ORIGIN=https://yourdomain.com
   
   # Logging
   LOG_LEVEL=info
   LOG_FORMAT=json
   
   # Security
   RATE_LIMIT=100
   RATE_LIMIT_WINDOW=1m
   ```

4. **Run the application:**
```bash
go run cmd/main.go
```

The server will start on `http://localhost:8082`

## 🔌 API Endpoints

### Authentication
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - User login

### Albums (Requires Authentication)
- `GET /api/albums/` - Get all albums for the user
- `GET /api/albums/search` - Search albums by title or artist
- `GET /api/albums/:id` - Get album by ID
- `POST /api/albums/` - Create a new album (artists only)
- `PUT /api/albums/:id` - Update an album
- `DELETE /api/albums/:id` - Delete an album

### Songs (Requires Authentication)
- `GET /api/songs/` - Get all songs for the user
- `GET /api/songs/search` - Search songs by title or artist
- `GET /api/songs/:id` - Get song by ID
- `POST /api/songs/` - Add a new song
- `PUT /api/songs/:id` - Update a song
- `DELETE /api/songs/:id` - Delete a song

### Playlists (Requires Authentication)
- `GET /api/playlists/` - Get all playlists for the user
- `GET /api/playlists/search` - Search playlists by name
- `GET /api/playlists/:id` - Get playlist by ID
- `POST /api/playlists/` - Create a new playlist
- `PUT /api/playlists/:id` - Update a playlist
- `DELETE /api/playlists/:id` - Delete a playlist

### Health Check
- `GET /api/ping` - Health check endpoint

## 📖 Swagger Documentation

The API documentation is automatically generated from code comments. To regenerate the documentation after making changes:

```bash
swag init -g cmd/main.go
```

### Environment-Specific Swagger Generation

Generate Swagger docs for different environments:

**Windows:**
```cmd
scripts\generate-swagger.bat production
scripts\generate-swagger.bat staging
scripts\generate-swagger.bat development
```

**Linux/Mac:**
```bash
./scripts/generate-swagger.sh production
./scripts/generate-swagger.sh staging
./scripts/generate-swagger.sh development
```

### Updating Swagger Port

If you need to change the port that Swagger uses for API requests:

**Windows:**
```cmd
update-swagger-port.bat 8082
```

**Linux/Mac:**
```bash
./update-swagger-port.sh 8082
```

Replace `8082` with your desired port number.

## 🏗️ Project Structure

```
music-lib-api/
├── cmd/
│   └── main.go                 # Application entry point
├── config/
│   ├── config.go              # Configuration and database setup
│   └── environment.go         # Environment variable management
├── controllers/               # HTTP request handlers
│   ├── albumsController.go    # Album management
│   ├── authController.go      # Authentication
│   ├── playistController.go   # Playlist management
│   └── songsContoller.go      # Song management
├── docs/                      # Generated Swagger documentation
├── middlewares/               # HTTP middlewares
│   └── authMiddleware.go      # JWT authentication middleware
├── models/                    # Data models
│   ├── album.go              # Album model
│   ├── playlist.go           # Playlist model
│   ├── song.go               # Song model
│   └── user.go               # User model
├── routes/                    # Route definitions
│   └── routes.go             # API route configuration
├── services/                  # Business logic layer
├── utils/                     # Utility functions
│   └── debug.go              # Debug utilities
├── scripts/                   # Build and deployment scripts
├── Dockerfile                 # Docker configuration
├── koyeb.yaml                 # Koyeb deployment config
└── env.example               # Environment variables template
```

## 🔐 Authentication

The API uses JWT tokens for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## 🗄️ Database Models

- **User**: Authentication and user management
- **Album**: Music album organization with artist information
- **Song**: Individual music tracks with metadata
- **Playlist**: Collections of songs with custom ordering

## 🐳 Docker Deployment

### Build and Run with Docker

```bash
# Build the Docker image
docker build -t music-lib-api .

# Run the container
docker run -p 8080:8080 --env-file .env music-lib-api
```

### Docker Compose (Optional)

Create a `docker-compose.yml` file for local development:

```yaml
version: '3.8'
services:
  app:
    build: .
    ports:
      - "8080:8080"
    environment:
      - ENVIRONMENT=development
    env_file:
      - .env
    depends_on:
      - db
  
  db:
    image: postgres:15
    environment:
      POSTGRES_DB: music_library
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: password
    ports:
      - "5432:5432"
    volumes:
      - postgres_data:/var/lib/postgresql/data

volumes:
  postgres_data:
```

## 🚀 Production Deployment

### Koyeb Deployment

The project includes a `koyeb.yaml` configuration file for easy deployment to Koyeb:

1. **Set up secrets in Koyeb:**
   - `db-uri`: Your PostgreSQL connection string
   - `jwt-secret`: Your JWT secret key

2. **Deploy using Koyeb CLI:**
```bash
koyeb app init music-lib-api --docker .
```

### Manual Production Setup

1. **Set Environment Variables:**
```bash
ENVIRONMENT=production
PROD_HOST=api.yourdomain.com
DB_URI=your-production-db-uri
JWT_SECRET=your-super-secret-jwt-key
```

2. **Generate Production Swagger Docs:**
```bash
./scripts/generate-swagger.sh production
```

3. **Build and Run:**
```bash
go build -ldflags="-s -w" -o music-api cmd/main.go
ENVIRONMENT=production ./music-api
```

For detailed deployment instructions, see:
- [DEPLOYMENT.md](./DEPLOYMENT.md) - Complete deployment guide
- [DEPLOYMENT-SIMPLE.md](./DEPLOYMENT-SIMPLE.md) - Simplified deployment guide

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

For support and questions:
- Check the [API Documentation](https://independent-carlene-tushar27x-a3461680.koyeb.app/swagger/index.html)
- Open an issue on GitHub
- Contact: support@swagger.io
