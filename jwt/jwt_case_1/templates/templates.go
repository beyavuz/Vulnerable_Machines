package templates

const ProfilePageHTML = `

<!DOCTYPE html>
<html>
<head>
    <title>Profile</title>
    <style>
        body { font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif; margin: 50px; background-color: #f4f4f4; }
        .card { background: white; padding: 30px; border-radius: 10px; max-width: 500px; margin: 0 auto; box-shadow: 0 4px 15px rgba(0,0,0,0.1); }
        h1 { margin-top: 0; color: #333; }
        
        .label { font-weight: bold; color: #555; display: block; margin-top: 15px; }
        .value { color: #000; display: block; margin-bottom: 5px; font-size: 1.1em; }
        
        /* Rozet Stilleri */
        .badge { padding: 4px 10px; border-radius: 15px; font-size: 0.85em; font-weight: bold; color: white; }
        .bg-green { background-color: #28a745; }
        .bg-red { background-color: #dc3545; }

        /* Buton Alanı Stilleri */
        .action-area { margin-top: 30px; padding-top: 20px; border-top: 1px solid #eee; display: flex; justify-content: space-between; align-items: center; }
        
        /* Genel Buton Stili */
        .btn { padding: 10px 20px; border-radius: 5px; text-decoration: none; font-size: 0.9em; border: none; cursor: pointer; transition: 0.3s; }
        
        /* Logout Linki */
        .link-logout { color: #dc3545; font-weight: bold; }
        .link-logout:hover { text-decoration: underline; }

        /* Admin Butonu (Aktif) */
        .btn-admin { background-color: #007bff; color: white; display: inline-block; }
        .btn-admin:hover { background-color: #0056b3; box-shadow: 0 2px 5px rgba(0,0,0,0.2); }

        /* Admin Butonu (Pasif/Disabled) */
        .btn-disabled { background-color: #e2e6ea; color: #6c757d; cursor: not-allowed; border: 1px solid #dae0e5; }
    </style>
</head>
<body>
    <div class="card">
        <h1>Welcome, {{.Username}}</h1>
        
        <div>
            <span class="label">User ID:</span>
            <span class="value">#{{.ID}}</span>
        </div>

        <div>
            <span class="label">Rol:</span>
            <span class="value">
                {{if eq .AccountRole "admin"}}
                    <span class="badge bg-red">Admin</span>
                {{else}}
                    <span class="badge bg-green">User</span>
                {{end}}
            </span>
        </div>

        <div>
            <span class="label">Token Created At (IAT):</span>
            <span class="value" style="font-family: monospace; font-size: 0.9em;">{{.RegisteredClaims.IssuedAt}}</span>
        </div>

        <div class="action-area">
            {{if eq .AccountRole "admin"}}
                <a href="/admin" class="btn btn-admin">Admin</a>
            {{else}}
                <button class="btn btn-disabled" disabled>Admin</button>
            {{end}}

            <a href="/logout" class="link-logout">Logout</a>
        </div>
    </div>
</body>
</html>
`

const LoginPageHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Login</title>
    <style>
        body { font-family: sans-serif; display: flex; justify-content: center; margin-top: 50px; background-color: #f0f2f5; }
        .container { background: white; padding: 30px; border-radius: 8px; box-shadow: 0 2px 10px rgba(0,0,0,0.1); width: 300px; }
        input { width: 100%; padding: 10px; margin: 10px 0; box-sizing: border-box; }
        button { width: 100%; padding: 10px; background: #007bff; color: white; border: none; border-radius: 4px; cursor: pointer; }
        button:hover { background: #0056b3; }
        .alert { color: red; font-size: 0.9em; text-align: center; }
    </style>
</head>
<body>
    <div class="container">
        <h2 style="text-align:center;">Login Page</h2>
        {{if .Error}}<div class="alert">{{.Error}}</div>{{end}}
        <form action="/login" method="POST">
            <input type="text" name="username" placeholder="Username" required>
            <input type="password" name="password" placeholder="Password" required>
            <button type="submit">Login</button>
        </form>
		<p style="font-size:0.8em; text-align:center; color:#666;">User => user:user123 </p>
    </div>
</body>
</html>
`

const AdminPageHTML = `
<!DOCTYPE html>
<html>
<head>
    <title>Admin Panel</title>
    <style>
        body { background-color: #1a1a1a; font-family: 'Courier New', Courier, monospace; color: #00ff00; display: flex; justify-content: center; align-items: center; height: 100vh; margin: 0; }
        .secret-card { border: 2px solid #00ff00; padding: 40px; max-width: 600px; text-align: center; box-shadow: 0 0 20px rgba(0, 255, 0, 0.2); background: #000; }
        h1 { font-size: 2em; margin-bottom: 10px; border-bottom: 1px dashed #00ff00; padding-bottom: 10px; text-transform: uppercase; letter-spacing: 3px; }
        .warning { color: #ff3333; font-weight: bold; margin-bottom: 30px; font-size: 0.9em; }
        
        .flag-container { margin: 30px 0; padding: 20px; border: 1px dotted #00ff00; background-color: #0a0a0a; }
        .flag-label { display: block; font-size: 0.8em; margin-bottom: 10px; color: #888; }
        .flag-text { font-size: 1.5em; font-weight: bold; word-break: break-all; color: #fff; text-shadow: 0 0 5px #fff; }
        
        .btn-back { display: inline-block; margin-top: 20px; color: #000; background-color: #00ff00; padding: 10px 20px; text-decoration: none; font-weight: bold; transition: 0.3s; }
        .btn-back:hover { background-color: #fff; box-shadow: 0 0 15px #fff; }
    </style>
</head>
<body>
    <div class="secret-card">
        <h1>Admin Page</h1>

        <div class="flag-container">
            <span class="flag-label">CTF FLAG:</span>
            <span class="flag-text">{{.Flag}}</span>
        </div>

        <a href="/profile" class="btn-back">Profile</a>
    </div>
</body>
</html>
`
