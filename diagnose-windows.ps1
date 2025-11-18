# Jellyfin Windows Diagnostic Script
# Run this in PowerShell as Administrator on Windows

Write-Host "=====================================================" -ForegroundColor Cyan
Write-Host "  Jellyfin Windows Connection Diagnostic" -ForegroundColor Cyan
Write-Host "=====================================================" -ForegroundColor Cyan
Write-Host ""

# Check if Jellyfin process is running
Write-Host "[1] Checking if Jellyfin is running..." -ForegroundColor Yellow
$jellyfinProcess = Get-Process | Where-Object {$_.ProcessName -like "*jellyfin*"}
if ($jellyfinProcess) {
    Write-Host "    ✓ Jellyfin is running" -ForegroundColor Green
    Write-Host "      Process: $($jellyfinProcess.ProcessName) (PID: $($jellyfinProcess.Id))" -ForegroundColor Gray
} else {
    Write-Host "    ✗ Jellyfin is NOT running" -ForegroundColor Red
    Write-Host "      → Start Jellyfin and run this script again" -ForegroundColor Yellow
    exit
}

# Check what port 8096 is bound to
Write-Host ""
Write-Host "[2] Checking port 8096 binding..." -ForegroundColor Yellow
$port8096 = netstat -ano | Select-String ":8096"
if ($port8096) {
    Write-Host "    Port 8096 status:" -ForegroundColor Gray
    $port8096 | ForEach-Object {
        $line = $_.Line
        Write-Host "      $line" -ForegroundColor Gray

        # Check if it's bound to localhost only
        if ($line -match "127\.0\.0\.1:8096") {
            Write-Host "    ✗ PROBLEM: Jellyfin is bound to localhost only (127.0.0.1)" -ForegroundColor Red
            Write-Host "      → This means it only accepts connections from Windows itself" -ForegroundColor Yellow
            Write-Host "      → Fix: Go to Jellyfin Dashboard → Networking" -ForegroundColor Yellow
            Write-Host "      → Set 'Bind to local network address' to: 0.0.0.0" -ForegroundColor Yellow
            Write-Host "      → Then restart Jellyfin" -ForegroundColor Yellow
        }
        elseif ($line -match "0\.0\.0\.0:8096") {
            Write-Host "    ✓ Good: Jellyfin is bound to all interfaces (0.0.0.0)" -ForegroundColor Green
        }
    }
} else {
    Write-Host "    ✗ Port 8096 is not listening" -ForegroundColor Red
    Write-Host "      → Check Jellyfin configuration" -ForegroundColor Yellow
}

# Check Windows IP addresses
Write-Host ""
Write-Host "[3] Windows IP Addresses..." -ForegroundColor Yellow
$ipAddresses = Get-NetIPAddress -AddressFamily IPv4 | Where-Object {$_.IPAddress -notlike "127.*"}
foreach ($ip in $ipAddresses) {
    Write-Host "    $($ip.InterfaceAlias): $($ip.IPAddress)" -ForegroundColor Gray
}

# Check firewall rules for port 8096
Write-Host ""
Write-Host "[4] Checking firewall rules for port 8096..." -ForegroundColor Yellow
$firewallRules = Get-NetFirewallRule | Where-Object {
    $_.Enabled -eq $true -and
    $_.Direction -eq "Inbound"
} | ForEach-Object {
    $portFilter = $_ | Get-NetFirewallPortFilter
    if ($portFilter.LocalPort -contains "8096") {
        $_
    }
}

if ($firewallRules) {
    Write-Host "    ✓ Found firewall rules allowing port 8096:" -ForegroundColor Green
    foreach ($rule in $firewallRules) {
        Write-Host "      - $($rule.DisplayName)" -ForegroundColor Gray
        Write-Host "        Enabled: $($rule.Enabled), Action: $($rule.Action)" -ForegroundColor Gray
    }
} else {
    Write-Host "    ✗ No firewall rules found for port 8096" -ForegroundColor Red
    Write-Host "      → Creating firewall rule..." -ForegroundColor Yellow

    try {
        New-NetFirewallRule -DisplayName "Jellyfin WSL Access" `
            -Direction Inbound `
            -Protocol TCP `
            -LocalPort 8096 `
            -Action Allow `
            -Profile Any `
            -Enabled True
        Write-Host "    ✓ Firewall rule created successfully" -ForegroundColor Green
    } catch {
        Write-Host "    ✗ Failed to create firewall rule (run as Administrator)" -ForegroundColor Red
    }
}

# Test local access
Write-Host ""
Write-Host "[5] Testing local access to Jellyfin..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8096" -UseBasicParsing -TimeoutSec 3
    Write-Host "    ✓ Jellyfin responds on localhost:8096" -ForegroundColor Green
    Write-Host "      Status: $($response.StatusCode)" -ForegroundColor Gray
} catch {
    Write-Host "    ✗ Cannot access Jellyfin on localhost:8096" -ForegroundColor Red
    Write-Host "      → Is Jellyfin actually running?" -ForegroundColor Yellow
}

# Test access from WSL IP
Write-Host ""
Write-Host "[6] Testing access from network interface..." -ForegroundColor Yellow
$mainIP = (Get-NetIPAddress -AddressFamily IPv4 | Where-Object {$_.IPAddress -like "10.*" -or $_.IPAddress -like "192.168.*" -or $_.IPAddress -like "172.*"} | Select-Object -First 1).IPAddress
if ($mainIP) {
    Write-Host "    Testing: http://${mainIP}:8096" -ForegroundColor Gray
    try {
        $response = Invoke-WebRequest -Uri "http://${mainIP}:8096" -UseBasicParsing -TimeoutSec 3
        Write-Host "    ✓ Jellyfin responds on ${mainIP}:8096" -ForegroundColor Green
        Write-Host "      → WSL should be able to connect using this IP" -ForegroundColor Green
    } catch {
        Write-Host "    ✗ Cannot access Jellyfin on ${mainIP}:8096" -ForegroundColor Red
        Write-Host "      → Jellyfin is likely bound to localhost only" -ForegroundColor Yellow
    }
}

# WSL Integration Check
Write-Host ""
Write-Host "[7] WSL Integration Info..." -ForegroundColor Yellow
Write-Host "    Your bot (in WSL) will try to connect to:" -ForegroundColor Gray
Write-Host "      http://10.255.255.254:8096" -ForegroundColor Cyan
Write-Host ""
Write-Host "    To test from WSL, run:" -ForegroundColor Gray
Write-Host "      curl -I http://10.255.255.254:8096" -ForegroundColor Cyan

Write-Host ""
Write-Host "=====================================================" -ForegroundColor Cyan
Write-Host "  Summary" -ForegroundColor Cyan
Write-Host "=====================================================" -ForegroundColor Cyan

# Determine most likely issue
$issuefound = $false
if ($port8096 -match "127\.0\.0\.1:8096") {
    Write-Host ""
    Write-Host "LIKELY ISSUE: Jellyfin bound to localhost only" -ForegroundColor Red
    Write-Host ""
    Write-Host "FIX:" -ForegroundColor Yellow
    Write-Host "  1. Open Jellyfin: http://localhost:8096" -ForegroundColor White
    Write-Host "  2. Go to: Dashboard → Networking" -ForegroundColor White
    Write-Host "  3. Find: 'Bind to local network address'" -ForegroundColor White
    Write-Host "  4. Set to: 0.0.0.0 (or leave empty)" -ForegroundColor White
    Write-Host "  5. Save and restart Jellyfin" -ForegroundColor White
    Write-Host "  6. Run this script again to verify" -ForegroundColor White
    $issuefound = $true
}

if (-not $firewallRules -and -not $issuefound) {
    Write-Host ""
    Write-Host "LIKELY ISSUE: No firewall rule for port 8096" -ForegroundColor Red
    Write-Host "Run this script as Administrator to create the rule" -ForegroundColor Yellow
    $issuefound = $true
}

if (-not $issuefound -and $port8096 -match "0\.0\.0\.0:8096") {
    Write-Host ""
    Write-Host "Configuration looks good! Try testing from WSL:" -ForegroundColor Green
    Write-Host "  curl -I http://10.255.255.254:8096" -ForegroundColor Cyan
}

Write-Host ""
