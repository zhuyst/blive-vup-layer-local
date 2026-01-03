# 巫女酱子弹幕姬 🎀

<p align="center">
  <img src="build/appicon.png" alt="Logo" width="120" height="120">
</p>

<p align="center">
  <strong>一款为虚拟主播设计的 B 站直播弹幕互动桌面应用</strong>
</p>

<p align="center">
  <a href="#功能特性">功能特性</a> •
  <a href="#快速开始">快速开始</a> •
  <a href="#配置说明">配置说明</a> •
  <a href="#使用指南">使用指南</a> •
  <a href="#开发指南">开发指南</a>
</p>

---

## ✨ 功能特性

### 🎯 核心功能

| 功能模块 | 描述 |
|---------|------|
| **弹幕监听** | 实时接收 B 站直播间弹幕、礼物、SC、大航海等消息 |
| **TTS 语音播报** | 使用阿里云 TTS 自动朗读弹幕、礼物感谢语、欢迎语等 |
| **大模型智能回复** | 接入多家 LLM（百度/智谱/豆包/千问），对弹幕进行智能回复 |
| **语音识别** | 支持阿里云语音识别，实现语音转文字功能 |
| **多窗口展示** | 独立透明窗口展示各类消息，可用于 OBS 窗口捕获 |
| **用户记忆** | 记录粉丝牌等级、舰长信息，实现个性化互动 |

### 🖥️ 独立窗口

应用提供多个独立的透明窗口，方便在直播画面中叠加展示：

- **弹幕窗口** - 实时弹幕流展示
- **礼物窗口** - 礼物赠送记录
- **进入直播间** - 观众入场提示
- **关注直播间** - 新增关注提示
- **醒目留言(SC)** - Super Chat 展示
- **大航海** - 舰长/提督/总督 上船提示
- **大模型回复** - AI 智能回复展示

### ⌨️ 全局快捷键

| 快捷键 | 功能 |
|-------|------|
| `Ctrl + Alt + F6` | 录音状态切换 |
| `Ctrl + Alt + F7` | 开始录音 |
| `Ctrl + Alt + F8` | 停止录音 |

### 🤖 支持的大模型

| 厂商 | 模型 |
|-----|------|
| 百度 | ERNIE 4.5 / DeepSeek |
| 智谱 | GLM-4 |
| 字节豆包 | Doubao Seed |
| 阿里千问 | Qwen Plus |

---

## 🚀 快速开始

### 环境要求

- **Go** 1.24+
- **Node.js** 18+
- **Wails** v3 (alpha)

### 安装 Wails

```bash
go install github.com/wailsapp/wails/v3/cmd/wails3@latest
```

### 克隆项目

```bash
git clone https://github.com/your-username/blive-vup-layer-local.git
cd blive-vup-layer-local
```

### 安装依赖

```bash
# 安装前端依赖
cd frontend
npm install
cd ..
```

### 配置文件

复制配置模板并填写必要信息：

```bash
cp etc/config.toml etc/config-dev.toml
```

编辑 `etc/config-dev.toml`，填写以下配置：

- B 站开放平台凭证
- 阿里云 TTS 密钥
- LLM API Key

详细配置请参考 [配置说明](#配置说明)。

### 开发模式运行

```bash
wails3 dev
```

### 生产构建

```bash
# Windows
PRODUCTION=true wails3 build
```

构建产物位于 `bin/` 目录。

---

## ⚙️ 配置说明

配置文件位于 `etc/config.toml`（生产环境）或 `etc/config-dev.toml`（开发环境）。

### 数据库配置

```toml
db_path = "./data/blive-vup-layer.db"
```

### B 站直播开放平台

前往 [B 站直播开放平台](https://open-live.bilibili.com/) 创建应用，获取以下凭证：

```toml
[bilibili]
access_key = "你的 Access Key"
secret_key = "你的 Secret Key"
app_id = 123456  # 应用 ID
```

### 阿里云 TTS（语音合成）

前往 [阿里云智能语音交互](https://nls-portal.console.aliyun.com/) 开通服务：

```toml
[aliyun_tts]
access_key = "阿里云 AccessKey ID"
secret_key = "阿里云 AccessKey Secret"
app_key = "智能语音应用 AppKey"
```

### 阿里云语音识别

```toml
[speech_recognition]
access_key = "阿里云 AccessKey ID"
secret_key = "阿里云 AccessKey Secret"
app_key = "智能语音应用 AppKey"
```

### 大模型配置

根据需要配置一个或多个大模型：

#### 百度（ERNIE / DeepSeek）

```toml
[llm.model.baidu]
base_url = "https://qianfan.baidubce.com/v2/ai_search/"
api_key = "你的 API Key"
ernie_model = "ernie-4.5-turbo-128k"
deepseek_model = "deepseek-v3.2"
```

#### 智谱 GLM

```toml
[llm.model.glm]
base_url = "https://open.bigmodel.cn/api/paas/v4/"
api_key = "你的 API Key"
glm_model = "glm-4.7"
```

#### 字节豆包

```toml
[llm.model.doubao]
base_url = "https://ark.cn-beijing.volces.com/api/v3"
api_key = "你的 API Key"
doubao_model = "doubao-seed-1-6-251015"
```

#### 阿里千问

```toml
[llm.model.qwen]
base_url = "https://dashscope.aliyuncs.com/compatible-mode/v1"
api_key = "你的 API Key"
qwen_model = "qwen-plus-latest"
```

### 自定义 LLM Prompt

可以在配置文件中自定义大模型的系统提示词：

```toml
[llm]
prompt = """
# 你的自定义 Prompt
...
"""
```

---

## 📖 使用指南

### 1. 获取身份码

1. 打开 B 站直播间
2. 在直播控制台中找到「互动玩法」-「身份码」
3. 复制身份码

### 2. 连接直播间

1. 启动应用，打开总控台
2. 在输入框中粘贴身份码
3. 点击「连接」按钮

### 3. 配置互动选项

在总控台中可以配置：

- **TTS 开关** - 是否启用语音播报
- **大模型回复** - 是否启用 AI 智能回复
- **欢迎限制** - 是否限制欢迎语播报（仅高等级粉丝）
- **空闲 TTS** - 是否启用直播间空闲时的随机语音
- **模型选择** - 选择使用的大模型

### 4. 窗口管理

- 通过系统托盘菜单显示/隐藏各个窗口
- 子窗口为无边框透明窗口，可直接用于 OBS 窗口捕获
- 可通过托盘菜单「显示所有窗口」快速恢复

### 5. 智能回复规则

大模型回复触发条件：

| 条件 | 效果 |
|-----|------|
| 佩戴主播粉丝牌且等级 ≥ 10 | 概率触发回复 |
| 舰长/提督/总督 | 强制触发回复 |
| 发送 Super Chat | 强制触发回复 |

为避免刷屏，系统会根据当前弹幕密度动态调整回复概率。

---

## 🛠️ 开发指南

### 项目结构

```
blive-vup-layer-local/
├── build/                # 构建配置和资源
├── config/               # Go 配置定义
├── dao/                  # 数据访问层（SQLite）
├── etc/                  # 配置文件目录
├── frontend/             # Vue 3 前端项目
│   ├── src/
│   │   ├── components/   # Vue 组件
│   │   ├── view/         # 页面视图
│   │   └── ...
│   └── package.json
├── llm/                  # 大模型接入层
├── speechrecognition/    # 语音识别模块
├── tts/                  # TTS 语音合成模块
├── util/                 # 工具函数
├── main.go               # 应用入口
├── service.go            # 核心业务逻辑
└── go.mod
```

### 技术栈

| 层级 | 技术 |
|-----|------|
| 桌面框架 | [Wails v3](https://wails.io/) |
| 后端语言 | Go 1.24 |
| 前端框架 | Vue 3 + Vite |
| 状态管理 | Pinia |
| 数据库 | SQLite (GORM) |
| B 站 API | [bianka](https://github.com/vtb-link/bianka) |

### 开发命令

```bash
# 开发模式（热重载）
wails3 dev

# 构建生产版本
PRODUCTION=true wails3 build

# 仅构建前端
cd frontend && npm run build
```

### 添加新的 LLM 支持

1. 在 `llm/` 目录下创建新的模型文件（参考 `qwen.go`）
2. 在 `config/define.go` 中添加配置结构
3. 在 `llm/llm.go` 中注册新模型

---

## 📝 更新日志

### v1.0.0

- ✅ B 站直播间弹幕监听
- ✅ 阿里云 TTS 语音播报
- ✅ 多大模型智能回复支持
- ✅ 多窗口透明展示
- ✅ 全局快捷键支持
- ✅ 用户数据持久化

---

## 📄 开源协议

本项目仅供学习交流使用。