# Music Library API

A RESTful API for managing music library with albums, songs, and playlists built with Go, Gin, and GORM.

## Features

- User authentication with JWT
- Album management (artists only)
- Song management
- Playlist creation and management
- PostgreSQL database with GORM
- Swagger API documentation

## API Documentation

The API is fully documented using Swagger/OpenAPI 3.0. Once the server is running, you can access the interactive API documentation at:

**Swagger UI**: `http://localhost:8082/swagger/index.html`

## Getting Started

### Prerequisites

- Go 1.23.5 or higher
- PostgreSQL database

### Installation

1. Clone the repository:
```bash
git clone <repository-url>
cd music-lib-api
```

2. Install dependencies:
```bash
go mod tidy
```

3. Set up environment variables in a `.env` file:
```env
PORT=8080
DB_HOST=localhost
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_PORT=5432
JWT_SECRET=your_jwt_secret
```

4. Run the application:
```bash
go run cmd/main.go
```

The server will start on `http://localhost:8082`

### API Endpoints

#### Authentication
- `POST /api/auth/register` - Register a new user
- `POST /api/auth/login` - User login

#### Albums (Requires Authentication)
- `GET /api/albums/getAllAlbums` - Get all albums for the user
- `POST /api/albums/addAlbum` - Create a new album (artists only)

#### Songs (Requires Authentication)
- `GET /api/songs/getAllSongs` - Get all songs for the user
- `POST /api/songs/addSong` - Add a new song

#### Playlists (Requires Authentication)
- `GET /api/playlists/getAllPlaylists` - Get all playlists for the user
- `POST /api/playlists/addPlaylist` - Create a new playlist

## Swagger Documentation

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

If you need to change the port that Swagger uses for API requests, you can use the provided scripts:

**Windows:**
```cmd
update-swagger-port.bat 8082
```

**Linux/Mac:**
```bash
./update-swagger-port.sh 8082
```

Replace `8082` with your desired port number. The script will automatically update all necessary files and regenerate the Swagger documentation.

## Project Structure

```
music-lib-api/
├── cmd/
│   └── main.go          # Application entry point
├── config/
│   └── config.go        # Configuration and database setup
├── controllers/          # HTTP request handlers
├── docs/                # Generated Swagger documentation
├── middlewares/         # HTTP middlewares (auth, etc.)
├── models/              # Data models
├── routes/              # Route definitions
└── utils/               # Utility functions
```

## Authentication

The API uses JWT tokens for authentication. Include the token in the Authorization header:

```
Authorization: Bearer <your-jwt-token>
```

## Database Models

- **User**: Authentication and user management
- **Album**: Music album organization
- **Song**: Individual music tracks
- **Playlist**: Collections of songs

## Contributing

1. Fork the repository
2. Create a feature branch
3. Commit your changes
4. Push to the branch
5. Create a Pull Request

## 🚀 **Production Deployment**

For production deployment instructions, see [DEPLOYMENT.md](./DEPLOYMENT.md).

**Don't have a staging environment?** See [DEPLOYMENT-SIMPLE.md](./DEPLOYMENT-SIMPLE.md) for a simplified guide.

### **Quick Production Setup**

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

## License

This project is licensed under the MIT License.
