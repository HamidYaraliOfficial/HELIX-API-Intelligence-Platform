# HELIX — API Intelligence Platform

**A self-hosted API Catalog & Discovery platform — Go backend, Next.js control center, PostgreSQL storage.**

<p align="center">
  <img alt="Go" src="https://img.shields.io/badge/backend-Go%201.22-00ADD8">
  <img alt="Next.js" src="https://img.shields.io/badge/frontend-Next.js%2014-000000">
  <img alt="PostgreSQL" src="https://img.shields.io/badge/storage-PostgreSQL%2016-336791">
  <img alt="License" src="https://img.shields.io/badge/license-MIT-informational">
</p>

---

## 🇬🇧 English

### Overview

HELIX is a working starter platform for **API discovery and cataloging**. Point it at an OpenAPI/Swagger document — upload a file or give it a URL — and HELIX parses the contract, registers the API, and stores every endpoint, parameter and response it finds. From there you get a searchable, filterable catalog of every API your organization owns: who owns it, what environment it runs in, its lifecycle state, and a live "operating hours" badge showing whether it's currently within its supported window.

This is the real, functioning core of a larger API Intelligence Platform vision — built to be extended with dependency graphs, traffic intelligence, breaking-change detection, and the rest of the HELIX roadmap.

### Features

- **API Discovery Engine** — upload an OpenAPI/Swagger 2.0 or 3.x document (JSON or YAML), or import directly from a URL. HELIX parses it into a normalized contract and stores every endpoint, parameter, and response.
- **API Catalog** — full-text search, plus filters by protocol, status, environment, lifecycle state, and team.
- **Endpoint Explorer** — expand any endpoint to see its parameters and possible responses.
- **Manual Registration** — register APIs by hand for services that don't (yet) have a spec.
- **Team Management** — group APIs by owning team, with live API counts.
- **Operating Hours / Availability** — set the days and hours an API is expected to be supported, in any IANA timezone. The catalog shows a **live badge**: "Open now, closes in 2h 15m" or "Closed, opens in 6h 40m" — computed client-side and updated every minute.
- **Dashboard** — total APIs, endpoints, teams, and active-API counts at a glance.
- **5 Fluent/Windows-11-inspired themes** — Light, Dark, AMOLED, Windows Blue, Windows Red — switchable at any time, persisted per browser.
- **3 languages, correct text direction** — English (LTR), Persian (RTL), Chinese (LTR), switchable at any time, persisted per browser.
- **Multi-tenant ready** — every table is organization-scoped; a default organization is seeded automatically.

### Tech stack

| Layer     | Technology                                                            |
|-----------|------------------------------------------------------------------------|
| Backend   | Go 1.22 (standard-library `net/http` router, no framework), `lib/pq`, `gopkg.in/yaml.v3` |
| Database  | PostgreSQL 16                                                          |
| Frontend  | Next.js 14 (App Router), React 18, TypeScript, Tailwind CSS            |
| Packaging | Docker, Docker Compose                                                 |

### Project structure

```
helix/
├── backend/                 # Go REST API
│   ├── cmd/server/          # main.go — entrypoint
│   ├── internal/
│   │   ├── config/          # environment configuration
│   │   ├── db/               # PostgreSQL connection + embedded migrations
│   │   ├── models/           # domain models & enums
│   │   ├── openapi/          # OpenAPI/Swagger 2.0 & 3.x parser
│   │   ├── repository/       # SQL data-access layer
│   │   ├── handlers/         # HTTP handlers + router
│   │   └── middleware/       # CORS, logging, panic recovery
│   └── Dockerfile
├── frontend/                 # Next.js control center
│   └── src/
│       ├── app/               # pages: dashboard, API detail, new/edit, discovery, teams
│       ├── components/        # UI components (table, forms, badges, switchers...)
│       ├── lib/                # api-client, i18n, theme, availability engine
│       ├── locales/            # en.json, fa.json, zh.json
│       └── styles/             # themes.css (5 themes as CSS variables)
├── docker-compose.yml
└── .env.example
```

### Installation

#### Option A — Docker Compose (recommended)

**Prerequisites:** Docker and Docker Compose installed.

1. Extract this project's files into a folder.
2. From that folder, run:
   ```
   docker compose up --build
   ```
3. Wait for all three containers (`postgres`, `backend`, `frontend`) to report healthy. First build takes a few minutes while Go modules and npm packages download.
4. Open the app:
   - Frontend (Control Center): **http://localhost:3000**
   - Backend (REST API): **http://localhost:8080**
   - Health check: **http://localhost:8080/health**

Database schema and a default organization/team are created automatically on first boot — no manual setup step required.

#### Option B — Manual / local development

**Prerequisites:** Go 1.22+, Node.js 20+, PostgreSQL 16 running locally.

1. Create a database and user matching `.env.example`, or edit `HELIX_DATABASE_URL` to match your own instance.
2. **Backend:**
   ```
   cd backend
   go mod tidy
   go run ./cmd/server
   ```
   The server listens on `:8080` and applies migrations automatically on startup.
3. **Frontend:**
   ```
   cd frontend
   npm install
   npm run dev
   ```
   The dev server listens on `:3000` and talks to the backend at `http://localhost:8080` by default (override with `NEXT_PUBLIC_API_BASE_URL`).

### Using Discovery

1. Open **Discovery** in the sidebar.
2. Either drag & drop an OpenAPI/Swagger file (`.json`, `.yaml`, `.yml`), or switch to the **From URL** tab and paste a link to a live spec (e.g. `https://api.example.com/openapi.json`).
3. Pick a target environment and (optionally) a team.
4. Click **Run discovery**. HELIX parses the contract, creates the API entry, imports every endpoint it found, and links straight to the new catalog page.

### Setting operating hours

When registering or editing an API, the **Operating hours** section lets you set:
- an opening and closing time,
- the days of the week it applies to,
- an IANA timezone (e.g. `Asia/Tehran`, `Europe/London`, `America/New_York`).

The catalog and API detail page then show a live badge computed in the visitor's browser — no server polling required — telling you whether the API is currently within its supported window, and how long until that changes.

### Environment variables

See `.env.example` for the full list. The most important ones:

| Variable                    | Where    | Default                                                     |
|------------------------------|----------|---------------------------------------------------------------|
| `HELIX_DATABASE_URL`         | Backend  | `postgres://helix:helix@localhost:5432/helix?sslmode=disable` |
| `HELIX_PORT`                  | Backend  | `8080`                                                        |
| `HELIX_CORS_ORIGIN`           | Backend  | `*`                                                           |
| `NEXT_PUBLIC_API_BASE_URL`    | Frontend | `http://localhost:8080`                                       |

### License

MIT — use it, extend it, ship it.

---

## 🇮🇷 فارسی

### معرفی

هلیکس (HELIX) یک پلتفرم واقعی و کارکردی برای **کشف و کاتالوگ‌سازی API** است. کافیست یک فایل OpenAPI/Swagger را آپلود کنید یا آدرس آن را بدهید؛ هلیکس آن را تجزیه می‌کند، API را ثبت می‌کند و تمام Endpointها، پارامترها و پاسخ‌های آن را ذخیره می‌سازد. نتیجه، یک کاتالوگ قابل جستجو و فیلتر از تمام APIهای سازمان شماست: مالک هرکدام، محیط اجرا، وضعیت چرخه عمر، و یک نشان زنده‌ی «ساعات کاری» که نشان می‌دهد آیا API در حال حاضر در بازه‌ی پشتیبانی‌شده قرار دارد یا نه.

این نسخه، هسته‌ی واقعی و کاملاً کارکردی از یک چشم‌انداز بزرگ‌تر برای پلتفرم هوشمندی API است؛ پایه‌ای که می‌توان با Dependency Graph، Traffic Intelligence، تشخیص Breaking Change و بقیه‌ی نقشه‌راه HELIX گسترشش داد.

### ویژگی‌ها

- **موتور کشف API** — آپلود فایل OpenAPI/Swagger نسخه‌ی 2.0 یا 3.x (با فرمت JSON یا YAML)، یا وارد کردن مستقیم از یک آدرس اینترنتی. هلیکس آن را به یک قرارداد استاندارد تبدیل کرده و تمام Endpointها، پارامترها و پاسخ‌ها را ذخیره می‌کند.
- **کاتالوگ API** — جستجوی متنی کامل، به‌همراه فیلتر بر اساس پروتکل، وضعیت، محیط، چرخه عمر و تیم.
- **کاوشگر Endpoint** — با کلیک روی هر Endpoint، پارامترها و پاسخ‌های احتمالی آن را ببینید.
- **ثبت دستی** — برای سرویس‌هایی که هنوز مشخصات رسمی ندارند، امکان ثبت دستی API وجود دارد.
- **مدیریت تیم‌ها** — گروه‌بندی APIها بر اساس تیم مالک، همراه با شمارش زنده‌ی تعداد APIها.
- **ساعات کاری / در دسترس بودن** — روزها و ساعاتی که انتظار می‌رود یک API پشتیبانی شود را در هر منطقه‌ی زمانی IANA تنظیم کنید. کاتالوگ یک **نشان زنده** نمایش می‌دهد: «اکنون باز است، تا ۲ ساعت و ۱۵ دقیقه دیگر بسته می‌شود» یا «بسته است، تا ۶ ساعت و ۴۰ دقیقه دیگر باز می‌شود» — که در مرورگر محاسبه و هر دقیقه به‌روزرسانی می‌شود.
- **داشبورد** — نمایش سریع مجموع APIها، Endpointها، تیم‌ها و APIهای فعال.
- **۵ پوسته با الهام از فلوئنت/ویندوز ۱۱** — روشن، تاریک، امولد، آبی ویندوز، قرمز ویندوز — هر زمان قابل تغییر و در مرورگر ذخیره می‌شود.
- **۳ زبان با جهت متن صحیح** — انگلیسی (چپ‌به‌راست)، فارسی (راست‌به‌چپ)، چینی (چپ‌به‌راست) — هر زمان قابل تغییر و در مرورگر ذخیره می‌شود.
- **آماده برای چند‌سازمانی (Multi-Tenant)** — تمام جدول‌ها بر اساس سازمان ایزوله هستند؛ یک سازمان پیش‌فرض به‌طور خودکار ساخته می‌شود.

### پشته‌ی فناوری

| لایه        | فناوری                                                                 |
|-------------|--------------------------------------------------------------------------|
| بک‌اند       | Go 1.22 (روتر استاندارد `net/http`، بدون فریمورک اضافه)، `lib/pq`، `gopkg.in/yaml.v3` |
| پایگاه‌داده  | PostgreSQL 16                                                             |
| فرانت‌اند    | Next.js 14 (App Router)، React 18، TypeScript، Tailwind CSS               |
| بسته‌بندی    | Docker، Docker Compose                                                   |

### نصب و راه‌اندازی

#### روش الف — Docker Compose (پیشنهادی)

**پیش‌نیاز:** نصب بودن Docker و Docker Compose.

۱. فایل‌های این پروژه را در یک پوشه استخراج کنید.
۲. از همان پوشه، این دستور را اجرا کنید:
   ```
   docker compose up --build
   ```
۳. صبر کنید تا هر سه کانتینر (`postgres`، `backend`، `frontend`) سالم (healthy) اعلام شوند. ساخت اولیه به‌دلیل دانلود ماژول‌های Go و بسته‌های npm ممکن است چند دقیقه طول بکشد.
۴. برنامه را باز کنید:
   - فرانت‌اند (پنل مدیریت): **http://localhost:3000**
   - بک‌اند (REST API): **http://localhost:8080**
   - بررسی سلامت سرویس: **http://localhost:8080/health**

طرح پایگاه‌داده و یک سازمان/تیم پیش‌فرض به‌طور خودکار در اولین اجرا ساخته می‌شوند و نیازی به تنظیم دستی نیست.

#### روش ب — راه‌اندازی محلی / دستی

**پیش‌نیاز:** Go نسخه‌ی ۱.۲۲ به بالا، Node.js نسخه‌ی ۲۰ به بالا، و PostgreSQL 16 که به‌صورت محلی در حال اجراست.

۱. یک پایگاه‌داده و کاربر متناسب با `.env.example` بسازید، یا مقدار `HELIX_DATABASE_URL` را متناسب با نمونه‌ی خودتان ویرایش کنید.
۲. **بک‌اند:**
   ```
   cd backend
   go mod tidy
   go run ./cmd/server
   ```
   سرور روی پورت `8080` اجرا می‌شود و Migrationها را به‌طور خودکار در زمان راه‌اندازی اعمال می‌کند.
۳. **فرانت‌اند:**
   ```
   cd frontend
   npm install
   npm run dev
   ```
   سرور توسعه روی پورت `3000` اجرا شده و به‌طور پیش‌فرض با بک‌اند در آدرس `http://localhost:8080` صحبت می‌کند (با متغیر `NEXT_PUBLIC_API_BASE_URL` قابل تغییر است).

### استفاده از بخش کشف (Discovery)

۱. گزینه‌ی **کشف** را از نوار کناری باز کنید.
۲. یا یک فایل OpenAPI/Swagger (با پسوند `.json`، `.yaml`، `.yml`) را Drag & Drop کنید، یا به تب **از طریق آدرس** بروید و لینک یک مشخصات زنده را وارد کنید (مثلاً `https://api.example.com/openapi.json`).
۳. یک محیط مقصد و در صورت نیاز یک تیم انتخاب کنید.
۴. روی **اجرای کشف** کلیک کنید. هلیکس قرارداد را تجزیه کرده، رکورد API را می‌سازد، تمام Endpointهای یافت‌شده را وارد می‌کند و مستقیماً به صفحه‌ی کاتالوگ آن API جدید پیوند می‌دهد.

### تنظیم ساعات کاری

هنگام ثبت یا ویرایش یک API، بخش **ساعات کاری** به شما اجازه می‌دهد این موارد را تنظیم کنید:
- ساعت شروع و پایان،
- روزهای هفته‌ای که اعمال می‌شود،
- یک منطقه‌ی زمانی IANA (مثلاً `Asia/Tehran`، `Europe/London`، `America/New_York`).

سپس کاتالوگ و صفحه‌ی جزئیات API یک نشان زنده نمایش می‌دهند که کاملاً در مرورگر کاربر محاسبه می‌شود — بدون نیاز به Polling سرور — و نشان می‌دهد آیا API در حال حاضر در بازه‌ی پشتیبانی‌شده است یا خیر، و چه مدت تا تغییر این وضعیت باقی مانده است.

### متغیرهای محیطی

فهرست کامل در `.env.example` موجود است. مهم‌ترین‌ها:

| متغیر                        | محل      | مقدار پیش‌فرض                                                  |
|------------------------------|----------|-----------------------------------------------------------------|
| `HELIX_DATABASE_URL`         | بک‌اند   | `postgres://helix:helix@localhost:5432/helix?sslmode=disable`  |
| `HELIX_PORT`                  | بک‌اند   | `8080`                                                          |
| `HELIX_CORS_ORIGIN`           | بک‌اند   | `*`                                                             |
| `NEXT_PUBLIC_API_BASE_URL`    | فرانت‌اند | `http://localhost:8080`                                        |

### مجوز

MIT — آزادانه استفاده، توسعه و منتشر کنید.

---

## 🇨🇳 中文

### 概述

HELIX 是一个真实可用的 **API 发现与目录管理** 平台。只需上传一个 OpenAPI/Swagger 文档，或提供其访问地址，HELIX 就会解析该契约、注册对应的 API，并存储其中的每一个接口、参数与响应。由此你将获得一个可搜索、可筛选的组织级 API 目录：每个 API 的负责人、运行环境、生命周期状态，以及一个实时的"运营时间"徽章，用于显示该 API 当前是否处于受支持的时间窗口内。

这是更宏大的 API 智能平台愿景中真实、完整可运行的核心部分，后续可以在此基础上扩展依赖关系图谱、流量智能分析、破坏性变更检测等 HELIX 路线图中的其他能力。

### 功能特性

- **API 发现引擎** — 上传 OpenAPI/Swagger 2.0 或 3.x 文档（JSON 或 YAML 格式），或直接通过网址导入。HELIX 会将其解析为统一的内部契约模型，并存储所有接口、参数和响应。
- **API 目录** — 支持全文搜索，并可按协议、状态、环境、生命周期状态和团队进行筛选。
- **接口浏览器** — 展开任意接口即可查看其参数与可能的响应。
- **手动注册** — 对于尚无规范文档的服务，也可以手动注册其 API 信息。
- **团队管理** — 按所属团队对 API 进行分组，并实时显示每个团队拥有的 API 数量。
- **运营时间 / 可用性** — 可为任意 IANA 时区设置 API 预期受支持的天数与时间段。目录页面会显示一个**实时徽章**：例如"当前开放，2 小时 15 分钟后关闭"或"已关闭，6 小时 40 分钟后开放"——完全在浏览器端计算，每分钟自动更新。
- **仪表盘** — 一目了然地查看 API 总数、接口总数、团队数量与活跃 API 数量。
- **5 种 Fluent / Windows 11 风格主题** — 浅色、深色、纯黑（AMOLED）、Windows 蓝、Windows 红，可随时切换并保存在浏览器中。
- **3 种语言，文字方向正确** — 英语（从左到右）、波斯语（从右到左）、中文（从左到右），可随时切换并保存在浏览器中。
- **原生支持多租户** — 所有数据表均按组织隔离，系统启动时会自动创建一个默认组织。

### 技术栈

| 层级       | 技术                                                                    |
|------------|---------------------------------------------------------------------------|
| 后端       | Go 1.22（使用标准库 `net/http` 路由，不依赖额外框架）、`lib/pq`、`gopkg.in/yaml.v3` |
| 数据库     | PostgreSQL 16                                                             |
| 前端       | Next.js 14（App Router）、React 18、TypeScript、Tailwind CSS               |
| 打包部署   | Docker、Docker Compose                                                    |

### 安装步骤

#### 方式一：Docker Compose（推荐）

**前置条件：** 已安装 Docker 与 Docker Compose。

1. 将本项目的文件解压到一个文件夹中。
2. 在该文件夹下执行：
   ```
   docker compose up --build
   ```
3. 等待 `postgres`、`backend`、`frontend` 三个容器均显示为健康（healthy）状态。首次构建时需要下载 Go 模块和 npm 依赖包，可能需要几分钟时间。
4. 打开应用：
   - 前端（控制中心）：**http://localhost:3000**
   - 后端（REST API）：**http://localhost:8080**
   - 健康检查接口：**http://localhost:8080/health**

数据库结构以及默认组织/团队会在首次启动时自动创建，无需任何手动配置步骤。

#### 方式二：本地手动开发

**前置条件：** Go 1.22 及以上版本、Node.js 20 及以上版本、本地已运行的 PostgreSQL 16。

1. 创建一个与 `.env.example` 相匹配的数据库和用户，或修改 `HELIX_DATABASE_URL` 以匹配你自己的实例。
2. **后端：**
   ```
   cd backend
   go mod tidy
   go run ./cmd/server
   ```
   服务将监听 `8080` 端口，并在启动时自动执行数据库迁移。
3. **前端：**
   ```
   cd frontend
   npm install
   npm run dev
   ```
   开发服务器将监听 `3000` 端口，默认通过 `http://localhost:8080` 与后端通信（可通过 `NEXT_PUBLIC_API_BASE_URL` 覆盖）。

### 使用「发现」功能

1. 在侧边栏打开 **发现** 页面。
2. 可以直接拖放一个 OpenAPI/Swagger 文件（`.json`、`.yaml`、`.yml`），或切换到 **通过网址** 标签页，粘贴一个在线规范文档的地址（例如 `https://api.example.com/openapi.json`）。
3. 选择目标环境，并可选地指定一个团队。
4. 点击 **运行发现**。HELIX 会解析该契约、创建对应的 API 条目、导入发现的所有接口，并直接跳转至新建的目录详情页面。

### 设置运营时间

在注册或编辑某个 API 时，**运营时间** 部分允许你设置：
- 开始与结束时间，
- 适用的星期几，
- 一个 IANA 时区（例如 `Asia/Shanghai`、`Europe/London`、`America/New_York`）。

之后，目录页面与 API 详情页面会显示一个完全在访问者浏览器中计算得出的实时徽章——无需服务器轮询——用于告知该 API 当前是否处于受支持的时间窗口内，以及距离状态变化还有多长时间。

### 环境变量

完整列表请参见 `.env.example`。其中最重要的几项：

| 变量名                        | 所属     | 默认值                                                          |
|--------------------------------|----------|--------------------------------------------------------------------|
| `HELIX_DATABASE_URL`           | 后端     | `postgres://helix:helix@localhost:5432/helix?sslmode=disable`     |
| `HELIX_PORT`                    | 后端     | `8080`                                                             |
| `HELIX_CORS_ORIGIN`             | 后端     | `*`                                                                |
| `NEXT_PUBLIC_API_BASE_URL`      | 前端     | `http://localhost:8080`                                            |

### 许可证

MIT — 可自由使用、扩展与发布。
