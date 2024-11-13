package async

// 发起周报请求
func SendWeekReport() error {
	cronSpec := "*/2 * * * *"
	body := map[string]interface{}{
		"type":     "week_report_msg",
		"cronSpec": cronSpec,
		"ttl":      1,
		"body":     map[string]interface{}{},
		"retry":    0,
	}
	err := SendPeriodTask(body)
	if err != nil {
		return err
	}
	return nil
}
