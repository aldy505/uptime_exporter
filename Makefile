build:
	rm -rf out
	mkdir out

	GOOS=darwin GOARCH=amd64 go build -o uptime_exporter .
	tar -czf out/uptime_exporter-darwin-amd64.tar.gz uptime_exporter LICENSE README.md
	rm uptime_exporter
	sha256sum out/uptime_exporter-darwin-amd64.tar.gz

	GOOS=darwin GOARCH=arm64 go build -o uptime_exporter .
	tar -czf out/uptime_exporter-darwin-arm64.tar.gz uptime_exporter LICENSE README.md
	rm uptime_exporter
	sha256sum out/uptime_exporter-darwin-arm64.tar.gz

	GOOS=linux GOARCH=386 go build -o uptime_exporter .
	tar -czf out/uptime_exporter-linux-386.tar.gz uptime_exporter LICENSE README.md
	rm uptime_exporter
	sha256sum out/uptime_exporter-linux-386.tar.gz

	GOOS=linux GOARCH=amd64 go build -o uptime_exporter .
	tar -czf out/uptime_exporter-linux-amd64.tar.gz uptime_exporter LICENSE README.md
	rm uptime_exporter
	sha256sum out/uptime_exporter-linux-amd64.tar.gz

	GOOS=linux GOARCH=arm go build -o uptime_exporter .
	tar -czf out/uptime_exporter-linux-arm.tar.gz uptime_exporter LICENSE README.md
	rm uptime_exporter
	sha256sum out/uptime_exporter-linux-arm.tar.gz

	GOOS=linux GOARCH=arm64 go build -o uptime_exporter .
	tar -czf out/uptime_exporter-linux-arm64.tar.gz uptime_exporter LICENSE README.md
	rm uptime_exporter
	sha256sum out/uptime_exporter-linux-arm64.tar.gz

	GOOS=windows GOARCH=386 go build -o uptime_exporter.exe .
	zip out/uptime_exporter-windows-386.zip uptime_exporter.exe LICENSE README.md
	rm uptime_exporter.exe
	sha256sum out/uptime_exporter-windows-386.zip

	GOOS=windows GOARCH=amd64 go build -o uptime_exporter.exe .
	zip out/uptime_exporter-windows-amd64.zip uptime_exporter.exe LICENSE README.md
	rm uptime_exporter.exe
	sha256sum out/uptime_exporter-windows-amd64.zip

	GOOS=windows GOARCH=arm go build -o uptime_exporter.exe .
	zip out/uptime_exporter-windows-arm.zip uptime_exporter.exe LICENSE README.md
	rm uptime_exporter.exe
	sha256sum out/uptime_exporter-windows-arm.zip