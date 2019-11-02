sudo docker system prune -f
sudo docker image build --no-cache --label single-app -t single-app .
sudo docker container stop $(sudo docker container ls -aq)
sudo docker run -it --log-driver none -p 8080:8080 --restart always single-app