# Four-in-Row: Real-time Multiplayer Game Platform

![Go](https://img.shields.io/badge/Backend-Go-00ADD8?style=flat&logo=go)
![React](https://img.shields.io/badge/Frontend-React%20%2B%20Vite-61DAFB?style=flat&logo=react)
![Kafka](https://img.shields.io/badge/Data-Apache%20Kafka-231F20?style=flat&logo=apachekafka)
![Docker](https://img.shields.io/badge/Infra-Docker%20Compose-2496ED?style=flat&logo=docker)
![TypeScript](https://img.shields.io/badge/Language-TypeScript-3178C6?style=flat&logo=typescript)

A production-ready, full-stack implementation of the classic "Four in a Row" game. This project demonstrates a scalable architecture featuring real-time WebSocket communication, an event-driven analytics pipeline, and a heuristic-based CPU opponent.

> **Live Demo:** [[Insert Link Here](https://fourinrow-emittr.onrender.com)] *(Optional: Add link if hosted)*

---

## 🏗 System Architecture

This project goes beyond a simple game by implementing an event-driven architecture suitable for data analysis and high-concurrency environments.

### High-Level Design
* **Frontend (SPA):** Built with React, TypeScript, and Tailwind CSS (Shadcn/UI). Serves as a single-page application managed by the Go router.
* **Game Server (Go):** Handles game state, validation, matchmaking, and WebSocket connections.
* **Event Bus (Kafka):** The application acts as a Producer, emitting game events (`GameStart`, `MoveMade`, `GameEnd`) to a Kafka topic (`game-events`) for downstream processing.
* **Bot Engine:** A server-side CPU opponent using a priority-based heuristic algorithm.

---

## ✨ Key Features

* **Real-time Multiplayer:** Instant state synchronization using WebSockets.
* **Intelligent Bot:** Server-side CPU opponent that calculates best moves based on win/block heuristics (`game/bot`).
* **Event Streaming:** Integrated `IBM/sarama` Kafka producer to stream game metrics for analytics.
* **Fallback Mechanisms:** Intelligent fallback to a "Stub Producer" if the Kafka broker is unreachable (Resiliency).
* **Production Grade UI:** polished interface using `radix-ui` primitives and Tailwind CSS.
* **Containerization:** Full `docker-compose` setup for zookeeper, kafka, and the application services.

---

## 🛠 Tech Stack

### Backend
* **Language:** Golang (1.22+)
* **Communication:** WebSockets (Real-time state), HTTP (REST API)
* **Streaming:** Apache Kafka (via Sarama library)
* **Architecture:** Hexagonal-inspired (handlers, service logic, repository, analytics layers isolated).

### Frontend
* **Framework:** React 18 + Vite
* **Language:** TypeScript
* **Styling:** Tailwind CSS + Shadcn/UI
* **State:** React Query (Tanstack)

### DevOps
* **Containerization:** Docker & Docker Compose
* **Build Pipeline:** Multi-stage builds (Frontend build embedded into Go binary).

---

## 🚀 Getting Started

The easiest way to run the application is via Docker Compose, which spins up the Kafka broker, Zookeeper, and the Game Server.

### Prerequisites
* Docker & Docker Compose
* Node.js (for local frontend dev only)
* Go 1.22+ (for local backend dev only)

### Running with Docker (Recommended)
```bash
# Clone the repository
git clone [https://github.com/your-username/four-in-row.git](https://github.com/your-username/four-in-row.git)
cd four-in-row

# Start all services (Kafka, Zookeeper, App)
docker-compose up --build
