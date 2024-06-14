# API Documentation

## 注册

**URL:** `http://10.33.74.157:8080/user/register`

**方法:** `POST`

**BODY:**

```json
{
    "username": "username",
    "password": "passwd"
}

```

**状态码**

成功: `200` 

json格式错误: `400` 

用户名已存在: `409` 

---

##  登录

**URL:**  `http://10.33.74.157:8080/user/login`

**方法:** `POST`

**BODY:**

```json
{
    "username": "username",
    "password": "newpasswd",
    "repassword": "repasswd"
}


```

**状态码**

成功: `200` 

json格式错误: `400` 

用户名不存在: `404` 

密码错误 : `401`


---

## 积分排行榜

**URL:**  `http://10.33.74.157:8080/user/rank`

**方法:** `GET`

**状态码**

成功: `200` 

(返回一个json作为BODY,JSON按照用户积分降序排序展示整张users表)

## MySQL :: game_db数据库 :: users表如下示例

| username | password        | score |
|----------|-----------------|-------|
| colacola | newpassword     |     0 |
| haha     | newpassword     |    50 |
| john_doe | secure_password |   100 |
| maomao   | lueluelue       |   200 |
| newuser  | newpassword     |     0 |
| test     | test            |   123 |
