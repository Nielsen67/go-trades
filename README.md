# API Manajemen Inventaris

Dokumentasi ini akan memandu Anda dalam menyiapkan dan menjalankan API Manajemen Inventaris.

## Requirement

Sebelum memulai, pastikan Anda telah menginstal perangkat lunak berikut:

*   **VS Code:** Editor kode yang direkomendasikan. Anda dapat mengunduhnya di [https://code.visualstudio.com/](https://code.visualstudio.com/).
*   **MySQL:** Sistem manajemen basis data yang akan digunakan. Anda dapat mengunduhnya di [https://www.mysql.com/](https://www.mysql.com/).  Pastikan MySQL server sudah berjalan sebelum menjalankan aplikasi.
*   **Go:** Bahasa pemrograman yang digunakan untuk membangun API. Anda dapat mengunduhnya di [https://go.dev/](https://go.dev/). Pastikan Go sudah di set di *environment variables*

## Installation

Berikut adalah langkah-langkah untuk menginstal dan menyiapkan proyek:

1.  **Repository Clone**

    ```bash
    git clone https://github.com/Nielsen67/go-trades
    ```

    Perintah ini akan mengunduh kode proyek dari repositori GitHub ke komputer Anda.  Pastikan `git` sudah terinstall di komputer anda.

2.  **Project Directory**

    ```bash
    cd inventory-management
    ```

    Perintah ini akan membawa Anda ke direktori proyek yang baru saja diunduh.

3.  **Configuration `.env`**

    ```bash
    cp .env.example .env
    ```

    Perintah ini akan menyalin file contoh konfigurasi `.env.example` menjadi file `.env`. File `.env` digunakan untuk menyimpan variabel enviroment termasuk koneksi database.

    Buka file `.env` lalu sesuaikan konfigurasi sesuai dengan pengaturan MySQL Anda. Berikut adalah contoh isinya:

    ```
    DB_USER=your_mysql_username  # Ganti dengan username MySQL Anda
    DB_PASSWORD=your_mysql_password # Ganti dengan password MySQL Anda
    DB_NAME=inventory_db       # Ganti dengan nama database yang akan Anda gunakan (buat dulu di MySQL)
    DB_HOST=localhost          # Biasanya localhost jika MySQL berjalan di komputer yang sama
    DB_PORT=3306             # Port default MySQL, sesuaikan jika berbeda
    ```

    *   **Penting:** Buat database dengan nama `go_trades` (atau nama yang Anda tentukan di `.env`) di MySQL Anda sebelum menjalankan aplikasi. Anda bisa menggunakan tools seperti phpMyAdmin, MySQL Workbench, atau command line MySQL untuk membuat database.

4.  **Dependencies Installation**

    ```bash
    go mod tidy
    ```

    Perintah ini akan mengunduh dan menginstal semua dependensi Go yang dibutuhkan oleh proyek. Dependensi ini didefinisikan dalam file `go.mod`.

## Running Application

Setelah konfigurasi selesai, Anda dapat menjalankan aplikasi dengan perintah berikut:

```bash
go run .

```
## API Guides
Dokumentasi Request Body, Response dan Endpoint terdapat pada Collection https://documenter.getpostman.com/view/37120004/2sB2izDYaK
