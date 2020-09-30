# seckillSystem  

**架构分析**
>Proxy层 --- 接收用户的秒杀请求，`初审`通过的请求会被放入`redis请求队列`中；同时从`redis响应队列`中取出秒杀结果，并响应给用户
>Layer层 --- 从`redis请求队列`中取出请求，进行秒杀相关的逻辑处理，最后将秒杀结果放入`redis响应队列`中   
>WebAdmin管理后台  --- 管理秒杀商品和秒杀活动 发布秒杀活动到etcd配置中  

+ mysql 数据库  
- redis 请求和响应队列  商品数量同步更新  
* etcd 商品配置 系统配置  

---

*为了时时同步更新秒杀活动商品的剩余数量，使用了redis hash结构来存储数据；  
Layer层定时的向redis中更新商品的剩余数，Admin后台则定时的从redis中获取最新的商品剩余数，并将其更新到数据库和etcd秒杀配置中  
Proxy层和Layer层也会时时的`watch`etcd中对应的秒杀配置，一旦etcd秒杀配置有修改就会进行对应的处理 (比如Proxy层更新etcd配置后 如果发现秒杀商品的剩余数量已经为0，则会停止接收用户请求，直接进行返回)*

