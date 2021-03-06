# API
## 名词约定
client：喧喧客户端
xxd：GO 聊天服务器
rzs：后台然之服务器

## 数据格式和约定
常见的请求对象格式
```json
{
    userID, // 用户id，xxd -> rzs 非登录时必须
    module, // 模块名称,必须
    method, // 方法名称,必须
    test,   // 可选参数，bool,默认为false。
    params, // 参数对象,可选
    data    // 请求数据,可选,与params配合使用,通常data传输是对象
}
```

常见的响应数据格式
```json
{
    module,            // 模块名称,必须
    method,            // 方法名称,必须
    users[],           // 该数据响应给哪些用户，users为空表示所有在线用户
    params,            // 参数对象,可选
    result: 'success', // 响应状态,可为'success'（成功）, 'fail'(失败), 'denied'(拒绝,需要登录),
    message: '',       // 消息,可选,当result不为success时,使用此字段来解释原因
    data               // 数据 
}
```

## 登录
### 请求  
#### 方向：client --> xxd
```json
{
    module: 'chat',
    method: 'login',
    params: [
		serverName, //多然之时客户端登录的服务器名称
        account,
        password,   // 已加密
        status      // 登录后设置的状态,包括online,away,busy
    ]
 }
```
#### 方向：xxd --> rzs
xxd服务器根据module、method和serverName把请求发送给指定的rzs

### 响应
#### 方向：rzs --> xxd
```json
{
    module: 'chat',
    method: 'login',
    result,
    users[]，
    data: 
    {             // 当前登录的用户数据
        id,
        account,  // 用户名
        realname, // 真实姓名
        avatar,   // 头像URL
        role,     // 角色
        dept,     // 部门ID
        status    // 当前状态
    }
}
```
登录成功以后xxd主动从rzs服务器获取用户列表和用户所参与的会话信息发送送给当前客户端。最后把rzs服务器响应给xxd服务器的登录信息去掉users字段后，发送给此会话包含的所有在线用户。

## 登出
### 请求
#### 方向：client --> xxd
```json
{
    userID, //登出用户的id号
    module: 'chat',
    method: 'logout',
}
```

#### 方向：xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
#### 方向 rzs --> xxd
```json
{
    module: 'chat',
    method: 'logout',
    result,
    users[],
    data:
    {
        id,
        account,
        realname,
        avatar,
        role,
        dept,
        status,
    }
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的登出信息去掉users字段后，发送给此会话包含的所有在线用户。

## 重复登录
>当同一用户重复登录时,系统会向前一个登录的用户推送一条特殊的消息,客户端接收到该消息后应该将用户登出，并关闭相关的网络连接。该消息不需要响应或返回结果。

#### 方向：xxd --> client
```json
{
    module:  'null',
    method:  'null',
    message: 'This account logined in another place.'
}
```

## 获取所有用户列表
### 请求
#### 方向： client --> xxd
```json
{
    userID, //用户的id号
    module: 'chat',
    method: 'userGetlist'
}
```

#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
#### 方向：rzs --> xxd
```json
{
    module: 'chat',
    method: 'userGetlist',
    result,
    users[],
    data:
    [           // 所有用户状态数组
        {       // 其中一个用户数据
            id,
            account,
            realname,
            avatar,
            role,
            dept
            status
        },
        // 更多用户数据...
    ]
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。

## 获取当前登录用户所有会话数据
### 请求
#### 方向：client --> xxd
```json
{
    userID,
    module: 'chat',
    method: 'getList',
}
```
#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
#### 方向：rzs --> xxd
```json
{
    module: 'chat',
    method: 'getList',
    result,
    users[],
    data:
    [                       // 所有会话信息数组
        {                   // 其中一个会话信息
            id,             // 会话在服务器数据保存的id
            gid,            // 会话的全局id,
            name,           // 会话的名称
            type,           // 会话的类型
            admins,         // 会话允许发言的用户列表
            subject,        // 主题会话的关联主题ID
            public,         // 是否公共会话
            createdBy,      // 创建者用户名
            createdDate,    // 创建时间
            editedBy,       // 编辑者用户名
            editedDate,     // 编辑时间
            lastActiveTime, // 会话最后一次发送消息的时间
            star,           // 当前登录用户是否收藏此会话
            hide,           // 当前登录用户是否隐藏此会话
            members: 
            [               // 当前会话中包含的所有用户信息,只需要包含id即可
                {
                    id,     //用户id
                },
                // 更多用户...
            ],
        },
        // 更多会话数据...
    ]
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。

## 更改当前登录用户的信息
### 请求
#### 方向：client --> xxd
```json
{
    userID,
    module: 'chat',
    method: 'userChange',
    params:
    {          //更改后的信息
        name,
        account,
        realname,
        avatar,
        role,
        dept,
        status //要设置的新状态,包括online, away, busy
    }
}
```
#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
#### 方向：rzs --> xxd
```json
{
    module: 'chat',
    method: 'userChange',
    result,
    users[],
    data: 
    {          //当前登录用户数据
        id,
        name,  // 更改后的姓名
        account,
        realname,
        avatar,
        role,
        dept,
        status //要设置的新状态,包括online, away, busy
    }
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。

## 创建聊天会话
### 请求
#### 方向：client --> xxd
```json
{
    userID,
    module: 'chat',
    method: 'create',
    params:
    {
        gid,     // 会话的全局id,
        name,    // 会话的名称
        type,    // 会话的类型
        members: [{id}, {id}...] // 会话的成员列表 
        subject, //可选,主题会话的关联主题ID,默认为0
        pulic    //可选,是否公共会话,默认为false
    }
}
```
#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
>服务器在创建会话时应该先检查gid是否已经存在，如果存在则直接为当前登录用户返回已存在的会话信息。

#### 方向：rzs --> xxd
```json
{
    module: 'chat',
    method: 'create',
    result,
    users[],
    data: // 新创建的会话完整信息
    {
        id,             // 会话在服务器数据保存的id
        gid,            // 会话的全局id,
        name,           // 会话的名称
        type,           // 会话的类型
        admins,         // 会话允许发言的用户列表
        subject,        // 主题会话的关联主题ID
        public,         // 是否公共会话
        createdBy,      // 创建者用户名
        createdDate,    // 创建时间
        editedBy,       // 编辑者用户名
        editedDate,     // 编辑时间
        lastActiveTime, // 会话最后一次发送消息的时间
        members: [{id}, {id}...] // 会话的成员列表 
    }
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。

## 加入或退出聊天会话
>用户可以加入类型为group并且公共的会话；用户可以退出类型为group的会话。

### 请求
#### 方向：client --> xxd
```json
{
    userID,
    module: 'chat',
    method: 'joinchat',
    params: 
    {
        gid, // 要加入或退出的会话id
        join // 可选, true加入会话, false退出会话, 默认为true
    }
}
```
#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
#### 方向：rzs --> xxd
```json
{
    module: 'chat',
    method: 'joinchat',
    result,
    users[],
    data:
    {                   // 会话的完整信息
        id,             // 会话在服务器数据保存的id
        gid,            // 会话的全局id,
        name,           // 会话的名称
        type,           // 会话的类型
        admins,         // 会话允许发言的用户列表
        subject,        // 主题会话的关联主题ID
        public,         // 是否公共会话
        createdBy,      // 创建者用户名
        createdDate,    // 创建时间
        editedBy,       // 编辑者用户名
        editedDate,     // 编辑时间
        lastActiveTime, // 会话最后一次发送消息的时间
        members: [{id}, {id}...] // 会话的成员列表 
    }
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。

## 更改会话名称
>用户可以更改类型为group的会话的名称。

### 请求
#### 方向：client --> xxd
```json
{
    userID,
    module: 'chat',
    method: 'changename',
    params:
    {
        id,  // 要更改的会话id
        name // 新的名称
    }
}
```
#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
#### 方向：rzs --> xxd
```json
{
    module: 'chat',
    method: 'changename',
    result,
    users[],
    data:
    {                   // 会话的完整信息
        id,             // 会话在服务器数据保存的id
        gid,            // 会话的全局id,
        name,           // 会话的名称
        type,           // 会话的类型
        admins,         // 会话允许发言的用户列表
        subject,        // 主题会话的关联主题ID
        public,         // 是否公共会话
        createdBy,      // 创建者用户名
        createdDate,    // 创建时间
        editedBy,       // 编辑者用户名
        editedDate,     // 编辑时间
        lastActiveTime, // 会话最后一次发送消息的时间
        members: [{id}, {id}...] // 会话的成员列表 
    }
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。

## 收藏或取消收藏会话
>每个用户都可以单独决定收藏或取消收藏会话（加星标记）。

### 请求
#### 方向：client --> xxd
```json
{
    userID,
    module: 'chat',
    method: 'star',
    params: 
    {
        gid, // 要收藏会话id
        star // 可选, true收藏会话, false取消收藏会话, 默认为true
    }
}
```
#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
#### 方向：rzs --> xxd
```json
{
    module: 'chat',
    method: 'star',
    result,
    users[],
    data:
    {
        gid // 要收藏或取消收藏的会话id
    }
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。

## 邀请新的用户到会话或者将用户踢出会话
### 请求
>用户可以邀请一个或多个用户到类型为group的已有会话中；会话管理员可以将一个或多个用户踢出类型为group的会话。

#### 方向：client --> xxd
```json
{
    userID,
    module: 'chat',
    method: 'addmember',
    params: 
    {
        gid,     // 要操作的会话id
        members, // 用户id数组
        join     // 可选, true邀请用户加入会话, false将用户踢出会话, 默认为true
    }
}
```
#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
>当新用户被添加到会话之后或者用户被踢出会话后,服务器应该主动推送此会话的信息给此会话的所有在线成员；此响应与chat/create/响应的结果一致。

#### 方向：rzs --> xxd
```json
{
    module: 'chat',
    method: 'addmember',
    result,
    users[],
    data: // 会话的完整信息
    {
        id,             // 会话在服务器数据保存的id
        gid,            // 会话的全局id,
        name,           // 会话的名称
        type,           // 会话的类型
        admins,         // 会话允许发言的用户列表
        subject,        // 主题会话的关联主题ID
        public,         // 是否公共会话
        createdBy,      // 创建者用户名
        createdDate,    // 创建时间
        editedBy,       // 编辑者用户名
        editedDate,     // 编辑时间
        lastActiveTime, // 会话最后一次发送消息的时间
        members: [{id}, {id}...] // 会话的成员列表 
    }
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。

## 向会话发送消息
### 请求
>用户向一个或多个会话中发送一条或多条消息,服务器推送此消息给此会话中的所有在线成员；当前不在线的成员会在下次上线时通过离线消息送达。

#### 方向：client --> xxd
```json
{
    userID,
    module: 'chat',
    method: 'message',
    params: 
    [        // 一个包含一条或多条新消息的数组
        {
            gid,         // 此消息的gid
            cgid,        // 此消息关联的会话的gid
            user,        // 如果为空,则为发送此请求的用户
            date,        // 如果为空,则已服务器处理时间为准
            type,        // 消息的类型
            contentType, // 消息内容的类型
            content      // 消息内容
        },
        // 可以在一个请求中发送多个消息
    ]
}
```
#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
>当有新的消息收到时,服务器会所有消息,并发送给对应会话的所有在线成员

#### 方向：rzs --> xxd
```json
{
    module: 'chat',
    method: 'message',
    result,
    users[],
    data:  // 一个包含一条或多条新消息的数组
    [
        {
            id,          // 消息在服务器保存的id
            gid,         // 此消息的gid
            cgid,        // 此消息关联的会话的gid
            user,        // 消息发送的用户名
            date,        // 消息发送的时间
            type,        // 消息的类型
            contentType, // 消息内容的类型
            content,     // 消息内容
        },
        // 可以有更多的消息
    ]
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。

## 获取会话的所有消息记录
### 请求
#### 方向：client --> xxd
```json
{
    userID,
    module: 'chat',
    method: 'history',
    params: 
    {
        gid // 要获取消息记录的会话gid
    }
}
```
#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
#### 方向：rzs --> xxd
```json
{
    module: 'chat',
    method: 'history',
    result,
    users[],
    data: 
    [
        {
        // 一条历史消息...
        },
        // 更多历史消息
    ]
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。

## 获取会话的所有成员信息
### 请求
#### 方向：client --> xxd
```json
{
    userID,
    module: 'chat',
    method: 'members',
    params:
    {
        gid // 要获取成员列表的会话gid
    }
}
```
#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
#### 方向：rzs --> xxd
```json
{
    module: 'chat',
    method: 'members',
    result,
    users[],
    data: 
    [
        {
            gid, // 此消息的gid
            members: 
            [
                {
                    id,
                    // ...
                } 
            ]
        },
    ]
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。

## 隐藏或显示会话
>每个用户都可以单独决定隐藏或显示已参与的会话。

### 请求
#### 方向：client --> xxd
```json
{
    userID,
    module: 'chat',
    method: 'hide',
    params: 
    {
        gid, // 要收藏会话id
        hide // 可选, true隐藏会话, false显示会话, 默认为true
    }
}
```
#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
#### 方向：rzs --> xxd
```json
{
    module: 'chat',
    method: 'hide',
    result,
    users[],
    data:
    {
        gid // 要隐藏或显示的会话id
    }
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。

## 将会话设置为公共会话或者取消设置公共会话
>用户可以将一个非主题会话设置为公共会话或者取消设置公共会话。

### 请求
#### 方向：client --> xxd
```json
{
    userID,
    module: 'chat',
    method: 'changepublic',
    params: 
    {
        gid,
        public, // 可选,true设置公共会话,false取消设置公共会话,默认为true
    }
}
```
#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
#### 方向：rzs --> xxd
```json
{
    module: ' chat';
    method: 'changepublic',
    result,
    users[],
    data: 
    {           // 会话的完整信息
        id,             // 会话在服务器数据保存的id
        gid,            // 会话的全局id,
        name,           // 会话的名称
        type,           // 会话的类型
        admins,         // 会话允许发言的用户列表
        subject,        // 主题会话的关联主题ID
        public,         // 是否公共会话
        createdBy,      // 创建者用户名
        createdDate,    // 创建时间
        editedBy,       // 编辑者用户名
        editedDate,     // 编辑时间
        lastActiveTime, // 会话最后一次发送消息的时间
        members: [{id}, {id}...] // 会话的成员列表 
    }
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。

## 获取所有公共会话列表
### 请求
#### 方向：client --> xxd
```json
{
    userID,
    module: 'chat',
    method: 'getPublicList'
}
```
#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
#### 方向：rzs --> xxd
```json
{
    module: 'chat',
    method: 'getPublicList',
    result,
    users[],
    data: 
    [         // 所有公共会话信息数组
        {     // 其中一个会话信息
            id,
            gid,
            name,
            type,
            admins, 
            subject,
            public,
            createdBy,
            createdDate,
            editedBy,
            editedDate,
            lastActiveTime
        },
        // 更多会话数据...
    ]
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。

## 获取指定会话的所有文件列表
### 请求
#### 方向：client --> xxd
```json
{
    userID,
    module: 'attach',
    method: 'getList',
    params: 
    {
        gid // 要获取文件列表的会话的全局id
    }
}
```
#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
#### 方向：rzs --> xxd
```json
{
    module: 'attach',
    method: 'getList',
    result,
    users[],
    data: 
    [      // 所有文件列表
        {
            id,      // 文件id, 下载文件需要此id值
            title,   // 文件标题
            size,    // 文件大小
            addedBy, // 上传人
        },
        /// 更多文件
    ]
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。

## 设置会话管理员
### 请求
#### 方向：client --> xxd
```json
{
    userID,
    module: 'chat',
    method: 'setAdmin',
    params: 
    {
        gid,  
        admins: [{id},{id}...], // 指定的用户列表
        isAdmin, //可选, true允许指定用户发言, false禁止指定用户发言, 默认为true 
    }
}
```
#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
#### 方向：rzs --> xxd
```json
{
    module: 'chat',
    method: 'setAdmin',
    result,
    users[],
    data: 
    {
        id,
        gid,
        name,
        type,
        admins,
        subject,
        public,
        createdBy,
        createdDate,
        editedBy,
        editedDate
        lastActiveTime
    }
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。

## 设置会话允许发言的人
>通过此功能可以设置会话白名单。

### 请求
#### 方向：client --> xxd
```json
{
    userID,
    module: 'chat',
    method: 'setCommitters',
    params: 
    {
        gid,  
        committers: [{id},{id}...] // 指定的用户列表
    }
}
```
#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
#### 方向：rzs --> xxd
```json
{
    module: 'chat',
    method: 'setCommitters',
    result,
    users[],
    data: 
    {
        id,
        gid,
        name,
        type,
        admins,
        subject,
        public,
        createdBy,
        createdDate,
        editedBy,
        editedDate
        lastActiveTime
    }
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。

## 上传下载用户在客户端的配置信息
### 请求
#### 方向：client --> xxd
```json
{
    userID,
    module: 'chat',
    method: 'settings',
    params: 
    {
        account, //用户名
        settings //用户配置信息, 可选, 为空表示下载用户配置信息, 不为空表示上传用户配置信息, 默认为空 
    }
}
```
#### 方向： xxd --> rzs
xxd把client发送的数据转发给rzs。

### 响应
#### 方向：rzs --> xxd
```json
{
    module: 'chat',
    method: 'settings',
    users[],
    result, 
    data // 用户配置信息
}
```
#### 方向：xxd --> client
把rzs服务器响应给xxd服务器的信息去掉users字段后，发送给此会话包含的所有在线用户。
