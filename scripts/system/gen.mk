codegen-update:
	@codegen-cli upload-http-server --service auth-id --source docs/openapi/auth-id.yaml
	@codegen-cli gen-http-server --service auth-id
