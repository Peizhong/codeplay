FROM jenkins/jenkins:lts
USER root
RUN apt update && apt -y upgrade
RUN groupadd docker -g 999
RUN usermod -aG docker jenkins
USER jenkins

# docker build . -t peizhong/jenkins-docker:lts