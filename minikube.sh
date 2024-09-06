eval $(minikube docker-env)
docker build -t chat-controller:latest ./chat_controller

eval $(minikube docker-env)
docker build -t chat-backend:latest ./chat_backend

eval $(minikube docker-env)
docker build -t chat-consumer:latest ./chat_consumer


kubectl apply -f chat-apps-manifest.yaml

kubectl rollout restart deployment chat-controller-deployment chat-backend-deployment chat-consumer-deployment

kubectl get pods