# INSTRUCTION TO DEPLOY A K3S CLUSTER USING ORANGE PIs

## **Introduction**
- This is a project from course ***Wireless Embedded System (NT131)*** at University of ***Information Technology (VNUHCM)***
- In this project, our team deployed a high availability K3S cluster using 5 Orange Pis model 2+. We demonstrated Kubernetes features by built on top of it a lightweight web application written in Golang, applied Kubernetes resources such as Deployment, Service, HorizontalPodAutoScaler, and paused random nodes to test for HA. We also deployed a Kubernetes Dashboard for easy management purpose
- This documentation only mention about the Notebook web application deployment. Details about hardware infrastructure and K3S cluster set up will be update later
## **System Architecture**
> Will be update later
## **Requirements** 
- To follow this instruction, you need:
    + An AWS account has permissions to interact with Elastic Container Registry (ECR)
    + A repository named **notebook**
    + Docker installed
    + K3s installed

> Now let's get started!

### **2. AWS Account Configuration**
- Configure AWS credentials, then enter your access key and secret access key:
```bash
aws configure
```
- Login to ECR repository:
```bash
aws ecr get-login-password --region <your-region> | docker login --username AWS --password-stdin <your-account-id>.dkr.ecr.<your-region>.amazonaws.com
```

### **2. Docker Image Preparation**
- First of all, clone this repository to your host:
``` bash
git clone https://github.com/tramcandoit/nt131-notebook
cd nt131-notebook
```
- Build the provided Dockerfile:
```bash
docker build --platform linux/arm/v7 -t notebook .
```
- Push that image to ECR:
```bash
docker tag notebook:latest <your-account-id>.dkr.ecr.<your-region>.amazonaws.com/notebook:<version>
docker push <your-account-id>.dkr.ecr.<your-region>.amazonaws.com/notebook:<version>
```

### **3. Notebook Deployment**
- Move to **infra**:
```bash
cd infra
```
- Apply PV and PVC:
```bash
kubectl apply -f redis-pv.yaml
kubectl apply -f redis-pvc.yaml
```
- Run Redis in advance:
```bash
kubectl apply -f redis-master-depl.yaml
kubectl apply -f redis-master-service.yaml
```
- Run Notebook:
```bash
kubectl apply -f notebook-depl.yaml
kubectl apply -f notebook-service.yaml
```
- Check Notebook Endpoint:
```bash
kubectl get svc
```
- Check the Notebook web content:
```bash
curl <node-ip>:<port>
```
