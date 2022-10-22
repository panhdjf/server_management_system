# sever_management_system
1. Khởi tạo file app.env giống trong file app.txt
2. Khởi tạo container docker :
    - Mở docker
    - Dùng lệnh docker-compose up -d
3. Kết nối Database (PostgresSQL/PgAdmin4)
    - Run file migrate.go
    - Mở PgAdmin4, tạo một server mới với các thông tin giống trong app.env :
        + Home/Address: localhost
        + Port: 6500
        + Maintainance DB: golang-gorm
        + Username: Postgres
4. Chạy chương trình:
    - Run lệnh "air"
5. Sử dụng Postman để test:
<img image.png>


