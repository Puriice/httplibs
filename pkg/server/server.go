package server

import "github.com/jackc/pgx/v5/pgxpool"

type Server struct {
	Host     string
	Port     string
	Database *pgxpool.Pool
}

func NewServer(host string, port string, db *pgxpool.Pool) *Server {
	return &Server{
		Host:     host,
		Port:     port,
		Database: db,
	}
}
