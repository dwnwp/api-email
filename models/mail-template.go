package models

func CreateMailBodyTemplate(bodySubject, BodyContent string) string {
	template := "<strong>"+ bodySubject +"</strong>" + "<br><p>&emsp;"+ BodyContent +"</p>"
	img := "<img src='https://img.freepik.com/free-vector/thank-you-lettering_1262-6963.jpg?semt=ais_hybrid&w=740' alt='Thank You' width='100' height='100'>"
	footer := "<br><strong>Best Regards,</strong><br>"+ img +"<br>Trainee Backend<br>ขอบคุณครับ"
	template = template + footer
	return template
}