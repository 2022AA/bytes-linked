# 前端部分

2022AA 前端部分

## 架构

Vite + Vue3 + Typescript + Element Plus + Axios + Vue Router 4 + Pinia状态管理

## 项目启动

```bash
# 安装依赖
$ pnpm install

# 启动服务
$ pnpm dev          # 本地环境开发
$ pnpm build        # 生产环境打包
```

## 目录

```md
├── dist/                          # 构建产物
├── public/
├── src/                           # 源码路径
│   ├── api/                       # 封装请求
│   ├── assets/                    # 素材库
│   ├── components/                # 功能组件
│   ├── router/                    # 路由配置
│   ├── store/                     # 维护状态
│   ├── utils/                     # [可选] 工具库
│   ├── views/                     # 视图页面
│   ├── main.ts                    # 应用主脚本
│   └── App.vue                    # 应用根页面
├── vite.config.ts                 # 工程配置
├── windi.config.ts                # windi配置
├── README.md
│   index.html                     # 应用入口HTML
├── package.json                   # 依赖包文件
├── .env                           # 环境配置
├── .eslintignore
├── .eslintrc.[j,t]s
├── .gitignore
└── [j,t]sconfig.json
```
