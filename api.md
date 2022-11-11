
## 接口
接口域名 `http://127.0.0.1:web_port`

| 接口        | 说明   |
| --------   | :-----  |
| [/api/v0/tweet/user/timeline](#api-v0-tweet-user-timeline)| 我的推特列表 |
| [/api/v0/tweet/release](#api-v0-tweet-release)| 发送推特 |
| [/api/v0/tweet/explore](#api-v0-tweet-explore)| 推特广场 |
| [/api/v0/tweet/user/timeline/:id](#api-v0-tweet-user-timeline-id)| 指定用户发送推特列表 |
| [/api/v0/user/:id](#api-v0-user-id)| 指定用户信息 |
| [/api/v0/tweet/forward](#api-v0-tweet-forward)| 转发推特 |
| [/api/v0/tweet/user/forward](#api-v0-tweet-user-forward)| 我的转发 |
| [/api/v0/user/profile](#api-v0-user-profile)| 当前用户信息 |

### /api/v0/tweet/user/timeline

| 接口说明 | HTTP请求方式 |
| :--- | --- |
| 我的推特列表 | GET |

#### 请求参数
| 名称 | 必须 | 类型及范围 | 说明   |
| --- | --- | --- | --- |
| page | false | int | 页码 |

#### 返回结果
```json
{
    "code": 0,
    "msg": "获取成功",
    "data": [
        {
            "Id": "264zpU8iQdgJ45ZyAkdhVdYxjMxpLFZPGrrJGtwPhbzCG2J2kjrhm3wsq9V4W6NdQxXwxsg24mxTP8oiGV9ZrErc",
            "UserId": "12D3KooWGhtFCSa7zN4AJvp2HDwwP59cNuHRBqPpcF9TncdGUS5V",
            "Content": "123123123",
            "Attachment": "",
            "Nonce": 1,
            "Sign": "4Fywn2aqqiAjqQXLw1B4TLBokXitGuJkvooCtamfa11ZAn6wqVpuWytBQuPohFJxKvTpmDL5dzxRJcuSQBbrwuz7",
            "UserInfo": null,
            "CreatedAt": "2021-04-20T16:56:29.834016+08:00",
            "UpdatedAt": "2021-04-20T16:56:29.834016+08:00"
        }
    ]
}
```



### /api/v0/tweet/release

| 接口说明 | HTTP请求方式 |
| :--- | --- |
| 发送推特 | POST |

#### 请求参数
| 名称 | 必须 | 类型及范围 | 说明   |
| --- | --- | --- | --- |
| forward_id | false | string | 转发内容的id 如果存在则会忽略content和附件 |
| content | true | string | 内容 |
| attachment | false | string | 内容 |

#### 返回结果
```json
{
  "code": 0,
  "msg": "发送成功",
  "data": {
    "Id": "xuqpZci8VPPpfZ5YMVxJMHbSABAeY3xhbCQ2w4nxBnTCXA49m1rtA82FfszbJ7nJor3jtDbUA688sgkYjukxB18",
    "UserId": "12D3KooWGhtFCSa7zN4AJvp2HDwwP59cNuHRBqPpcF9TncdGUS5V",
    "Content": "123123123",
    "Attachment": "",
    "Nonce": 2,
    "Sign": "41CWRMMt2FADv5Wx6XGoHiaHGgd3LEShzytmYPdqhmq794qKrofcbt2djobb1zqgUGhL4pxiaEjZomwMLbt2NzNB",
    "CreatedAt": "2021-04-20T16:57:04.20815+08:00",
    "UpdatedAt": "2021-04-20T16:57:04.20815+08:00"
  }
}
```

### /api/v0/tweet/explore

| 接口说明 | HTTP请求方式 |
| :--- | --- |
| 推特广场 | GET |

#### 请求参数
| 名称 | 必须 | 类型及范围 | 说明   |
| --- | --- | --- | --- |
| page | false | int | 页码 |

#### 返回结果
```json
{
  "code": 0,
  "msg": "获取成功",
  "data": [
    {
      "Id": "264zpU8iQdgJ45ZyAkdhVdYxjMxpLFZPGrrJGtwPhbzCG2J2kjrhm3wsq9V4W6NdQxXwxsg24mxTP8oiGV9ZrErc",
      "UserId": "12D3KooWGhtFCSa7zN4AJvp2HDwwP59cNuHRBqPpcF9TncdGUS5V",
      "Content": "123123123",
      "Attachment": "",
      "Nonce": 1,
      "Sign": "4Fywn2aqqiAjqQXLw1B4TLBokXitGuJkvooCtamfa11ZAn6wqVpuWytBQuPohFJxKvTpmDL5dzxRJcuSQBbrwuz7",
      "UserInfo": null,
      "CreatedAt": "2021-04-20T16:56:29.834016+08:00",
      "UpdatedAt": "2021-04-20T16:56:29.834016+08:00"
    }
  ]
}
```

### /api/v0/tweet/user/timeline/:id

| 接口说明 | HTTP请求方式 |
| :--- | --- |
| 指定用户发送推特列表 | GET |

#### 请求参数
| 名称 | 必须 | 类型及范围 | 说明   |
| --- | --- | --- | --- |
| page | false | int | 页码 |

#### 返回结果
```json
{
  "code": 0,
  "msg": "获取成功",
  "data": [
    {
      "Id": "264zpU8iQdgJ45ZyAkdhVdYxjMxpLFZPGrrJGtwPhbzCG2J2kjrhm3wsq9V4W6NdQxXwxsg24mxTP8oiGV9ZrErc",
      "UserId": "12D3KooWGhtFCSa7zN4AJvp2HDwwP59cNuHRBqPpcF9TncdGUS5V",
      "Content": "123123123",
      "Attachment": "",
      "Nonce": 1,
      "Sign": "4Fywn2aqqiAjqQXLw1B4TLBokXitGuJkvooCtamfa11ZAn6wqVpuWytBQuPohFJxKvTpmDL5dzxRJcuSQBbrwuz7",
      "UserInfo": null,
      "CreatedAt": "2021-04-20T16:56:29.834016+08:00",
      "UpdatedAt": "2021-04-20T16:56:29.834016+08:00"
    }
  ]
}
```

### /api/v0/user/:id

| 接口说明 | HTTP请求方式 |
| :--- | --- |
| 指定用户信息 | GET |

#### 返回结果
```json
{
  "code": 0,
  "msg": "获取成功",
  "data": {
    "Id": "12D3KooWSj3G2XTFscfaVEwkf5C6UDmY9MGLJyk8PskfhA211Mph",
    "Name": "OeFSzp",
    "Desc": "",
    "LatestCid": "QmX16PFRRvTDLgkVXjbUhRGzm16RfmD2F3pStivW7hYNo5",
    "Avatar": "",
    "Nonce": 0,
    "LocalNonce": 0,
    "Sign": "3NEtNSHGDPM1AihLyL3nfzJzDaWXeGWZVuVfDJSQaLZ8LKv1kkH7gbweptbbe6wU8B66psbUpZqmcMbsokr9tiHJMLbRZ7sUVGHApg9FMfNFgZu",
    "PubKey": "4XTTMJ3iusfvLUyoEWuQAn75vcpRCd2FnSCMN86ARjPdXy8bo",
    "CreatedAt": "0001-01-01T00:00:00Z",
    "UpdatedAt": "0001-01-01T00:00:00Z"
  }
}
```

### /api/v0/tweet/user/forward

| 接口说明 | HTTP请求方式 |
| :--- | --- |
| 推特转发 | GET |

#### 请求参数
| 名称 | 必须 | 类型及范围 | 说明   |
| --- | --- | --- | --- |
| page | false | int | 页码 |

#### 返回结果
```json
{
  "code": 0,
  "msg": "success",
  "data": [
    {
      "Id": "24uAkYuaLefzKwGVu9ZVHmPy3ANSx1G6SbfEwuXsWbzY1djvJGARkGcEwx43LHSCCpnzqbFKdAT8VbfiHoUwyfZa",
      "UserId": "12D3KooWNJcjhN4c3Uv74eC2ceTjRTEVpGCGdRkfAQc6CXcMS5HH",
      "Content": "手机发送好像有问题",
      "Attachment": "",
      "Nonce": 13,
      "Sign": "4SuX1VEPtJPUEz3fwGTxxrb5y4GMW5Ciy9cbvSM46JqiLq7sJxPG19qE7oqDB7zfjHWvWt2kWWtDDkfyuap6qenq",
      "CreatedAt": "2021-04-26T17:59:20.2950209+08:00",
      "UpdatedAt": "2021-04-26T17:59:20.2950209+08:00",
      "OriginUserId": "12D3KooWR6UqF4z1Qbe3MzcMYR7kGQBEEqb1MtszwDzx2khwki4o",
      "ForwardAt": "2021-04-26T17:59:18.4155214+08:00",
      "OriginTwId": "2y73tDb3VtdscSd47aABTgkuXbmmwvtwY2Mo82sWXCpfKztyXPJ3D8rn5oQ2UcZ6D98HXHMzjS9FvvuvoJbDGu2F"
    },
    {
      "Id": "29YCTMsgs5Pauib9PXH2cxyGYucnCY6v4EHRgxJX3sXJcX7MHSR4eLTSHjo9QPifCbYaFrTPkcE8uZ8jfLNA1odi",
      "UserId": "12D3KooWNJcjhN4c3Uv74eC2ceTjRTEVpGCGdRkfAQc6CXcMS5HH",
      "Content": "手机发送好像有问题",
      "Attachment": "",
      "Nonce": 12,
      "Sign": "37pX8eRRiCDAkW5yjYxC8SrUYkTphkDLwg5kQ24Uxe6BGLpz3KuWBk2PgXHxPqyxBXcARVpPJ71GuHVYShyHHinD",
      "CreatedAt": "2021-04-26T17:58:25.1323159+08:00",
      "UpdatedAt": "2021-04-26T17:58:25.1323159+08:00",
      "OriginUserId": "12D3KooWR6UqF4z1Qbe3MzcMYR7kGQBEEqb1MtszwDzx2khwki4o",
      "ForwardAt": "2021-04-26T17:58:23.2830045+08:00",
      "OriginTwId": "2y73tDb3VtdscSd47aABTgkuXbmmwvtwY2Mo82sWXCpfKztyXPJ3D8rn5oQ2UcZ6D98HXHMzjS9FvvuvoJbDGu2F"
    },
    {
      "Id": "28BxAaMpzZ3LbG2sUrbT4PtwGkTJhQukXZKBspoZDmBFffzdXn5E7tczLuGNdiKykwGPoKjEzwdYwR9GGCsmH54u",
      "UserId": "12D3KooWNJcjhN4c3Uv74eC2ceTjRTEVpGCGdRkfAQc6CXcMS5HH",
      "Content": "手机发送好像有问题",
      "Attachment": "",
      "Nonce": 11,
      "Sign": "2vasiD9q7rFXhkBty698hqCaMUVK19wXPEbWLLdjinUsXpKHmVtkosyVMHFFxYwt8te13s5cnqcmoVPB8ptXK4yZ",
      "CreatedAt": "2021-04-26T17:55:29.0511492+08:00",
      "UpdatedAt": "2021-04-26T17:55:29.0511492+08:00",
      "OriginUserId": "12D3KooWR6UqF4z1Qbe3MzcMYR7kGQBEEqb1MtszwDzx2khwki4o",
      "ForwardAt": "2021-04-26T17:55:27.1086883+08:00",
      "OriginTwId": "2y73tDb3VtdscSd47aABTgkuXbmmwvtwY2Mo82sWXCpfKztyXPJ3D8rn5oQ2UcZ6D98HXHMzjS9FvvuvoJbDGu2F"
    }
  ]
}
```

### /api/v0/user/profile

| 接口说明 | HTTP请求方式 |
| :--- | --- |
| 指定用户信息 | GET |

#### 返回结果
```json
{
  "code": 0,
  "msg": "获取成功",
  "data": {
    "Id": "12D3KooWSj3G2XTFscfaVEwkf5C6UDmY9MGLJyk8PskfhA211Mph",
    "Name": "OeFSzp",
    "Desc": "",
    "LatestCid": "QmX16PFRRvTDLgkVXjbUhRGzm16RfmD2F3pStivW7hYNo5",
    "Avatar": "",
    "Nonce": 0,
    "LocalNonce": 0,
    "Sign": "3NEtNSHGDPM1AihLyL3nfzJzDaWXeGWZVuVfDJSQaLZ8LKv1kkH7gbweptbbe6wU8B66psbUpZqmcMbsokr9tiHJMLbRZ7sUVGHApg9FMfNFgZu",
    "PubKey": "4XTTMJ3iusfvLUyoEWuQAn75vcpRCd2FnSCMN86ARjPdXy8bo",
    "CreatedAt": "0001-01-01T00:00:00Z",
    "UpdatedAt": "0001-01-01T00:00:00Z"
  }
}
```

