# Personal Storage Server

> A lightweight self-hosted Website to store files on a local machine, and safely be able to download, delete, upload the files across the internet.

[**Project Portfolio**](https://m7md2323.github.io/Portfolio/pages/personal_cloud_sotrage.html) | [**Report Bug**](https://github.com/m7md2323/Personal-Storage-Server/issues)

---
## 📖 Table of Contents

* [Tech Stack](#-tech-stack)
* [Features](#-features)
* [Project Structure](#-project-structure)
* [Getting Started](#-getting-started)
* [Screenshots](#-screenshots)
* [Contributing](#-contributing)
* [License](#-license)
* [Contact](#-contact)

---

## 🛠️ Tech Stack

### Backend, OS & Database
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-05122A?style=for-the-badge&logo=go)
![SQLite](https://img.shields.io/badge/sqlite-%2307405e.svg?style=for-the-badge&logo=sqlite&logoColor=white)
![GORM](https://img.shields.io/badge/GORM-blueviolet?style=for-the-badge)
![Debian](https://img.shields.io/badge/Debian-A81D33?style=for-the-badge&logo=debian&logoColor=white)

> **Note:** This system is designed to be hosted locally on a **Linux (Debian 13)** machine to serve as a private cloud.

### Frontend
![JavaScript](https://img.shields.io/badge/javascript-%23323330.svg?style=for-the-badge&logo=javascript&logoColor=%23F7DF1E)
![HTML5](https://img.shields.io/badge/html5-%23E34F26.svg?style=for-the-badge&logo=html5&logoColor=white)
![CSS3](https://img.shields.io/badge/css3-%231572B6.svg?style=for-the-badge&logo=css3&logoColor=white)

### Networking & Deployment
![Tailscale](https://img.shields.io/badge/Tailscale-%23FFFFFF.svg?style=for-the-badge&logo=tailscale&logoColor=black)

---

## ✨ Features
* **Security:**: Middlewares to make sure no untrusted devices can access the website, and login Authentication.
* **UI:**: Simple and easy to use user interface, a Gallery-like page to display the user Images.
* **Files Management:** The abilty to Upload/Delete/Retrieve all kind of files, every user has his own Directory.
* **Secure Remote Access:**: Avaliable through Tailscale(serve) service, a Mesh peer-to-peer network.

---

## 📂 Project Structure

.
├── backend/
│   ├── database/         # SQLite connection & GORM initialization
│   ├── handlers/         # Gin route handlers (Upload, Download, Delete)
│   ├── middlewares/      # Auth & device trust logic
│   └── models/           # Data structures & schemas
├── frontend/             # Web interface
│   ├── css/              # Styling
│   ├── js/               # Client-side logic
│   └── index.html        # Main dashboard
├── main.go               # Application entry point
├── go.mod                # Go dependencies
├── go.sum                # Cryptographic checksums
├── .gitignore            # Excludes sensitive data & binaries
└── README.md             # Project documentation

---

## 🚀 Getting Started

### Prerequisites
* Go (1.21+)
* Git
* Tailscale (Required if you want to access the server outside your local network)

### Installation & Setup

1. **Clone the repository**
   ```bash
   git clone https://github.com/m7md2323/Personal-Storage-Server.git
   cd Personal-Storage-Server
   ```

2. **Install dependencies**
   ```bash
   go mod tidy
   ```

3. **Configure Environment**

* Database Configuration (SQLite)
DATABASE_FILE_PATH="YOUR_DATABASE_FILE_PATH"

* Storage Configuration
UPLOADS="YOUR_FILES_UPLOAD_PATH"


4. **Run the server**
   ```bash
   go run main.go
   ```
   
---

## 📸 Screenshots

| Dashboard View|
![Dashboard](https://m7md2323.github.io/Portfolio/images/upload.PNG)

---

## 🤝 Contributing

1. Fork the Project
2. Create your Feature Branch (`git checkout -b feature/AmazingFeature`)
3. Commit your Changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the Branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## ⚖️ License

Distributed under the no License.

---

## 📬 Contact

**Mohammad K. Al Harahsheh** - [GitHub Profile](https://github.com/m7md2323) | **Email** [Send Email](mohammadalharahsheh04@gmail.com)
