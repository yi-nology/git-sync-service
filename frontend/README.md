# Git Sync Service Frontend

基于 Vue 3 + Vite + Element Plus 构建的 Git 代码同步管理系统前端。

## 技术栈

- Vue 3.4
- Vite 5.x
- TypeScript
- Element Plus
- Pinia
- ECharts
- Axios

## 开发

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev

# 构建生产版本
npm run build
```

## 项目结构

```
frontend/
├── src/
│   ├── views/          # 页面组件
│   │   ├── dashboard/ # 仪表盘
│   │   ├── sync/      # 同步任务
│   │   ├── webhook/   # Webhook
│   │   └── settings/  # 系统设置
│   ├── components/     # 通用组件
│   │   └── layout/   # 布局组件
│   ├── router/        # 路由配置
│   ├── stores/        # Pinia 状态管理
│   ├── api/           # API 接口
│   ├── styles/        # 样式文件
│   ├── utils/         # 工具函数
│   └── main.ts      # 入口文件
├── index.html
├── vite.config.ts
├── tsconfig.json
└── package.json
```

## 功能特性

- 📊 数据仪表盘 - 同步状态统计图表
- 🔄 同步任务管理 - 任务创建、编辑、执行
- 📜 同步历史记录 - 查看详细的同步日志
- 🎣 Webhook 规则管理 - 自动触发同步
- 🔐 Git 凭证管理 - 多平台 Token 管理
- ⚙️ 系统设置 - 全局配置管理
