package models

func CreateMailBodyTemplate(bodySubject, BodyContent string) string {
	template := "<strong>"+ bodySubject +"</strong>" + "<br><p>&emsp;"+ BodyContent +"</p>"
	img := "<img src='https://img.freepik.com/premium-vector/thank-you-caligraphic-vector-illustration-concept_1253202-40937.jpg?semt=ais_hybrid&w=100' alt='Thank You'>"
	footer := "<br><strong>Best Regards,</strong><br>"+ img +"<br>Trainee Backend<br>ขอบคุณครับ"
	return template + footer
}