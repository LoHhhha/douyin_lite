# DouYin_lite

## tables in database

* users

    * desc
      
        ```
        +-----------------+--------+------+-----+---------+----------------+
        | Field           | Type   | Null | Key | Default | Extra          |
        +-----------------+--------+------+-----+---------+----------------+
        | id              | bigint | NO   | PRI | NULL    | auto_increment |
        | name            | text   | NO   |     | NULL    |                |
        | follow_count    | bigint | NO   |     | 0       |                |
        | follower_count  | bigint | NO   |     | 0       |                |
        | total_favorited | bigint | NO   |     | 0       |                |
        | password        | text   | NO   |     | NULL    |                |
        | work_count      | bigint | NO   |     | 0       |                |
        | favorite_count  | bigint | NO   |     | 0       |                |
        +-----------------+--------+------+-----+---------+----------------+
        ```
    
    * code

        ```mysql
        CREATE TABLE `users` (
        `id` bigint NOT NULL AUTO_INCREMENT,
        `name` text NOT NULL,
        `follow_count` bigint NOT NULL DEFAULT '0',
        `follower_count` bigint NOT NULL DEFAULT '0',
        `total_favorited` bigint NOT NULL DEFAULT '0',
        `password` text NOT NULL,
        `work_count` bigint NOT NULL DEFAULT '0',
        `favorite_count` bigint NOT NULL DEFAULT '0',
        PRIMARY KEY (`id`)
        ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
        ```

* videos

    * desc
      
        ```
        +----------------+--------------+------+-----+---------+----------------+
        | Field          | Type         | Null | Key | Default | Extra          |
        +----------------+--------------+------+-----+---------+----------------+
        | vid            | bigint       | NO   | PRI | NULL    | auto_increment |
        | uid            | bigint       | NO   | MUL | NULL    |                |
        | play_url       | varchar(255) | NO   |     | NULL    |                |
        | cover_url      | varchar(255) | YES  |     | NULL    |                |
        | favorite_count | bigint       | YES  |     | 0       |                |
        | comment_count  | bigint       | YES  |     | 0       |                |
        | uploadtime     | datetime     | NO   |     | NULL    |                |
        | title          | varchar(255) | NO   |     | NULL    |                |
        +----------------+--------------+------+-----+---------+----------------+
        ```
    
    * code

        ```mysql
        CREATE TABLE `videos` (
        `vid` bigint NOT NULL AUTO_INCREMENT,
        `uid` bigint NOT NULL,
        `play_url` varchar(255) NOT NULL COMMENT '视频url',
        `cover_url` varchar(255) DEFAULT NULL COMMENT '封面url',
        `favorite_count` bigint DEFAULT '0',
        `comment_count` bigint DEFAULT '0',
        `uploadtime` datetime NOT NULL COMMENT '上传时间',
        `title` varchar(255) NOT NULL,
        PRIMARY KEY (`vid`),
        KEY `fk_vid_uid` (`uid`)
        ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
        ```

* comments

    * desc
      
        ```
        +-------------+---------------+------+-----+---------+----------------+
        | Field       | Type          | Null | Key | Default | Extra          |
        +-------------+---------------+------+-----+---------+----------------+
        | id          | bigint        | NO   | PRI | NULL    | auto_increment |
        | vid         | bigint        | YES  | MUL | NULL    |                |
        | uid         | bigint        | YES  | MUL | NULL    |                |
        | content     | varchar(1024) | NO   |     | NULL    |                |
        | commentdate | datetime      | NO   |     | NULL    |                |
        +-------------+---------------+------+-----+---------+----------------+
        ```
    
    * code

        ```mysql
        CREATE TABLE `comments` (
        `id` bigint NOT NULL AUTO_INCREMENT,
        `vid` bigint DEFAULT NULL,
        `uid` bigint DEFAULT NULL,
        `content` varchar(1024) NOT NULL,
        `commentdate` datetime NOT NULL,
        PRIMARY KEY (`id`),
        KEY `fk_cvid_vvid` (`vid`),
        KEY `fk_vuid_uuid` (`uid`)
        ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
        ```

* favorites

    * desc
      
        ```
        +-------+--------+------+-----+---------+----------------+
        | Field | Type   | Null | Key | Default | Extra          |
        +-------+--------+------+-----+---------+----------------+
        | id    | bigint | NO   | PRI | NULL    | auto_increment |
        | uid   | bigint | NO   | MUL | NULL    |                |
        | vid   | bigint | NO   |     | NULL    |                |
        +-------+--------+------+-----+---------+----------------+

        ```
    
    * code

        ```mysql
        CREATE TABLE `favorites` (
        `id` bigint NOT NULL AUTO_INCREMENT,
        `uid` bigint NOT NULL,
        `vid` bigint NOT NULL,
        PRIMARY KEY (`id`),
        KEY `favorites_uid_vid_index` (`uid`,`vid`)
        ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
        ```

* follows

    * desc

        ```
        +--------------+--------+------+-----+---------+----------------+
        | Field        | Type   | Null | Key | Default | Extra          |
        +--------------+--------+------+-----+---------+----------------+
        | uid          | bigint | YES  | MUL | NULL    |                |
        | id           | bigint | NO   | PRI | NULL    | auto_increment |
        | follower_uid | bigint | YES  | MUL | NULL    |                |
        +--------------+--------+------+-----+---------+----------------+
        ```

    * code

        ```mysql
        CREATE TABLE `follows` (
        `uid` bigint DEFAULT NULL,
        `id` bigint NOT NULL AUTO_INCREMENT,
        `follower_uid` bigint DEFAULT NULL,
        PRIMARY KEY (`id`),
        KEY `follow_follower_uid_index` (`follower_uid`),
        KEY `follow_uid_index` (`uid`)
        ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
        ```

* messages

    * desc

        ```
        +--------------+--------+------+-----+---------+----------------+
        | Field        | Type   | Null | Key | Default | Extra          |
        +--------------+--------+------+-----+---------+----------------+
        | id           | bigint | NO   | PRI | NULL    | auto_increment |
        | to_user_id   | bigint | YES  | MUL | NULL    |                |
        | from_user_id | bigint | YES  | MUL | NULL    |                |
        | content      | text   | YES  |     | NULL    |                |
        | create_time  | int    | YES  |     | NULL    |                |
        +--------------+--------+------+-----+---------+----------------+
        ```

    * code

        ```mysql
        CREATE TABLE `messages` (
        `id` bigint NOT NULL AUTO_INCREMENT,
        `to_user_id` bigint DEFAULT NULL,
        `from_user_id` bigint DEFAULT NULL,
        `content` text,
        `create_time` int DEFAULT NULL,
        PRIMARY KEY (`id`),
        KEY `messages_from_user_id_index` (`from_user_id`),
        KEY `messages_to_user_id_index` (`to_user_id`)
        ) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
        ```

## test case

```mysql
insert into users (name,password) values
("lohhhha","123456"),
("testuser1","123456"),
("testuser2","123456"),
("testuser3","123456"),
("testuser4","123456"),
("testuser5","123456");

insert into videos (uid,play_url,cover_url,favorite_count,comment_count,uploadtime,title) values
(1,"/static/videos/1691829728_cat.mp4","/static/covers/1691829728_cat.jpg",0,0,"2023-08-12 16:42:09","cat");

insert into follows (uid,follower_uid) values
(1,2),
(1,3),
(1,4),
(1,5),
(1,6),
(2,1);
```