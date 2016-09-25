#!/usr/bin/env powershell


$blacklight = "blacklight.exe"
Write-Host " -- building $blacklight"

$srcfiles = (Get-ChildItem -Path "src\*.go" -Recurse).FullName
go build -v -x -o "$blacklight" $srcfiles

if ( ! ($LastExitCode -eq 0)) {
  Write-Host " -- something went wrong!"
  Write-Host " -- go was unable to build $blacklight"
  exit 2
}

Get-ChildItem -Path "examples\*.bl" -exclude _* | % {
  Write-Host " -- running: " $_.FullName
    $output = (echo "test" | & .\$blacklight $_.FullName 2>&1)
    $blexit = $LastExitCode
    Write-Host $output
# FIXME: verify program output vs expected

    if ( ! ($blexit -eq 0)) {
      Write-Host " -- something went wrong!"
      exit 1
    }
}

rm $blacklight
Write-Host " -- done"
