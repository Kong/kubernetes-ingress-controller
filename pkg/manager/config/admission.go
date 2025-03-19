package config

// AdmissionServerConfig defines parameters that configure the Admission Server run by the controller.
type AdmissionServerConfig struct {
	ListenAddr string

	CertPath string
	Cert     string

	KeyPath string
	Key     string
}
