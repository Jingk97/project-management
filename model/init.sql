/*
    创建数据库后创建对应服务模块连接sql账号，并且开启远程数据库连接，然后针对数据库赋权
    SHOW DATABASES;
*/
CREATE DATABASE projectManagement DEFAULT CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

CREATE USER 'app-user'@'%' IDENTIFIED BY 'user@EMG1Fv6+!';

GRANT SELECT,INSERT,UPDATE,DELETE,CREATE,ALTER,INDEX ON projectManagement.* TO 'app-user'@'%';

FLUSH PRIVILEGES;