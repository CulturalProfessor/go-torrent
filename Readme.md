# 🚀 Go Torrent 📂

## 📦 Project Overview

The **Go Torrent Project** is a web-based application that enables users to upload torrent files, monitor download progress, and download completed files. Built with a **Next.js client** 🌐 and a **Go server** 🐹, this project is containerized with Docker, making deployment simple and efficient across different environments.

## 🚀 Features

- **Upload Torrent Files**: Easily upload torrent files to start downloads.
- **Monitor Download Progress**: Track the real-time download progress of torrents.
- **Download Completed Files**: Download fully completed torrent files directly.
- **Docker Support**: Deploy the entire system with Docker for a streamlined setup.

---

## 📂 Project Structure

```bash
├── client/                     # Next.js frontend for upload and progress tracking
├── server/                     # Go backend for torrent downloading and serving files
├── .gitignore                  # Git ignore file for unnecessary files and folders
├── docker-compose.yml          # Docker Compose setup for client and server
├── Dockerfile                  # Dockerfile for building server and client images
├── makefile                    # Makefile for easier development
└── README.md                   # Project documentation
```

## 🛠️ Installation
### Prerequisites

- **Docker**: Make sure Docker is installed and running on your machine.

### Steps

1. **Clone the repository**:

   ```bash
   git clone https://github.com/CulturalProfessor/go-torrent.git
   cd go-torrent
   ```

2. **Build and run the Docker container**:

   ```bash
   make build
   # or alternatively
   docker-compose up --build
   ```

3. **Access the system**:

Once the containers are up and running, open your browser and navigate to:

- **Client (Next.js): http://localhost:3000 to access the web interface.**
- **Server (Go): API available at http://localhost:8080.**

### 🎉 Sample Torrent File

For testing purposes, you can download a sample torrent file here:  
📥 [Download Sample Torrent File](https://drive.google.com/file/d/1At-GhxB_7PGoYAh7IsrTaEw63uiyRkJt/view?usp=sharing)

