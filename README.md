# مستندات پروژه Reality

## مقدمه

پروژه Reality یک سرویس API جامع است که برای مدیریت کاربران، محصولات، اعلان‌ها و دیگر بخش‌های سیستم طراحی شده است. این مستند به شما نحوه نصب، راه‌اندازی، و اجرای پروژه را با استفاده از `make` و تنظیمات متغیرهای محیطی (`.env`) توضیح می‌دهد.

این پروژه از سرویس‌های خارجی مانند **PayPing** برای پردازش پرداخت و **Liara** برای زیرساخت ابری استفاده می‌کند. این ابزارها به افزایش قابلیت‌های پلتفرم کمک می‌کنند، پردازش پرداخت‌های امن و مقیاس‌پذیری زیرساخت‌ها را فراهم می‌کنند.

---

## فهرست مطالب

1. [پیش‌نیازها](#پیش‌نیازها)
2. [نصب و راه‌اندازی](#نصب-و-راه‌اندازی)
3. [اجرای پروژه](#اجرای-پروژه)
4. [ساختار پروژه](#ساختار-پروژه)
5. [متغیرهای محیطی](#متغیرهای-محیطی)
6. [راهنمایی مشارکت](#راهنمایی-مشارکت)
7. [سرویس‌های خارجی](#سرویس‌های-خارجی)
8. [مجوز](#مجوز)

---

## پیش‌نیازها

برای اجرای این پروژه، به ابزارهای زیر نیاز دارید:

- **Go**: نسخه 1.19 یا بالاتر
- **PostgreSQL**: برای پایگاه داده
- **Docker**: برای اجرای کانتینرها (اختیاری)
- **Make**: برای اجرای دستورات سریع
- **RabbitMQ**: برای پیام‌رسانی داخلی
- **Liara**: برای میزبانی و مقیاس‌پذیری سرویس‌های ابری
- **PayPing**: برای پرداخت‌ها و مدیریت تراکنش‌ها

---

## نصب و راه‌اندازی

### 1. کلون کردن مخزن

ابتدا باید مخزن پروژه را کلون کنید:

```bash
git clone https://github.com/ARmo-BigBang/reality.git
cd reality
```

### 2. تنظیم متغیرهای محیطی

یک فایل `.env` ایجاد کنید و اطلاعات زیر را در آن قرار دهید. این اطلاعات شامل تنظیمات پایگاه داده، پیام‌رسانی، احراز هویت، و دیگر تنظیمات پروژه است:

```env
NAME=Reality
ENVIRONMENT=development
PORT=5454
VERSION=0.0.1
APP_URL=0.0.0.0

MAX_LOGIN_DEVICE_COUNT=3
AUTO_DELETE_DEVICE=true

ACCESS_TOKEN_SIGNING_KEY=your_access_token_signing_key
REFRESH_TOKEN_SIGNING_KEY=your_refresh_token_signing_key
ACCESS_TOKEN_EXPIRATION=9999
REFRESH_TOKEN_EXPIRATION=888888

RABBITMQ_USERNAME=admin
RABBITMQ_PASSWORD=""
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672

DATABASE_PORT=5432
DATABASE_NAME=your_database_name
DATABASE_HOST=localhost
DATABASE_USERNAME=admin
DATABASE_PASSWORD=""
DATABASE_SSL_MODE=disable

SMTP_HOST=your_smtp_host
SMTP_PORT=587
SMTP_USERNAME=your_smtp_username
SMTP_PASSWORD=your_smtp_password
SMTP_FROM=info@armogroup.tech

POSTGRES_PASSWORD="your_postgres_password"

PAYPING_TOKEN=your_payping_token
PAYPING_CALLBACK_TARGET_WEBSITE="http://localhost:8081/dashboard"

AWS_S3_CONNECTION_TARGET=CUSTOM
AWS_S3_PRIVATE_BUCKET_NAME=your_private_bucket_name
AWS_S3_PUBLIC_BUCKET_NAME=your_public_bucket_name
AWS_S3_HOST=http://localhost:9000
AWS_S3_PUBLIC_HOST_URI=bytebase.armogroup.tech
AWS_REGION=us-east-1
AWS_ACCESS_KEY_ID=arash
AWS_SECRET_ACCESS_KEY="your_aws_secret_key"

STATIC_DIR_PATH=./public
```

### 3. نصب وابستگی‌ها

پس از نصب Go و کلون کردن پروژه، وابستگی‌های پروژه را نصب کنید:

```bash
go mod tidy
```

### 4. راه‌اندازی پایگاه داده

در صورتی که از Docker استفاده می‌کنید، می‌توانید PostgreSQL و RabbitMQ را با دستور زیر راه‌اندازی کنید:

```bash
docker-compose up -d
```

---

## اجرای پروژه

### اجرای پروژه با `make`

برای اجرای پروژه می‌توانید از دستورات زیر استفاده کنید:

- **برای محیط توسعه:**

```bash
make dev
```

این دستور پروژه را در حالت توسعه اجرا می‌کند.

- **تازه‌سازی و راه‌اندازی محیط توسعه:**

```bash
make fresh-dev
```

این دستور پایگاه داده را تازه‌سازی کرده و سپس پروژه را اجرا می‌کند.

- **اجرای مهاجرت پایگاه داده در محیط توسعه:**

```bash
make migrate-dev
```

- **برای محیط تولید:**

```bash
make prod
```

پروژه را در حالت تولید اجرا می‌کند.

- **تازه‌سازی و راه‌اندازی محیط تولید:**

```bash
make fresh-prod
```

- **اجرای مهاجرت پایگاه داده در محیط تولید:**

```bash
make migrate-prod
```

---

## ساختار پروژه

ساختار پوشه‌های پروژه به صورت زیر است:

```
reality/
│
├── cmd/            # شامل فایل‌های اجرایی اصلی مانند http و chef
├── services/       # شامل تمام سرویس‌های پروژه
├── models/         # مدل‌های مربوط به پایگاه داده
├── migrations/     # فایل‌های مهاجرت پایگاه داده
└── .env.example    # مثال از فایل محیطی
```

---

## متغیرهای محیطی

### پایگاه داده (PostgreSQL)

- `DATABASE_PORT`: پورت پایگاه داده
- `DATABASE_NAME`: نام پایگاه داده
- `DATABASE_HOST`: آدرس پایگاه داده
- `DATABASE_USERNAME`: نام کاربری پایگاه داده
- `DATABASE_PASSWORD`: رمز عبور پایگاه داده
- `DATABASE_SSL_MODE`: وضعیت SSL پایگاه داده

### پیام‌رسانی (RabbitMQ)

- `RABBITMQ_USERNAME`: نام کاربری RabbitMQ
- `RABBITMQ_PASSWORD`: رمز عبور RabbitMQ
- `RABBITMQ_HOST`: آدرس میزبان RabbitMQ
- `RABBITMQ_PORT`: پورت RabbitMQ

### احراز هویت (JWT)

- `ACCESS_TOKEN_SIGNING_KEY`: کلید امضای توکن دسترسی
- `REFRESH_TOKEN_SIGNING_KEY`: کلید امضای توکن تجدید
- `ACCESS_TOKEN_EXPIRATION`: زمان انقضای توکن دسترسی
- `REFRESH_TOKEN_EXPIRATION`: زمان انقضای توکن تجدید

### ذخیره‌سازی (AWS S3)

- `AWS_S3_PRIVATE_BUCKET_NAME`: نام باکت خصوصی
- `AWS_S3_PUBLIC_BUCKET_NAME`: نام باکت عمومی
- `AWS_S3_HOST`: آدرس هاست S3
- `AWS_S3_PUBLIC_HOST_URI`: آدرس عمومی هاست

---

## سرویس‌های خارجی

### PayPing

پروژه از PayPing برای مدیریت پرداخت‌ها و تراکنش‌های مالی استفاده می‌کند. PayPing یک درگاه پرداخت آنلاین است که به کاربران اجازه می‌دهد به‌راحتی پرداخت‌های خود را از طریق کارت‌های بانکی ایرانی انجام دهند. این سرویس از طریق توکن امنیتی احراز هویت می‌شود که باید در متغیر محیطی `PAYPING_TOKEN` تعریف شود.

### Liara

برای میزبانی و مقیاس‌پذیری، پروژه از Liara استفاده می‌کند. Liara یک سرویس ابری است که امکان میزبانی برنامه‌های کاربردی با استفاده از زیرساخت‌های مقیاس‌پذیر را فراهم می‌کند. می‌توانید از Liara برای اجرای این پروژه به صورت آنلاین و بدون نگرانی از مدیریت زیرساخت‌ها استفاده کنید.

---

## راهنمایی مشارکت

برای مشارکت در این پروژه، موارد زیر را رعایت کنید:

1. مخزن را فورک کنید.
2. تغییرات خود را اعمال کرده و درخواست ادغام (Pull Request) ایجاد کنید.
3. حتماً کدهای خود را تست کرده و مستندات مربوطه را به‌روزرسانی کنید.

---

## مجوز

این پروژه تحت مجوز MIT منتشر شده است.

---

