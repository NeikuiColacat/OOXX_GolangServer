# API Documentation

## 注册

**URL:** `http://10.33.74.157:8080/user/register`

**方法:** `POST`

**BODY:**

```json
{
    "username": "username",
    "password": "passwd"
    "repassword": "repasswd"
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

(返回一个json作为BODY,JSON按照用户积分降序排序展示整张users表，如下为返回json BODY示例)

```json
[
    {
        "username": "maomao",
        "password": "lueluelue",
        "score": 200
    },
    {
        "username": "test",
        "password": "test",
        "score": 123
    },
    {
        "username": "john_doe",
        "password": "secure_password",
        "score": 100
    },
    {
        "username": "haha",
        "password": "newpassword",
        "score": 50
    },
    {
        "username": "colacola",
        "password": "newpassword",
        "score": 0
    },
    {
        "username": "newuser",
        "password": "newpassword",
        "score": 0
    }
]
```

## MySQL :: game_db数据库 :: users表如下示例

| username | password        | score |
|----------|-----------------|-------|
| colacola | newpassword     |     0 |
| haha     | newpassword     |    50 |
| john_doe | secure_password |   100 |
| maomao   | lueluelue       |   200 |
| newuser  | newpassword     |     0 |
| test     | test            |   123 |


---

## 积分更新 

**URL:**  `http://10.33.74.157:8080/user/score`

**方法:** `POST`

**BODY:**

```json
{
    "username": "123",
    "score": 100
}

```

原有的用户积分上加上scroe

**状态码**

成功: `200` 

---

## 查询单个用户积分

**URL:**  `http://10.33.74.157:8080/user/query`

**方法:** `POST`

**BODY:**

```json
{
    "username": "maomao"
}

```
**返回**
```json
{
    "score": 200
}

```

---

## 请求加入匹配队列

**URL:**  `http://10.33.74.157:8080/user/query_queue`

**方法:** `POST`

**BODY:**

```json
{
  "username": "匹配玩家1"
}
```

该请求会将用户名为 `匹配玩家1` 的玩家加入到匹配队列中 ， 等待匹配

---

## 查询匹配结果


**URL:**  `http://10.33.74.157:8080/user/query_match`

**方法:** `POST`

**BODY:**

```json
{
  "username": "匹配玩家1"
}
```

**返回BODY:**

```json
{
    "username": "匹配玩家1",
    "match": "匹配玩家2"
}
```

该请求会返回一个json 查询 `匹配玩家1`  的对战玩家 `匹配玩家2`

**状态码**

尚未匹配成功 : `404`

匹配成功 : `200`

---

## 移除玩家匹配

**URL:**  `http://10.33.74.157:8080/user/query_remove`

**方法:** `POST`

**BODY:**

```json
{
  "username": "匹配玩家1"
}
```

该请求会将 `匹配玩家1`  和配对的 `匹配玩家2`  移除匹配完成状态 ， 代表对战结束

**状态码**

移除成功 : `200`

`匹配玩家1` 当前尚未匹配成功 : `404`


---

## 插入时间

**URL:**  `http://10.33.74.157:8080/user/time`

**方法:** `POST`

**BODY:**

```json
{
  "username": "testUser",
  "time": "2023-04-01T12:00:00Z"
}
```

向map中用户名为 `testuser` 的用户映射一个时间字符串

**状态码**

成功 : `200`

---

## 查询时间

**URL:**  `http://10.33.74.157:8080/user/query_time`

**方法:** `POST`

**BODY:**

```json
{
  "username": "testUser"
}
```

**返回BODY**

```json
{
    "time": "2023-04-01T12:00:00Z"
}
```

返回一个JSON作为用户的提交时间

**状态码**

map中键为空 ， 查询失败返回 : `404`