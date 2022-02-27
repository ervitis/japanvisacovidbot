package email

type (
	Message struct {
		ContentType string
		Body        string
	}
)

const (
	ContentType = "text/plain"
	Charset     = "UTF-8"

	layoutDateEmailSend = "2006-01-02 15:04:05 MST"

	MessageFormat = `%s: %s\r\n`

	locationEmail = "Europe/Madrid"
)

var (
	MessageBody = `%s,
Me llamo Victor Martin Alonso. El motivo de mi correo es pedir el visado de trabajo de Japón.

Tengo en posesión los documentos necesarios que aparecen en la página web de la embajada. 

- DNI
- CoE enviado por mi empresa
- Pasaporte en vigor
- Fotografías
- Aplicación para el visado
- Aplicación registrada en el sistema Entrants, Returnees Follow-up System (ERFS) patrocinado por mi empresa

¿Cuales serían los siguientes pasos para obtener el visado?

Un saludo y muchas gracias.
`

	MessageConfirmation = `
Before sending the email, check the following list of documentation to prepare:
- Valid passport with two blank pages
- Valid NIF
- Photo face
- Filled the visa documents %s or %s
- Certificate of Eligibility, check the expiration date
- Application for ERFS system issued by the company

Remember that it's possible that the embassy would ask for more papers and documentation
`

	PriorityHeaders = []struct {
		Header string
		Value  string
	}{
		{
			Header: "X-Priority", Value: "1 (Highest)",
		},
		{
			Header: "X-MSMail-Priority", Value: "High",
		},
		{
			Header: "Importance", Value: "High",
		},
	}
)
