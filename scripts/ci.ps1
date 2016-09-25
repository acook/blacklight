#!/usr/bin/env powershell

Write-Host " -- OHAI"
Write-Host " -- Meesa build yousa blacklight!"

$blacklight = "blacklight.exe"
$srcfiles = (Get-ChildItem -Path "src\*.go" -Recurse).FullName
go build -v -x -o "$blacklight" $srcfiles

Write-Host " -- Weesa running blacklight programs now!"

Get-ChildItem -Path "examples\*.bl" -exclude _* | % {
  Write-Host " -- Running: " $_.FullName
  $output = (echo "test" | & .\$blacklight $_.FullName 2>&1)
  $blexit = $LastExitCode
  Write-Host $output
  # FIXME: verify program output vs expected

  if ( ! ($blexit -eq 0)) {
    Write-Host " -- Encountered error!"
    exit 1
  }
}

rm $blacklight
Write-Host " -- All done!"
