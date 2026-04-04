# 🌐 Personal Storage Server

> A lightweight Website to store files on a local machine, and safely be able to download, delete, upload the files across the internet.

[**Project Portfolio**](https://m7md2323.github.io/Portfolio/pages/personal_cloud_sotrage.html) | [**Report Bug**](https://github.com/m7md2323/Personal-Storage-Server/issues)

---

## 📖 Table of Contents

* [Features](#-features)
* [Tech Stack](#-tech-stack)
* [Getting Started](#-getting-started)
* [Screenshots](#-screenshots)
* [Contributing](#-contributing)
* [License](#-license)
* [Contact](#-contact)

---

## 🛠️ Tech Stack

### Backend & Database
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-05122A?style=for-the-badge&logo=go)
![SQLite](https://img.shields.io/badge/sqlite-%2307405e.svg?style=for-the-badge&logo=sqlite&logoColor=white)
![GORM](https://img.shields.io/badge/GORM-blueviolet?style=for-the-badge)

### Frontend
![JavaScript](https://img.shields.io/badge/javascript-%23323330.svg?style=for-the-badge&logo=javascript&logoColor=%23F7DF1E)
![HTML5](https://img.shields.io/badge/html5-%23E34F26.svg?style=for-the-badge&logo=html5&logoColor=white)
![CSS3](https://img.shields.io/badge/css3-%231572B6.svg?style=for-the-badge&logo=css3&logoColor=white)

### Networking & Deployment
![Tailscale](https://img.shields.io/badge/Tailscale-%23FFFFFF.svg?style=for-the-badge&logo=tailscale&logoColor=black)

---

## ✨ Features

* **Full File Lifecycle:** Seamlessly upload, download, and delete files through a clean web interface.
* **Secure Remote Access:** Access your files from anywhere via a private mesh VPN using Tailscale.
* **High Performance:** Built with a Go-based backend for fast file processing and low memory footprint.
* **Persistent Storage:** Uses SQLite and GORM to manage file metadata and organization efficiently.
* **Minimalist UI:** A responsive vanilla frontend designed for ease of use without heavy frameworks.

---

## 🚀 Getting Started

### Prerequisites
* Go (1.21+)
* Git
* Tailscale (Required if you want to access the server outside your local network)

### Installation & Setup

1. **Clone the repository**
   ```bash
   git clone [https://github.com/m7md2323/Personal-Storage-Server.git](https://github.com/m7md2323/Personal-Storage-Server.git)
   cd Personal-Storage-Server
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure Environment**
   Ensure your storage directory is set up or defined in the configuration.

4. **Run the server**
   ```bash
   go run main.go
   ```

---

## 📸 Screenshots

| Dashboard View | File Management |
| :--- | :--- |
| ![Dashboard](https://via.placeholder.com/400x250?text=Storage+Dashboard) | ![Files](https://via.placeholder.com/400x250?text=File+List+View) |

---

## 🤝 Contributing

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## ⚖️ License

Distributed under the **MIT License**. See `LICENSE` for more information.

---

## 📬 Contact

**Mohammad K. Al Harahsheh** - [GitHub Profile](https://github.com/m7md2323)

Project Link: [https://github.com/m7md2323/Personal-Storage-Server](https://github.com/m7md2323/Personal-Storage-Server)
