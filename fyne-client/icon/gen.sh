# capitialse name for export
rm start_icon.go stop_icon.go meter_icon.go
fyne bundle -name Metericon -package icon meter.png > meter_icon.go
fyne bundle -name Starticon -package icon startbttn.png > start_icon.go
fyne bundle -name Stopicon -package icon stopbttn.png > stop_icon.go