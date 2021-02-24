FROM flant/shell-operator:v1.0.0-rc.1
ADD hooks /hooks
ADD deploy /deploy
RUN chmod -R 777 /hooks
RUN chmod -R 777 /var/run/

USER 10001