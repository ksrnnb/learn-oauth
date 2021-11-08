FROM mariadb:10.6

RUN apt-get update \
    && apt-get -y install locales-all

ENV LANG ja_JP.UTF-8
ENV LANGUAGE ja_JP:ja
ENV LC_ALL ja_JP.UTF-8

COPY ./docker/db/mysqld_charset.cnf /etc/mysql/conf.d/mysqld_charset.cnf