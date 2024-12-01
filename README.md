# Community Event Management Platform

## DevOps Learning Project

This project demonstrates DevOps principles and modern tooling through a Community Event Management Platform.

### Architecture Overview

- Microservices Architecture
- Containerized Deployment
- CI/CD Pipeline
- Infrastructure as Code
- Monitoring and Observability

### Services

1. Event Service (Python/Flask)
2. User Service (Node.js)
3. Notification Service (Go)

### DevOps Tools and Technologies

- Containerization: Docker
- Orchestration: Kubernetes
- CI/CD: GitHub Actions
- Infrastructure as Code: Terraform
- Monitoring: Prometheus, Grafana
- Logging: ELK Stack

## Prerequisites

Before you begin, ensure you have the following installed:
- Docker
- Kubernetes (minikube or kind)
- Terraform
- AWS CLI
- kubectl
- Helm

## Local Development

### 1. Clone the Repository
```bash
git clone https://github.com/yourusername/devops-event-platform.git
cd devops-event-platform
```

### 2. Local Service Development

#### Event Service (Python/Flask)
```bash
cd services/event-service
python -m venv venv
source venv/bin/activate
pip install -r requirements.txt
python app.py
```

#### User Service (Node.js)
```bash
cd services/user-service
npm install
npm start
```

#### Notification Service (Go)
```bash
cd services/notification-service
go mod download
go run main.go
```

### 3. Docker Containerization

Build Docker images:
```bash
docker build -t event-service ./services/event-service
docker build -t user-service ./services/user-service
docker build -t notification-service ./services/notification-service
```

### 4. Kubernetes Deployment

#### Local Cluster Setup
```bash
minikube start
kubectl apply -f infrastructure/k8s/
```

### 5. Infrastructure as Code

#### Terraform AWS EKS Deployment
```bash
cd infrastructure/terraform
terraform init
terraform plan
terraform apply
```

### 6. Monitoring Setup

#### Prometheus and Grafana
```bash
helm repo add prometheus-community https://prometheus-community.github.io/helm-charts
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update

helm install prometheus prometheus-community/prometheus
helm install grafana grafana/grafana
```

## DevOps Workflow

1. **Continuous Integration**: GitHub Actions runs tests on every push
2. **Containerization**: Docker images built for each service
3. **Deployment**: Kubernetes manages service orchestration
4. **Monitoring**: Prometheus and Grafana track system health

## Real-World Use Cases

1. **Event Management**
   - Create, list, and manage community events
   - User registration and authentication
   - Real-time event notifications

2. **Scalable Microservices**
   - Independent service deployment
   - Horizontal scaling
   - Fault isolation

3. **DevOps Automation**
   - CI/CD pipeline
   - Infrastructure as Code
   - Comprehensive monitoring
