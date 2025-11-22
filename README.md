# JWT Authentication dan RBAC Implementation

## Setup Dependencies

Jalankan perintah berikut untuk menginstall dependencies yang diperlukan:

\`\`\`bash
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
\`\`\`

## Database Setup

1. Jalankan script SQL untuk membuat tabel users:
\`\`\`sql
-- File: scripts/001_create_users_table.sql
-- Tabel Users untuk autentikasi
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    role VARCHAR(20) DEFAULT 'user' CHECK (role IN ('admin', 'user')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert sample users (password: "123456")
INSERT INTO users (username, email, password_hash, role) VALUES 
('admin', 'admin@university.com', 
'$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'admin'),
('user1', 'user1@university.com', 
'$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi', 'user')
ON CONFLICT (username) DO NOTHING;
\`\`\`

## API Endpoints dengan RBAC

### Public Endpoints
- `POST /api/login` - Login untuk mendapatkan JWT token

### Protected Endpoints (Memerlukan Authentication)
- `GET /api/profile` - Melihat profile user yang sedang login

### Alumni Endpoints
- `GET /api/alumni` - Ambil semua data alumni (Admin dan User)
- `GET /api/alumni/:id` - Ambil data alumni berdasarkan ID (Admin dan User)
- `POST /api/alumni` - Tambah alumni baru (Admin Only)
- `PUT /api/alumni/:id` - Update data alumni (Admin Only)
- `DELETE /api/alumni/:id` - Hapus data alumni (Admin Only)

### Pekerjaan Alumni Endpoints
- `GET /api/pekerjaan` - Ambil semua data pekerjaan alumni (Admin dan User)
- `GET /api/pekerjaan/:id` - Ambil data pekerjaan berdasarkan ID (Admin dan User)
- `GET /api/pekerjaan/alumni/:alumni_id` - Ambil semua pekerjaan berdasarkan alumni (Admin Only)
- `POST /api/pekerjaan` - Tambah pekerjaan baru (Admin Only)
- `PUT /api/pekerjaan/:id` - Update data pekerjaan (Admin Only)
- `DELETE /api/pekerjaan/:id` - Hapus data pekerjaan (Admin Only)

## Testing dengan Postman

### 1. Login untuk mendapatkan token
\`\`\`http
POST http://localhost:3000/api/login
Content-Type: application/json

{
  "username": "admin",
  "password": "123456"
}
\`\`\`

Response:
\`\`\`json
{
  "success": true,
  "message": "Login berhasil",
  "data": {
    "user": {
      "id": 1,
      "username": "admin",
      "email": "admin@university.com",
      "role": "admin"
    },
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
  }
}
\`\`\`

### 2. Menggunakan token untuk akses protected endpoints
\`\`\`http
GET http://localhost:3000/api/alumni
Authorization: Bearer YOUR_TOKEN_HERE
\`\`\`

### 3. Akses profile
\`\`\`http
GET http://localhost:3000/api/profile
Authorization: Bearer YOUR_TOKEN_HERE
\`\`\`

### 4. Operasi Admin (Create, Update, Delete)
\`\`\`http
POST http://localhost:3000/api/alumni
Authorization: Bearer ADMIN_TOKEN_HERE
Content-Type: application/json

{
  "nim": "2023001",
  "nama": "Test Alumni",
  "jurusan": "Informatika",
  "angkatan": 2023,
  "tahun_lulus": 2027,
  "email": "test@email.com"
}
\`\`\`

## Sample Users
- **Admin**: username: `admin`, password: `123456`
- **User**: username: `user1`, password: `123456`

## Security Features
- Password hashing dengan bcrypt
- JWT token dengan expiration (24 jam)
- Role-based access control (RBAC)
- Input validation
- Logging untuk audit trail
