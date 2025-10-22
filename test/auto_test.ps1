# Automated test for error display
Write-Host "Starting automated test..." -ForegroundColor Green

# Start sqlcmder in background
$process = Start-Process -FilePath ".\sqlcmder.exe" -PassThru -WindowStyle Hidden

Write-Host "SQLCmder started with PID: $($process.Id)"

# Wait for startup
Start-Sleep -Seconds 3

# Check if process is still running
if ($process.HasExited) {
    Write-Host "FAIL: Process exited during startup" -ForegroundColor Red
    exit 1
}

Write-Host "Process is running, waiting 5 seconds..." -ForegroundColor Yellow
Start-Sleep -Seconds 5

# Check if process is still running
if ($process.HasExited) {
    Write-Host "FAIL: Process crashed" -ForegroundColor Red
    exit 1
} else {
    Write-Host "SUCCESS: Process is stable" -ForegroundColor Green
}

# Cleanup
Stop-Process -Id $process.Id -Force
Write-Host "Test completed" -ForegroundColor Green

